package notifications

import (
	"api/env"
	"api/models"
	"api/utilities"
	"crypto/tls"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)





type EmailSystem interface {
	SendAdminEmail(c *gin.Context, disputeID int64, resEmail string, title string, summary string)
	NotifyDisputeStateChanged(c *gin.Context, disputeID int64, disputeStatus , description string)
	SendDefaultUserEmail(c *gin.Context, email string, pass string, title string, summary string)
}

type emailImpl struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) EmailSystem {
	return &emailImpl{db: db}
}

func (e *emailImpl) SendAdminEmail(c *gin.Context, disputeID int64, resEmail string, title string, summary string) {
	logger := utilities.NewLogger().LogWithCaller()
	envReader := env.NewEnvLoader()

	//get the dispute details using ID from request body
	var dbDispute models.Dispute
	e.db.Where("id = ?", disputeID).First(&dbDispute)

	var respondentEmail = resEmail

	companyEmail, err := envReader.Get("COMPANY_EMAIL")
	if err != nil {
		utilities.InternalError(c)
		return
	}

	email := models.Email{
		From:    companyEmail,
		To:      respondentEmail,
		Subject: "Notification of formal dispute",
		Body: "Dear valued respondent,\r\n\r\n We hope this email finds you well. A dispute has arisen between you and a user of our system.\r\n\r\n The dispute details are as followed:\r\n" +
			"Title: " + title + ".\r\n" +
			"Summary of dispute: " + summary + ".\r\n" +
			"Please login to your DRE account and review it .",
	}

	go SendMail(email)
	logger.Info("Admin email notification sent successfully")
}

func (e *emailImpl) SendDefaultUserEmail(c *gin.Context, email string, pass string, title string, summary string) {
	logger := utilities.NewLogger().LogWithCaller()
	envReader := env.NewEnvLoader()

	companyEmail, err := envReader.Get("COMPANY_EMAIL")
	if err != nil {
		utilities.InternalError(c)
		return
	}

	emailBody := models.Email{
		From:    companyEmail,
		To:      email,
		Subject: "Default DRE Account",
		Body: "Dear valued respondent,\r\n\r\n We hope this email finds you well. A dispute has arisen between you and a user of our system.\r\n\r\n The dispute details are as followed: \r\n" +
			"Title: " + title + ".\r\n\r\n" +
			"Summary of dispute: " + summary + ".\r\n\r\n" +
			"Please login to your DRE account and review it, use this inbox's email address along with this password: " + pass + " .",
	}

	go SendMail(emailBody)
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

	if err := SendMail(email1); err != nil {
		logger.WithError(err).Error("Failed to send respondent email")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal server error notifying"})
	}

	if err := SendMail(email2); err != nil {
		logger.WithError(err).Error("Failed to send complainant email")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal server error notifying"})
	}
	logger.Info("Email notifications sent successfully")
	c.JSON(http.StatusOK, models.Response{Error: "Email notifications sent successfully"})
}*/

func (e *emailImpl) NotifyDisputeStateChanged(c *gin.Context, disputeID int64, disputeStatus, description string) {
	logger := utilities.NewLogger().LogWithCaller()
	envLoader := env.NewEnvLoader()




	var dbDispute models.Dispute
	e.db.Where("\"workflow\" = ?", disputeID).First(&dbDispute)

	logger.Info("dispute:", dbDispute)
	
	var respondent models.User
	var complainant models.User
	err := e.db.First(&respondent, *dbDispute.Respondant)
	if err.Error != nil {
		logger.Error("Failed to get the respondent details:", err.Error)
		return
	}
	err = e.db.First(&complainant, dbDispute.Complainant)
	if err.Error != nil {
		logger.WithError(err.Error).Error("Failed to get the complainant details")
		return
	}
	body := "Dear valued user,\r\n\r\nWe hope this email finds you well. The status of a dispute you are involved with has changed. \r\n\r\n"
	body += "The dispute details are as follows:\r\n"
	body += "Current status: " + disputeStatus + ".\r\n"
	body += description + "\r\n\r\n"
	body += "Please visit DRE and check your emails regularly for future updates."

	companyEmail, err2 := envLoader.Get("COMPANY_EMAIL")
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
	logger.Info("Sending mail to: ", respondent.Email, " and ", complainant.Email)
	go SendMail(emailComplainant)
	go SendMail(emailRespondent)
	logger.Info("Emails sent out")
}

func SendMail(email models.Email) error {
	envReader := env.NewEnvLoader()
	companyEmail, err := envReader.Get("COMPANY_EMAIL")
	if err != nil {
		return err
	}
	companyAuth, err := envReader.Get("COMPANY_AUTH")
	if err != nil {
		return err
	}

	logger := utilities.NewLogger().LogWithCaller()
	d := gomail.NewDialer("smtp.gmail.com", 587, companyEmail, companyAuth)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", email.From)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain", email.Body)
	logger.WithField("email", email).Info("Sending email")

	if err := d.DialAndSend(m); err != nil {
		logger.WithError(err).Error("Error sending email")
		return err
	}
	return nil
}