//coverage:ignore file
package dispute

import (
	"api/auditLogger"
	"api/env"
	"api/handlers/notifications"
	"api/handlers/ticket"
	"api/handlers/workflow"
	mediatorassignment "api/mediatorAssignment"
	"api/middleware"
	"api/models"
	"api/utilities"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type DisputeModel interface {
	UploadWriteup(userId, disputeId int64, path string, file io.Reader) error
	UploadEvidence(userId, disputeId int64, path string, file io.Reader) (uint, error)
	GetEvidenceByDispute(disputeId int64) ([]models.Evidence, error)
	GetDisputeExperts(disputeId int64) ([]models.Expert, error)

	GetAdminDisputes(searchTerm *string, limit *int, offset *int, sort *models.Sort, filters *[]models.Filter, dateFilter *models.DateFilter) ([]models.AdminDisputeSummariesResponse, int64, error)
	GetDisputesByUser(userId int64) ([]models.Dispute, error)
	GetDispute(disputeId int64) (models.Dispute, error)

	GetUserByEmail(email string) (models.User, error)
	GetUserById(userId int64) (models.User, error)
	CreateDispute(dispute models.Dispute) (int64, error)
	UpdateDisputeStatus(disputeId int64, status string) error

	ObjectExpert(disputeId, expertId, ticketId int64) error
	ReviewExpertObjection(objectionId int64, approved models.ExpObjStatus) error
	GetDisputeIDByTicketID(ticketID int64) (int64, error)
	GetExpertRejections(expertID, disputeID *int64, limit, offset *int) ([]models.ExpertObjectionsView, error)

	CreateDefaultUser(email string, fullName string, pass string) error
	AssignExpertsToDispute(disputeID int64) ([]models.User, error)
	AssignExpertswithDisputeAndExpertIDs(disputeID int64, expertIDs []int) error

	GetWorkflowRecordByID(id uint64) (*models.Workflow, error)
	CreateActiverWorkflow(workflow *models.ActiveWorkflows) (int, error)
	DeleteActiveWorkflow(workflow *models.ActiveWorkflows) error

	GetExperts(disputeID int64) ([]models.AdminDisputeExperts, error)
	GenerateAISummary(disputeID int64, disputeDesc string, apiKey string)

	GetUser(UserID int64) (models.UserDetails, error)
	GetAdminDisputeDetails(disputeId int64) (models.AdminDisputeDetailsResponse, error)
}

type Dispute struct {
	Model              DisputeModel
	TicketModel        ticket.TicketModel
	Email              notifications.EmailSystem
	JWT                middleware.Jwt
	Env                env.Env
	AuditLogger        auditLogger.DisputeProceedingsLoggerInterface
	OrchestratorEntity WorkflowOrchestrator
	MediatorAssignment mediatorassignment.AlgorithmAssignment
	WorkflowModel      workflow.WorkflowDBModel
}

type OrchestratorRequest struct {
	ID int64 `json:"id"`
}

type Trigger struct {
	Id      int64
	Trigger string
}

type WorkflowOrchestrator interface {
	MakeRequestToOrchestrator(endpoint string, payload OrchestratorRequest) (string, error)
	SendTriggerToOrchestrator(endpoint string, activerWfId int64, trigger string) (int, error)
}

type OrchestratorReal struct {
}

func (w OrchestratorReal) SendTriggerToOrchestrator(endpoint string, activerWfId int64, trigger string) (int, error) {
	logger := utilities.NewLogger().LogWithCaller()
	payload := Trigger{
		Id:      activerWfId,
		Trigger: trigger,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("marshal error: ", err)
		return http.StatusInternalServerError, err
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Error("post error: ", err)
		return http.StatusInternalServerError, err
	}
	//read body
	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.Error("read body error: ", err)
		return http.StatusInternalServerError, err
	}

	// log the response body for debugging
	logger.Info("Response Body: ", string(responseBody))

	return resp.StatusCode, nil
}

