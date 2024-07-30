package middleware_test

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

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

	// Call the GenerateJWT function
	token, err := middleware.GenerateJWT(user)

	// Assert that no error is returned
	suite.NoError(err)
	suite.NotEmpty(token)
}

// TestGenerateJWTError tests the GenerateJWT function when an error occurs
func (suite *JWTTestSuite) TestGenerateJWTError() {
	// Create a mock user
	user := models.User{
		Email:    "test@example.com",
	}

	// Call the GenerateJWT function with an invalid user
	token, err := middleware.GenerateJWT(user)

	// Assert that an error is returned
	suite.Error(err)
	suite.Empty(token)
}



// mocks

type JWTTestSuite struct {
	suite.Suite
	router *gin.Engine
	logger *utilities.Logger
}

type LoggerMock interface{
	LogWithCaller() LoggerMock
	WithError(err error) LoggerMock
	Error(msg string) 
	Info(msg string)
}

type Env interface {
	LoadFromFile(files ...string)
	Register(key string)
	RegisterDefault(key, fallback string)
	Get(key string) (string, error)
}

type loggerMockImpl struct{
	mock.Mock
}

func (m *loggerMockImpl) LogWithCaller() LoggerMock {
	args := m.Called()
	return args.Get(0).(LoggerMock)
}

func (m *loggerMockImpl) WithError(err error) LoggerMock {
	args := m.Called(err)
	return args.Get(0).(LoggerMock)
}

func (m *loggerMockImpl) Error(msg string) {
	m.Called(msg)
}

func (m *loggerMockImpl) Info(msg string) {
	m.Called(msg)
}