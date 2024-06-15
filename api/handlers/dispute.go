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
}

func (h Handler) getSummaryListOfDisputes(w http.ResponseWriter, r *http.Request) {
	// Logic for handling /disputes
	w.Write([]byte("Disputes Endpoint"))
}

func (h Handler) getDispute(w http.ResponseWriter, r *http.Request) {
	// Logic for handling /disputes/{id}
	vars := mux.Vars(r)
	id := vars["id"]

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "Dispute Detail Endpoint for ID: " + id})
}

func (h Handler) patchDispute(w http.ResponseWriter, r *http.Request) {
	// Logic for handling /disputes/{id}
	vars := mux.Vars(r)
	id := vars["id"]

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "Dispute Patch Endpoint for ID: " + id})
}