func (w OrchestratorReal) MakeRequestToOrchestrator(endpoint string, payload OrchestratorRequest) (string, error) {
	logger := utilities.NewLogger().LogWithCaller()

	// Marshal the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("marshal error: ", err)
		return "", fmt.Errorf("internal server error")
	}
	logger.Info("Payload: ", string(payloadBytes))

	// Send the POST request to the orchestrator
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Error("post error: ", err)
		return "", fmt.Errorf("internal server error")
	}
	defer resp.Body.Close()

	// Check for a successful status code (200 OK)

	if resp.StatusCode == http.StatusInternalServerError {
		logger.Error("status code error: ", resp.StatusCode)
		return "", fmt.Errorf("Check that you gave the correct state name if resetting")
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("status code error: ", resp.StatusCode)
		rsponseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("read body error: ", err)
			return "", fmt.Errorf("internal server error")
		}

		return string(rsponseBody), fmt.Errorf("internal server error")
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read body error: ", err)
		return "", fmt.Errorf("internal server error")
	}

	// Convert the response body to a string
	responseBody := string(bodyBytes)

	// Log the response body for debugging
	logger.Info("Response Body: ", responseBody)

	return responseBody, nil
}

type disputeModelReal struct {
	db  *gorm.DB
	env env.Env
}

func NewHandler(db *gorm.DB, envReader env.Env) Dispute {

	return Dispute{
		Email:              notifications.NewHandler(db),
		JWT:                middleware.NewJwtMiddleware(),
		Env:                env.NewEnvLoader(),
		Model:              &disputeModelReal{db: db, env: env.NewEnvLoader()},
		AuditLogger:        auditLogger.NewDisputeProceedingsLogger(db, envReader),
		OrchestratorEntity: OrchestratorReal{},
		MediatorAssignment: mediatorassignment.DefaultAlorithmAssignment(db),
		TicketModel:        ticket.NetTicketModelReal(db, envReader),
		WorkflowModel:      &workflow.WorkflowModelReal{DB: db},
	}
}

func (m *disputeModelReal) GetUserById(userId int64) (models.User, error) {
	logger := utilities.NewLogger().LogWithCaller()
	user := models.User{}
	err := m.db.Where("\"id\"= ?", userId).First(&user).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving user")
	}
	return user, err
}

func (m *disputeModelReal) UploadWriteup(userId, disputeId int64, path string, file io.Reader) error {
	env := env.NewEnvLoader()
	fileStorageRoot, err := env.Get("FILESTORAGE_ROOT")
	if err != nil {
		return err
	}
	fileStorageUrl, err := env.Get("FILESTORAGE_URL")
	if err != nil {
		return err
	}

	logger := utilities.NewLogger().LogWithCaller()
	dir, name := filepath.Split(path)
	dest := filepath.Join(fileStorageRoot, path)

	if err := os.MkdirAll(filepath.Join(fileStorageRoot, dir), 0755); err != nil {
		logger.WithError(err).Error("failed to create folder for file upload")
		return err
	}

	url := fmt.Sprintf("%s/%s", fileStorageUrl, path)

	// Open the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		logger.WithError(err).Errorf("failed to create file in storage: '%s'", dest)
		return errors.New("failed to create file in storage")
	}
	defer destFile.Close()

	// Copy file content to destination
	_, err = io.Copy(destFile, file)
	if err != nil {
		logger.WithError(err).Error("failed to copy file content")
		return errors.New("failed to copy file content")
	}

	// Add file entry to Database
	fileRow := models.File{
		FileName: name,
		FilePath: url,
		Uploaded: time.Now(),
	}

	if err := m.db.Create(&fileRow).Error; err != nil {
		logger.WithError(err).Error("error adding file to database")
		return errors.New("error adding file to database")
	}

	//add entry to dispute decisions table
	disputeDecision := models.DisputeDecisions{
		DisputeID: disputeId,
		ExpertID:  userId,
		WriteUpID: int64(*fileRow.ID),
	}

	if err := m.db.Create(&disputeDecision).Error; err != nil {
		logger.WithError(err).Error("error creating dispute decision")
		return errors.New("error creating dispute decision")
	}

	return nil
}

func (m *disputeModelReal) GetWorkflowRecordByID(id uint64) (*models.Workflow, error) {
	logger := utilities.NewLogger().LogWithCaller()
	workflow := models.Workflow{}
	err := m.db.Where("id = ?", id).First(&workflow).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving workflow record")
		return nil, err
	}
	return &workflow, nil
}

func (m *disputeModelReal) CreateActiverWorkflow(workflow *models.ActiveWorkflows) (int, error) {
	result := m.db.Create(workflow)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(workflow.ID), nil
}

func (m *disputeModelReal) DeleteActiveWorkflow(workflow *models.ActiveWorkflows) error {
	result := m.db.Delete(workflow)

	if result.Error != nil {
		return result.Error
	}

	return nil

}

