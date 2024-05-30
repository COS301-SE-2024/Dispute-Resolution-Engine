package main

import (
	// "fmt"
	// "log"
	// "net/http"

	"api/database"
	"encoding/json"
	"log"
	"net/http"

	// "api/handlers"

	"github.com/gorilla/mux"
	// "github.com/rs/cors"
)

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}


type apiFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Error string
}


func makeHTTPHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
}

func newServer(liseningADD string) *APIServer {
	return &APIServer{
		listenAddr: liseningADD,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/api", makeHTTPHandler(s.handleRequests))
	log.Println("Server is running on port ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)

}


func (s *APIServer) handleRequests(w http.ResponseWriter, r *http.Request) error {
	return nil
}



func main() {
    db.ConnectDB()
}
