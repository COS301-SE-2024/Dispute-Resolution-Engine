package handlers

import (
	"api/models"
	"api/utilities"
	"net/http"

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
	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "Archive Summary Endpoint"})
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

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "Archive Detail Endpoint for ID: " + id})
}
