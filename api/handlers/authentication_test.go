package handlers_test

import (
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
