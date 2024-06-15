package handlers

import (
	"api/models"
	"api/utilities"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupUserRoutes(router *mux.Router, h Handler) {
	router.HandleFunc("/update", h.updateUser).Methods(http.MethodPut)
}

func (h Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	//get the user id from the request
	var updateUser struct {
		FirstName string `json:"first_name"`
		Surname   string `json:"surname"`
		Email     string `json:"email"`

		Code        *string `json:"code"` //This is the country code
		Country     *string `json:"country"`
		Province    *string `json:"province"`
		City        *string `json:"city"`
		Street3     *string `json:"street3"`
		Street2     *string `json:"street2"`
		Street      *string `json:"street"`
		AddressType *int    `json:"address_type"`
	}

	//Unmarshal the body into the local variable
	err = json.Unmarshal(body, &updateUser)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	//retrieve the user from the database
	var dbUser models.User
	h.DB.Where("email = ?", updateUser.Email).First(&dbUser)

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
		dbAddress.Country = updateUser.Country
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
	h.DB.Save(&dbUser)
	h.DB.Save(&dbAddress)

	utilities.WriteJSON(w, http.StatusOK, models.Response{Data: "User updated successfully"})
}
