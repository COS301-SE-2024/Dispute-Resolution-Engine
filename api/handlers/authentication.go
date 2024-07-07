package handlers

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type StringWrapper struct {
	Data string `json:"Data"`
}

// Define Credentials struct globally
type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func SetupAuthRoutes(group *gin.RouterGroup, h Auth) {
	group.POST("/signup", h.CreateUser)
	group.POST("/login", h.LoginUser)
	group.POST("/verify", h.Verify)
	/*
		group.Handle("/reset-password", middleware.RoleMiddleware(http.AuthFunc(h.ResetPassword), 0)).Methods(http.MethodPost)
		// router.Handle("/verify", middleware.RoleMiddleware(http.AuthFunc(h.Verify), 0)).Methods(http.MethodPost)
	*/
}

// @Summary Reset a user's password
// @Description Reset a user's password
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Password reset not available yet..."
// @Router /auth/reset-password [post]
/*
func (h Auth) ResetPassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	c.JSON(http.StatusOK, models.Response{Data: "password reset not available yet..."})
}
*/

// @Summary Create a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "User Details"
// @Success 201 {object} models.User
// @Failure 400 {object} models.Response "Bad Request"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Router /auth/signup [post]
func (h Auth) CreateUser(c *gin.Context) {
	hasher := utilities.NewArgon2idHash(1, 12288, 4, 32, 16)

	var reqUser models.CreateUser
	if err := c.BindJSON(&reqUser); err != nil {
		return
	}

	//stub timezone
	zone, _ := time.Now().Zone()
	timezone := zone
	reqUser.Timezone = &timezone
	//Now put stuff in the actual user object
	date, _ := time.Parse("2006-01-02", reqUser.Birthdate)
	user := models.User{
		FirstName:         reqUser.FirstName,
		Surname:           reqUser.Surname,
		Birthdate:         date,
		Nationality:       reqUser.Nationality,
		Email:             reqUser.Email,
		PasswordHash:      reqUser.Password,
		PhoneNumber:       reqUser.PhoneNumber,
		AddressID:         nil,
		Status:            "Unverified",
		Gender:            reqUser.Gender,
		PreferredLanguage: reqUser.PreferredLanguage,
		Timezone:          reqUser.Timezone,
	}

	//Check if there is an existing email
	duplicate := h.checkUserExists(user.Email)

	if duplicate {
		c.JSON(http.StatusConflict, models.Response{Error: "Email already in use"})
		return
	}

	//Hash the password
	hashAndSalt := hasher.HashPassword(user.PasswordHash)
	user.PasswordHash = base64.StdEncoding.EncodeToString(hashAndSalt.Hash)
	user.Salt = base64.StdEncoding.EncodeToString(hashAndSalt.Salt)

	//update log metrics
	user.CreatedAt = utilities.GetCurrentTime()
	user.UpdatedAt = utilities.GetCurrentTimePtr()
	user.Status = "Active"

	//Small user preferences
	user.Role = "user"
	stubbedPref := "en-US"
	user.PreferredLanguage = &stubbedPref
	user.LastLogin = nil

	if result := h.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating user"})
		return
	}
	sendOTP(user.Email)

	jwt, err := middleware.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating token"})
		return
	}
	c.JSON(http.StatusCreated, models.Response{Data: jwt})
}

// @Summary Login a user
// @Description Login an existing user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body Credentials true "User Credentials"
// @Success 200 {string} string "token"
// @Failure 400 {object} string "Bad Request"
// @Failure 401 {object} string "Unauthorized"
// @Router /auth/login [post]
func (h Auth) LoginUser(c *gin.Context) {
	hasher := utilities.NewArgon2idHash(1, 12288, 4, 32, 16)

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		return
	}

	if !h.checkUserExists(user.Email) {
		c.JSON(http.StatusNotFound, models.Response{Error: "User does not exist"})
		return
	}

	var dbUser models.User
	h.DB.Where("email = ?", user.Email).First(&dbUser)

	realSalt, err := base64.StdEncoding.DecodeString(dbUser.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong processing your request..."})
		return
	}
	checkHash, err := hasher.GenerateHash([]byte(user.PasswordHash), realSalt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong processing your request..."})
		return
	}

	if dbUser.PasswordHash != base64.StdEncoding.EncodeToString(checkHash.Hash) {
		print(dbUser.PasswordHash)
		print(base64.StdEncoding.EncodeToString(checkHash.Hash))
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Invalid credentials"})
		return
	}

	dbUser.LastLogin = utilities.GetCurrentTimePtr()
	h.DB.Where("email = ?", user.Email).Update("last_login", utilities.GetCurrentTime())

	token, err := middleware.GenerateJWT(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: token})

}

