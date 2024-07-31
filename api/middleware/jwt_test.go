package middleware_test

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type JWTTestSuite struct {
	suite.Suite
	router *gin.Engine
	logger *utilities.Logger
}

func TestJWT(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}

func (suite *JWTTestSuite) SetupTest() {
	suite.router = gin.Default()
}

// TestGenerateJWT tests the GenerateJWT function
func (suite *JWTTestSuite) TestGenerateJWTSuccess() {
	// Create a mock user

	user := models.User{
		Email: "test@example.com",
	}
	jwtMiddleware := middleware.NewJwtMiddleware()
	// Call the GenerateJWT function
	token, err := jwtMiddleware.GenerateJWT(user)

	// Assert that no error is returned
	suite.NoError(err)
	suite.NotEmpty(token)
}

// TestGenerateJWTError tests the GenerateJWT function when an error occurs
func (suite *JWTTestSuite) TestGenerateJWTError() {
	// Create a mock user
	user := models.User{
		Email: "test@example.com",
	}

	// Call the GenerateJWT function with an invalid user
	jwtMiddleware := middleware.NewJwtMiddleware()
	token, err := jwtMiddleware.GenerateJWT(user)

	// Assert that an error is returned
	suite.Error(err)
	suite.Empty(token)
}

// mocks
