package handlers

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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

func SetupAuthRoutes(router *mux.Router, h Handler) {
	router.HandleFunc("/signup", h.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/login", h.LoginUser).Methods(http.MethodPost)
	router.Handle("/reset-password", middleware.RoleMiddleware(http.HandlerFunc(h.ResetPassword), 0)).Methods(http.MethodPost)
	// router.Handle("/verify", middleware.RoleMiddleware(http.HandlerFunc(h.Verify), 0)).Methods(http.MethodPost)
	router.HandleFunc("/verify", h.Verify).Methods(http.MethodPost)
}

// @Summary Reset a user's password
// @Description Reset a user's password
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "Password reset not available yet..."
// @Router /auth/reset-password [post]
func (h Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "password reset not available yet..."})
}

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
func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	hasher := utilities.NewArgon2idHash(1, 12288, 4, 32, 16)
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	//create a local variable to store the user details
	var reqUser models.CreateUser
	//Unmarshal the body into the local variable
	err = json.Unmarshal(body, &reqUser)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	//stub timezone
	zone, _ := time.Now().Zone()
	timezone := zone
	reqUser.Timezone = &timezone
	//Now put stuff in the actual user object
	date, err := time.Parse("2006-01-02", reqUser.Birthdate)
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

	// address := models.Address{
	// 	Code:        nil, //to be filled in a later request
	// 	Country:     reqUser.Country,
	// 	Province:    reqUser.Province,
	// 	City:        reqUser.City,
	// 	Street3:     reqUser.Street3,
	// 	Street2:     reqUser.Street2,
	// 	Street:      reqUser.Street,
	// 	AddressType: reqUser.AddressType,
	// }

	//Check if there is an existing email
	duplicate := h.checkUserExists(user.Email)

	if duplicate {
		utilities.WriteJSON(w, http.StatusConflict, models.Response{Error: "Email already in use"})
		return
	}
	//get country code
	// var country models.Country
	// h.DB.Where("country_name = ?", reqUser.Country).First(&country)
	// address.Code = &country.CountryCode

	//create the address
	// if result := h.DB.Create(&address); result.Error != nil {
	// 	utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Something went wrong creating the address..."})
	// 	return
	// }
	//get the address id
	// user.AddressID = &address.ID

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
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Error creating user"})
		return
	}
	sendOTP(user.Email)
	utilities.WriteJSON(w, http.StatusCreated, models.Response{Data: "User created successfully"})
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
func (h Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	hasher := utilities.NewArgon2idHash(1, 12288, 4, 32, 16)
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request, please check request body.", http.StatusBadRequest)
		return
	}

	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	if !h.checkUserExists(user.Email) {
		utilities.WriteJSON(w, http.StatusNotFound, models.Response{Error: "User does not exist"})
		return
	}

	var dbUser models.User
	h.DB.Where("email = ?", user.Email).First(&dbUser)

	realSalt, err := base64.StdEncoding.DecodeString(dbUser.Salt)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Something went wrong processing your request..."})
		return
	}
	checkHash, err := hasher.GenerateHash([]byte(user.PasswordHash), realSalt)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Something went wrong processing your request..."})
		return
	}

	if dbUser.PasswordHash != base64.StdEncoding.EncodeToString(checkHash.Hash) {
		print(dbUser.PasswordHash)
		print(base64.StdEncoding.EncodeToString(checkHash.Hash))
		utilities.WriteJSON(w, http.StatusUnauthorized, models.Response{Error: "Invalid credentials"})
		return
	}

	dbUser.LastLogin = utilities.GetCurrentTimePtr()
	h.DB.Where("email = ?", user.Email).Update("last_login", utilities.GetCurrentTime())

	token, err := middleware.GenerateJWT(dbUser)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Error generating token"})
		return
	}

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: token})

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
func (h Handler) Verify(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	var pinReq models.VerifyUser
	err = json.Unmarshal(body, &pinReq)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}
	valid, err := utilities.RemoveFromFile("stubbedStorage/verify.txt", pinReq.Pin)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Error verifying pin"})
		return
	}
	if !valid {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: "Invalid pin"})
		return
	}
	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "Email verified successfully"})
}
