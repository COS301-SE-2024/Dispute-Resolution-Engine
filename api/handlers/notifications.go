package handlers

import (
	"api/middleware"
	"api/models"
	"net/http"
	"os"

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

	email := models.Email{
		From:    os.Getenv("COMPANY_EMAIL"),
		To:      respondentEmail,
		Subject: "Notification of formal dispute",
		Body:    "Dear valued respondent,\n We hope this email finds you well. A dispute has arisen between you and a user of our system. Please login to your DRE account and review it, if you do not have an account you may create one.",
	}

	if err := sendMail(email); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error sending notification email."})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Email notification send successfully"})

}
