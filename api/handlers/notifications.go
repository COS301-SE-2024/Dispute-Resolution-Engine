package handlers

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// func SetupNotificationRoutes(g *gin.RouterGroup, h Notification) {
// 	g.POST("/invite", h.sendAdminNotification)
// }

func (h Handler) sendAdminNotification(c *gin.Context, resEmail string) {
	logger := utilities.NewLogger().LogWithCaller()
	//get claims
	jwtClaims := middleware.GetClaims(c)
	if jwtClaims == nil {
		logger.Error("Unauthorized")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	//parse the request body
	var reqInv models.DisputeNotify
	if err := c.BindJSON(&reqInv); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
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
		logger.WithError(err).Error("Failed to send admin email")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error sending admin notification email."})
		return
	}
	logger.Info("Admin email notification sent successfully")
	c.JSON(http.StatusOK, models.Response{Data: "Admin email notification sent successfully"})
}

func (h Notification) AcceptanceNotification(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	jwtClaims := middleware.GetClaims(c)
	if jwtClaims == nil {
		logger.Error("Unauthorized")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
	}

	var respondentEmail = jwtClaims.Email
	var reqNotif models.DisputeNotify

	if err := c.BindJSON(&reqNotif); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Bad request body"})
	}

	var disputeDetails models.Dispute
	h.DB.Where("id = ?", reqNotif.DisputeID).First(&disputeDetails)

	var dbComplainant models.User
	h.DB.Where("id = ?", disputeDetails.Complainant).First(&dbComplainant)

	var complainantEmail = dbComplainant.Email

	email1 := models.Email{
		From:    os.Getenv("COMPANY_EMAIL"),
		To:      respondentEmail,
		Subject: "Formal Dispute Accepted",
		Body:    "Dear users,\n A dispute has been accepted from both parties and will now commence, please stay active on DRE to always be up to date.",
	}

	email2 := models.Email{
		From:    os.Getenv("COMPANY_EMAIL"),
		To:      complainantEmail,
		Subject: "Formal Dispute Accepted",
		Body:    "Dear users,\n A dispute has been accepted from both parties and will now commence, please stay active on DRE to always be up to date.",
	}

	if err := sendMail(email1); err != nil {
		logger.WithError(err).Error("Failed to send respondent email")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal server error notifying"})
	}

	if err := sendMail(email2); err != nil {
		logger.WithError(err).Error("Failed to send complainant email")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal server error notifying"})
	}
	logger.Info("Email notifications sent successfully")
	c.JSON(http.StatusOK, models.Response{Error: "Email notifications sent successfully"})
}
