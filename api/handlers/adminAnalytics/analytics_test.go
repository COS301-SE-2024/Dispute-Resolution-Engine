package adminanalytics_test

import (
	adminanalytics "api/handlers/adminAnalytics"
	"api/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type mockAdminAnalyticsModel struct {
	throwErrors bool
}

type mockJwtModel struct {
	throwErrors bool
	returnUser  models.UserInfoJWT
}

type mockEnv struct {
	throwErrors bool
	Error       error
}

type AdminAnalyticsErrorTestSuite struct {
	suite.Suite
	adminAnalyticsMock *mockAdminAnalyticsModel
	jwtMock            *mockJwtModel
	envMock            *mockEnv
	router             *gin.Engine
}

func (suite *AdminAnalyticsErrorTestSuite) SetupTest() {
	suite.adminAnalyticsMock = &mockAdminAnalyticsModel{}
	suite.jwtMock = &mockJwtModel{}
	suite.envMock = &mockEnv{}
	router := gin.Default()

	handler := adminanalytics.AdminAnalyticsHandler{DB: suite.adminAnalyticsMock, JWT: suite.jwtMock, EnvReader: suite.envMock}
	router.Use(suite.jwtMock.JWTMiddleware)
	router.GET("/time/estimation", handler.GetTimeEstimation)
	router.GET("/dispute/countries", handler.GetDisputeGrouping) //by status or country
	router.POST("/stats/:table", handler.GetTableStats)
	suite.router = router
}

func TestAdminAnalyticsErrors(t *testing.T) {
	suite.Run(t, new(AdminAnalyticsErrorTestSuite))
}

/*-------------------------------MOCK MODELS-----------------------------------------*/
//Analytics model mocks

func (a *mockAdminAnalyticsModel) CalculateAverageResolutionTime() (float64, error) {
	if a.throwErrors {
		return 0, errors.ErrUnsupported
	}

	return 0, nil
}

// DisputeStatusCount represents the count of disputes grouped by status.
type DisputeStatusCount struct {
	Status string
	Count  int64
}

func (h *mockAdminAnalyticsModel) CountDisputesByMonth(
	tableName string,
	dateColumn string,
) (map[string]int64, error) {
	return nil, nil
}


func (h *mockAdminAnalyticsModel) CalculateAverageResolutionTimeByMonth() (map[string]float64, error) {
	return nil, nil
}

// GetDisputeGroupingByStatus counts the number of disputes grouped by their statuses.
func (h *mockAdminAnalyticsModel) CountRecordsWithGroupBy(
	tableName string,
	column *string,
	value *interface{},
	groupBy *string,
) (map[string]int64, error) {

	if h.throwErrors {
		return nil, errors.ErrUnsupported
	}

	return nil, nil
}

// GetDisputeGroupingByCountry counts the number of disputes grouped by the country of the complainant.
func (h *mockAdminAnalyticsModel) GetDisputeGroupingByCountry() (map[string]int, error) {
	if h.throwErrors {
		return nil, errors.ErrUnsupported
	}
	return nil, nil
}

