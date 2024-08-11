package handlers_test

import (
	"api/env"
	"api/handlers"
	"bytes"

	// "api/middleware"
	"api/models"
	// "api/utilities"
	// "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	// "time"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	mockUserCount = 10
)

type UserTestSuite struct {
	suite.Suite

	mock   sqlmock.Sqlmock
	db     *gorm.DB
	router *gin.Engine
	envLoader *mockenvLoader
}

type mockenvLoader struct {
}

func (m * mockenvLoader) LoadFromFile(files ...string) {

}
func (m * mockenvLoader) Register(key string) {

}
func (m * mockenvLoader) RegisterDefault(key, fallback string) {

}
func (m * mockenvLoader) Get(key string) (string, error) {
	return "secret", nil
}


func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) SetupTest() {
	conn, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       conn,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})

	handler := handlers.NewUserHandler(db)
	gin.SetMode("release")
	router := gin.Default()
	router.PUT("/user/profile", handler.UpdateUser)
	router.GET("/user/profile", handler.GetUser)
	router.PUT("/user/profile/address", handler.UpdateUserAddress)
	router.DELETE("/user/remove", handler.RemoveAccount)
	router.POST("/user/analytics", handler.UserAnalyticsEndpoint)

	suite.mock = mock
	suite.db = db
	suite.router = router
	suite.envLoader = &mockenvLoader{}

	handler.EnvReader = suite.envLoader

}

func initUserTestRows() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
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

	for i := 1; i <= mockUserCount; i++ {
		rows.AddRow(
			i,
			fmt.Sprintf("FirstName%d", i),
			fmt.Sprintf("Surname%d", i),
			fmt.Sprintf("1990-01-%02d", i), // Mock birthdate
			"Nationality",
			"User",
			fmt.Sprintf("user%d@example.com", i),
			"mocked_hash",
			fmt.Sprintf("123-456-789%d", i),
			i,                   // address_id is null
			"2023-01-01 00:00:00", // created_at
			"2023-01-01 00:00:00", // updated_at
			"2023-01-01 00:00:00", // last_login
			"Active",
			"Male",
			"English",
			"UTC",
			"mocked_salt",
		)
	}
	return rows
}

func initAddressRows() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id",
		"code",
		"country",
		"province",
		"city",
		"street3",
		"street2",
		"street",
		"address_type",
		"last_updated",
	})
	// Add mock data here
	for i := 1; i <= mockUserCount; i++ {
		rows.AddRow(
			i,
			"ZA",
			fmt.Sprintf("Country%d", i),
			fmt.Sprintf("Province%d", i),
			fmt.Sprintf("City%d", i),
			fmt.Sprintf("Street3%d", i),
			fmt.Sprintf("Street2%d", i),
			fmt.Sprintf("Street%d", i),
			"Postal",
			"2023-01-01 00:00:00",
		)
	}
	return rows
}

func (suite *UserTestSuite) TestGetUser() {
	envLoader := env.NewEnvLoader()
	envLoader.RegisterDefault("JWT_SECRET", "secret")
	rows := initUserTestRows()
	suite.mock.ExpectQuery("SELECT (.+) FROM \"users\"").WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/user/profile", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	var result models.Response
	assert.NoError(suite.T(),json.Unmarshal(w.Body.Bytes(), &result))
	assert.NotEmpty(suite.T(), result.Error)
}

