package handlers

import (
	"api/models"
	"api/utilities"
	"net/http"
)

func (h handler) GetCountries(w http.ResponseWriter, r *http.Request) {
	var countries []models.Country
	err := h.DB.Find(&countries).Error
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: countries})
}
