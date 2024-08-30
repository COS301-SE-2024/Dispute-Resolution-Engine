package handlers_test

import (
	"api/handlers"
	"api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)


type UtilityTestSuite struct {
	suite.Suite

	mock   sqlmock.Sqlmock
	db     *gorm.DB
	router *gin.Engine
	mockenvLoader *mockenvLoader
	mockJwt *mockJwt
}

func TestCountries(t *testing.T) {
	suite.Run(t, new(UtilityTestSuite))
}

// Runs before every test to set up the DB and routers
func (suite *UtilityTestSuite) SetupTest() {
	mock, db, _ := mockDatabase()

	suite.mock = mock
	suite.db = db
	suite.mockenvLoader = &mockenvLoader{}
	suite.mockJwt = &mockJwt{}

	handler := handlers.Utility{
		Handler: handlers.Handler{
			DB: suite.db,
			EnvReader: suite.mockenvLoader,
			Jwt: suite.mockJwt,
		},
	}
	
	gin.SetMode("release")
	router := gin.Default()
	router.GET("/countries", handler.GetCountries)
	router.GET("/statuses", handler.GetDisputeStatuses)
	suite.router = router

}

func createCountryRows(count int) (*sqlmock.Rows, []models.Country) {
	rows := sqlmock.NewRows([]string{
		"country_code",
		"country_name",
	})
	countries := make([]models.Country, count)
	for i := 0; i < count; i++ {
		rows = rows.AddRow(
			fmt.Sprint(i),
			fmt.Sprintf("Country %d", i),
		)
		countries[i] = models.Country{
			CountryCode: fmt.Sprint(i),
			CountryName: fmt.Sprintf("Country %d", i),
		}
	}
	return rows, countries
}

func createStatusRows(count int) (*sqlmock.Rows, []string) {
	rows := sqlmock.NewRows([]string{
		"enum_value",
	})
	countries := make([]string, count)
	for i := 0; i < count; i++ {
		status := fmt.Sprintf("Status %d", i)
		rows = rows.AddRow(status)

		countries[i] = status
	}
	return rows, countries
}

func (suite *UtilityTestSuite) TestReturnsCorrectCountries() {
	rows, data := createCountryRows(5)
	suite.mock.ExpectQuery("^SELECT (.+) FROM \"?countries\"?.*").WillReturnRows(rows)

	// Send request
	req, _ := http.NewRequest("GET", "/countries", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert properties
	var result struct {
		Data []models.Country `json:"data"`
	}

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.NoError(suite.T(), json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(suite.T(), data, result.Data)
}

func (suite *UtilityTestSuite) TestDisputeStatusReturnsCorrectStatus() {
	rows, data := createStatusRows(10)
	suite.mock.ExpectQuery("^SELECT (.+) FROM \"?pg_type\"?.*").WillReturnRows(rows)

	// Send request
	req, _ := http.NewRequest("GET", "/statuses", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert properties
	var result struct {
		Data []string `json:"data"`
	}

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.NoError(suite.T(), json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(suite.T(), data, result.Data)
}