func (m *disputeModelReal) UploadEvidence(userId, disputeId int64, path string, file io.Reader) (uint, error) {
	env := env.NewEnvLoader()
	fileStorageRoot, err := env.Get("FILESTORAGE_ROOT")
	if err != nil {
		return 0, err
	}
	fileStorageUrl, err := env.Get("FILESTORAGE_URL")
	if err != nil {
		return 0, err
	}

	logger := utilities.NewLogger().LogWithCaller()

	dir, name := filepath.Split(path)
	dest := filepath.Join(fileStorageRoot, path)

	if err := os.MkdirAll(filepath.Join(fileStorageRoot, dir), 0755); err != nil {
		logger.WithError(err).Error("failed to create folder for file upload")
		return 0, err
	}

	url := fmt.Sprintf("%s/%s", fileStorageUrl, path)

	// Open the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		logger.WithError(err).Errorf("failed to create file in storage: '%s'", dest)
		return 0, errors.New("failed to create file in storage")
	}
	defer destFile.Close()

	// Copy file content to destination
	_, err = io.Copy(destFile, file)
	if err != nil {
		logger.WithError(err).Error("failed to copy file content")
		return 0, errors.New("failed to copy file content")
	}

	// Add file entry to Database
	fileRow := models.File{
		FileName: name,
		FilePath: url,
		Uploaded: time.Now(),
	}

	if err := m.db.Create(&fileRow).Error; err != nil {
		logger.WithError(err).Error("error adding file to database")
		return 0, errors.New("error adding file to database")
	}

	//add entry to dispute evidence table
	disputeEvidence := models.DisputeEvidence{
		Dispute: disputeId,
		FileID:  int64(*fileRow.ID),
		UserID:  userId,
	}

	if err := m.db.Create(&disputeEvidence).Error; err != nil {
		logger.WithError(err).Error("error creating dispute evidence")
		return 0, errors.New("error creating dispute evidence")
	}

	return *fileRow.ID, nil
}

func (m *disputeModelReal) GetUser(userID int64) (models.UserDetails, error) {
	logger := utilities.NewLogger().LogWithCaller()
	var user models.UserDetails
	err := m.db.Raw("SELECT CONCAT(u.first_name, ' ', u.surname) AS full_name, u.email, CONCAT(a.street, ' ',a.street2,' ',a.street3, CHR(10), a.city, CHR(10), a.province, CHR(10), a.country)  AS address FROM users u JOIN addresses a ON a.id = u.id WHERE u.id = ?", userID).Scan(&user).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving user")
		return user, err
	}
	return user, nil
}

func (m *disputeModelReal) GetDispute(disputeId int64) (dispute models.Dispute, err error) {
	logger := utilities.NewLogger().LogWithCaller()
	err = m.db.Model(&models.Dispute{}).Where("id = ?", disputeId).First(&dispute).Error

	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute")
	}
	return dispute, err
}

func (m *disputeModelReal) GetAdminDisputeDetails(disputeId int64) (models.AdminDisputeDetailsResponse, error) {

	logger := utilities.NewLogger().LogWithCaller()
	dispute, err := m.GetDispute(disputeId)
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute")
		return models.AdminDisputeDetailsResponse{}, err
	}

	var adminIntermed models.AdminIntermediate
	err = m.db.Raw("SELECT id, title, status, case_date, date_resolved FROM disputes WHERE id = ?", disputeId).Scan(&adminIntermed).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute")
		return models.AdminDisputeDetailsResponse{}, err
	}
	var workflow models.WorkflowResp
	err = m.db.Raw("SELECT wf.id, wf.name FROM disputes d JOIN active_workflows aw ON d.workflow = aw.id JOIN workflows wf ON wf.id = aw.workflow WHERE d.id = ?", adminIntermed.Id).First(&workflow).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving workflow for dispute with ID: " + strconv.Itoa(int(adminIntermed.Id)))
		return models.AdminDisputeDetailsResponse{}, err
	}

	experts, err := m.GetExperts(disputeId)
	if err != nil {
		logger.WithError(err).Error("Error retrieving experts for dispute with ID: " + strconv.Itoa(int(adminIntermed.Id)))
		return models.AdminDisputeDetailsResponse{}, err
	}

	var evidence []models.Evidence
	evidence, err = m.GetEvidenceByDispute(disputeId)
	if err != nil {
		logger.WithError(err).Error("Error retrieving evidence for dispute with ID: " + strconv.Itoa(int(adminIntermed.Id)))
		return models.AdminDisputeDetailsResponse{}, err
	}

	var complainant models.UserDetails
	complainant, err = m.GetUser(dispute.Complainant)
	if err != nil {
		logger.WithError(err).Error("Error retrieving complainant for dispute with ID: " + strconv.Itoa(int(adminIntermed.Id)))
		return models.AdminDisputeDetailsResponse{}, err
	}

	var respondent models.UserDetails
	respondent, err = m.GetUser(*dispute.Respondant)
	if err != nil {
		logger.WithError(err).Error("Error retrieving respondent for dispute with ID: " + strconv.Itoa(int(adminIntermed.Id)))
		return models.AdminDisputeDetailsResponse{}, err
	}

	adminDisputeDetails := models.AdminDisputeDetailsResponse{
		AdminDisputeSummariesResponse: models.AdminDisputeSummariesResponse{
			Id:        strconv.Itoa(int(adminIntermed.Id)),
			Title:     adminIntermed.Title,
			Status:    adminIntermed.Status,
			Workflow:  workflow,
			DateFiled: adminIntermed.CaseDate.Format("2006-01-02"),
			Experts:   experts,
		},
		Description: dispute.Description,
		Evidence:    evidence,
		Complainant: complainant,
		Respondent:  respondent,
	}
	if dispute.DateResolved != nil {
		dateResolved := dispute.DateResolved.Format("2006-01-02")
		adminDisputeDetails.DateResolved = &dateResolved
	}
	if adminDisputeDetails.Experts == nil {
		adminDisputeDetails.Experts = []models.AdminDisputeExperts{}
	}
	if adminDisputeDetails.Evidence == nil {
		adminDisputeDetails.Evidence = []models.Evidence{}
	}

	return adminDisputeDetails, nil
}

