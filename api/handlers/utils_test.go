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

type CountriesTestSuite struct {
	suite.Suite

	mock   sqlmock.Sqlmock
	db     *gorm.DB
	router *gin.Engine
}

func TestCountries(t *testing.T) {
	suite.Run(t, new(CountriesTestSuite))
}

// Runs before every test to set up the DB and routers
func (suite *CountriesTestSuite) SetupTest() {
	mock, db, _ := mockDatabase()

	handler := handlers.NewUtilitiesHandler(db)
	gin.SetMode("release")
	router := gin.Default()
	router.GET("/countries", handler.GetCountries)

	suite.mock = mock
	suite.db = db
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

func (suite *CountriesTestSuite) TestReturnsCorrectCountries() {
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
