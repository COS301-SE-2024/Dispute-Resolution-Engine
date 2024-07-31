package handlers_test

import (
	"api/env"
	"api/handlers"
	"api/models"
	"api/utilities"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthTestSuite struct {
	suite.Suite

	mock   sqlmock.Sqlmock
	db     *gorm.DB
	router *gin.Engine

	userRows *sqlmock.Rows
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func createUserRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id",
		"first_name",
		"surname",
		"birthdate",
		"nationality",
		"role",
		"email",
		"password_hash",
		"phone_number",
		"address_id",
		"created_at",
		"updated_at",
		"last_login",
		"status",
		"gender",
		"preferred_language",
		"timezone",
		"salt",
	})
}

// Runs before every test to set up the DB and routers
func (suite *AuthTestSuite) SetupTest() {
	mock, db, _ := mockDatabase()

	handler := handlers.NewAuthHandler(db)
	gin.SetMode("release")
	router := gin.Default()
	router.POST("/login", handler.LoginUser)
	router.POST("/signup", handler.CreateUser)

	suite.mock = mock
	suite.db = db
	suite.router = router

	suite.userRows = createUserRows()
}

func initUserRows() *sqlmock.Rows {
	salt := []byte("salt")
	hash := utilities.HashPasswordWithSalt("pass", salt)

	return createUserRows().AddRow(
		0,                                       // id
		"Test",                                  // first_name
		"User",                                  // surname
		time.Now(),                              // birthdate
		"ZA",                                    // nationality
		"user",                                  // role
		"test@example.com",                      // email
		base64.StdEncoding.EncodeToString(hash), // password_hash
		"0123456789",                            // phone_number
		0,                                       // address_id
		time.Now(),                              // created_at
		time.Now(),                              // updated_at
		time.Now(),                              // last_login
		"Active",                                // status
		"Male",                                  // gender
		"English",                               // preferred_language
		"timezone",                              // timezone
		base64.StdEncoding.EncodeToString(salt), // salt

	)
}

// creates a request with the the passed-in payload as the body
func createJSONRequest(method string, url string, req any) (*http.Request, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, url, bytes.NewReader(body))
}

func (suite *AuthTestSuite) TestSignupWithValidInformation() {
	env.RegisterDefault("JWT_SECRET", "your_secret_key")
	// Create the expected result for the INSERT query
	result := sqlmock.NewResult(1, 1) // 1 is the last insert ID, 1 is the number of rows affected

	// Set up the expectation for the INSERT query
	suite.mock.ExpectExec("INSERT INTO \"users\" \\(first_name, surname, email, phone_number, password, birthdate, gender, nationality, timezone, preferred_language\\) VALUES \\(.+\\)").
		WithArgs("John", "Doe", "john.doe@example.com", "1234567890", "hashedpassword", "1990-01-01", "male", "USA", "UTC", "en-US").
		WillReturnResult(result)

	// Create the request
	req, _ := createJSONRequest("POST", "/signup", gin.H{
		"first_name":         "John",
		"surname":            "Doe",
		"email":              "john.doe@example.com",
		"phone_number":       "1234567890",
		"password":           "password",
		"birthdate":          "1990-01-01",
		"gender":             "male",
		"nationality":        "USA",
		"timezone":           "UTC",
		"preferred_language": "en-US",
	})

	// Perform the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Check the response
	var response models.Response
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.NoError(suite.T(), json.Unmarshal(w.Body.Bytes(), &response))
	assert.Empty(suite.T(), response.Error)
}

func (suite *AuthTestSuite) TestSignupWithExistingEmail() {
	result := sqlmock.NewResult(1, 1)

	suite.mock.ExpectExec("INSERT INTO \"users\" \\(first_name, surname, email, phone_number, password, birthdate, gender, nationality, timezone, preferred_language\\) VALUES \\(.+\\)").
		WithArgs("John", "Doe", "john.doe@example.com", "1234567890", "hashedpassword", "1990-01-01", "male", "USA", "UTC", "en-US").
		WillReturnResult(result)

	// Create the request
	req, _ := createJSONRequest("POST", "/signup", gin.H{
		"first_name":         "John",
		"surname":            "Doe",
		"email":              "john.doe@example.com",
		"phone_number":       "1234567890",
		"password":           "password",
		"birthdate":          "1990-01-01",
		"gender":             "male",
		"nationality":        "USA",
		"timezone":           "UTC",
		"preferred_language": "en-US",
	})

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response models.Response
	assert.Equal(suite.T(), http.StatusConflict, w.Code)
	assert.NoError(suite.T(), json.Unmarshal(w.Body.Bytes(), &response))
	assert.NotEmpty(suite.T(), response.Error)
}

func (suite *AuthTestSuite) TestSignupWithBadRequestBody() {
	result := sqlmock.NewResult(1, 1)

	suite.mock.ExpectExec("INSERT INTO \"users\" \\(first_name, surname, email, phone_number, password, birthdate, gender, nationality, timezone, preferred_language\\) VALUES \\(.+\\)").
		WithArgs("John", "Doe", "john.doe@example.com", "1234567890", "hashedpassword", "1990-01-01", "male", "USA", "UTC", "en-US").
		WillReturnResult(result)

	// Create the request
	req, _ := createJSONRequest("POST", "/signup", gin.H{
		"first_name":         "John",
		"surname":            "Doe",
		"email":              "john.doe@example.com",
		"phone_number":       "1234567890",
		"birthdate":          "1990-01-01",
		"gender":             "male",
		"nationality":        "USA",
		"timezone":           "UTC",
		"preferred_language": "en-US",
	})

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response models.Response
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.NoError(suite.T(), json.Unmarshal(w.Body.Bytes(), &response))
	assert.NotEmpty(suite.T(), response.Error)
}

func (suite *AuthTestSuite) TestLoginEmailDoesNotExist() {
	rows := createUserRows()
	suite.mock.ExpectQuery("SELECT (.+) FROM \"users\" .*").WithArgs("test", 1).WillReturnRows(rows)

	req, _ := createJSONRequest("POST", "/login", gin.H{
		"email":    "test",
		"password": "pass",
	})

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response models.Response
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	assert.NoError(suite.T(), json.Unmarshal(w.Body.Bytes(), &response))
	assert.NotEmpty(suite.T(), response.Error)
}

func (suite *AuthTestSuite) TestLoginPasswordIncorrect() {
	suite.mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE email =.*").WithArgs("test", 1).WillReturnRows(initUserRows())
	suite.mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE email =.*").WithArgs("test", 1).WillReturnRows(initUserRows())

	req, _ := createJSONRequest("POST", "/login", gin.H{
		"email":    "test",
		"password": "incorrect",
	})

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response models.Response
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.NoError(suite.T(), json.Unmarshal(w.Body.Bytes(), &response))
	assert.NotEmpty(suite.T(), response.Error)
}

func (suite *AuthTestSuite) TestLoginWithCorrectCredentials() {
	rows := initUserRows()
	suite.mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE email =.*").WithArgs("test", 1).WillReturnRows(rows)
	suite.mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE email =.*").WithArgs("test", 1).WillReturnRows(rows)

	req, _ := createJSONRequest("POST", "/login", gin.H{
		"email":    "test",
		"password": "pass",
	})

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response struct {
		Data string `json:"data"`
	}
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.NoError(suite.T(), json.Unmarshal(w.Body.Bytes(), &response))
	assert.NotEmpty(suite.T(), response.Data)
}
