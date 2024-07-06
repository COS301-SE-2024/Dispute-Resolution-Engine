package handlers

import (
	//"api/middleware"
	"api/middleware"
	"api/models"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupDisputeRoutes(g *gin.RouterGroup, h Dispute) {
	//dispute routes
	g.Use(middleware.JWTMiddleware)

	g.GET("", h.getSummaryListOfDisputes)
	g.POST("/create", h.createDispute)
	g.GET("/:id", h.getDispute)

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
// @Router /dispute [get]
func (h Dispute) getSummaryListOfDisputes(c *gin.Context) {
	var disputes []models.DisputeSummaryResponse
	err := h.DB.Raw("SELECT id, title, description, status FROM disputes").Scan(&disputes).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Data: disputes})
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

	var DisputeDetailsResponse models.DisputeDetailsResponse
	err := h.DB.Raw("SELECT id, title, description, status, case_date FROM disputes WHERE id = ?", id).Scan(&DisputeDetailsResponse).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}

	err = h.DB.Raw("SELECT file_path FROM files WHERE id IN (SELECT file_id FROM dispute_evidence WHERE dispute = ?)", id).Scan(&DisputeDetailsResponse.Evidence).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}

	DisputeDetailsResponse.Experts = []string{"Expert 1", "Expert 2", "Expert 3"}

	c.JSON(http.StatusOK, models.Response{Data: DisputeDetailsResponse})
	// c.JSON(http.StatusOK, models.Response{Data: "Dispute Detail Endpoint for ID: " + id})
}

func (h Dispute) createDispute(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
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
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid full name"})
			return
		}

	} else if err != nil {
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
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating dispute"})
		return
	}

	//get the id of the created dispute
	var disputeFromDbInserted models.Dispute
	err = h.DB.Where("title = ? AND case_date = ? AND status = ? AND description = ? AND complainant = ? AND resolved = ? AND decision = ?", title, time.Now(), "Awaiting Respondant", description, complainantID, false, models.Unresolved).First(&disputeFromDbInserted).Error
	if err != nil {
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
		fileLocation := filepath.Join("/app/filestorage", fileName) // Assuming '/files' is where Docker mounts its storage

		// Create the file in Docker (or any storage system you use)
		f, err := os.Create(fileLocation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to create file in storage"})
			return
		}
		defer f.Close()

		// Copy file content to destination
		_, err = io.Copy(f, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to copy file content"})
			return
		}

		// Generate URL for accessing the file
		fileURL := fmt.Sprintf("https://your-domain.com%s", fileLocation)
		fileURLs = append(fileURLs, fileURL)
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
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating file"})
			return
		}

		//get id of the created file enrty
		var fileFromDbInserted models.File
		err = h.DB.Where("file_name = ? AND file_path = ?", fileNames[i], fileURL).First(&fileFromDbInserted).Error
		if err != nil {
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
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating dispute evidence"})
			return
		}
	}

	// Respond with success message
	c.JSON(http.StatusCreated, models.Response{Data: "Dispute created successfully"})
	log.Printf("Dispute created successfully: %s", title)
}

// Function to get MIME type from file header
func getFileType(fh *multipart.FileHeader) string {
	file, err := fh.Open()
	if err != nil {
		return "application/octet-stream" // default to octet-stream if cannot determine type
	}
	defer file.Close()

	// Determine file type based on MIME type
	buffer := make([]byte, 512) // Read the first 512 bytes to detect MIME type
	_, err = file.Read(buffer)
	if err != nil {
		return "application/octet-stream" // default to octet-stream if cannot determine type
	}

	mimeType := http.DetectContentType(buffer)
	return mimeType
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

func getMockArchiveDisputeSummaries() []models.ArchivedDisputeSummary {
	return []models.ArchivedDisputeSummary{
		{
			ID:           1,
			Title:        "Dispute 1: Contract Disagreement",
			Summary:      "A contractual dispute between parties over payment terms.",
			Category:     []string{"Legal"},
			DateFiled:    time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			DateResolved: time.Date(2021, time.January, 15, 0, 0, 0, 0, time.UTC),
			Resolution:   "Settlement reached with revised terms.",
		},
		{
			ID:           2,
			Title:        "Dispute 2: Product Quality Issue",
			Summary:      "Customer complained about product defects; manufacturer's response needed.",
			Category:     []string{"Customer Service", "Quality Control"},
			DateFiled:    time.Date(2021, time.February, 10, 0, 0, 0, 0, time.UTC),
			DateResolved: time.Date(2021, time.February, 28, 0, 0, 0, 0, time.UTC),
			Resolution:   "Product recall initiated; replacements provided.",
		},
		{
			ID:           3,
			Title:        "Dispute 3: Employment Dispute",
			Summary:      "Employee termination dispute due to performance issues.",
			Category:     []string{"Human Resources", "Legal"},
			DateFiled:    time.Date(2021, time.March, 5, 0, 0, 0, 0, time.UTC),
			DateResolved: time.Date(2021, time.March, 20, 0, 0, 0, 0, time.UTC),
			Resolution:   "Settled with severance package and agreement.",
		},
	}
}

func paginateSummaries(summaries []models.ArchivedDisputeSummary, offset int, limit int) []models.ArchivedDisputeSummary {
	start := offset
	end := offset + limit
	if start >= len(summaries) {
		return []models.ArchivedDisputeSummary{}
	}
	if end > len(summaries) {
		end = len(summaries)
	}
	return summaries[start:end]
}

func sortSummaries(summaries []models.ArchivedDisputeSummary, sorting string, order string) {
	switch sorting {
	case "title":
		if order == "asc" {
			sort.Slice(summaries, func(i, j int) bool {
				return summaries[i].Title < summaries[j].Title
			})
		} else {
			sort.Slice(summaries, func(i, j int) bool {
				return summaries[i].Title > summaries[j].Title
			})
		}
	case "date_filed":
		if order == "asc" {
			sort.Slice(summaries, func(i, j int) bool {
				return summaries[i].DateFiled.Before(summaries[j].DateFiled)
			})
		} else {
			sort.Slice(summaries, func(i, j int) bool {
				return summaries[i].DateFiled.After(summaries[j].DateFiled)
			})
		}
	case "date_resolved":
		if order == "asc" {
			sort.Slice(summaries, func(i, j int) bool {
				return summaries[i].DateResolved.Before(summaries[j].DateResolved)
			})
		} else {
			sort.Slice(summaries, func(i, j int) bool {
				return summaries[i].DateResolved.After(summaries[j].DateResolved)
			})
		}
	}
}

func filterSummariesBySearch(summaries []models.ArchivedDisputeSummary, searchTerm string) []models.ArchivedDisputeSummary {
	if searchTerm == "" {
		return summaries
	}
	var filteredSummaries []models.ArchivedDisputeSummary
	for _, summary := range summaries {
		// Example of case-insensitive search
		if strings.Contains(strings.ToLower(summary.Title), strings.ToLower(searchTerm)) {
			filteredSummaries = append(filteredSummaries, summary)
		}
	}
	return filteredSummaries
}
