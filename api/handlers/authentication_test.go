package handlers_test

import (
	"api/handlers"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
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

	suite.userRows = sqlmock.NewRows([]string{
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

func (suite *AuthTestSuite) TestLogin() {

}
