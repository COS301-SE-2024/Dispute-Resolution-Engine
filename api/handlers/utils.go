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
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to retrieve countries"})
		return
	}
	logger.Info("Countries retrieved successfully")
	c.JSON(http.StatusOK, models.Response{Data: countries})
}

func (h Utility) GetDisputeStatuses(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	var statuses []string
	err := h.DB.Raw("SELECT e.enumlabel AS enum_value FROM pg_type t JOIN pg_enum e ON t.oid = e.enumtypid JOIN pg_namespace n ON t.typnamespace = n.oid WHERE t.typname = 'dispute_status' ORDER BY e.enumsortorder;").Scan(&statuses).Error
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve dispute statuses")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to retrieve dispute statuses"})
		return
	}
	logger.Info("Dispute statuses retrieved successfully")
	c.JSON(http.StatusOK, models.Response{Data: statuses})
}