func (suite *UserTestSuite) TestUpdateUser() {
	envLoader := env.NewEnvLoader()
	envLoader.RegisterDefault("JWT_SECRET", "secret")
	updatePayload := map[string]interface{}{
		"first_name":         "NewFirstName",
		"surname":            "NewSurname",
		"phone_number":       "123-456-7891",
		"gender":             "Male",
		"nationality":        "NewNationality",
		"timezone":           "UTC",
		"preferred_language": "English",
	}
	payloadBytes, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest("PUT", "/user/profile", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	var result models.Response
	assert.Empty(suite.T(), result.Error)
}

func (suite *UserTestSuite) TestUpdateUserAddress() {
	envLoader := env.NewEnvLoader()
	envLoader.RegisterDefault("JWT_SECRET", "secret")
	updatePayload := map[string]interface{}{
		"country":  "NewCountry",
		"province": "NewProvince",
		"city":     "NewCity",
		"street3":  "NewStreet3",
		"street2":  "NewStreet2",
		"street":   "NewStreet",
	}
	payloadBytes, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest("PUT", "/user/profile/address", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ik5ld0FkZHJlc3MiLCJpYXQiOjE1MTYyMzkwMjJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ik5ld0FkZHJlc3MiLCJpYXQiOjE1MTYyMzkwMjJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ik5ld0FkZHJlc3MiLCJpYXQiOjE1MTYyMzkwMjJ9")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *UserTestSuite) TestRemoveAccount() {
	envLoader := env.NewEnvLoader()
	envLoader.RegisterDefault("JWT_SECRET", "secret")

	// Case 1: User does not exist
	removePayload := map[string]interface{}{
		"email":    "nonexistentuser@example.com",
		"password": "wrongpassword",
	}
	payloadBytes, _ := json.Marshal(removePayload)
	req, _ := http.NewRequest("DELETE", "/user/remove", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer valid_jwt_token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Case 2: Invalid credentials
	rows := initUserTestRows()
	suite.mock.ExpectQuery("SELECT (.+) FROM \"users\"").WillReturnRows(rows)

	removePayload = map[string]interface{}{
		"email":    "user1@example.com",
		"password": "wrongpassword",
	}
	payloadBytes, _ = json.Marshal(removePayload)
	req, _ = http.NewRequest("DELETE", "/user/remove", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer valid_jwt_token")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	var result models.Response
	assert.NoError(suite.T(), json.Unmarshal(w.Body.Bytes(), &result))
	assert.NotEmpty(suite.T(), result.Error)
}

func (suite *UserTestSuite) TestUserAnalyticsEndpoint() {
	envLoader := env.NewEnvLoader()
	envLoader.RegisterDefault("JWT_SECRET", "secret")

	// Case 1: Count users
	analyticsPayload := map[string]interface{}{
		"count": true,
	}
	payloadBytes, _ := json.Marshal(analyticsPayload)
	req, _ := http.NewRequest("POST", "/user/analytics", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer valid_jwt_token")
	w := httptest.NewRecorder()

	// Mock the SQL query for counting users
	suite.mock.ExpectQuery(`SELECT count\(\*\) FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	// Serve the HTTP request
	suite.router.ServeHTTP(w, req)

	// Check the status code and response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), float64(10), result["count"])

	// Case 2: Filter users by nationality and gender
	analyticsPayload = map[string]interface{}{
		"columnvalue_comparisons": []map[string]interface{}{
			{"column": "nationality", "value": "Nationality"},
			{"column": "gender", "value": "Male"},
		},
	}
	payloadBytes, _ = json.Marshal(analyticsPayload)
	req, _ = http.NewRequest("POST", "/user/analytics", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer valid_jwt_token")
	w = httptest.NewRecorder()

	// Mock the SQL query for filtering users
	userRows := sqlmock.NewRows([]string{
		"id", "first_name", "surname", "birthdate", "nationality", "role", "email", "password_hash",
		"phone_number", "address_id", "created_at", "updated_at", "last_login", "status", "gender",
		"preferred_language", "timezone", "salt",
	}).
		AddRow(1, "John", "Doe", "1990-01-01", "Nationality", "User", "user1@example.com", "mocked_hash", "123-456-7891",
			1, "2023-01-01 00:00:00", "2023-01-01 00:00:00", "2023-01-01 00:00:00", "Active", "Male", "English", "UTC", "mocked_salt").
		AddRow(2, "Jane", "Doe", "1991-02-02", "Nationality", "User", "user2@example.com", "mocked_hash", "123-456-7892",
			2, "2023-02-01 00:00:00", "2023-02-01 00:00:00", "2023-02-01 00:00:00", "Active", "Female", "English", "UTC", "mocked_salt")

	// Adjusted mock to expect the exact query
	suite.mock.ExpectQuery(`SELECT \* FROM "users" WHERE nationality LIKE \$1 AND gender LIKE \$2`).
		WithArgs("%Nationality%", "%Male%").
		WillReturnRows(userRows)

	// Serve the HTTP request
	suite.router.ServeHTTP(w, req)

	// Check the status code and response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	// Since the response contains a map with "data" key, extract users
	usersData, ok := responseBody["data"].([]interface{})
	assert.True(suite.T(), ok)
	assert.NotEmpty(suite.T(), usersData)
}



