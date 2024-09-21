package handlers

import (
	"api/models"
	"api/utilities"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupArchiveRoutes(g *gin.RouterGroup, h Archive) {
	g.POST("/search", h.SearchArchive)
	g.GET("/highlights", h.Highlights)
	g.GET("/:id", h.getArchive)
}

// @Summary Get a list of highlight disputes from the archive
// @Description Get a list of highlight disputes from the archive
// @Tags archive
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Archive Highlights Endpoint"
// @Router /archive [post]
func (h Archive) Highlights(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	limitParam := c.DefaultQuery("limit", "3")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		logger.WithError(err).Error("Error parsing query parameter")
		c.JSON(http.StatusInternalServerError, models.Response{Error: fmt.Sprintf("Invalid query parameter: %s", limitParam)})
		return
	}

	// Query the database
	var disputes []models.Dispute
	if err := h.DB.Model(&models.Dispute{}).Where("date_resolved IS NOT NULL").Limit(limit).Scan(&disputes).Error; err != nil {
		logger.WithError(err).Error("Error retrieving disputes")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving disputes"})
		return
	}

	// Transform the results to ArchivedDisputeSummary
	summaries := make([]models.ArchivedDisputeSummary, len(disputes))
	for i, dispute := range disputes {
		disputeSummary := models.DisputeSummaries{}
		err := h.DB.Model(models.DisputeSummaries{}).Where("dispute = ?", *dispute.ID).First(&disputeSummary).Error
		if err != nil {
			logger.WithError(err).Error("Could not get dispute for id:" + fmt.Sprint(*dispute.ID))
		}

		summaries[i] = models.ArchivedDisputeSummary{
			ID:           *dispute.ID,
			Title:        dispute.Title,
			Description:  dispute.Description,
			Summary:      disputeSummary.Summary,
			Category:     []string{"Dispute"}, // Assuming a default category for now
			DateFiled:    dispute.CaseDate.Format("2006-01-02"),
			DateResolved: dispute.DateResolved.Format("2006-01-02"), // Placeholder for resolved date
			Resolution:   string(dispute.Status),
		}
	}

	// Return the response
	logger.Info("Successfully retrieved disputes")
	c.JSON(http.StatusOK, models.Response{Data: models.ArchiveSearchResponse{
		Archives: summaries,
		Total:    int64(limit),
	}})
}

// @Summary Get a summary list of archives
// @Description Get a summary list of archives
// @Tags archive
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Archive Summary Endpoint"
// @Router /archive [post]
func (h Archive) SearchArchive(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	var body models.ArchiveSearchRequest

	if err := c.BindJSON(&body); err != nil {
		logger.WithError(err).Error("Error binding request body")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request body"})
		return
	}

	// Handle the request
	searchTerm := ""
	limit := 10
	offset := 0
	order := "asc"
	sort := "id"

	if body.Search != nil {
		searchTerm = *body.Search
	}
	if body.Limit != nil {
		limit = *body.Limit
	}
	if body.Offset != nil {
		offset = *body.Offset
	}
	if body.Order != nil {
		order = *body.Order
		if strings.ToLower(order) != "asc" && strings.ToLower(order) != "desc" {
			logger.Error("Invalid order value")
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid order value"})
			return
		}
	}
	if body.Sort != nil {
		sort = string(*body.Sort)
	}

	// Query the database
	var disputes []models.Dispute
	query := h.DB.Model(&models.Dispute{})

	// Apply search filter
	if searchTerm != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%")
	}

	query = query.Where("date_resolved IS NOT NULL", true)

	var count int64
	query = query.Count(&count)

	// Apply sorting
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	// Apply pagination
	query = query.Offset(offset).Limit(limit)

	// Execute the query
	if err := query.Find(&disputes).Error; err != nil {
		logger.WithError(err).Error("Error retrieving disputes")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving disputes"})
		return
	}

	if len(disputes) == 0 {
		logger.Info("No disputes found")
		c.JSON(http.StatusOK, models.Response{Data: models.ArchiveSearchResponse{
			Archives: []models.ArchivedDisputeSummary{},
			Total:    0,
		}})
		return
	}

	// Transform the results to ArchivedDisputeSummary
	var archiveDisputeSummaries []models.ArchivedDisputeSummary
	for _, dispute := range disputes {
		disputeSummary := models.DisputeSummaries{}
		err := h.DB.Model(models.DisputeSummaries{}).Where("dispute = ?", *dispute.ID).First(&disputeSummary).Error
		if err != nil {
			logger.WithError(err).Error("Could not get dispute for id:" + fmt.Sprint(*dispute.ID))
		}
		archiveDisputeSummaries = append(archiveDisputeSummaries, models.ArchivedDisputeSummary{
			ID:           *dispute.ID,
			Title:        dispute.Title,
			Description:  dispute.Description,
			Summary:      disputeSummary.Summary,
			Category:     []string{"Dispute"}, // Assuming a default category for now
			DateFiled:    dispute.CaseDate.Format("2006-01-02"),
			DateResolved: dispute.DateResolved.Format("2006-01-02"), // Placeholder for resolved date
			Resolution:   string(dispute.Status),
		})
	}

	// Return the response
	logger.Info("Successfully retrieved disputes")
	c.JSON(http.StatusOK, models.Response{Data: models.ArchiveSearchResponse{
		Archives: archiveDisputeSummaries,
		Total:    count,
	}})
}

