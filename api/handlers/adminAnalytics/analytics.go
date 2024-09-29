package adminanalytics

import (
	"api/models"
	"api/utilities"

	"github.com/gin-gonic/gin"
)

func SetupAnalyticsRoute(router *gin.RouterGroup, h AdminAnalyticsHandler){
	router.GET("/time/estimation", h.GetTimeEstimation)
	router.POST("/dispute/grouping", h.GetDisputeSummary) //by status or country
	router.GET("/stats/tickets", h.GetTicketStats)
	router.POST("/stats/tickets/summary", h.GetTicketSummary)
	router.GET("/stats/disputes", h.GetDisputeStats)
	router.POST("/stats/disputes/summary", h.GetDisputeSummary)
	router.GET("/stats/experts", h.GetExpertStats)
	router.POST("/stats/experts/summary", h.GetExpertSummary)

}

func (h AdminAnalyticsHandler) GetTimeEstimation(c *gin.Context){
	logger := utilities.NewLogger().LogWithCaller()
	if !h.IsAuthorized(c, "admin", logger) {
		return
	}
	
}

func (h AdminAnalyticsHandler) GetDisputeSummary(c *gin.Context){
	logger := utilities.NewLogger().LogWithCaller()
	if !h.IsAuthorized(c, "admin", logger) {
		return
	}

}

func (h AdminAnalyticsHandler) GetTicketStats(c *gin.Context){
	logger := utilities.NewLogger().LogWithCaller()
	if !h.IsAuthorized(c, "admin", logger) {
		return
	}
	
}

func (h AdminAnalyticsHandler) GetTicketSummary(c *gin.Context){
	logger := utilities.NewLogger().LogWithCaller()
	if !h.IsAuthorized(c, "admin", logger) {
		return
	}
	
}


func (h AdminAnalyticsHandler) GetDisputeStats(c *gin.Context){
	logger := utilities.NewLogger().LogWithCaller()
	if !h.IsAuthorized(c, "admin", logger) {
		return
	}
	
}


func (h AdminAnalyticsHandler) GetExpertStats(c *gin.Context){
	logger := utilities.NewLogger().LogWithCaller()
	if !h.IsAuthorized(c, "admin", logger) {
		return
	}
	
}

func (h AdminAnalyticsHandler) GetExpertSummary(c *gin.Context){
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