func (m *disputeModelReal) GetEvidenceByDispute(disputeId int64) (evidence []models.Evidence, err error) {
	logger := utilities.NewLogger().LogWithCaller()
	err = m.db.Table("dispute_evidence").Select(`files.id, file_name, uploaded, file_path,  CASE
            WHEN disputes.complainant = dispute_evidence.user_id THEN 'Complainant'
            WHEN disputes.respondant = dispute_evidence.user_id THEN 'Respondent'
            ELSE 'Other'
        END AS uploader_role`).Joins("JOIN files ON dispute_evidence.file_id = files.id").Joins("JOIN disputes ON dispute_evidence.dispute = disputes.id").Where("dispute = ?", disputeId).Find(&evidence).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute evidence")
		return
	}
	if evidence == nil {
		evidence = []models.Evidence{}
	}
	return evidence, err
}
func (m *disputeModelReal) GetDisputeExperts(disputeId int64) (experts []models.Expert, err error) {
	logger := utilities.NewLogger().LogWithCaller()
	err = m.db.Table("dispute_experts_view").
		Select("users.id, users.first_name || ' ' || users.surname AS full_name, email, users.phone_number AS phone, role").
		Joins("JOIN users ON dispute_experts_view.expert = users.id").
		Where("dispute = ?", disputeId).
		Where("dispute_experts_view.status = 'Approved'").
		Where("role = 'Mediator' OR role = 'Arbitrator' OR role = 'Conciliator' OR role = 'expert'").
		Find(&experts).Error

	if err != nil && err.Error() != "record not found" {
		logger.WithError(err).Error("Error retrieving dispute experts")
		return
	}

	if experts == nil {
		experts = []models.Expert{}
	}
	return experts, err
}

func (m *disputeModelReal) GetDisputesByUser(userId int64) (disputes []models.Dispute, err error) {
	logger := utilities.NewLogger().LogWithCaller()
	logger.Infof("Starting GetDisputesByUser with userId: %d", userId)

	err = m.db.Where("complainant = ? OR respondant = ?", userId, userId).Find(&disputes).Error
	if err != nil {
		logger.WithError(err).Errorf("Failed to find disputes of user with ID %d", userId)
		return nil, err
	}

	if len(disputes) == 0 {
		logger.Infof("No disputes found for user with ID %d, checking expert disputes", userId)
		var disputesExpert []models.DisputeExpert
		err = m.db.Raw("SELECT * FROM public.dispute_experts WHERE public.dispute_experts.user = ?", userId).Scan(&disputesExpert).Error
		if err != nil {
			logger.WithError(err).Error("Failed to find expert with disputes")
			return nil, err
		}

		for _, disputeExpert := range disputesExpert {
			dispute := models.Dispute{}
			err1 := m.db.Where("id = ?", disputeExpert.Dispute).First(&dispute).Error
			if err1 != nil {
				logger.WithError(err1).Errorf("Failed to find dispute with expert and id: %d", disputeExpert.Dispute)
				continue
			}
			disputes = append(disputes, dispute)
		}
		logger.Infof("Disputes after processing experts: %+v", disputes)
	}

	return disputes, err
}

