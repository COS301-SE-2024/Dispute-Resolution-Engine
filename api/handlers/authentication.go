package handlers

import (
	"api/env"
	"api/middleware"
	"api/models"
	"api/redisDB"
	"api/utilities"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
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
	group.POST("/reset-password/send-email", h.ResetPassword)
	group.POST("/reset-password/reset", h.ActivateResetPassword)
	group.POST("/resend-otp", h.ResendOTP)
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
	logger := utilities.NewLogger().LogWithCaller()

	var reqUser models.CreateUser
	if err := c.BindJSON(&reqUser); err != nil {
		logger.WithError(err).Error("Invalid Request")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Request"})
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
		logger.Error("Email already in use")
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
	user.LastLogin = nil

	// if result := h.DB.Create(&user); result.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating user"})
	// 	return
	// }

	jwt, err := middleware.GenerateJWT(user)
	if err != nil {
		logger.WithError(err).Error("Error getting user jwt.")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating token"})
		return
	}

	// generate a pin code
	pin := utilities.GenerateVerifyEmailToken()

	// create UserVerify struct in redis
	userkey := user.Email + user.Surname
	userVerify := models.ConvertUserToUserVerify(user, pin)

	//convert to json
	userVerifyJSON, err := json.Marshal(userVerify)
	if err != nil {
		logger.WithError(err).Error("Error marshalling user-verify JSON.")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating token"})
	}

	//store in redis cacher
	err = redisDB.RDB.Set(context.Background(), userkey, userVerifyJSON, 24*time.Hour).Err()
	if err != nil {
		logger.WithError(err).Error("Error storing OTP in redis.")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating token"})
		return
	}
	logger.Info("OTP generated")
	//send OTP
	go sendOTP(user.Email, pin)

	logger.Info("User created successfully")
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
	logger := utilities.NewLogger().LogWithCaller()

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		logger.WithError(err).Error("Invalid Request")
		return
	}

	if !h.checkUserExists(user.Email) {
		logger.Error("User does not exist")
		c.JSON(http.StatusNotFound, models.Response{Error: "User does not exist"})
		return
	}

	var dbUser models.User
	h.DB.Where("email = ?", user.Email).First(&dbUser)

	realSalt, err := base64.StdEncoding.DecodeString(dbUser.Salt)
	if err != nil {
		logger.WithError(err).Error("Error decoding salt")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong processing your request..."})
		return
	}
	checkHash, err := hasher.GenerateHash([]byte(user.PasswordHash), realSalt)
	if err != nil {
		logger.WithError(err).Error("Error hashing password")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong processing your request..."})
		return
	}

	if dbUser.PasswordHash != base64.StdEncoding.EncodeToString(checkHash.Hash) {
		print(dbUser.PasswordHash)
		print(base64.StdEncoding.EncodeToString(checkHash.Hash))
		logger.Error("Invalid credentials")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Invalid credentials"})
		return
	}

	dbUser.LastLogin = utilities.GetCurrentTimePtr()
	h.DB.Where("email = ?", user.Email).Update("last_login", utilities.GetCurrentTime())

	token, err := middleware.GenerateJWT(dbUser)
	if err != nil {
		logger.WithError(err).Error("Error generating token")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating token"})
		return
	}

	logger.Info("User logged in successfully")
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
	logger := utilities.NewLogger().LogWithCaller()

	var pinReq models.VerifyUser
	if err := c.BindJSON(&pinReq); err != nil {
		logger.WithError(err).Error("Invalid Request")
		return
	}
	var valid bool
	valid = false
	jwtClaims := middleware.GetClaims(c)
	if jwtClaims == nil {
		logger.Error("No claims found")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Request"})
		return
	}
	userkey := jwtClaims.Email + jwtClaims.User.Surname
	// valid, err := utilities.RemoveFromFile("stubbedStorage/verify.txt", pinReq.Pin)
	userVerifyJSON, err := redisDB.RDB.Get(context.Background(), userkey).Result()
	if err != nil {
		logger.WithError(err).Error("Error getting userVerify from redis")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error verifying pin"})
		return
	}
	//unmarshal result
	var userVerify models.UserVerify
	err = json.Unmarshal([]byte(userVerifyJSON), &userVerify)
	if err != nil {
		logger.WithError(err).Error("Error unmarshalling userVerify")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error verifying pin"})
		return
	}

	//check if the pin is correct
	pin := userVerify.Pin
	if pin == pinReq.Pin {
		valid = true
	}
	if err != nil {
		logger.WithError(err).Error("Error verifying pin")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error verifying pin"})
		return
	}
	if !valid {
		logger.Error("Invalid pin")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid pin"})
		return
	}
	logger.Info("Email verified successfully")

	// insert the user into database with updated status status to verified
	user := models.ConvertUserVerifyToUser(userVerify)
	user.Status = "Active"

	if result := h.DB.Create(&user); result.Error != nil {
		logger.WithError(result.Error).Error("Error creating user")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating user"})
		return
	}
	logger.Info("User added to the Database")

	//create new jwt from the claims
	var updatedUser models.User
	h.DB.Where("email = ?", jwtClaims.Email).First(&updatedUser)
	newJWT, err := middleware.GenerateJWT(updatedUser)
	if err != nil {
		logger.WithError(err).Error("Error generating token")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating token"})
		return

	}
	c.JSON(http.StatusOK, models.Response{Data: newJWT})
}

func (h Handler) checkUserExists(email string) bool {
	var user models.User
	result := h.DB.Where("email = ?", email).First(&user)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false // User does not exist
	}
	return user.Email != "" // Check if email is not empty
}