// @Summary Get an archive
// @Description Get an archive
// @Tags archive
// @Accept json
// @Produce json
// @Param id path string true "Archive ID"
// @Success 200 {object} models.Response "Archive Detail Endpoint"
// @Router /archive/{id} [get]
func (h Archive) getArchive(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	id := c.Param("id")

	intID, err := strconv.Atoi(id)
	if err != nil {
		logger.WithError(err).Error("Invalid ID")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid ID"})
		return
	}

	//mock body
	// body := models.ArchivedDispute{
	// 	ArchivedDisputeSummary: models.ArchivedDisputeSummary{
	// 		ID:           int64(intID),
	// 		Title:        "Dispute " + id,
	// 		Summary:      "Summary " + id,
	// 		Category:     []string{"Category " + id},
	// 		DateFiled:    time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
	// 		DateResolved: time.Date(2021, time.January, 2, 0, 0, 0, 0, time.UTC),
	// 		Resolution:   "Resolution " + id,
	// 	},
	// 	Events: []models.Event{
	// 		{
	// 			Timestamp:   "2021-01-01T00:00:00Z",
	// 			Type:        "Type 1",
	// 			Description: "Details 1",
	// 		},
	// 		{
	// 			Timestamp:   "2021-01-02T00:00:00Z",
	// 			Type:        "Type 2",
	// 			Description: "Details 2",
	// 		},
	// 	},
	// }

	//request to db
	var dispute models.Dispute

	err = h.DB.Where("id = ?", intID).First(&dispute).Error
	if err != nil && err.Error() == "record not found" {
		logger.WithError(err).Error("Dispute not found")
		c.JSON(http.StatusNotFound, models.Response{Data: ""})
		return
	} else if err != nil {
		logger.WithError(err).Error("Error retrieving dispute")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving dispute"})
		return
	}

	//transform to archive dispute
	var archiveDispute models.ArchivedDispute
	if *dispute.ID != 0 {
		disputeSummary := models.DisputeSummaries{}
		err := h.DB.Model(models.DisputeSummaries{}).Where("dispute = ?", *dispute.ID).First(&disputeSummary).Error
		if err != nil {
			logger.WithError(err).Error("Could not get dispute for id:" + fmt.Sprint(*dispute.ID))
		}
		archiveDispute = models.ArchivedDispute{

			ArchivedDisputeSummary: models.ArchivedDisputeSummary{
				ID:           *dispute.ID,
				Title:        dispute.Title,
				Description:  dispute.Description,
				Summary:      disputeSummary.Summary,
				Category:     []string{"Dispute"}, // Assuming a default category for now
				DateFiled:    dispute.CaseDate.Format("2006-01-02"),
				DateResolved: dispute.DateResolved.Format("2006-01-02"), // Placeholder for resolved date
				Resolution:   string(dispute.Status),
			},
			Events: []models.Event{},
		}
		logger.Info("Successfully retrieved dispute")
		c.JSON(http.StatusOK, models.Response{Data: archiveDispute})
		return
	} else {
		logger.Info("Dispute not found")
		c.JSON(http.StatusNotFound, models.Response{Data: ""})
	}
}