func (m *disputeModelReal) GetUserByEmail(email string) (user models.User, err error) {
	logger := utilities.NewLogger().LogWithCaller()
	err = m.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		logger.WithError(err).Errorf("Failed to find user with email '%s'", email)
	}
	return user, err
}
func (m *disputeModelReal) CreateDispute(dispute models.Dispute) (int64, error) {
	logger := utilities.NewLogger().LogWithCaller()
	disputeCloned := dispute
	if err := m.db.Create(&disputeCloned).Error; err != nil {
		logger.WithError(err).Error("Error creating dispute in database")
		return 0, err
	}
	return *disputeCloned.ID, nil
}
func (m *disputeModelReal) UpdateDisputeStatus(disputeId int64, status string) error {
	logger := utilities.NewLogger().LogWithCaller()
	if status != "Closed" {
		err := m.db.Model(&models.Dispute{}).Where("id = ?", disputeId).Update("status", status).Error
		if err != nil {
			logger.WithError(err).Errorf("Failed to update dispute (ID = %d) status to '%s'", disputeId, status)
		}
		return err
	} else {
		additionalFields := map[string]interface{}{
			"status":   status,
			"resolved": true,
		}
		err := m.db.Model(&models.Dispute{}).Where("id = ?", disputeId).Updates(additionalFields).Error
		if err != nil {
			logger.WithError(err).Errorf("Failed to update dispute (ID = %d) status to '%s'", disputeId, status)
		}
		var dbDispute models.Dispute
		err2 := m.db.Model(&models.Dispute{}).Where("id = ?", disputeId).First(&dbDispute).Error
		if err2 != nil {
			logger.WithError(err2).Error("There was an error trying to retrieve the dispute related to the summary")
		}
		apiKey, err4 := m.env.Get("OPENAI_KEY")
		if err4 != nil {
			logger.WithError(err4).Error("Something went wrong getting the API key.")
		}
		go m.GenerateAISummary(disputeId, dbDispute.Description, apiKey)
		return err
	}
}
func (m *disputeModelReal) ObjectExpert(disputeId, expertId, ticketId int64) error {
	logger := utilities.NewLogger().LogWithCaller()

	//add entry to expert objections table
	expertObjection := models.ExpertObjection{
		ExpertID: expertId,
		TicketID: ticketId,
		Status:   models.ObjectionReview,
	}

	if err := m.db.Create(&expertObjection).Error; err != nil {
		logger.WithError(err).Error("Error creating expert objection")
		return err
	}
	return nil
}
func (m *disputeModelReal) ReviewExpertObjection(rejectionID int64, approved models.ExpObjStatus) error {
	logger := utilities.NewLogger().LogWithCaller()

	var expertObjections models.ExpertObjection
	logger.Infof("Rejection ID: %d", rejectionID)
	if err := m.db.Where("\"id\" = ?", rejectionID).First(&expertObjections).Error; err != nil {
		logger.WithError(err).Error("Error retrieving expert objections")
		return err
	}

	// Update status
	expertObjections.Status = approved

	if err := m.db.Save(&expertObjections).Error; err != nil {
		logger.WithError(err).Error("Error updating expert objections")
		return err
	}

	return nil
}

func (m *disputeModelReal) CreateDefaultUser(email string, fullName string, pass string) error {
	logger := utilities.NewLogger().LogWithCaller()
	nameSplit := strings.Split(fullName, " ")
	//stub timezone
	zone, _ := time.Now().Zone()
	timezone := zone
	actualTimezone := &timezone
	//Now put stuff in the actual user object
	date, _ := time.Parse("2006-01-02", time.Now().String())
	user := models.User{
		FirstName:         nameSplit[0],
		Surname:           nameSplit[1],
		Birthdate:         date,
		Nationality:       "",
		Email:             email,
		PasswordHash:      pass,
		PhoneNumber:       nil,
		AddressID:         nil,
		Status:            "Unverified",
		Gender:            "Other",
		PreferredLanguage: nil,
		Timezone:          actualTimezone,
	}
	//create a default entry for the user
	//Hash the password
	hash, salt, err := utilities.HashPassword(user.PasswordHash)
	if err != nil {
		logger.WithError(err).Error("Error hashing the default password.")
		return err
	}

	user.PasswordHash = base64.StdEncoding.EncodeToString(hash)
	user.Salt = base64.StdEncoding.EncodeToString(salt)

	//update log metrics
	user.Status = "Active"
	user.Role = "user"

	if result := m.db.Create(&user); result.Error != nil {
		logger.WithError(result.Error).Error("Error creating default user")
		return result.Error
	}
	logger.Info("User added to the Database")
	return nil
}