func (h Auth) ResendOTP(c *gin.Context) {
	pin := utilities.GenerateVerifyEmailToken()
	jwtClaims := middleware.GetClaims(c)
	if jwtClaims == nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "IDIOT"})
		return
	}
	userkey := jwtClaims.Email + jwtClaims.User.Surname
	fmt.Println(userkey)
	userVerifyJSON, err := redisDB.RDB.Get(context.Background(), userkey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	var userVerify models.UserVerify
	err = json.Unmarshal([]byte(userVerifyJSON), &userVerify)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}
	userVerify.Pin = pin
	newVerifyJSON, err := json.Marshal(userVerify)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating token"})
	}

	err = redisDB.RDB.Set(context.Background(), userkey, newVerifyJSON, 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating token"})
		return
	}
	go sendOTP(userVerify.Email, pin)
	c.JSON(http.StatusOK, models.Response{Data: "Pin resent!"})
}

func sendOTP(userInfo string, pin string) {
	logger := utilities.NewLogger().LogWithCaller()
	// SMTP server configuration for Gmail
	smtpServer := "smtp.gmail.com"
	smtpPort := 587

	smtpUser, err := env.Get("COMPANY_EMAIL")
	if err != nil {
		// I'm sorry, this is the way it has to be
		return
	}

	smtpPassword, err := env.Get("COMPANY_AUTH") // Use app password if 2-factor authentication is enabled
	if err != nil {
		// I'm sorry, this is the way it has to be
		return
	}

	// Recipient email address
	to := userInfo
	// Email subject and body
	subject := "Verify Account"
	body := "Hello,\nPlease verify your DRE account using this pin: " + pin + "\n\nThanks,\nTeam Techtonic."

	// Initialize the SMTP dialer
	d := gomail.NewDialer(smtpServer, smtpPort, smtpUser, smtpPassword)
	d.TLSConfig = &tls.Config{ServerName: smtpServer, InsecureSkipVerify: true}
	// Create a new email message
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		logger.WithError(err).Error("Error sending OTP email")
	}
	logger.Info("OTP Email sent successfully")
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
	logger := utilities.NewLogger().LogWithCaller()
	//check the body of the request
	var body models.SendResetRequest
	if err := c.BindJSON(&body); err != nil {
		logger.WithError(err).Error("Invalid Request")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Request"})
		return
	}

	//check if the email exists in the database
	var user models.User
	h.DB.Where("email = ?", body.Email).First(&user)
	if !h.checkUserExists(user.Email) {
		logger.Error("User does not exist")
		c.JSON(http.StatusNotFound, models.Response{Error: "User does not exist"})
		return
	}

	//generate a jwt token with the user's email
	tempUser := models.User{
		Email: user.Email,
	}

	jwt, err := middleware.GenerateJWT(tempUser)
	if err != nil {
		logger.WithError(err).Error("Error generating token")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error generating token"})
		return
	}

	frontendBase, err := env.Get("FRONTEND_BASE_URL")
	if err != nil {
		utilities.InternalError(c)
		return
	}

	companyEmail, err := env.Get("COMPANY_EMAIL")
	if err != nil {
		utilities.InternalError(c)
		return
	}

	//send an email to the user with a temporary link to reset the password
	linkURL := fmt.Sprintf("%s/reset-password/%s", frontendBase, jwt)
	email := models.Email{
		From:    companyEmail,
		To:      user.Email,
		Subject: "Password Reset",
		Body:    "Click here to reset your password: " + linkURL,
	}
	log.Println(email)
	if err := sendMail(email); err != nil {
		logger.WithError(err).Error("Error sending reset email")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error sending reset email"})
		return
	}
	logger.Info("Reset Email sent successfully")
	c.JSON(http.StatusOK, models.Response{Data: "Reset Email sent successfully"})
}

// activate reset password
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
	logger := utilities.NewLogger().LogWithCaller()
	//get the token from the url
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		logger.Error("Missing token")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Missing token"})
		return
	}
	claims := middleware.GetClaims(c)
	if claims == nil {
		logger.Error("Missing token")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Missing token"})
		return
	}

	//check the body of the request
	var body models.ResetPassword
	if err := c.BindJSON(&body); err != nil {
		logger.WithError(err).Error("Invalid Request")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid Request"})
		return
	}

	//get user form the database
	var user models.User
	h.DB.Where("email = ?", claims.Email).First(&user)
	if !h.checkUserExists(user.Email) {
		logger.Error("User does not exist")
		c.JSON(http.StatusNotFound, models.Response{Error: "User does not exist"})
		return
	}

	//hash the new password
	// WARNING: this function changes the salt of the user.
	hashedPassword, err := hashPassword(body.NewPassword, &user)
	if err != nil {
		logger.WithError(err).Error("Error hashing password")
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
	companyEmail, err := env.Get("COMPANY_EMAIL")
	if err != nil {
		return err
	}
	companyAuth, err := env.Get("COMPANY_AUTH")
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
	m.SetBody("text/html", email.Body)
	logger.WithField("email", email).Info("Sending email")

	if err := d.DialAndSend(m); err != nil {
		logger.WithError(err).Error("Error sending email")
		return err
	}
	return nil
}

func hashPassword(newPassword string, user *models.User) (string, error) {
	logger := utilities.NewLogger().LogWithCaller()
	hasher := utilities.NewArgon2idHash(1, 12288, 4, 32, 16)

	hashAndSalt := hasher.HashPassword(newPassword)
	user.Salt = base64.StdEncoding.EncodeToString(hashAndSalt.Salt)
	logger.Info("Password hashed successfully")
	return base64.StdEncoding.EncodeToString(hashAndSalt.Hash), nil
}
