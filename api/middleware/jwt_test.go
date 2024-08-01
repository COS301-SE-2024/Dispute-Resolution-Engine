package middleware_test

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type JWTTestSuite struct {
	suite.Suite
	logger *utilities.Logger
	mockEnvLoader *MockEnvLoader
	jwtMiddleware *middleware.JwtMiddleware
}

func TestJWT(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}

func (suite *JWTTestSuite) SetupTest() {
	suite.logger = utilities.NewLogger().LogWithCaller()
	suite.mockEnvLoader = new(MockEnvLoader)
	suite.jwtMiddleware = &middleware.JwtMiddleware{
		EnvLoader: suite.mockEnvLoader,
		Logger:    suite.logger,
	}
}

// // TestGenerateJWT tests the GenerateJWT function
// func (suite *JWTTestSuite) TestGenerateJWTSuccess() {
// 	// Create a mock user

// 	user := models.User{
// 		Email: "test@example.com",
// 	}
// 	jwtMiddleware := middleware.NewJwtMiddleware()
// 	// Call the GenerateJWT function
// 	token, err := jwtMiddleware.GenerateJWT(user)

// 	// Assert that no error is returned
// 	suite.NoError(err)
// 	suite.NotEmpty(token)
// }

// TestGenerateJWTSuccess tests the GenerateJWT function when successful
func (suite *JWTTestSuite) TestGenerateJWTSuccess() {
	// Create a mock user
	user := models.User{
		Email: "test@example.com",
	}

	// Set up the mock to return a valid JWT secret
	suite.mockEnvLoader.On("Get", "JWT_SECRET").Return("test_secret", nil)

	// Call the GenerateJWT function
	token, err := suite.jwtMiddleware.GenerateJWT(user)

	// Assert that no error is returned
	suite.NoError(err)
	suite.NotEmpty(token)

	// Assert that the mock was called as expected
	suite.mockEnvLoader.AssertCalled(suite.T(), "Get", "JWT_SECRET")
}

// TestGenerateJWTError tests the GenerateJWT function when an error occurs
func (suite *JWTTestSuite) TestGenerateJWTError() {
	// Create a mock user
	user := models.User{
		Email: "test@example.com",
	}

	// Set up the mock to return an error when retrieving the JWT secret
	suite.mockEnvLoader.On("Get", "JWT_SECRET").Return("", errors.New("failed to get JWT secret"))

	// Call the GenerateJWT function
	token, err := suite.jwtMiddleware.GenerateJWT(user)

	// Assert that an error is returned
	suite.Error(err)
	suite.Empty(token)

	// Assert that the mock was called as expected
	suite.mockEnvLoader.AssertCalled(suite.T(), "Get", "JWT_SECRET")
}

// mocks

type MockEnvLoader struct {
    mock.Mock
}

func (m *MockEnvLoader) Get(key string) (string, error) {
    args := m.Called(key)
    return args.String(0), args.Error(1)
}
func (m *MockEnvLoader) Register(key string) {
	m.Called(key)
}

func (m *MockEnvLoader) RegisterDefault(key, fallback string) {

	m.Called(key, fallback)
}

func (m *MockEnvLoader) LoadFromFile(files ...string) {
	m.Called(files)
}
