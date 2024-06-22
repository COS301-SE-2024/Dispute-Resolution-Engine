package handlers

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"encoding/json"
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

	//archive routes
	archiveRouter := router.PathPrefix("/archive").Subrouter()
	archiveRouter.HandleFunc("/search", h.getSummaryListOfArchives).Methods(http.MethodPost)
	archiveRouter.HandleFunc("/{id}", h.getArchive).Methods(http.MethodPost)
}

// @Summary Get a summary list of disputes
// @Description Get a summary list of disputes
// @Tags dispute
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Dispute Summary Endpoint"
// @Router /dispute [get]
func (h Handler) getSummaryListOfDisputes(w http.ResponseWriter, r *http.Request) {
	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "Dispute Summary Endpoint"})
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

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "Dispute Detail Endpoint for ID: " + id})
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
	//get the request body
	var body models.ArchiveSearchRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: "Invalid request body, could not parse JSON"})
		return
	}

	//handle the request
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
	}
	if body.Sort != nil {
		sort = string(*body.Sort)
	}

	//mock response
	archiveDisputeSummaries := getMockArchiveDisputeSummaries()

	//filter the summaries
	archiveDisputeSummaries = filterSummariesBySearch(archiveDisputeSummaries, searchTerm)

	//sort the summaries
	sortSummaries(archiveDisputeSummaries, sort, order)

	//paginate the summaries
	archiveDisputeSummaries = paginateSummaries(archiveDisputeSummaries, offset, limit)

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
	body := models.ArchivedDispute{
		ArchivedDisputeSummary: models.ArchivedDisputeSummary{
			ID:           int64(intID),
			Title:        "Dispute " + id,
			Summary:      "Summary " + id,
			Category:     []string{"Category " + id},
			DateFiled:    time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			DateResolved: time.Date(2021, time.January, 2, 0, 0, 0, 0, time.UTC),
			Resolution:   "Resolution " + id,
		},
		Events: []models.Event{
			{
				Timestamp:   "2021-01-01T00:00:00Z",
				Type:        "Type 1",
				Description: "Details 1",
			},
			{
				Timestamp:   "2021-01-02T00:00:00Z",
				Type:        "Type 2",
				Description: "Details 2",
			},
		},
	}

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: body})
}
