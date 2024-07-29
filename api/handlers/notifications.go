package handlers

import (
	"api/env"
	"api/models"
	"api/utilities"

	"github.com/gin-gonic/gin"
)

// func SetupNotificationRoutes(g *gin.RouterGroup, h Notification) {
// 	g.POST("/invite", h.sendAdminNotification)
// }

func (h Handler) sendAdminNotification(c *gin.Context, disputeID int64, resEmail string) {
	logger := utilities.NewLogger().LogWithCaller()

	//get the dispute details using ID from request body
	var dbDispute models.Dispute
	h.DB.Where("id = ?", disputeID).First(&dbDispute)

	var respondentEmail = resEmail

	companyEmail, err := env.Get("COMPANY_EMAIL")
	if err != nil {
		utilities.InternalError(c)
		return
	}

	email := models.Email{
		From:    companyEmail,
		To:      respondentEmail,
		Subject: "Notification of formal dispute",
		Body:    "Dear valued respondent,\n We hope this email finds you well. A dispute has arisen between you and a user of our system. Please login to your DRE account and review it, if you do not have an account you may create one.",
	}

	go sendMail(email)
	logger.Info("Admin email notification sent successfully")
}

/*func (h Notification) AcceptanceNotification(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	jwtClaims := middleware.GetClaims(c)
	if jwtClaims == nil {
		logger.Error("Unauthorized")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
	}

	var respondentEmail = jwtClaims.Email

	var disputeDetails models.Dispute
	h.DB.Where("id = ?", reqNotif.DisputeID).First(&disputeDetails)

	var dbComplainant models.User
	h.DB.Where("id = ?", disputeDetails.Complainant).First(&dbComplainant)

	var complainantEmail = dbComplainant.Email

	email1 := models.Email{
		From:    env.Get("COMPANY_EMAIL"),
		To:      respondentEmail,
		Subject: "Formal Dispute Accepted",
		Body:    "Dear users,\n A dispute has been accepted from both parties and will now commence, please stay active on DRE to always be up to date.",
	}

	email2 := models.Email{
		From:    env.Get("COMPANY_EMAIL"),
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
}*/

func (h Handler) StateChangeNotifications(c *gin.Context, disputeID int64, disputeStatus string) {
	logger := utilities.NewLogger().LogWithCaller()

	var dbDispute models.Dispute
	h.DB.Where("id = ?", disputeID).First(&dbDispute)

	var respondent models.User
	var complainant models.User
	err := h.DB.Where("id = ?", dbDispute.Respondant).First(&respondent)
	if err != nil {
		logger.WithError(err.Error).Error("Failed to get the respondent details")
		return
	}
	err = h.DB.Where("id = ?", dbDispute.Complainant).First(&complainant)
	if err != nil {
		logger.WithError(err.Error).Error("Failed to get the complainant details")
		return
	}
	body := "Dear valued user,\n We hope this email finds you well. The status of a dispute you are involved with has changed to " + disputeStatus + ". Please visit DRE and check your emails regularly for future updates."
	companyEmail, err2 := env.Get("COMPANY_EMAIL")
	if err2 != nil {
		utilities.InternalError(c)
		return
	}
	emailRespondent := models.Email{
		From:    companyEmail,
		To:      respondent.Email,
		Subject: "Dispute Status Change",
		Body:    body,
	}
	emailComplainant := models.Email{
		From:    companyEmail,
		To:      complainant.Email,
		Subject: "Dispute Status Change",
		Body:    body,
	}
	go sendMail(emailComplainant)
	go sendMail(emailRespondent)
	logger.Info("Emails sent out")
}
