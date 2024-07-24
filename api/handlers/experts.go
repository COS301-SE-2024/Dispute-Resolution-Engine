package handlers

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupExpertRoutes(g *gin.RouterGroup, h Expert) {
	g.Use(middleware.JWTMiddleware)

	g.POST("/assign", h.recommendExpert)
	g.POST("/reject", h.rejectExpert)
	g.GET("/dispute/:dispute_id", h.getDisputeExperts)
	g.GET("/dispute", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.Response{Data: "Dispute experts"})
	})
}

func (h Expert) recommendExpert(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	
	//get the dispute id from the request
	var recommendexpert models.RecommendExpert
	if err := c.BindJSON(&recommendexpert); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	// use algorithm to find recommended expert
	// currently select the first 4

	var users []models.User
	roles := []string{"Mediator", "Adjudicator", "Arbitrator"}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Query for users with the specified roles
	h.DB.Where("role IN ?", roles).Find(&users)

	// Shuffle the results and take the first 4
	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })

	// Select the first 4 users after shuffle
	selectedUsers := users
	if len(users) > 4 {
		selectedUsers = users[:4]
	}

	// insert the selected experts into the dispute_experts table
	for _, expert := range selectedUsers {
		// insert the expert into the dispute_experts table

		h.DB.Create(&models.DisputeExpert{
			Dispute: int64(recommendexpert.DisputeId),
			User:    expert.ID,
			ComplainantVote: "Approved",
			RespondantVote:  "Approved",
			ExpertVote:      "Approved",
			Status:  "Approved",
		})
	}
	logger.Info("Recommended experts successfully")
	c.JSON(http.StatusOK, models.Response{Data: selectedUsers})
}

func (h Expert) rejectExpert(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	//get dispute ID and expert ID from the request
	var rejectexpert models.RejectExpert
	if err := c.BindJSON(&rejectexpert); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	//set status to rejected

	h.DB.Model(&models.DisputeExpert{}).Where("dispute = ? AND dispute_experts.user = ?", rejectexpert.DisputeId, rejectexpert.ExpertId).Update("status", "rejected")
	logger.Info("Expert rejected successfully")
	c.JSON(http.StatusOK, models.Response{Data: "Expert rejected"})
}

func (h Expert) getDisputeExperts(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	//get dispute ID from the request
	disputeID := c.Param("dispute_id")

	//get the experts assigned to the dispute
	var disputeExperts []models.DisputeExpert
	h.DB.Where("dispute = ?", disputeID).Find(&disputeExperts)
	logger.Info("Dispute experts retrieved successfully")
	c.JSON(http.StatusOK, models.Response{Data: disputeExperts})
}
