package handlers

import (
	"api/middleware"
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func SetupNotificationRoutes(g *gin.RouterGroup, h Notification) {
// 	g.POST("/invite", h.sendAdminNotification)
// }

func (h Notification) sendAdminNotification(c *gin.Context, resEmail string) {

	//get claims
	jwtClaims := middleware.GetClaims(c)
	if jwtClaims == nil {
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	//parse the request body
	var reqInv models.DisputeNotifyInvite
	if err := c.BindJSON(&reqInv); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request"})
		return
	}

	//get the dispute details using ID from request body
	var dbDispute models.Dispute
	h.DB.Where("id = ?", reqInv.DisputeID).First(&dbDispute)

	var respondentEmail = resEmail
	//For now will stub it, but here we will finalize and send the appropriate email to invite a party
	
}
