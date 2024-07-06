package handlers

import (
	"api/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Archive struct {
	DB *gorm.DB
}

func NewArchiveHandler(db *gorm.DB) Archive {
	return Archive{DB: db}
}

func SetupArchiveRoutes(g *gin.RouterGroup, h Archive) {
	g.POST("/search", h.getSummaryListOfArchives)
	g.GET("/:id", h.getArchive)
}

// @Summary Get a summary list of archives
// @Description Get a summary list of archives
// @Tags archive
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Archive Summary Endpoint"
// @Router /archive [post]
func (h Archive) getSummaryListOfArchives(c *gin.Context) {
	var body models.ArchiveSearchRequest
	if err := c.BindJSON(&body); err != nil {
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

	query = query.Where("resolved = ?", true)

	// Apply sorting
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	// Apply pagination
	query = query.Offset(offset).Limit(limit)

	// Execute the query
	if err := query.Find(&disputes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving disputes"})
		return
	}

	if len(disputes) == 0 {
		c.JSON(http.StatusOK, models.Response{Data: []models.ArchivedDisputeSummary{}})
		return
	}

	// Transform the results to ArchivedDisputeSummary
	var archiveDisputeSummaries []models.ArchivedDisputeSummary
	for _, dispute := range disputes {
		archiveDisputeSummaries = append(archiveDisputeSummaries, models.ArchivedDisputeSummary{
			ID:           *dispute.ID,
			Title:        dispute.Title,
			Summary:      dispute.Description,
			Category:     []string{"Dispute"}, // Assuming a default category for now
			DateFiled:    dispute.CaseDate,
			DateResolved: dispute.CaseDate.Add(48 * time.Hour), // Placeholder for resolved date
			Resolution:   string(dispute.Decision),
		})
	}

	// Return the response
	c.JSON(http.StatusOK, models.Response{Data: archiveDisputeSummaries})
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

	id := c.Param("id")

	intID, err := strconv.Atoi(id)
	if err != nil {
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
		c.JSON(http.StatusNotFound, models.Response{Data: ""})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving dispute"})
		return
	}

	//transform to archive dispute
	var archiveDispute models.ArchivedDispute
	if *dispute.ID != 0 {
		archiveDispute = models.ArchivedDispute{
			ArchivedDisputeSummary: models.ArchivedDisputeSummary{
				ID:           *dispute.ID,
				Title:        dispute.Title,
				Summary:      dispute.Description,
				Category:     []string{"Dispute"}, // Assuming a default category for now
				DateFiled:    dispute.CaseDate,
				DateResolved: dispute.CaseDate.Add(48 * time.Hour), // Placeholder for resolved date
				Resolution:   string(dispute.Decision),
			},
			Events: []models.Event{},
		}
		c.JSON(http.StatusOK, models.Response{Data: archiveDispute})
		return
	} else {
		c.JSON(http.StatusNotFound, models.Response{Data: ""})
	}

}