// JWT Mocks
func (m *mockJwtModel) GenerateJWT(user models.User) (string, error) {
	if m.throwErrors {
		return "", errors.ErrUnsupported
	}
	return "mock", nil
}
func (m *mockJwtModel) StoreJWT(email string, jwt string) error {
	if m.throwErrors {
		return errors.ErrUnsupported
	}
	return nil
}
func (m *mockJwtModel) GetJWT(email string) (string, error) {
	if m.throwErrors {
		return "", errors.ErrUnsupported
	}
	return "", nil
}
func (m *mockJwtModel) JWTMiddleware(c *gin.Context) {
	if m.throwErrors {
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func (m *mockJwtModel) GetClaims(c *gin.Context) (models.UserInfoJWT, error) {
	if m.throwErrors {
		return models.UserInfoJWT{}, errors.ErrUnsupported
	}
	return m.returnUser, nil

}

//mock env

func (m *mockEnv) LoadFromFile(files ...string) {
}

func (m *mockEnv) Register(key string) {
}

func (m *mockEnv) RegisterDefault(key, fallback string) {
}

func (m *mockEnv) Get(key string) (string, error) {
	if m.throwErrors {
		return "", m.Error
	}
	return "", nil
}

// ---------------------------------------------------------------- GET TIME ESTIMATION ----------------------------------------------------------------

func (suite *AdminAnalyticsErrorTestSuite) TestGetTimeEstimationUnauthorized() {
	req, _ := http.NewRequest("GET", "/time/estimation", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.jwtMock.throwErrors = true
	var result struct {
		Error string `json:"error"`
	}

	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
}

func (suite *AdminAnalyticsErrorTestSuite) TestGetTimeEstimationError() {
	suite.adminAnalyticsMock.throwErrors = true
	suite.jwtMock.returnUser.Role = "admin"
	req, _ := http.NewRequest("GET", "/time/estimation", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Equal("No disputes Have been resolved yet", result.Error)
}

func (suite *AdminAnalyticsErrorTestSuite) TestGetTimeEstimationSuccess() {
	suite.adminAnalyticsMock.throwErrors = false
	suite.jwtMock.returnUser.Role = "admin"
	req, _ := http.NewRequest("GET", "/time/estimation", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]int `json:"data"`
	}

	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

// ---------------------------------------------------------------- GET DISPUTE GROUPING ----------------------------------------------------------------

// func (suite *AdminAnalyticsErrorTestSuite) TestGetDisputeGroupingUnauthorized() {
// 	req, _ := http.NewRequest("GET", "/dispute/countries", nil)
// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)
// 	suite.jwtMock.throwErrors = true
// 	var result struct {
// 		Error string `json:"error"`
// 	}

// 	suite.Equal(http.StatusOK, w.Code)
// 	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
// }

func (suite *AdminAnalyticsErrorTestSuite) TestGetDisputeGroupingError() {
	suite.adminAnalyticsMock.throwErrors = true
	suite.jwtMock.returnUser.Role = "admin"
	req, _ := http.NewRequest("GET", "/dispute/countries", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Equal("Failed to get dispute grouping by country", result.Error)
}

func (suite *AdminAnalyticsErrorTestSuite) TestGetDisputeGroupingSuccess() {
	suite.adminAnalyticsMock.throwErrors = false
	suite.jwtMock.returnUser.Role = "admin"
	req, _ := http.NewRequest("GET", "/dispute/countries", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]int `json:"data"`
	}

	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
}

// ---------------------------------------------------------------- GET TABLE STATS ----------------------------------------------------------------

// func (suite *AdminAnalyticsErrorTestSuite) TestGetTableStatsUnauthorized() {
// 	req, _ := http.NewRequest("POST", "/stats/disputes", nil)
// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)
// 	suite.jwtMock.throwErrors = true
// 	var result struct {
// 		Error string `json:"error"`
// 	}

// 	suite.Equal(http.StatusBadRequest, w.Code)
// 	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
// 	suite.NotEmpty(result.Error)
// }

func (suite *AdminAnalyticsErrorTestSuite) TestGetTableStatsInvalidTable() {
	suite.jwtMock.returnUser.Role = "admin"
	req, _ := http.NewRequest("POST", "/stats/$", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *AdminAnalyticsErrorTestSuite) TestGetTableStatsError() {
	suite.adminAnalyticsMock.throwErrors = true
	suite.jwtMock.returnUser.Role = "admin"

	body := `{"group": "country", "where": {"": {"before": "2021-01-01", "after": "2020-01-01"}}}"}`
	req, _ := http.NewRequest("POST", "/stats/disputes", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *AdminAnalyticsErrorTestSuite) TestGetTableStatsSuccess() {
	suite.adminAnalyticsMock.throwErrors = false
	suite.jwtMock.returnUser.Role = "admin"
	body := `{"group": "country", "where": {"resolved": {"before": "2021-01-01", "after": "2020-01-01"}}}"}`
	req, _ := http.NewRequest("POST", "/stats/disputes", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]int `json:"data"`
	}

	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
}
