package auditLogger_test

import (
	"api/auditLogger"
	"api/env"
	"api/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type DisputeProceedingsLoggerTestSuite struct {
	suite.Suite
	DBMock  *auditLogger.MockDBDisputeProceedingsLogger
	EnvMock *env.MockEnv
}

func (suite *DisputeProceedingsLoggerTestSuite) SetupTest() {
	suite.DBMock = &auditLogger.MockDBDisputeProceedingsLogger{}
	suite.EnvMock = &env.MockEnv{}

}
func (suite *DisputeProceedingsLoggerTestSuite) TestNewDisputeProceedingsLogger() {
	db := &gorm.DB{}
	envLoader := &env.MockEnv{}

	logger := auditLogger.NewDisputeProceedingsLogger(db, envLoader)

	suite.NotNil(logger)
	suite.IsType(auditLogger.DisputeProceedingsLogger{}, logger)
	suite.Equal(envLoader, logger.(auditLogger.DisputeProceedingsLogger).EnvReader)
}



func TestDisputeProceedingsLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeProceedingsLoggerTestSuite))
}
func (suite *DisputeProceedingsLoggerTestSuite) TestLogDisputeProceedings_Success() {
	logger := auditLogger.DisputeProceedingsLogger{DB: suite.DBMock, EnvReader: suite.EnvMock}

	eventType := models.EventTypes("TestEvent")
	eventData := map[string]interface{}{"key": "value"}

	err := logger.LogDisputeProceedings(eventType, eventData)

	suite.NoError(err)
}

func (suite *DisputeProceedingsLoggerTestSuite) TestLogDisputeProceedings_Error() {
	suite.DBMock.Error = errors.New("db error")
	suite.DBMock.ThrowError = true
	logger := auditLogger.DisputeProceedingsLogger{DB: suite.DBMock, EnvReader: suite.EnvMock}


	eventType := models.EventTypes("TestEvent")
	eventData := map[string]interface{}{"key": "value"}

	expectedError := errors.New("db error")

	err := logger.LogDisputeProceedings(eventType, eventData)

	suite.Error(err)
	suite.Equal(expectedError, err)
}