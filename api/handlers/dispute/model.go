package dispute

import (
	"api/auditLogger"
	"api/env"
	"api/handlers/notifications"
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
	"strings"
	"time"

	"gorm.io/gorm"
)

type DisputeModel interface {
	UploadEvidence(userId, disputeId int64, path string, file io.Reader) (uint, error)
	GetEvidenceByDispute(disputeId int64) ([]models.Evidence, error)
	GetDisputeExperts(disputeId int64) ([]models.Expert, error)

	GetDisputesByUser(userId int64) ([]models.Dispute, error)
	GetDispute(disputeId int64) (models.Dispute, error)

	GetUserByEmail(email string) (models.User, error)
	CreateDispute(dispute models.Dispute) (int64, error)
	UpdateDisputeStatus(disputeId int64, status string) error

	ObjectExpert(userId, disputeId, expertId int64, reason string) error
	ReviewExpertObjection(userId, disputeId, expertId int64, approved bool) error

	CreateDefaultUser(email string, fullName string, pass string) error
	AssignExpertsToDispute(disputeID int64) ([]models.User, error)

	GenerateAISummary(disputeID int64, disputeDesc string, apiKey string)
}

type Dispute struct {
	Model       DisputeModel
	Email       notifications.EmailSystem
	JWT         middleware.Jwt
	Env         env.Env
	AuditLogger auditLogger.DisputeProceedingsLoggerInterface
}
type disputeModelReal struct {
	db  *gorm.DB
	env env.Env
}

func NewHandler(db *gorm.DB, envReader env.Env) Dispute {
	return Dispute{
		Email:       notifications.NewHandler(db),
		JWT:         middleware.NewJwtMiddleware(),
		Env:         env.NewEnvLoader(),
		Model:       &disputeModelReal{db: db, env: env.NewEnvLoader()},
		AuditLogger: auditLogger.NewDisputeProceedingsLogger(db, envReader),
	}
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

func (m *disputeModelReal) GetDispute(disputeId int64) (dispute models.Dispute, err error) {
	logger := utilities.NewLogger().LogWithCaller()
	err = m.db.Model(&models.Dispute{}).Where("id = ?", disputeId).First(&dispute).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute")
	}
	return dispute, err
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
	err = m.db.Table("dispute_experts").Select("users.id, users.first_name || ' ' || users.surname AS full_name, email, users.phone_number AS phone, role").Joins("JOIN users ON dispute_experts.user = users.id").Where("dispute = ?", disputeId).Where("dispute_experts.status = 'Approved'").Where("role = 'Mediator' OR role = 'Arbitrator' OR role = 'Conciliator' OR role = 'expert'").Find(&experts).Error
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
func (m *disputeModelReal) ObjectExpert(userId, disputeId, expertId int64, reason string) error {
	logger := utilities.NewLogger().LogWithCaller()

	//add entry to expert objections table
	expertObjection := models.ExpertObjection{
		DisputeID: disputeId,
		ExpertID:  expertId,
		UserID:    userId,
		Reason:    reason,
	}

	if err := m.db.Create(&expertObjection).Error; err != nil {
		logger.WithError(err).Error("Error creating expert objection")
		return err
	}
	return nil
}
func (m *disputeModelReal) ReviewExpertObjection(userId, disputeId, expertId int64, approved bool) error {
	logger := utilities.NewLogger().LogWithCaller()

	var expertObjections models.ExpertObjection
	if err := m.db.Where("dispute_id = ? AND expert_id = ? AND status = ?", disputeId, expertId, models.ReviewStatus).First(&expertObjections).Error; err != nil {
		logger.WithError(err).Error("Error retrieving expert objections")
		return err
	}

	// Update status
	if approved {
		expertObjections.Status = models.Sustained
	} else {
		expertObjections.Status = models.Overruled
	}

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
		if err := m.db.Create(&models.DisputeExpert{
			Dispute: disputeID,
			User:    expert.ID,
		}).Error; err != nil {
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
