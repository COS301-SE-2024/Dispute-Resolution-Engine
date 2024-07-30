package handlers

import (
	"api/env"
	"api/middleware"
	"api/models"
	"api/utilities"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

// Uploads a multipart file to the file storage, returning the id of the file entry in the database
func uploadFile(db *gorm.DB, path string, header *multipart.FileHeader) (uint, error) {
	fileStorageRoot, err := env.Get("FILESTORAGE_ROOT")
	if err != nil {
		return 0, err
	}
	fileStorageUrl, err := env.Get("FILESTORAGE_URL")
	if err != nil {
		return 0, err
	}

	logger := utilities.NewLogger().LogWithCaller()

	fileName := filepath.Base(header.Filename)
	storePath := filepath.Join(fileStorageRoot, path, fileName)

	if err := os.MkdirAll(filepath.Join(fileStorageRoot, path), 0755); err != nil {
		logger.WithError(err).Error("failed to create folder for file upload")
		return 0, err
	}

	var storeUrl string
	if path != "" {
		storeUrl = strings.Join([]string{fileStorageUrl, path, fileName}, "/")
	} else {
		storeUrl = strings.Join([]string{fileStorageUrl, fileName}, "/")
	}

	// Open the form file
	formFile, err := header.Open()
	if err != nil {
		logger.WithError(err).Error("failed to open form file")
		return 0, errors.New("failed to open form file")
	}
	defer formFile.Close()

	// Open the destination file
	storeFile, err := os.Create(storePath)
	if err != nil {
		logger.WithError(err).Error("failed to create file in storage")
		return 0, errors.New("failed to create file in storage")
	}
	defer storeFile.Close()

	// Copy file content to destination
	_, err = io.Copy(storeFile, formFile)
	if err != nil {
		logger.WithError(err).Error("failed to copy file content")
		return 0, errors.New("failed to copy file content")
	}

	// Add file entry to Database
	file := models.File{
		FileName: fileName,
		FilePath: storeUrl,
		Uploaded: time.Now(),
	}

	if err := db.Create(&file).Error; err != nil {
		logger.WithError(err).Error("error adding file to database")
		return 0, errors.New("error adding file to database")
	}

	return *file.ID, nil
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
		id, err := uploadFile(h.DB, folder, fileHeader)
		if err != nil {
			logger.WithError(err).Error("Error uploading evidence")
			c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
			return
		}

		//add enrty to dispute evidence table
		disputeEvidence := models.DisputeEvidence{
			Dispute: int64(disputeId),
			FileID:  int64(id),
			UserID:  claims.User.ID,
		}

		if err := h.DB.Create(&disputeEvidence).Error; err != nil {
			logger.WithError(err).Error("Error creating dispute evidence")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating dispute evidence"})
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

	var disputes []models.Dispute
	err := h.DB.Where("complainant = ? OR respondant = ?", userID, userID).Find(&disputes).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving disputes")
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	logger.Info("Retrieving disputes: ", disputes)

	var disputeSummaries []models.DisputeSummaryResponse
	for _, dispute := range disputes {
		var role string = ""
		if dispute.Complainant == userID {
			role = "Complainant"
		} else if *(dispute.Respondant) == userID {
			role = "Respondant"
		}
		summary := models.DisputeSummaryResponse{
			ID:          *dispute.ID,
			Title:       dispute.Title,
			Description: dispute.Description,
			Status:      dispute.Status,
			Role:        &role,
		}
		disputeSummaries = append(disputeSummaries, summary)
	}
	logger.Info("Dispute summaries retrieved successfully")
	c.JSON(http.StatusOK, models.Response{Data: disputeSummaries})
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
	id := c.Param("id")
	logger := utilities.NewLogger().LogWithCaller()
	var disputes models.Dispute
	err := h.DB.Raw("SELECT id, title, description, status, case_date, respondant, complainant FROM disputes WHERE id = ?", id).Scan(&disputes).Error
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

	if userId == disputes.Complainant {
		role = "Complainant"
	} else if userId == *(disputes.Respondant) {
		role = "Respondent"
	}

	DisputeDetailsResponse := models.DisputeDetailsResponse{
		ID:          *disputes.ID,
		Title:       disputes.Title,
		Description: disputes.Description,
		Status:      disputes.Status,
		DateCreated: disputes.CaseDate,
		Role:        role,
	}

	var evidence []models.Evidence
	err = h.DB.Table("dispute_evidence").Select(`files.id, file_name, uploaded, file_path,  CASE
            WHEN disputes.complainant = dispute_evidence.user_id THEN 'Complainant'
            WHEN disputes.respondant = dispute_evidence.user_id THEN 'Respondent'
            ELSE 'Other'
        END AS uploader_role`).Joins("JOIN files ON dispute_evidence.file_id = files.id").Joins("JOIN disputes ON dispute_evidence.dispute = disputes.id").Where("dispute = ?", id).Find(&evidence).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute evidence")
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	if evidence == nil {
		evidence = []models.Evidence{}
	}

	var experts []models.Expert
	err = h.DB.Table("dispute_experts").Select("users.id, users.first_name || ' ' || users.surname AS full_name, email, users.phone_number AS phone, role").Joins("JOIN users ON dispute_experts.user = users.id").Where("dispute = ?", id).Where("dispute_experts.status = 'Approved'").Where("role = 'Mediator' OR role = 'Arbitrator' OR role = 'Conciliator'").Find(&experts).Error
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
	telephone := form.Value["respondent[telephone]"][0]

	//get complainants id
	complainantID := claims.User.ID

	//check if respondant is in database by email and phone number
	var respondantID *int64
	var respondent models.User
	err = h.DB.Where("email = ? AND phone_number = ?", email, telephone).First(&respondent).Error
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
	dispute := models.Dispute{
		Title:       title,
		CaseDate:    time.Now(),
		Workflow:    nil,
		Status:      "Awaiting Respondant",
		Description: description,
		Complainant: complainantID,
		Respondant:  respondantID,
		Resolved:    false,
		Decision:    models.Unresolved,
	}

	err = h.DB.Create(&dispute).Error
	if err != nil {
		logger.WithError(err).Error("Error creating dispute")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating dispute"})
		return
	}

	// Store files in Docker and retrieve URLs
	files := form.File["files"]
	folder := fmt.Sprintf("%d", *dispute.ID)
	for _, fileHeader := range files {
		fileId, err := uploadFile(h.DB, folder, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Error: fmt.Sprintf("Failed to upload file: %s", fileHeader.Filename)})
			return
		}

		disputeEvidence := models.DisputeEvidence{
			Dispute: *dispute.ID,
			FileID:  int64(fileId),
			UserID:  complainantID,
		}

		if err := h.DB.Create(&disputeEvidence).Error; err != nil {
			logger.WithError(err).Error("Error creating dispute evidence")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating dispute evidence"})
			return
		}
	}
	disputeID := *dispute.ID
	// Respond with success message
	go h.sendAdminNotification(c, disputeID, email)
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

	var dbDispute models.Dispute
	h.DB.Where("id = ?", disputeStatus.DisputeID).First(&dbDispute)

	dbDispute.Status = disputeStatus.Status

	h.DB.Model(&dbDispute).Where("id = ?", dbDispute.ID).Updates(dbDispute)
	go h.StateChangeNotifications(c, disputeStatus.DisputeID, disputeStatus.Status)
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

	//update dispute experts table
	var disputeExpert models.DisputeExpert
	err = h.DB.Where("dispute = ? AND dispute_experts.user = ?", int64(disputeIdInt), req.ExpertID).First(&disputeExpert).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute expert")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error filing objection"})
		return
	}

	disputeExpert.Status = models.ReviewStatus
	if err := h.DB.Save(&disputeExpert).Error; err != nil {
		logger.WithError(err).Error("Error updating dispute expert")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error filing objection"})
		return
	}

	//add entry to expert objections table
	expertObjection := models.ExpertObjection{
		DisputeID: int64(disputeIdInt),
		ExpertID:  req.ExpertID,
		UserID:    claims.User.ID,
		Reason:    req.Reason,
	}

	if err := h.DB.Create(&expertObjection).Error; err != nil {
		logger.WithError(err).Error("Error creating expert objection")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error filing objection"})
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

	// Get expert objections
	var expertObjections []models.ExpertObjection
	err = h.DB.Where("dispute_id = ? AND expert_id = ? AND status = ?", disputeIdInt, req.ExpertID, models.ReviewStatus).Find(&expertObjections).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving expert objections")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error updating expert objections"})
		return
	}

	var disputeExpert models.DisputeExpert
	err = h.DB.Where("dispute = ? AND dispute_experts.user = ? AND status = ?", disputeId, req.ExpertID, models.ReviewStatus).First(&disputeExpert).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute expert")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error updating expert objections"})
		return
	}

	// Start a transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		logger.WithError(tx.Error).Error("Error starting transaction")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error starting transaction"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.Error("Recovered in deferred function, rolled back transaction")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal server error"})
		}
	}()

	// Update expert objections
	if req.Accepted {
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
			logger.WithError(err).Error("Error updating expert objections")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error updating expert objections"})
			return
		}
	}

	// Save the dispute expert
	if err := tx.Save(&disputeExpert).Error; err != nil {
		tx.Rollback()
		logger.WithError(err).Error("Error updating dispute expert")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error updating expert objections"})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		logger.WithError(err).Error("Error committing transaction")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error committing transaction"})
		return
	}

	logger.Info("Expert objections reviewed successfully")
	c.JSON(http.StatusOK, models.Response{Data: "Expert objections reviewed successfully"})
}

