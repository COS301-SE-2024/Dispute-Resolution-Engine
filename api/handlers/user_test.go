package handlers_test

import (
	"api/handlers"
	"api/middleware"
	"api/models"
	"api/utilities"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

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
	router.PUT("/user/profile",handler.UpdateUser)
	router.GET("/user/profile", handler.GetUser)
	router.PUT("/user/profile/address", handler.UpdateUserAddress)
	router.DELETE("/user/remove", handler.RemoveAccount)
	router.POST("/user/analytics", handler.UserAnalyticsEndpoint)

	suite.mock = mock
	suite.db = db
	suite.router = router
}

func initUserRows() *sqlmock.Rows {
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
	// Add mock data here
	// Add mock data here
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
			nil, // address_id is null
			"2023-01-01 00:00:00", // created_at
			"2023-01-01 00:00:00", // updated_at
			nil, // last_login is null
			"Active",
			"Male",
			"English",
			"UTC",
			"mocked_salt",
		)
	}
	return rows
}

