package handlers_test

import (
	"api/handlers"
	"api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	mockDisputeCount = 10
)

// Initializes a GORM instance with a mocked SQL database
func initMockGorm() (sqlmock.Sqlmock, *gorm.DB, error) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		Conn:       conn,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	return mock, db, err
}

func initCountRow(count int) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"count"}).AddRow(count)
}

func initRows() *sqlmock.Rows {
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

func initArchiveRouter(db *gorm.DB) *gin.Engine {
	archiveHandler := handlers.Archive{DB: db}
	r := gin.Default()
	r.POST("/archive/search", archiveHandler.SearchArchive)
	return r
}

// Creates a new POST request to /archive/search using the passed-in payload
func createSearchRequest(req models.ArchiveSearchRequest) (*http.Request, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return http.NewRequest("POST", "/archive/search", bytes.NewReader(body))
}

func TestSearchInvalidJSONShouldReturnError(t *testing.T) {
	// Initialize mock database
	_, db, _ := initMockGorm()
	router := initArchiveRouter(db)

	// Set up request + response
	req, _ := http.NewRequest("POST", "/archive/search", strings.NewReader(""))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert properties
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var result models.Response
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.NotEmpty(t, result.Error)
}

func TestSearchReturnsValidJSON(t *testing.T) {
	// Initialize mock database
	mock, db, _ := initMockGorm()
	rows := initRows()

	// Set up API route
	router := initArchiveRouter(db)

	// Set up request + response
	searchTerm := "Hello"
	req, _ := createSearchRequest(models.ArchiveSearchRequest{
		Search: &searchTerm,
	})

	// Mock SQL queries
	mock.ExpectQuery("^SELECT count(.+) FROM \"?disputes\"?.*").WillReturnRows(initCountRow(mockDisputeCount))
	mock.ExpectQuery("^SELECT (.+) FROM \"?disputes\"?.*").WillReturnRows(rows)

	// Send request to router
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert properties
	assert.Equal(t, 200, w.Code)

	var result models.Response
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))

}
