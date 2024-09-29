package adminanalytics

import (
	"api/models"
	"api/utilities"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupAnalyticsRoute(router *gin.RouterGroup, h AdminAnalyticsHandler) {
	router.GET("/time/estimation", h.GetTimeEstimation)
	router.POST("/dispute/grouping", h.GetDisputeGrouping) //by status or country
	router.POST("/stats/:table", h.GetTableStats)
	router.GET("/stats/:table/summary", h.GetTableSummary)

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

	// Parse the request body
	var groupingRequest models.AdminGroupingAnalytics
	if err := c.BindJSON(&groupingRequest); err != nil {
		logger.WithError(err).Error("Failed to parse request body")
		c.JSON(400, models.Response{Error: "Failed to parse request body"})
		return
	}

	if groupingRequest.Group == nil || (*groupingRequest.Group != "status" && *groupingRequest.Group != "country") {
		logger.Error("Invalid request body")
		c.JSON(400, models.Response{Error: "Invalid request body"})
		return
	}

	if *groupingRequest.Group == "status" {
		grouping, err := h.DB.GetDisputeGroupingByStatus()
		if err != nil {
			logger.WithError(err).Error("Failed to get dispute grouping by status")
			c.JSON(500, models.Response{Error: "Failed to get dispute grouping by status"})
			return
		}

		c.JSON(200, models.Response{Data: grouping})
		return
	} else {
		grouping, err := h.DB.GetDisputeGroupingByCountry()
		if err != nil {
			logger.WithError(err).Error("Failed to get dispute grouping by country")
			c.JSON(500, models.Response{Error: "Failed to get dispute grouping by country"})
			return
		}

		c.JSON(200, models.Response{Data: grouping})
		return
	}
}

func (h AdminAnalyticsHandler) GetTableStats(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	if !h.IsAuthorized(c, "admin", logger) {
		return
	}

}

func (h AdminAnalyticsHandler) GetTableSummary(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	if !h.IsAuthorized(c, "admin", logger) {
		return
	}

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