// bandaid fix, will be removed in future

func (m *disputeModelReal) AssignExpertswithDisputeAndExpertIDs(disputeID int64, expertIDs []int) error {
	logger := utilities.NewLogger().LogWithCaller()
	for _, expertID := range expertIDs {
		if err := m.db.Exec("INSERT INTO dispute_experts_view VALUES (?, ?)", disputeID, expertID).Error; err != nil {
			logger.WithError(err).Error("Error inserting expert into dispute_experts table")
			return err
		}
	}
	return nil
}

func (m disputeModelReal) AssignExpertsToDispute(disputeID int64) ([]models.User, error) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Define the roles to select from
	roles := []string{"Mediator", "Adjudicator", "Arbitrator", "expert"}

	// Query for users with the specified roles
	var users []models.User
	if err := m.db.Where("role IN ?", roles).Find(&users).Error; err != nil {
		return nil, err
	}

	// Shuffle the results and take the first 4
	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })

	// Select the first 4 users after shuffle
	selectedUsers := users
	if len(users) > 4 {
		selectedUsers = users[:4]
	}

	// Insert the selected experts into the dispute_experts table
	for _, expert := range selectedUsers {

		// A raw query is used here because GORM tries to insert the Status field of the struct,
		// despite the field being marked as read-only. Why did we use an ORM?
		if err := m.db.Exec("INSERT INTO dispute_experts_view VALUES (?, ?)", disputeID, expert.ID).Error; err != nil {
			return nil, err
		}
	}

	return selectedUsers, nil
}

func (m *disputeModelReal) GenerateAISummary(disputeID int64, disputeDesc string, apiKey string) {
	logger := utilities.NewLogger().LogWithCaller()

	// Define the messages
	messages := []map[string]string{
		{
			"role":    "system",
			"content": "You are an AI agent specialized in Alternative Dispute Resolution. Your role is to generate concise and informative summaries of resolved dispute cases for archival purposes. These summaries will help future users understand the nature of the disputes, the evidence presented, the domain in which the dispute occurred, and the final outcome. Your summaries should focus on key details, such as: Type of Dispute: Clearly identify the nature of the dispute (e.g., contract disagreement, service complaint, intellectual property issue). Domain: Specify the context or industry relevant to the dispute (e.g., e-commerce, real estate, software development). Your summaries should be clear, neutral, about 200 words in length and useful for guiding future decisions and actions related to similar disputes. Provide all output as plaintext and in paragraph form such that it looks valid in an archive section.",
		},
		{
			"role":    "user",
			"content": disputeDesc,
		},
	}

	// Create the request payload
	payload := map[string]interface{}{
		"model":    "gpt-4-turbo",
		"user":     "dre1",
		"messages": messages,
	}

	// Convert payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		logger.WithError(err).Error("There was an error converting payload to JSON.")
		return
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadJSON))
	if err != nil {
		logger.WithError(err).Error("Something went wrong creating the request.")
		return
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Create a custom Transport with TLSClientConfig to skip verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Send the request using an http.Client with the custom Transport
	client := &http.Client{Transport: tr}
	response, err := client.Do(req)
	if err != nil {
		logger.WithError(err).Error("Something went wrong sending the request.")
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.WithError(err).Error("Something went wrong reading the response body.")
		return
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		logger.WithError(err).Error("There was an error parsing the response.")
		return
	}
	prettyJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
	if err != nil {
		logger.WithError(err).Error("Failed to marshal JSON response.")
	}

	choices, ok := jsonResponse["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		logger.Info("JSON response: " + string(prettyJSON))
		logger.Error("Unexpected response format or empty choices.")
		return
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		logger.Error("Unexpected response format for first choice.")
		return
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		logger.Error("Unexpected response format for message.")
		return
	}

	content, ok := message["content"].(string)
	if !ok {
		logger.Error("Unexpected response format for content.")
		return
	}

	archiveSummary := models.DisputeSummaries{
		ID:      disputeID,
		Summary: content,
	}
	if err = m.db.Create(&archiveSummary).Error; err != nil {
		logger.WithError(err).Error("Error inserting the summary.")
		return
	}
}

