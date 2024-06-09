package handlers

import (
	"api/models"
	"api/utilities"
	"encoding/json"
	"io"
	"net/http"
)

type StringWrapper struct {
	Data string `json:"Data"`
}

func (h handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	hasher := utilities.NewArgon2idHash(1,12288,4,32,16)
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

	duplicate := h.checkUserExists(user.Email)

	if duplicate {
		// errJson := StringWrapper{"User already exists"}
		// jsonData, err := json.Marshal(errJson)
		// if err != nil {
		// 	fmt.Println("Error marshalling JSON", err)
		// 	return
		// }
		
		utilities.WriteJSON(w, http.StatusConflict, models.Response{Error: "User already exists"})
		return
	}
	if result := h.DB.Create(&user); result.Error != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Error creating user"})
		return
	}
	
	utilities.WriteJSON(w, http.StatusCreated, models.Response{Data: user})
}

func (h handler) checkUserExists(email string) bool {
	var user models.User
	h.DB.Where("email = ?", email).First(&user)
	return user.Email != ""
}