// Verify verifies the user's email through a pin code
// @Summary Verify user email
// @Description Verifies the user's email by checking the provided pin code against stored values.
// @Tags auth
// @Accept json
// @Produce json
// @Param pinReq body models.VerifyUser true "Verify User"
// @Success 200 {object} interface{} "Email verified successfully - Example response: { 'message': 'Email verified successfully' }"
// @Failure 400 {object} interface{} "Invalid Request - Example error response: { 'error': 'Invalid Request' }"
// @Failure 400 {object} interface{} "Invalid pin - Example error response: { 'error': 'Invalid pin' }"
// @Failure 500 {object} interface{} "Error verifying pin - Example error response: { 'error': 'Error verifying pin' }"
// @Router /auth/verify [post]
func (h Auth) Verify(c *gin.Context) {
	var pinReq models.VerifyUser
	if err := c.BindJSON(&pinReq); err != nil {
		return
	}
	valid, err := utilities.RemoveFromFile("stubbedStorage/verify.txt", pinReq.Pin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error verifying pin"})
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid pin"})
		return
	}
	c.JSON(http.StatusOK, models.Response{Data: "Email verified successfully"})
}

func (h Handler) checkUserExists(email string) bool {
	var user models.User
	h.DB.Where("email = ?", email).First(&user)
	return user.Email != ""
}

func sendOTP(userInfo string) {
	// SMTP server configuration for Gmail
	smtpServer := "smtp.gmail.com"
	smtpPort := 587
	smtpUser := os.Getenv("COMPANY_EMAIL")
	smtpPassword := os.Getenv("COMPANY_AUTH") // Use app password if 2-factor authentication is enabled

	// Recipient email address
	to := userInfo
	pin := utilities.GenerateVerifyEmailToken()
	// Email subject and body
	subject := "Verify Account"
	body := "Hello,\nPlease verify your DRE account using this pin: " + pin + "\n\nThanks,\nTeam Techtonic."

	// Initialize the SMTP dialer
	d := gomail.NewDialer(smtpServer, smtpPort, smtpUser, smtpPassword)
	d.TLSConfig = &tls.Config{ServerName: smtpServer, InsecureSkipVerify: false}
	// Create a new email message
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email:", err)
		os.Exit(1)
	}
	err := utilities.WriteToFile(pin, "stubbedStorage/verify.txt")
	if err != nil {
		fmt.Println("Error writing to file: " + err.Error())
	}
	fmt.Println("Email sent successfully!")
}

// resetPassword sends an email to the user with a link to reset their password
// @Summary Reset a user's password
// @Description Reset a user's password
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Email sent successfully"
// @Failure 400 {object} models.Response "Invalid Request"
// @Router /auth/reset-password/send-email [post]

func (h Auth) ResetPassword(c *gin.Context) {
	defer c.Request.Body.Close()
	
	//check the body of the request

	var body models.SendResetRequest;
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Request"})
		return
	}

	//check if the email exists in the database

	var user models.User;
	h.DB.Where("email = ?", body.Email).First(&user);
	if h.checkUserExists(user.Email) {
		c.JSON(http.StatusNotFound, models.Response{Error: "User does not exist"})
		return
	}

	//generate a jwt token with the user's email
	tempUser:= models.User{
		Email: user.Email,
	}

	jwt, err := middleware.GenerateJWT(tempUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating token"})
		return
	}

	//send an email to the user with a temporary link to reset the password
	linkURL := "http://localhost:8080/reset-password/" + jwt
	email := models.Email{
		From:    os.Getenv("COMPANY_EMAIL"),
		To:      user.Email,
		Subject: "Password Reset",
		Body:    "Click here to reset your password: " + linkURL,
	}

	if err := sendMail(email); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error sending email"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Email sent successfully"})

}

//activate reset password
// @Summary Activate reset password
// @Description Activate reset password
// @Tags auth
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Success 200 {object} models.Response "Password reset link activated"
// @Failure 400 {object} models.Response "Invalid Request"
// @Router /auth/reset-password/reset [post]
func (h Auth) ActivateResetPassword(c *gin.Context) {
	defer c.Request.Body.Close()

	//get the token from the url
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, models.Response{Error: "missing token"})
		return
	}
	claims := middleware.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "missing token"})
		return
	}

	//check the body of the request

	var body models.ResetPassword;
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Request"})
		return
	}

	//get user form the database
	var user models.User;
	h.DB.Where("email = ?", claims.Email).First(&user);
	if h.checkUserExists(user.Email) {
		c.JSON(http.StatusNotFound, models.Response{Error: "User does not exist"})
		return
	}

	//hash the new password
	hashedPassword, err := hashPassword(body.NewPassword, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error hashing password"})
		return
	}

	//update the user's password
	h.DB.Model(&user).Updates(map[string]interface{}{
		"password_hash": hashedPassword,
		"salt":          user.Salt,
	})

	c.JSON(http.StatusOK, models.Response{Data: "Password reset successfully"})
}

func sendMail(email models.Email) error {
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("COMPANY_EMAIL"), os.Getenv("COMPANY_AUTH"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", email.From)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/html", email.Body)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func hashPassword(newPassword string, user *models.User) (string, error) {
	hasher := utilities.NewArgon2idHash(1, 12288, 4, 32, 16)

	hashAndSalt := hasher.HashPassword(newPassword)
	user.Salt = base64.StdEncoding.EncodeToString(hashAndSalt.Salt)
	return base64.StdEncoding.EncodeToString(hashAndSalt.Hash), nil
}
