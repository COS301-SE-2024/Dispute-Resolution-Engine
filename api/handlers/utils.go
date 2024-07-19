package handlers

import (
	"api/models"
	"api/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Utility) GetCountries(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	var countries []models.Country
	err := h.DB.Find(&countries).Error
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve countries")
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	logger.Info("Countries retrieved successfully")
	c.JSON(http.StatusOK, models.Response{Data: countries})
}
