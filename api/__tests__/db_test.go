package api

import (
	"api/db"
	"api/handlers"
	"net/http"
	"github.com/gorilla/mux"
	
	"testing"
)

type MockStorage struct{}

func TestInit(t *testing.T) {
	DB := db.Init()
	if DB == nil {
        t.Errorf("Expected db to be initialized, got nil")
    }
	h := handlers.New(DB)
	router := mux.NewRouter()

	router.HandleFunc("/createAcc", h.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/login", h.LoginUser).Methods(http.MethodPost)
	
}