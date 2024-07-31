package handlers

import (
	"api/env"
	"api/middleware"
	"api/models"
	"api/utilities"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
}

type EmailSystem interface {
	SendAdminEmail(c *gin.Context, disputeID int64, resEmail string)
	NotifyDisputeStateChanged(c *gin.Context, disputeID int64, disputeStatus string)
}

type Dispute struct {
	Model DisputeModel
	Email EmailSystem
}

type disputeModelReal struct {
	db gorm.DB
}

func NewDisputeHandler(db gorm.DB) Dispute {
	return Dispute{
		Model: &disputeModelReal{db: db},
	}
}

func SetupDisputeRoutes(g *gin.RouterGroup, h Dispute) {
	//dispute routes
	g.Use(middleware.JWTMiddleware)

	g.GET("", h.getSummaryListOfDisputes)
	g.POST("/create", h.createDispute)
	g.GET("/:id", h.getDispute)

	g.POST("/:id/experts/reject", h.expertObjection)
	g.POST("/:id/experts/review-rejection", h.expertObjectionsReview)
	g.POST("/:id/evidence", h.uploadEvidence)
	g.PUT("/dispute/status", h.updateStatus)

	//patch is not to be integrated yet
	// disputeRouter.HandleFunc("/{id}", h.patchDispute).Methods(http.MethodPatch)

	//create dispute

	//archive routes
}

// @Summary Get a summary list of disputes
// @Description Get a summary list of disputes
// @Tags dispute
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Dispute Summary Endpoint"
// @Router /dispute/:id/evidence [post]
func (h Dispute) uploadEvidence(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	claims := middleware.GetClaims(c)
	if claims == nil {
		logger.Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	disputeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.WithError(err).Error("Invalid Dispute ID")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Dispute ID"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logger.WithError(err).Error("Failed to parse form data")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Failed to parse form data"})
		return
	}

	files := form.File["files"]
	folder := fmt.Sprintf("%d", disputeId)
	for _, fileHeader := range files {
		path := filepath.Join(folder, fileHeader.Filename)
		file, err := fileHeader.Open()
		if err != nil {
			logger.WithError(err).Error("error opening file")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "error opening file"})
			return
		}

		_, err = h.Model.UploadEvidence(claims.User.ID, int64(disputeId), path, file)
		if err != nil {
			logger.WithError(err).Error("Error uploading evidence")
			c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
			return
		}
	}
	logger.Info("Evidence uploaded successfully")
	c.JSON(http.StatusCreated, models.Response{
		Data: "Files uploaded",
	})
}

// @Summary Get a summary list of disputes
// @Description Get a summary list of disputes
// @Tags dispute
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Dispute Summary Endpoint"
// @Router /dispute [get]
func (h Dispute) getSummaryListOfDisputes(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	jwtClaims := middleware.GetClaims(c)
	userID := jwtClaims.User.ID

	disputes, err := h.Model.GetDisputesByUser(userID)
	if err != nil {
		logger.WithError(err).Error("error retrieving disputes")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error while retrieving disputes"})
		return
	}

	summaries := make([]models.DisputeSummaryResponse, len(disputes))
	for i, dispute := range disputes {
		var role string = ""
		if dispute.Complainant == userID {
			role = "Complainant"
		} else if *(dispute.Respondant) == userID {
			role = "Respondant"
		}
		summaries[i] = models.DisputeSummaryResponse{
			ID:          *dispute.ID,
			Title:       dispute.Title,
			Description: dispute.Description,
			Status:      dispute.Status,
			Role:        &role,
		}
	}
	logger.Info("Dispute summaries retrieved successfully")
	c.JSON(http.StatusOK, models.Response{Data: summaries})
}

