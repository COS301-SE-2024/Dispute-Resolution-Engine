package handlers

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func SetupDisputeRoutes(router *mux.Router, h Handler) {
	//dispute routes
	disputeRouter := router.PathPrefix("").Subrouter()
	disputeRouter.Use(middleware.JWTMiddleware)
	disputeRouter.HandleFunc("", h.getSummaryListOfDisputes).Methods(http.MethodGet)
	disputeRouter.HandleFunc("/{id}", h.getDispute).Methods(http.MethodGet)
	disputeRouter.HandleFunc("/{id}", h.patchDispute).Methods(http.MethodPatch)
	disputeRouter.HandleFunc("/create", h.createDispute).Methods(http.MethodPost)

	//archive routes
	archiveRouter := router.PathPrefix("/archive").Subrouter()
	archiveRouter.HandleFunc("/search", h.getSummaryListOfArchives).Methods(http.MethodPost)
	archiveRouter.HandleFunc("/{id}", h.getArchive).Methods(http.MethodGet)
}

// @Summary Get a summary list of disputes
// @Description Get a summary list of disputes
// @Tags dispute
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Dispute Summary Endpoint"
// @Router /dispute [get]
func (h Handler) getSummaryListOfDisputes(w http.ResponseWriter, r *http.Request) {
	var disputes []models.DisputeSummaryResponse
	err := h.DB.Raw("SELECT id, title, description, status FROM disputes").Scan(&disputes).Error
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: disputes})
}

// @Summary Get a dispute
// @Description Get a dispute
// @Tags dispute
// @Accept json
// @Produce json
// @Param id path string true "Dispute ID"
// @Success 200 {object} models.Response "Dispute Detail Endpoint"
// @Router /dispute/{id} [get]
func (h Handler) getDispute(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	// var

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "Dispute Detail Endpoint for ID: " + id})
}

func (h Handler) createDispute(w http.ResponseWriter, r *http.Request) {

	formData := false
	jsonData := false

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		jsonData = true
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		formData = true
	}

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "Contains FormData: " + strconv.FormatBool(formData) + " Contains JSON: " + strconv.FormatBool(jsonData)})

	// if err != nil {
	// 	utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: "Invalid request body"})
	// 	return
	// }
	// var createReq models.CreateDispute
	// if err := json.Unmarshal(body, &createReq); err != nil {
	// 	utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: "Invalid request body, could not parse JSON"})
	// 	return
	// }

	// // Create a new dispute
	// dispute := models.Dispute{
	// 	Title:       createReq.Title,
	// 	Description: createReq.Description,
	// }

	// if err := h.DB.Create(&dispute).Error; err != nil {
	// 	utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Error creating dispute"})
	// 	return
	// }

	// utilities.WriteJSON(w, http.StatusCreated, models.Response{Data: "Success, id: " + string(dispute.ID)})
}

// @Summary Update a dispute
// @Description Update a dispute
// @Tags dispute
// @Accept json
// @Produce json
// @Param id path string true "Dispute ID"
// @Success 200 {object} models.Response "Dispute Patch Endpoint"
// @Router /dispute/{id} [patch]
func (h Handler) patchDispute(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "Dispute Patch Endpoint for ID: " + id})
}

// @Summary Get a summary list of archives
// @Description Get a summary list of archives
// @Tags archive
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Archive Summary Endpoint"
// @Router /archive [post]
func (h Handler) getSummaryListOfArchives(w http.ResponseWriter, r *http.Request) {
	// Get the request body
	var body models.ArchiveSearchRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: "Invalid request body, could not parse JSON"})
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
			utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: "Invalid order value"})
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
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Error retrieving disputes"})
		return
	}

	if len(disputes) == 0 {
		utilities.WriteJSON(w, http.StatusNotFound, models.Response{Data: []models.ArchivedDisputeSummary{}})
		return
	}

	// Transform the results to ArchivedDisputeSummary
	var archiveDisputeSummaries []models.ArchivedDisputeSummary
	for _, dispute := range disputes {
		archiveDisputeSummaries = append(archiveDisputeSummaries, models.ArchivedDisputeSummary{
			ID:           dispute.ID,
			Title:        dispute.Title,
			Summary:      dispute.Description,
			Category:     []string{"Dispute"}, // Assuming a default category for now
			DateFiled:    dispute.CaseDate,
			DateResolved: dispute.CaseDate.Add(48 * time.Hour), // Placeholder for resolved date
			Resolution:   string(dispute.Decision),
		})
	}

	// Return the response
	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: archiveDisputeSummaries})
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

// @Summary Get an archive
// @Description Get an archive
// @Tags archive
// @Accept json
// @Produce json
// @Param id path string true "Archive ID"
// @Success 200 {object} models.Response "Archive Detail Endpoint"
// @Router /archive/{id} [get]
func (h Handler) getArchive(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	intID, err := strconv.Atoi(id)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: "Invalid ID"})
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
		utilities.WriteJSON(w, http.StatusNotFound, models.Response{Data: ""})
		return
	} else if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Error retrieving dispute"})
		return
	}

	//transform to archive dispute
	var archiveDispute models.ArchivedDispute
	if dispute.ID != 0 {
		archiveDispute = models.ArchivedDispute{
			ArchivedDisputeSummary: models.ArchivedDisputeSummary{
				ID:           dispute.ID,
				Title:        dispute.Title,
				Summary:      dispute.Description,
				Category:     []string{"Dispute"}, // Assuming a default category for now
				DateFiled:    dispute.CaseDate,
				DateResolved: dispute.CaseDate.Add(48 * time.Hour), // Placeholder for resolved date
				Resolution:   string(dispute.Decision),
			},
			Events: []models.Event{},
		}
		utilities.WriteJSON(w, http.StatusOK, models.Response{Data: archiveDispute})
		return
	} else {
		utilities.WriteJSON(w, http.StatusNotFound, models.Response{Data: ""})
	}

}
