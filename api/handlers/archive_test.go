package handlers_test

import (
	"api/handlers"
	"api/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	mockDisputeCount = 10
)

type mockJwtModel struct {
	throwErrors bool
}

type mockAuditLogger struct {
}



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
func (m *mockJwtModel) JWTMiddleware(c *gin.Context) {}

func (m *mockJwtModel) GetClaims(c *gin.Context) (models.UserInfoJWT, error) {
	if m.throwErrors {
		return models.UserInfoJWT{}, errors.ErrUnsupported
	}
	return models.UserInfoJWT{
		ID:                0,
		FirstName:         "",
		Surname:           "",
		Birthdate:         time.Now(),
		Nationality:       "",
		Role:              "",
		Email:             "",
		PhoneNumber:       new(string),
		AddressID:         new(int64),
		Status:            "",
		Gender:            "",
		PreferredLanguage: new(string),
		Timezone:          new(string),
	}, nil

}

func TestArchive(t *testing.T) {
	suite.Run(t, new(ArchiveTestSuite))
}

//mock model auditlogger
func (m *mockAuditLogger) LogDisputeProceedings(proceedingType models.EventTypes, eventData map[string]interface{}) error {
	return nil
}

type ArchiveTestSuite struct {
	suite.Suite

	mock   sqlmock.Sqlmock
	db     *gorm.DB
	router *gin.Engine
	auditlogger *mockAuditLogger
}

// Runs before every test to set up the DB and routers
func (suite *ArchiveTestSuite) SetupTest() {
	mock, db, _ := mockDatabase()

	suite.mock = mock
	suite.db = db
	suite.auditlogger = &mockAuditLogger{}


	handler := handlers.Archive{
		Handler: handlers.Handler{
			DB: db,
			DisputeProceedingsLogger: suite.auditlogger,
		},
	}
	gin.SetMode("release")
	router := gin.Default()
	router.POST("/archive/search", handler.SearchArchive)
	suite.router = router

}



func mockDatabase() (sqlmock.Sqlmock, *gorm.DB, error) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		Conn:       conn,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, err
	}
	return mock, db, nil
}

func initCountRow(count int) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"count"}).AddRow(count)
}

func initDisputeRows() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id",
		"case_date",
		"workflow",
		"status",
		"title",
		"description",
		"complainant",
		"respondant",
		"resolved",
		"decision",
	})
	for i := 0; i < mockDisputeCount; i++ {
		rows = rows.AddRow(
			i,
			time.Now(),
			nil,
			"Awaiting Respondant",
			fmt.Sprintf("Dispute Title %d", i),
			fmt.Sprintf("Description %d", i),
			0,
			nil,
			true,
			"Unresolved",
		)
	}
	return rows
}

// Creates a new POST request to /archive/search using the passed-in payload
func createSearchRequest(req models.ArchiveSearchRequest) (*http.Request, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return http.NewRequest("POST", "/archive/search", bytes.NewReader(body))
}

func (suite *ArchiveTestSuite) TestBadRequestReturnsError() {
	req, _ := http.NewRequest("POST", "/archive/search", strings.NewReader(""))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert properties
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	var result models.Response
	assert.NoError(suite.T(), json.Unmarshal(w.Body.Bytes(), &result))
	assert.NotEmpty(suite.T(), result.Error)
}

func (suite *ArchiveTestSuite) TestReturnsValidJSON() {
	rows := initDisputeRows()
	suite.mock.ExpectQuery("^SELECT count(.+) FROM \"?disputes\"?.*").WillReturnRows(initCountRow(mockDisputeCount))
	suite.mock.ExpectQuery("^SELECT (.+) FROM \"?disputes\"?.*").WillReturnRows(rows)

	// Set up request + response
	searchTerm := "Hello"
	req, _ := createSearchRequest(models.ArchiveSearchRequest{
		Search: &searchTerm,
	})

	// Send request to router
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert properties
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	var result models.Response
	assert.Error(suite.T(), json.Unmarshal(w.Body.Bytes(), &result))
}
