package api

import (
    "testing"
	"api/api"
    //"github.com/stretchr/testify/assert"
)

type MockStorage struct{}

// func TestNewServer(t *testing.T) {
//     store, err := storage.NewPostgresStore()
// 	server := api.NewServer(":8080", store)
// 	server.Run()	
// }
func TestRandomSales(t *testing.T) {
	secret, err := api.RandomSalt(32)
	if err != nil {
		t.Errorf("RandomSalt() error = %v", err)
		return
	}
	if len(secret) != 32 {
		t.Errorf("RandomSalt() = %v, want %v", len(secret), 32)
	}	
}