// @Summary Get a dispute
// @Description Get a dispute
// @Tags dispute
// @Accept json
// @Produce json
// @Param id path string true "Dispute ID"
// @Success 200 {object} models.Response "Dispute Detail Endpoint"
// @Router /dispute/{id} [get]
func (h Dispute) getDispute(c *gin.Context) {

	logger := utilities.NewLogger().LogWithCaller()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Error: fmt.Sprintf("Invalid dispute id '%s'", idParam)})
		return
	}

	dispute, err := h.Model.GetDispute(int64(id))
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute")
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}

	jwtClaims := middleware.GetClaims(c)
	userId := jwtClaims.User.ID
	role := ""
	//name and email
	// var respondantData models.User
	// err = h.DB.Where("id = ?", disputes.Respondant).Scan(&respondantData).Error
	// if err!=nil {

	// }

	if userId == dispute.Complainant {
		role = "Complainant"
	} else if userId == *(dispute.Respondant) {
		role = "Respondent"
	}

	DisputeDetailsResponse := models.DisputeDetailsResponse{
		ID:          *dispute.ID,
		Title:       dispute.Title,
		Description: dispute.Description,
		Status:      dispute.Status,
		DateCreated: dispute.CaseDate,
		Role:        role,
	}

	evidence, err := h.Model.GetEvidenceByDispute(int64(id))
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute evidence")
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	if evidence == nil {
		evidence = []models.Evidence{}
	}

	experts, err := h.Model.GetDisputeExperts(int64(id))
	if err != nil && err.Error() != "record not found" {
		logger.WithError(err).Error("Error retrieving dispute experts")
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}

	if experts == nil {
		experts = []models.Expert{}
	}

	DisputeDetailsResponse.Evidence = evidence
	DisputeDetailsResponse.Experts = experts

	logger.Info("Dispute details retrieved successfully")
	c.JSON(http.StatusOK, models.Response{Data: DisputeDetailsResponse})
	// c.JSON(http.StatusOK, models.Response{Data: "Dispute Detail Endpoint for ID: " + id})
}

func (h Dispute) createDispute(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	claims := middleware.GetClaims(c)
	if claims == nil {
		logger.Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logger.WithError(err).Error("Error parsing form")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Failed to parse form data"})
		return
	}

	// Access form values
	title := form.Value["title"][0]
	description := form.Value["description"][0]
	fullName := form.Value["respondent[full_name]"][0]
	email := form.Value["respondent[email]"][0]
	// telephone := form.Value["respondent[telephone]"][0]

	//get complainants id
	complainantID := claims.User.ID

	//check if respondant is in database by email and phone number
	var respondantID *int64
	respondent, err := h.Model.GetUserByEmail(email)
	if err != nil && err.Error() == "record not found" {
		//create a default entry for the user
		nameSplit := strings.Split(fullName, " ")
		if len(nameSplit) < 2 {
			logger.Error("Invalid full name")
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid full name"})
			return
		} else {
			logger.WithError(err).Error("Error retrieving respondent")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving respondent"})
			return
		}
	} else if err != nil {
		logger.WithError(err).Error("Error retrieving respondent")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving respondent"})
		return
	} else {
		respondantID = &respondent.ID
	}

	//create entry into the dispute table
	disputeId, err := h.Model.CreateDispute(models.Dispute{
		Title:       title,
		CaseDate:    time.Now(),
		Workflow:    nil,
		Status:      "Awaiting Respondant",
		Description: description,
		Complainant: complainantID,
		Respondant:  respondantID,
		Resolved:    false,
		Decision:    models.Unresolved,
	})
	if err != nil {
		logger.WithError(err).Error("Error creating dispute")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating dispute"})
		return
	}

	// Store files in Docker and retrieve URLs
	files := form.File["files"]
	folder := fmt.Sprintf("%d", disputeId)
	for _, fileHeader := range files {
		path := filepath.Join(folder, fileHeader.Filename)
		file, err := fileHeader.Open()
		if err != nil {
			logger.WithError(err).Error("failed to open file")
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong."})
			return
		}

		_, err = h.Model.UploadEvidence(complainantID, disputeId, path, file)
		if err != nil {
			logger.WithError(err).Error("failed to upload evidence")
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong."})
			return
		}
	}

	// Respond with success message
	go h.Email.SendAdminEmail(c, disputeId, email)
	logger.Info("Admin email sent")
	c.JSON(http.StatusCreated, models.Response{Data: "Dispute created successfully"})
	logger.Info("Dispute created successfully: ", title)
}

