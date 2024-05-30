package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
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

func NewServer(liseningADD string) *APIServer {
	return &APIServer{
		listenAddr: liseningADD,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/api", makeHTTPHandler(s.HandleRequests))
	log.Println("Server is running on port ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)

}


func (s *APIServer) HandleRequests(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
		case "GET":
			return s.getAPI(w, r)
		case "POST":
			return s.postAPI(w, r)
		default:
			return fmt.Errorf("method not allowed: %s", r.Method)
	}
}

func (s *APIServer) getAPI(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, "GET API")
}

func (s *APIServer) postAPI(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, "POST API")
}