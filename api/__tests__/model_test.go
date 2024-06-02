package api

import (
    "testing"
	"api/model"
    //"github.com/stretchr/testify/assert"
)

// func TestNewServer(t *testing.T) {
//     store, err := storage.NewPostgresStore()
// 	server := api.NewServer(":8080", store)
// 	server.Run()	
// }
func TestAuthUserCreation(t *testing.T) {
	user := model.AuthUser()
	if user == nil {
		t.Fatalf("AuthUser() = nil, want &LoginUser{}")
	}
}