package handlers

import (
	"api/models"
	"api/utilities"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func SetupDisputeRoutes(router *mux.Router, h Handler) {
	router.HandleFunc("", h.getSummaryListOfDisputes).Methods(http.MethodGet)
	router.HandleFunc("/{id}", h.getDispute).Methods(http.MethodGet)
	router.HandleFunc("/{id}", h.patchDispute).Methods(http.MethodPatch)

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
// @Router /archive [get]
func (h Handler) getSummaryListOfArchives(w http.ResponseWriter, r *http.Request) {
	//get the request body
	var body models.ArchiveSearchRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: "Invalid request body, could not parse JSON"})
		return
	}

	//mock ArchiveDisputeSummary array
	archiveDisputeSummaries := []models.ArchivedDisputeSummary{
		{
			ID:           1,
			Title:        "Dispute 1",
			Summary:      "Summary 1",
			Category:     []string{"Category 1"},
			DateFiled:    time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			DateResolved: time.Date(2021, time.January, 2, 0, 0, 0, 0, time.UTC),
			Resolution:   "Resolution 1",
		},
		{
			ID:           2,
			Title:        "Dispute 2",
			Summary:      "Summary 2",
			Category:     []string{"Category 2"},
			DateFiled:    time.Date(2021, time.January, 3, 0, 0, 0, 0, time.UTC),
			DateResolved: time.Date(2021, time.January, 4, 0, 0, 0, 0, time.UTC),
			Resolution:   "Resolution 2",
		},
		{
			ID:           3,
			Title:        "Dispute 3",
			Summary:      "Summary 3",
			Category:     []string{"Category 3"},
			DateFiled:    time.Date(2021, time.January, 5, 0, 0, 0, 0, time.UTC),
			DateResolved: time.Date(2021, time.January, 6, 0, 0, 0, 0, time.UTC),
			Resolution:   "Resolution 3",
		},
	}


	//handle the request
	// limit := body.Limit
	// offset := body.Offset
	// sort := body.Sort
	// order := body.Order
	// filter := body.Filter
	// search := body.Search

	

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: archiveDisputeSummaries})
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
				Timestamp: "2021-01-01T00:00:00Z",
				Type:      "Type 1",
				Description:   "Details 1",
			},
			{
				Timestamp: "2021-01-02T00:00:00Z",
				Type:      "Type 2",
				Description:   "Details 2",
			},
		},
	};

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: body})
}