func (h Dispute) updateStatus(c *gin.Context) {
	var disputeStatus models.DisputeStatusChange
	logger := utilities.NewLogger().LogWithCaller()
	if err := c.BindJSON(&disputeStatus); err != nil {
		logger.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.Model.UpdateDisputeStatus(disputeStatus.DisputeID, disputeStatus.Status)
	if err != nil {
		logger.WithError(err).Error("failed to update dispute status")
		utilities.InternalError(c)
		return
	}
	go h.Email.NotifyDisputeStateChanged(c, disputeStatus.DisputeID, disputeStatus.Status)

	logger.Info("Dispute status updated successfully")
	c.JSON(http.StatusOK, models.Response{Data: "Dispute status update successful"})
}

// @Summary Update a dispute
// @Description Update a dispute
// @Tags dispute
// @Accept json
// @Produce json
// @Param id path string true "Dispute ID"
// @Success 200 {object} models.Response "Dispute Patch Endpoint"
// @Router /dispute/{id} [patch]
func (h Dispute) patchDispute(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	id := c.Param("id")
	logger.Info("Dispute Patch Endpoint for ID: ", id)
	c.JSON(http.StatusOK, models.Response{Data: "Dispute Patch Endpoint for ID: " + id})
}

func (h Dispute) expertObjection(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	//get dispute id
	disputeId := c.Param("id")
	disputeIdInt, err := strconv.Atoi(disputeId)
	if err != nil {
		logger.WithError(err).Error("Cannot convert dispute ID to integer")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Dispute ID"})
		return
	}

	//get info from body of post
	var req models.ExpertRejectRequest
	if err := c.BindJSON(&req); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	//get user properties from token
	claims := middleware.GetClaims(c)
	if claims == nil {
		logger.Error("Unauthorized access attempt in function expertObjection")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}
	err = h.Model.ObjectExpert(claims.User.ID, int64(disputeIdInt), req.ExpertID, req.Reason)
	if err != nil {
		logger.Error("Unauthorized access attempt in function expertObjection")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Unauthorized"})
		return
	}

	logger.Info("Expert rejected suggestion")
	c.JSON(http.StatusOK, models.Response{Data: "objection filed successfully"})
}

func (h Dispute) expertObjectionsReview(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	// Get dispute id
	disputeId := c.Param("id")
	disputeIdInt, err := strconv.Atoi(disputeId)
	if err != nil {
		logger.WithError(err).Error("Cannot convert dispute ID to integer")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Dispute ID"})
		return
	}

	// Get info from token
	claims := middleware.GetClaims(c)
	if claims == nil {
		logger.WithError(err).Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Get body of post
	var req models.RejectExpertReview
	if err := c.BindJSON(&req); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Request failed"})
		return
	}

	err = h.Model.ReviewExpertObjection(claims.User.ID, int64(disputeIdInt), req.ExpertID, req.Accepted)
	if err != nil {
		logger.WithError(err).Error("failed to review objection")
		c.JSON(http.StatusBadRequest, models.Response{Error: "failed to review objection"})
		return
	}

	logger.Info("Expert objections reviewed successfully")
	c.JSON(http.StatusOK, models.Response{Data: "Expert objections reviewed successfully"})
}

func (m *disputeModelReal) UploadEvidence(userId, disputeId int64, path string, file io.Reader) (uint, error) {
	fileStorageRoot, err := env.Get("FILESTORAGE_ROOT")
	if err != nil {
		return 0, err
	}
	fileStorageUrl, err := env.Get("FILESTORAGE_URL")
	if err != nil {
		return 0, err
	}
	logger := utilities.NewLogger().LogWithCaller()

	name := filepath.Base(path)
	dest := filepath.Join(fileStorageRoot, path)

	if err := os.MkdirAll(filepath.Join(fileStorageRoot, path), 0755); err != nil {
		logger.WithError(err).Error("failed to create folder for file upload")
		return 0, err
	}

	url := fmt.Sprintf("%s/%s", fileStorageUrl, path)

	// Open the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		logger.WithError(err).Error("failed to create file in storage")
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

	if err := m.db.Create(&file).Error; err != nil {
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

func (m *disputeModelReal) GetDisputeSummaries(userId int64) ([]models.DisputeSummaryResponse, error) {
	var disputes []models.Dispute
	err := m.db.Where("complainant = ? OR respondant = ?", userId, userId).Find(&disputes).Error
	if err != nil {
		return nil, err
	}

	var summaries []models.DisputeSummaryResponse
	for _, dispute := range disputes {
		var role string = ""
		if dispute.Complainant == userId {
			role = "Complainant"
		} else if *(dispute.Respondant) == userId {
			role = "Respondant"
		}
		summary := models.DisputeSummaryResponse{
			ID:          *dispute.ID,
			Title:       dispute.Title,
			Description: dispute.Description,
			Status:      dispute.Status,
			Role:        &role,
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (m *disputeModelReal) GetDispute(disputeId int64) (dispute models.Dispute, err error) {
	err = m.db.Model(&models.Dispute{}).Where("id = ?", disputeId).First(&dispute).Error
	return dispute, err
}

func (m *disputeModelReal) GetEvidenceByDispute(disputeId int64) (evidence []models.Evidence, err error) {
	err = m.db.Table("dispute_evidence").Select(`files.id, file_name, uploaded, file_path,  CASE
            WHEN disputes.complainant = dispute_evidence.user_id THEN 'Complainant'
            WHEN disputes.respondant = dispute_evidence.user_id THEN 'Respondent'
            ELSE 'Other'
        END AS uploader_role`).Joins("JOIN files ON dispute_evidence.file_id = files.id").Joins("JOIN disputes ON dispute_evidence.dispute = disputes.id").Where("dispute = ?", disputeId).Find(&evidence).Error
	if err != nil {
		return
	}
	if evidence == nil {
		evidence = []models.Evidence{}
	}
	return evidence, err
}
func (m *disputeModelReal) GetDisputeExperts(disputeId int64) (experts []models.Expert, err error) {
	err = m.db.Table("dispute_experts").Select("users.id, users.first_name || ' ' || users.surname AS full_name, email, users.phone_number AS phone, role").Joins("JOIN users ON dispute_experts.user = users.id").Where("dispute = ?", disputeId).Where("dispute_experts.status = 'Approved'").Where("role = 'Mediator' OR role = 'Arbitrator' OR role = 'Conciliator'").Find(&experts).Error
	if err != nil && err.Error() != "record not found" {
		return
	}

	if experts == nil {
		experts = []models.Expert{}
	}
	return experts, err
}

func (m *disputeModelReal) GetDisputesByUser(userId int64) (disputes []models.Dispute, err error) {
	err = m.db.Where("complainant = ? OR respondant = ?", userId, userId).Find(&disputes).Error
	return disputes, err
}

func (m *disputeModelReal) GetUserByEmail(email string) (user models.User, err error) {
	err = m.db.Where("email = ?", email).First(&user).Error
	return user, err
}
func (m *disputeModelReal) CreateDispute(dispute models.Dispute) (id int64, err error) {
	disputeCloned := dispute
	err = m.db.Create(&disputeCloned).Error
	return *disputeCloned.ID, nil
}
func (m *disputeModelReal) UpdateDisputeStatus(disputeId int64, status string) error {
	return m.db.Model(&models.Dispute{}).Where("id = ?", disputeId).Update("status", status).Error
}
func (m *disputeModelReal) ObjectExpert(userId, disputeId, expertId int64, reason string) error {
	//update dispute experts table
	var disputeExpert models.DisputeExpert
	if err := m.db.Where("dispute = ? AND dispute_experts.user = ?", disputeId, expertId).First(&disputeExpert).Error; err != nil {
		return err
	}

	disputeExpert.Status = models.ReviewStatus
	if err := m.db.Save(&disputeExpert).Error; err != nil {
		return err
	}

	//add entry to expert objections table
	expertObjection := models.ExpertObjection{
		DisputeID: disputeId,
		ExpertID:  expertId,
		UserID:    userId,
		Reason:    reason,
	}

	if err := m.db.Create(&expertObjection).Error; err != nil {
		return err
	}
	return nil
}
func (m *disputeModelReal) ReviewExpertObjection(userId, disputeId, expertId int64, approved bool) error {

	var expertObjections []models.ExpertObjection
	if err := m.db.Where("dispute_id = ? AND expert_id = ? AND status = ?", disputeId, expertId, models.ReviewStatus).Find(&expertObjections).Error; err != nil {
		return err
	}

	var disputeExpert models.DisputeExpert
	if err := m.db.Where("dispute = ? AND dispute_experts.user = ? AND status = ?", disputeId, expertId, models.ReviewStatus).First(&disputeExpert).Error; err != nil {
		return nil
	}

	// Start a transaction
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update expert objections
	if approved {
		for i := range expertObjections {
			expertObjections[i].Status = models.Sustained
		}
		disputeExpert.Status = models.RejectedStatus
	} else {
		for i := range expertObjections {
			expertObjections[i].Status = models.Overruled
		}
		disputeExpert.Status = models.ApprovedStatus
	}

	// Save the expert objections
	for _, objection := range expertObjections {
		if err := tx.Save(&objection).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Save the dispute expert
	if err := tx.Save(&disputeExpert).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
