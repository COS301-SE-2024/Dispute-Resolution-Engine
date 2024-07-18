package handlers

import (
	//"api/middleware"
	"api/middleware"
	"api/models"
	"api/utilities"
	"errors"
	"io"
	"log"
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
	g.POST("/:id/evidence", h.uploadEvidence)

	//patch is not to be integrated yet
	// disputeRouter.HandleFunc("/{id}", h.patchDispute).Methods(http.MethodPatch)

	//create dispute

	//archive routes
}

func uploadFile(db *gorm.DB, header *multipart.FileHeader) (uint, error) {
	logger := utilities.NewLogger().LogWithCaller()

	// Open the form file
	formFile, err := header.Open()
	if err != nil {
		logger.WithError(err).Error("Failed to open form file")
		return 0, errors.New("Failed to open form file")
	}
	defer formFile.Close()

	// Open the destination file
	fileName := filepath.Base(header.Filename)
	destPath := filepath.Join(os.Getenv("FILESTORAGE_ROOT"), fileName) // Assuming '/files' is where Docker mounts its storage
	destFile, err := os.Create(destPath)
	if err != nil {
		logger.WithError(err).Error("Failed to create file in storage")
		return 0, errors.New("Failed to create file in storage")
	}
	defer destFile.Close()

	// Copy file content to destination
	_, err = io.Copy(destFile, formFile)
	if err != nil {
		logger.WithError(err).Error("Failed to copy file content")
		return 0, errors.New("Failed to copy file content")
	}

	// TODO: Change this to a proper URL
	fileUrl := destPath
	//add file to Database
	file := models.File{
		FileName: fileName,
		Uploaded: time.Now(),

		FilePath: fileUrl,
	}

	if err := db.Create(&file).Error; err != nil {
		logger.WithError(err).Error("Error adding file to database")
		return 0, errors.New("Error adding file to database")
	}

	//get id of the created file enrty
	var fileFromDbInserted models.File
	err = db.Where("file_name = ? AND file_path = ?", fileName, fileUrl).First(&fileFromDbInserted).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving file entry from database")
		return 0, errors.New("Error retrieving file entry from database")
	}
	return *fileFromDbInserted.ID, nil
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

	disputeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Dispute ID"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Failed to parse form data"})
		return
	}

	files := form.File["files"]
	for _, fileHeader := range files {
		id, err := uploadFile(h.DB, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
			return
		}

		//add enrty to dispute evidence table
		disputeEvidence := models.DisputeEvidence{
			Dispute: int64(disputeId),
			FileID:  int64(id),
		}
		
        if err := h.DB.Create(&disputeEvidence).Error; err != nil {
			logger.WithError(err).Error("Error creating dispute evidence")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating dispute evidence"})
			return
		}
	}
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
	jwtClaims := middleware.GetClaims(c)
	var disputes []models.Dispute
	err := h.DB.Find(&disputes).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	var disputeSummaries []models.DisputeSummaryResponse
	userID := jwtClaims.User.ID
	for _, dispute := range disputes {
		var role string = ""
		if dispute.Complainant == userID {
			role = "Complainant"
		}
		if dispute.Respondant == &userID {
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

	var disputes models.Dispute
	err := h.DB.Raw("SELECT id, title, description, status, case_date, respondant, complainant FROM disputes WHERE id = ?", id).Scan(&disputes).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}

	//name and email
	// var respondantData models.User
	// err = h.DB.Where("id = ?", disputes.Respondant).Scan(&respondantData).Error
	// if err!=nil {

	// }

	DisputeDetailsResponse := models.DisputeDetailsResponse{
		ID:          *disputes.ID,
		Title:       disputes.Title,
		Description: disputes.Description,
		Status:      disputes.Status,
		DateCreated: disputes.CaseDate,
	}

	err = h.DB.Raw("SELECT file_name,uploaded,file_path FROM files WHERE id IN (SELECT file_id FROM dispute_evidence WHERE dispute = ?)", id).Scan(&DisputeDetailsResponse.Evidence).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}

	DisputeDetailsResponse.Experts = []string{"Expert 1", "Expert 2", "Expert 3"}

	c.JSON(http.StatusOK, models.Response{Data: DisputeDetailsResponse})
	// c.JSON(http.StatusOK, models.Response{Data: "Dispute Detail Endpoint for ID: " + id})
}

