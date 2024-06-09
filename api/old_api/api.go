package api

import (
	"api/old_model"
	"api/storage"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"reflect"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/argon2"
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
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Error: "invalid request method"})
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
		return writeJSON(w, http.StatusInternalServerError, model.Response{Status: 500, Error: err.Error()})
	}

	var wrappedResJSON any
	err = json.Unmarshal([]byte(wrappedResStr), &wrappedResJSON)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, model.Response{Status: 500, Error: err.Error()})
	}

	return writeJSON(w, http.StatusOK, wrappedResJSON)
}

func (s *APIServer) postAPI(w http.ResponseWriter, r *http.Request) error {
	req := new(model.BaseRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Error: err.Error()})
	}

	switch req.RequestType {
	case "create_account":
		return s.createAccount(w, req.Body)
	case "login":
		return s.login(w, req.Body)
	case "dispute_summary":
		return s.getDisputeSummary(w, req.Body)
	default:
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Error: "invalid request type"})
	}
}

func (s *APIServer) getDisputeSummary(w http.ResponseWriter, rawBody json.RawMessage) error {
	var body model.DisputeSummaryBody
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Error: err.Error()})
	}

	if body.UserID == "" {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Error: "missing required fields"})
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

	return writeJSON(w, http.StatusOK, model.Response{Status: 200, Data: disputeSummaryLists})
}

func (s *APIServer) createAccount(w http.ResponseWriter, rawBody json.RawMessage) error {
	var body model.CreateAccountBody
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Error: err.Error()})
	}

	if body.FirstName == "" || body.Surname == "" || body.Password == "" {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Error: "missing required fields"})
	}
	hasher := Argon2idHash{time: 1, memory: 12288, threads: 4, keylen: 32, saltlen: 16}
	hashAndSalt := hasher.hashPassword(body.Password)
	user := &model.User{
		First_name:         body.FirstName,
		Surname:            body.Surname,
		Birthdate:          "1990-01-01",
		Nationality:        "USA",
		Role:               "user",
		Email:              body.Email,
		Password_hash:      base64.StdEncoding.EncodeToString(hashAndSalt.Hash),
		Salt:               base64.StdEncoding.EncodeToString(hashAndSalt.Salt),
		Phone_number:       "",
		Address_id:         1,
		Created_at:         "xedolek889@kernuo.com",
		Updated_at:         "2024-05-30 10:00:00",
		Last_login:         "2024-05-30 10:00:00",
		Status:             "active",
		Gender:             "male",
		Preferred_language: "en",
		Timezone:           "UTC",
	}

	if err := s.store.CreateUser(user); err != nil {
		return writeJSON(w, http.StatusInternalServerError, model.Response{Status: 500, Error: err.Error()})
	}

	var bodyResponse = map[string]interface{}{
		"message": "account created",
	}

	sendOTP(body.Email)

	return writeJSON(w, http.StatusOK, model.Response{Status: 200, Data: bodyResponse})
}

func (s *APIServer) login(w http.ResponseWriter, rawBody json.RawMessage) error {
	var body model.LoginBody
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Error: err.Error()})
	}

	if body.Email == "" || body.Password == "" {
		return writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Error: "missing required fields"})
	}

	authUser := model.AuthUser()
	authUser.Email = body.Email

	salt, err := s.store.GetSalt(authUser)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, model.Response{Status: 500, Error: err.Error()})
	}

	realSalt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, model.Response{Status: 500, Error: err.Error()})
	}

	hasher := Argon2idHash{time: 1, memory: 12288, threads: 4, keylen: 32, saltlen: 16}
	checkHash, err := hasher.GenerateHash([]byte(body.Password), realSalt)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, model.Response{Status: 500, Error: err.Error()})
	}

	authUser.Password_hash = base64.StdEncoding.EncodeToString(checkHash.Hash)

	if err := s.store.AuthenticateUser(authUser); err != nil {
		return writeJSON(w, http.StatusInternalServerError, model.Response{Status: 500, Error: err.Error()})
	}

	bodyResponse := map[string]interface{}{
		"message": "login successful",
	}

	return writeJSON(w, http.StatusOK, model.Response{Status: 200, Data: bodyResponse})

}

func (a *Argon2idHash) hashPassword(password string) *HashSalt {
	salt, err := RandomSalt(16)
	if err != nil {
		return nil
	}
	hashSalt, err := a.GenerateHash([]byte(password), salt)
	if err != nil {
		return nil
	}
	return hashSalt
}

type Argon2idHash struct {
	time    uint32
	memory  uint32
	threads uint8
	keylen  uint32
	saltlen uint32
}

type HashSalt struct {
	Hash []byte
	Salt []byte
}

func NewArgon2idHash(time, saltLen uint32, memory uint32, threads uint8, keylen uint32) *Argon2idHash {
	return &Argon2idHash{time, memory, threads, keylen, saltLen}
}

func RandomSalt(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (a *Argon2idHash) GenerateHash(password, salt []byte) (*HashSalt, error) {
	var err error
	if len(salt) == 0 {
		salt, err = RandomSalt(a.saltlen)
	}
	if err != nil {
		return nil, err
	}
	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keylen)
	return &HashSalt{Hash: hash, Salt: salt}, nil

}

func (a *Argon2idHash) Compare(hash, salt, password []byte) error {
	// Generate hash for comparison.
	hashSalt, err := a.GenerateHash(password, salt)
	if err != nil {
		return err
	}
	// Compare the generated hash with the stored hash.
	// If they don't match return error.
	if !bytes.Equal(hash, hashSalt.Hash) {
		return errors.New("hash doesn't match")
	}
	return nil
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
			writeJSON(w, http.StatusBadRequest, model.Response{Status: 400, Error: err.Error()})
		}
	}
}

func sendOTP(userInfo string) {
	from := os.Getenv("COMPANY_EMAIL")
	to := []string{userInfo}
	password := os.Getenv("COMPANY_AUTH")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Subject: Verification Email \n\nPlease follow this link to verify your email: \n\n" + "http://localhost:8080/verify" + "\n" + "OTP: 123456" + "\n\n" + "Thank you!" + "\n\n" + "Techtonic Team")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")

}
