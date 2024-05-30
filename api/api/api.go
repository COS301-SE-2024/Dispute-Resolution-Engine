package api

import (
	"api/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
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
	router.HandleFunc("/api/{id}", makeHTTPHandler(s.HalderTestGet))
	log.Println("Server is running on port ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)

}

func (s *APIServer) HalderTestGet(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	return writeJSON(w, http.StatusOK, "GET_API_with_ID: "+id)
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
	address := model.NewAddress()
	address.Id = 1
	address.Code = "12345"
	address.Country = "USA"
	address.Province = "NY"
	address.City = "NYC"
	address.Street3 = "Street3"
	address.Street2 = "Street2"
	address.Street = "Street"
	address.Address_type = 1
	address.Last_updated = "2021-01-01"

	account := model.NewUser()
	account.ID = 1
	account.First_name = "John"
	account.Surname = "Doe"
	account.Birthdate = "1990-01-01"
	account.Nationality = "USA"
	account.Role = "Admin"
	account.Email = "j@d.com"
	account.Password_hash = "123456"
	account.Phone_number = "1234567890"
	account.Address_id = 1
	account.Created_at = "2021-01-01"
	account.Updated_at = "2021-01-01"
	account.Last_login = "2021-01-01"
	account.Status = "active"

	wrappedResStr, err := s.wrapInJSON(*address, *account)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, APIError{Error: err.Error()})
	}

	var wrappedResJSON any
	err = json.Unmarshal([]byte(wrappedResStr), &wrappedResJSON)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, APIError{Error: err.Error()})
	}

	return writeJSON(w, http.StatusOK, wrappedResJSON)
}

func (s *APIServer) postAPI(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, "POST API")
}

func (s *APIServer) wrapInJSON(objects ...interface{}) (string, error) {
	jsonData := make(map[string]interface{})

	for _, obj := range objects {
		objType := reflect.TypeOf(obj).String()

		jsonData[objType] = obj
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
