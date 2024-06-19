package handlers

import (
	"api/models"
	"api/utilities"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupUserRoutes(router *mux.Router, h Handler) {
	router.HandleFunc("/profile", h.updateUser).Methods(http.MethodPut)
	router.HandleFunc("/profile", h.getUser).Methods(http.MethodGet)
	router.HandleFunc("/remove", h.RemoveAccount).Methods(http.MethodDelete)
}

// @Summary Get user profile
// @Description Get user profile
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "User profile not available yet..."
// @Router /user/profile [get]
func (h Handler) getUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "User profile not available yet..."})
}

// @Summary Update user profile
// @Description Update user profile
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.UpdateUser true "User"
// @Success 200 {object} models.Response "User updated successfully"
// @Failure 400 {object} models.Response "Bad Request"
// @Router /user/profile [put]
func (h Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request, please check request body.", http.StatusBadRequest)
		return
	}

	//get the user id from the request
	var updateUser models.UpdateUser

	//Unmarshal the body into the local variable
	err = json.Unmarshal(body, &updateUser)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	//retrieve the user from the database
	var dbUser models.User
	h.DB.Where("id = ?", updateUser.Email).First(&dbUser)

	var dbAddress models.Address
	h.DB.Where("id = ?", dbUser.AddressID).First(&dbAddress)

	//check which fields to update
	if updateUser.FirstName != "" {
		dbUser.FirstName = updateUser.FirstName
	}
	if updateUser.Surname != "" {
		dbUser.Surname = updateUser.Surname
	}
	if updateUser.Code != nil {
		dbAddress.Code = updateUser.Code
	}
	if updateUser.Country != nil {
		var dbCountry models.Country
		dbAddress.Country = updateUser.Country
		h.DB.Where("country_name = ?", updateUser.Country).First(&dbCountry)
		dbAddress.Code = &dbCountry.CountryCode
	}
	if updateUser.Province != nil {
		dbAddress.Province = updateUser.Province
	}
	if updateUser.City != nil {
		dbAddress.City = updateUser.City
	}
	if updateUser.Street3 != nil {
		dbAddress.Street3 = updateUser.Street3
	}
	if updateUser.Street2 != nil {
		dbAddress.Street2 = updateUser.Street2
	}
	if updateUser.Street != nil {
		dbAddress.Street = updateUser.Street
	}
	if updateUser.AddressType != nil {
		dbAddress.AddressType = updateUser.AddressType
	}

	//now update the user and address
	h.DB.Model(&dbUser).Where("id = ?", dbUser.ID).Updates(dbUser)
	h.DB.Model(&dbAddress).Where("id = ?", dbUser.AddressID).Updates(dbAddress)

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "User updated successfully"})
}

// @Summary Remove user account
// @Description Remove user account
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.DeleteUser true "User"
// @Success 200 {object} models.Response "User account removed successfully"
// @Failure 400 {object} models.Response "Bad Request"
// @Router /user/remove [delete]
func (h Handler) RemoveAccount(w http.ResponseWriter, r *http.Request) {
	hasher := utilities.NewArgon2idHash(1, 12288, 4, 32, 16)
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var user models.DeleteUser
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
	checkHash, err := hasher.GenerateHash([]byte(user.Password), realSalt)
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

	h.DB.Where("email = ?", user.Email).Delete(&dbUser)
	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "User account removed successfully"})
}

func (h Handler) UpdateUserAddress(w http.ResponseWriter, r *http.Request) {
	//read request body into variable
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request, please check request body.", http.StatusBadRequest)
		return
	}
	//here we get the details of the request
	var UpdateUserAddress models.UpdateAddress
	err = json.Unmarshal(body, &UpdateUserAddress)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}
	//retrieve the record from the database
	var dbUser models.User
	h.DB.Where("email = ?", UpdateUserAddress.Email).First(&dbUser)

	//now we have to set the address parameters using the passed in data

}
