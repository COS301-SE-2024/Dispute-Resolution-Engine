package api

import (
	"api/model"
	"api/storage"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewServer(liseningADD string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: liseningADD,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/api", makeHTTPHandler(s.HandleRequests))
	router.HandleFunc("/api/{id}", makeHTTPHandler(s.HandlerTestGet))
	log.Println("Server is running on port ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)

}

func (s *APIServer) HandlerTestGet(w http.ResponseWriter, r *http.Request) error {
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
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Body: model.APIError{Error: "invalid request type"}})
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
		return writeJSON(w, http.StatusInternalServerError, model.Response{Status: 500, Body: model.APIError{Error: err.Error()}})
	}

	var wrappedResJSON any
	err = json.Unmarshal([]byte(wrappedResStr), &wrappedResJSON)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, model.Response{Status: 500, Body: model.APIError{Error: err.Error()}})
	}

	return writeJSON(w, http.StatusOK, wrappedResJSON)
}

func (s *APIServer) postAPI(w http.ResponseWriter, r *http.Request) error {
	req := new(model.BaseRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Body: model.APIError{Error: err.Error()}})
	}

	switch req.RequestType {
	case "create_account":
		return s.createAccount(w, req.Body)
	case "login":
		return s.login(w, req.Body)
	case "dispute_summary":
		return s.getDisputeSummary(w, req.Body)
	default:
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Body: model.APIError{Error: "invalid request type"}})
	}
}

func (s *APIServer) getDisputeSummary(w http.ResponseWriter, rawBody json.RawMessage) error {
	var body model.DisputeSummaryBody
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Body: model.APIError{Error: err.Error()}})
	}

	if body.UserID == "" {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Body: model.APIError{Error: "missing required fields"}})
	}

	//stubbed db access
	disputeSummaryLists := make([]model.DisputeSummary, 0)
	disputeSummaryLists = append(disputeSummaryLists, model.DisputeSummary{DisputeID: "1", DisputeTitle: "Dispute 1"})
	disputeSummaryLists = append(disputeSummaryLists, model.DisputeSummary{DisputeID: "2", DisputeTitle: "Dispute 2"})
	disputeSummaryLists = append(disputeSummaryLists, model.DisputeSummary{DisputeID: "3", DisputeTitle: "Dispute 3"})
	disputeSummaryLists = append(disputeSummaryLists, model.DisputeSummary{DisputeID: "4", DisputeTitle: "Dispute 4"})
	disputeSummaryLists = append(disputeSummaryLists, model.DisputeSummary{DisputeID: "5", DisputeTitle: "Dispute 5"})
	disputeSummaryLists = append(disputeSummaryLists, model.DisputeSummary{DisputeID: "6", DisputeTitle: "Dispute 6"})
	disputeSummaryLists = append(disputeSummaryLists, model.DisputeSummary{DisputeID: "7", DisputeTitle: "Dispute 7"})
	disputeSummaryLists = append(disputeSummaryLists, model.DisputeSummary{DisputeID: "8", DisputeTitle: "Dispute 8"})
	disputeSummaryLists = append(disputeSummaryLists, model.DisputeSummary{DisputeID: "9", DisputeTitle: "Dispute 9"})
	disputeSummaryLists = append(disputeSummaryLists, model.DisputeSummary{DisputeID: "10", DisputeTitle: "Dispute 10"})

	return writeJSON(w, http.StatusOK, model.Response{Status: 200, Body: disputeSummaryLists})
}

func (s *APIServer) createAccount(w http.ResponseWriter, rawBody json.RawMessage) error {
	var body model.CreateAccountBody
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Body: model.APIError{Error: err.Error()}})
	}

	if body.FirstName == "" || body.Surname == "" || body.PasswordHash == "" {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Body: model.APIError{Error: "missing required fields"}})
	}
	user := &model.User{
		First_name:         body.FirstName,
		Surname:            body.Surname,
		Birthdate:          "1990-01-01",
		Nationality:        "USA",
		Role:               "user",
		Email:              body.Email,
		Password_hash:      body.PasswordHash,
		Phone_number:       "",
		Address_id:         1,
		Created_at:         "2024-05-30 10:00:00",
		Updated_at:         "2024-05-30 10:00:00",
		Last_login:         "2024-05-30 10:00:00",
		Status:             "active",
		Gender:             "male",
		Preferred_language: "en",
		Timezone:           "UTC",
	}

	if err := s.store.CreateUser(user); err != nil {
		return writeJSON(w, http.StatusInternalServerError, model.Response{Status: 500, Body: model.APIError{Error: err.Error()}})
	}

	var bodyResponse = map[string]interface{}{
		"message": "account created",
	}

	return writeJSON(w, http.StatusOK, model.Response{Status: 200, Body: bodyResponse})
}

func (s *APIServer) login(w http.ResponseWriter, rawBody json.RawMessage) error {
	var body model.LoginBody
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Body: model.APIError{Error: err.Error()}})
	}

	if body.Email == "" || body.PasswordHash == "" {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Body: model.APIError{Error: "missing required fields"}})
	}

	authUser := model.AuthUser()
	authUser.Email = body.Email
	authUser.Password_hash = body.PasswordHash

	if err := s.store.AuthenticateUser(authUser); err != nil {
		return writeJSON(w, http.StatusInternalServerError, model.Response{Status: 500, Body: model.APIError{Error: err.Error()}})
	}

	bodyResponse := map[string]interface{}{
		"message": "login successful",
	}

	return writeJSON(w, http.StatusOK, bodyResponse)

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

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Body: model.APIError{Error: err.Error()}})
		}
	}
}
