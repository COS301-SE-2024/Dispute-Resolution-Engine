package handlers

import (
	"api/models"
	"api/utilities"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type StringWrapper struct {
	Data string `json:"Data"`
}

func (h handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	hasher := utilities.NewArgon2idHash(1, 12288, 4, 32, 16)
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var reqUser struct {
		//These are all the user details that are required to create a user
		FirstName         string  `json:"first_name"`
		Surname           string  `json:"surname"`
		Birthdate         string  `json:"birthdate"`
		Nationality       string  `json:"nationality"`
		Email             string  `json:"email"`
		Password          string  `json:"password"`
		PhoneNumber       *string `json:"phone_number"`
		Gender            string  `json:"gender"`
		PreferredLanguage *string `json:"preferred_language"`
		Timezone          *string `json:"timezone"`

		//These are the user's address details
		AddressID *int64  //something I need to derive from the address info
		Code      *string `json:"code"` //This is the country code
		Country   *string `json:"country"`
		Province  *string `json:"province"`
		City      *string `json:"city"`
		Street3   *string `json:"street3"`
		Street2   *string `json:"street2"`
		Street    *string `json:"street"`
	}
	//create an appropriate user object
	//var user models.User
	//parse the body into the user object
	//err = json.Unmarshal(body, &user)
	err = json.Unmarshal(body, &reqUser)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}
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
		AddressID:         reqUser.AddressID,
		Status:            "Active",
		Gender:            reqUser.Gender,
		PreferredLanguage: reqUser.PreferredLanguage,
		Timezone:          reqUser.Timezone,
	}

	//Check if there is an existing email
	duplicate := h.checkUserExists(user.Email)

	if duplicate {
		utilities.WriteJSON(w, http.StatusConflict, models.Response{Error: "Email already in use"})
		return
	}

	//Hash the password
	hashAndSalt := hasher.HashPassword(user.PasswordHash)
	user.PasswordHash = base64.StdEncoding.EncodeToString(hashAndSalt.Hash)
	user.Salt = base64.StdEncoding.EncodeToString(hashAndSalt.Salt)

	//update log metrics
	user.CreatedAt = utilities.GetCurrentTime()
	user.UpdatedAt = utilities.GetCurrentTime()
	user.Status = "Active"

	//Small user preferences
	user.Role = "user"
	user.PreferredLanguage = nil
	user.LastLogin = nil

	if result := h.DB.Create(&user); result.Error != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Error creating user"})
		return
	}

	utilities.WriteJSON(w, http.StatusCreated, models.Response{Data: "User created successfully"})
}

func (h handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	hasher := utilities.NewArgon2idHash(1, 12288, 4, 32, 16)
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: "Invalid Request"})
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
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Error decoding salt"})
		return
	}
	checkHash, err := hasher.GenerateHash([]byte(user.PasswordHash), realSalt)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Something went wrong..."})
		return
	}

	if dbUser.PasswordHash != base64.StdEncoding.EncodeToString(checkHash.Hash) {
		print(dbUser.PasswordHash)
		print(base64.StdEncoding.EncodeToString(checkHash.Hash))
		utilities.WriteJSON(w, http.StatusUnauthorized, models.Response{Error: "Invalid password"})
		return
	}

	dbUser.LastLogin = utilities.GetCurrentTimePtr()
	h.DB.Where("email = ?", user.Email).Update("last_login", utilities.GetCurrentTime())
	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "Login successful"})

}

func (h handler) checkUserExists(email string) bool {
	var user models.User
	h.DB.Where("email = ?", email).First(&user)
	return user.Email != ""
}