func (m *disputeModelReal) GetExperts(disputeID int64) ([]models.AdminDisputeExperts, error) {
	logger := utilities.NewLogger().LogWithCaller()
	var expert []models.AdminDisputeExperts
	err := m.db.Raw("SELECT u.id, CONCAT(u.first_name,' ' ,u.surname) AS full_name, de.status FROM users u JOIN dispute_experts_view de ON u.id = de.expert WHERE de.dispute = ?", disputeID).Scan(&expert).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving experts")
		return expert, err
	}
	return expert, nil
}

func (m *disputeModelReal) GetAdminDisputes(searchTerm *string, limit *int, offset *int, sort *models.Sort, filters *[]models.Filter, dateFilter *models.DateFilter) ([]models.AdminDisputeSummariesResponse, int64, error) {
	logger := utilities.NewLogger().LogWithCaller()
	var disputes []models.AdminDisputeSummariesResponse = []models.AdminDisputeSummariesResponse{}
	var queryString strings.Builder
	var countString strings.Builder
	var countParams []interface{}
	var queryParams []interface{}

	queryString.WriteString("SELECT id, title, status, case_date, date_resolved FROM disputes")
	countString.WriteString("SELECT COUNT(*) FROM disputes")
	if searchTerm != nil {
		queryString.WriteString(" WHERE disputes.title LIKE ?")
		countString.WriteString(" WHERE disputes.title LIKE ?")
		queryParams = append(queryParams, "%"+*searchTerm+"%")
		countParams = append(countParams, "%"+*searchTerm+"%")
	}

	if filters != nil && len(*filters) > 0 {
		if searchTerm != nil {
			queryString.WriteString(" AND ")
			countString.WriteString(" AND ")
		} else {
			queryString.WriteString(" WHERE ")
			countString.WriteString(" WHERE ")
		}
		for i, filter := range *filters {
			queryString.WriteString(filter.Attr + " = ?")
			countString.WriteString(filter.Attr + " = ?")
			queryParams = append(queryParams, filter.Value)
			countParams = append(countParams, filter.Value)
			if i < len(*filters)-1 {
				queryString.WriteString(" AND ")
				countString.WriteString(" AND ")
			}
		}
	}

	if dateFilter != nil {
		if searchTerm != nil || (filters != nil && len(*filters) > 0) {
			queryString.WriteString(" AND ")
			countString.WriteString(" AND ")
		} else {
			queryString.WriteString(" WHERE ")
			countString.WriteString(" WHERE ")
		}
		if dateFilter.Filed != nil {
			if dateFilter.Filed.Before != nil {
				queryString.WriteString("disputes.case_date < ? ")
				countString.WriteString("disputes.case_date < ? ")
				queryParams = append(queryParams, *dateFilter.Filed.Before)
				countParams = append(countParams, *dateFilter.Filed.Before)
			}
			if dateFilter.Filed.After != nil {
				if dateFilter.Filed.Before != nil {
					queryString.WriteString("AND ")
					countString.WriteString("AND ")
				}
				queryString.WriteString("disputes.case_date > ? ")
				countString.WriteString("disputes.case_date > ? ")
				queryParams = append(queryParams, *dateFilter.Filed.After)
				countParams = append(countParams, *dateFilter.Filed.After)
			}
		}
		if dateFilter.Resolved != nil {
			if dateFilter.Resolved.Before != nil {
				if dateFilter.Filed != nil {
					queryString.WriteString("AND ")
					countString.WriteString("AND ")
				}
				queryString.WriteString("disputes.date_resolved < ? ")
				countString.WriteString("disputes.date_resolved < ? ")
				queryParams = append(queryParams, *dateFilter.Resolved.Before)
				countParams = append(countParams, *dateFilter.Resolved.Before)
			}
			if dateFilter.Resolved.After != nil {
				if dateFilter.Filed != nil || dateFilter.Resolved.Before != nil {
					queryString.WriteString("AND ")
					countString.WriteString("AND ")
				}
				queryString.WriteString("disputes.date_resolved > ? ")
				countString.WriteString("disputes.date_resolved > ? ")
				queryParams = append(queryParams, *dateFilter.Resolved.After)
				countParams = append(countParams, *dateFilter.Resolved.After)
			}
		}
	}

	if sort != nil {
		validSortAttrs := map[string]bool{
			"id":            true,
			"case_date":     true,
			"workflow":      true,
			"status":        true,
			"title":         true,
			"description":   true,
			"complainant":   true,
			"respondant":    true,
			"date_resolved": true,
		}

		if _, valid := validSortAttrs[sort.Attr]; !valid {
			return disputes, 0, errors.New("invalid sort attribute")
		}

		if sort.Order != "asc" && sort.Order != "desc" {
			sort.Order = "asc"
		}

		queryString.WriteString(" ORDER BY " + sort.Attr + " " + sort.Order)
	}
	if limit != nil {
		queryString.WriteString(" LIMIT ?")
		queryParams = append(queryParams, *limit)
	}
	if offset != nil {
		queryString.WriteString(" OFFSET ?")
		queryParams = append(queryParams, *offset)
	}

	//get relevant disputes data
	var intermediateDisputes []models.AdminIntermediate
	err := m.db.Raw(queryString.String(), queryParams...).Scan(&intermediateDisputes).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving disputes")
		return disputes, 0, err
	}

	var countRows int64 = 0
	err = m.db.Raw(countString.String(), countParams...).Scan(&countRows).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute count")
		return disputes, 0, err
	}

	//take the intermediate disputes and fill in the admin response information
	for _, dispute := range intermediateDisputes {
		var disputeResp models.AdminDisputeSummariesResponse
		disputeResp.Id = strconv.Itoa(int(dispute.Id))
		disputeResp.Title = dispute.Title
		disputeResp.Status = dispute.Status
		disputeResp.DateFiled = dispute.CaseDate.Format("2006-01-02")
		//get the workflow
		var workflow models.WorkflowResp

		err = m.db.Raw("SELECT wf.id, wf.name FROM disputes d JOIN active_workflows aw ON d.workflow = aw.id JOIN workflows wf ON wf.id = aw.workflow WHERE d.id = ?", dispute.Id).First(&workflow).Error
		if err != nil {
			logger.WithError(err).Error("Error retrieving workflow for dispute with ID: " + strconv.Itoa(int(dispute.Id)))
		}

		experts, err := m.GetExperts(dispute.Id)
		if err != nil {
			logger.WithError(err).Error("Error retrieving experts for dispute with ID: " + strconv.Itoa(int(dispute.Id)))
		}

		disputeResp.Workflow = workflow
		disputeResp.Experts = experts
		//get the date resolved
		if dispute.DateResolved != nil {
			dateResolved := dispute.DateResolved.Format("2006-01-02")
			disputeResp.DateResolved = &dateResolved
		}

		disputes = append(disputes, disputeResp)
	}

	return disputes, countRows, err
}

