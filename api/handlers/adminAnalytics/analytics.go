package adminanalytics

import (
	"api/models"
	"api/utilities"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupAnalyticsRoute(router *gin.RouterGroup, h AdminAnalyticsHandler) {
	router.GET("/time/estimation", h.GetTimeEstimation)
	router.GET("/dispute/countries", h.GetDisputeGrouping) //by status or country
	router.POST("/stats/:table", h.GetTableStats)

}

func (h AdminAnalyticsHandler) GetTimeEstimation(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	if !h.IsAuthorized(c, "admin", logger) {
		return
	}

	days, err := h.DB.CalculateAverageResolutionTime()
	if err != nil {
		logger.WithError(err).Error("Failed to calculate average resolution time")
		c.JSON(500, models.Response{Error: "Failed to calculate average resolution time"})
		return
	}

	// Convert the average time (in days) to a time.Duration
	duration := time.Duration(days * float64(24*time.Hour))

	c.JSON(200, models.Response{Data: duration})
}

func (h AdminAnalyticsHandler) GetDisputeGrouping(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	if !h.IsAuthorized(c, "admin", logger) {
		return
	}
	grouping, err := h.DB.GetDisputeGroupingByCountry()
	if err != nil {
		logger.WithError(err).Error("Failed to get dispute grouping by country")
		c.JSON(500, models.Response{Error: "Failed to get dispute grouping by country"})
		return
	}

	c.JSON(200, models.Response{Data: grouping})
}

func (h AdminAnalyticsHandler) GetTableStats(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	if !h.IsAuthorized(c, "admin", logger) {
		return
	}

	table := c.Param("table")
	if table == "" {
		logger.Error("Invalid request")
		c.JSON(400, models.Response{Error: "Invalid request"})
		return
	}

	// Get body
	var body models.AdminTableStats
	if err := c.BindJSON(&body); err != nil {
		logger.WithError(err).Error("Failed to bind request body")
		c.JSON(400, models.Response{Error: "Failed to bind request body"})
		return
	}

	var value interface{}
	if body.Where != nil {
		// If Where is provided, use its value. Otherwise, leave value as nil.
		value = body.Where.Value
	}

	// Call the CountRecordsWithGroupBy function
	resCount, err := h.DB.CountRecordsWithGroupBy(table, &body.Where.Column, &value, body.Group)
	if err != nil {
		logger.WithError(err).Error("Failed to count records")
		c.JSON(500, models.Response{Error: "Failed to count records"})
		return
	}

	// Return the result as a JSON response
	c.JSON(200, models.Response{Data: resCount})
}



func (h AdminAnalyticsHandler) IsAuthorized(c *gin.Context, role string, logger *utilities.Logger) bool {
	claims, err := h.JWT.GetClaims(c)
	if err != nil || claims.Role != role {
		logger.WithError(err).Error("Unauthorized")
		c.JSON(401, models.Response{Error: "Unauthorized"})
		return false
	}
	return true
}