func (h Dispute) createDispute(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	form, err := c.MultipartForm()
	if err != nil {
		logger.WithError(err).Error("Error parsing form")
		return
	}

	// Access form values
	title := form.Value["title"][0]
	description := form.Value["description"][0]
	fullName := form.Value["respondent[full_name]"][0]
	email := form.Value["respondent[email]"][0]
	telephone := form.Value["respondent[telephone]"][0]

	//get complainants id
	claims := middleware.GetClaims(c)
	if claims == nil {
		logger.Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}
	complainantID := claims.User.ID

	//check if respondant is in database by email and phone number
	var respondantID *int64
	var respondent models.User
	err = h.DB.Where("email = ? AND phone_number = ?", email, telephone).First(&respondent).Error
	if err != nil && err.Error() == "record not found" {
		//create a deafult entry for the user
		nameSplit := strings.Split(fullName, " ")
		if len(nameSplit) < 2 {
			logger.Error("Invalid full name")
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid full name"})
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

	//get the id of the created dispute
	var disputeFromDbInserted models.Dispute
	err = h.DB.Where("title = ? AND case_date = ? AND status = ? AND description = ? AND complainant = ? AND resolved = ? AND decision = ?", title, time.Now(), "Awaiting Respondant", description, complainantID, false, models.Unresolved).First(&disputeFromDbInserted).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving dispute")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving dispute"})
		return
	}

	// Store files in Docker and retrieve URLs
	fileURLs := []string{}
	fileNames := []string{}
	files := form.File["files"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to open file"})
			return
		}
		defer file.Close()

		// Generate a unique filename
		fileName := filepath.Base(fileHeader.Filename)
		fileNames = append(fileNames, fileName)

		fileLocation := filepath.Join(os.Getenv("FILESTORAGE_ROOT"), fileName) // Assuming '/files' is where Docker mounts its storage

		// Create the file in Docker (or any storage system you use)
		f, err := os.Create(fileLocation)
		if err != nil {
			logger.WithError(err).Error("Failed to create file in storage")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to create file in storage"})
			return
		}
		defer f.Close()

		// Copy file content to destination
		_, err = io.Copy(f, file)
		if err != nil {
			logger.WithError(err).Error("Failed to copy file content")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to copy file content"})
			return
		}

		// Generate URL for accessing the file
		// fileURL := fmt.Sprintf("https://your-domain.com%s", fileLocation)
		fileURLs = append(fileURLs, fileLocation)
	}

	// Store file URLs in PostgreSQL database
	for i, fileURL := range fileURLs {
		//add file to Database
		file := models.File{
			FileName: fileNames[i],
			Uploaded: time.Now(),
			FilePath: fileURL,
		}

		err = h.DB.Create(&file).Error
		if err != nil {
			logger.WithError(err).Error("Error creating file")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating file"})
			return
		}

		//get id of the created file enrty
		var fileFromDbInserted models.File
		err = h.DB.Where("file_name = ? AND file_path = ?", fileNames[i], fileURL).First(&fileFromDbInserted).Error
		if err != nil {
			logger.WithError(err).Error("Error retrieving file")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving file"})
			return
		}

		//add enrty to dispute evidence table
		disputeEvidence := models.DisputeEvidence{
			Dispute: *disputeFromDbInserted.ID,
			FileID:  int64(*fileFromDbInserted.ID),
		}
		err = h.DB.Create(&disputeEvidence).Error
		if err != nil {
			logger.WithError(err).Error("Error creating dispute evidence")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating dispute evidence"})
			return
		}
	}

	// Respond with success message
	h.sendAdminNotification(c, email)
	c.JSON(http.StatusCreated, models.Response{Data: "Dispute created successfully"})
	log.Printf("Dispute created successfully: %s", title)
	logger.Info("Dispute created successfully: ", title)
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

	id := c.Param("id")
	c.JSON(http.StatusOK, models.Response{Data: "Dispute Patch Endpoint for ID: " + id})
}