func (m *disputeModelReal) GetExpertRejections(expertID, disputeID *int64, limit, offset *int) ([]models.ExpertObjectionsView, error) {
	logger := utilities.NewLogger().LogWithCaller()
	var rejections []models.ExpertObjectionsView = []models.ExpertObjectionsView{}
	var queryString strings.Builder
	var queryParams []interface{}

	queryString.WriteString("SELECT * FROM expert_objections_view")
	if expertID != nil || disputeID != nil {
		queryString.WriteString(" WHERE ")
		if expertID != nil {
			queryString.WriteString("expert_id = ?")
			queryParams = append(queryParams, *expertID)
		}
		if disputeID != nil {
			if expertID != nil {
				queryString.WriteString(" AND ")
			}
			queryString.WriteString("dispute_id = ?")
			queryParams = append(queryParams, *disputeID)
		}
	}

	if limit != nil {
		queryString.WriteString(" LIMIT ?")
		queryParams = append(queryParams, *limit)
	}
	if offset != nil {
		queryString.WriteString(" OFFSET ?")
		queryParams = append(queryParams, *offset)
	}
	fmt.Println(queryString.String())
	err := m.db.Raw(queryString.String(), queryParams...).Scan(&rejections).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving expert rejections")
		return rejections, err
	}

	return rejections, err
}

func (m *disputeModelReal) GetDisputeIDByTicketID(ticketID int64) (int64, error) {
	logger := utilities.NewLogger().LogWithCaller()
	var disputeID int64
	err := m.db.Raw(`SELECT 
	t.dispute_id
FROM 
	expert_objections eo
JOIN 
	tickets t ON eo.ticket_id = t.id
WHERE 
	eo.id = ?`, ticketID).Scan(&disputeID).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute ID by ticket ID")
		return 0, err
	}
	return disputeID, err
}
