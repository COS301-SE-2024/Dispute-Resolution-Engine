package handlers

import (
	"api/models"
	"api/utilities"
	"encoding/json"
	"fmt"
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

	duplicate := h.checkUserExists(user.Email)

	if duplicate {
		errJson := StringWrapper{"User already exists"}
		jsonData, err := json.Marshal(errJson)
		if err != nil {
			fmt.Println("Error marshalling JSON", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonData)
		return
	}

	hashAndSalt := hasher.HashPassword(user.PasswordHash)

	user.CreatedAt = utilities.GetCurrentTime()
	user.UpdatedAt = utilities.GetCurrentTime()
	user.LastLogin = nil
	user.PasswordHash = string(hashAndSalt.Hash)
	user.Salt = string(hashAndSalt.Salt)

	errJson, err := json.Marshal(StringWrapper{"Failed to create user"})
	if result := h.DB.Create(&user); result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(errJson)
		return
	}
	resbody, err := json.Marshal(StringWrapper{"Account created successfully"})

	// Send a 201 created response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resbody)
}

func (h handler) Login (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var user models.User
	err = json.Unmarshal(body, &user)

	var dbUser models.User
	h.DB.Where("email = ?", user.Email).First(&dbUser)

	if dbUser.Email == "" {
		errJson, err := json.Marshal(StringWrapper{"User not found"})
		if err != nil {
			fmt.Println("Error marshalling JSON", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(errJson)
		return
	}

	if dbUser.PasswordHash != user.PasswordHash {
		errJson, err := json.Marshal(StringWrapper{"Invalid password"})
		if err != nil {
			fmt.Println("Error marshalling JSON", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errJson)
		return
	}

	resbody, err := json.Marshal(StringWrapper{"Login successful"})

	// Send a 200 OK response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resbody)
}

func (h handler) checkUserExists(email string) bool {
	var user models.User
	h.DB.Where("email = ?", email).First(&user)
	return user.Email != ""
}
