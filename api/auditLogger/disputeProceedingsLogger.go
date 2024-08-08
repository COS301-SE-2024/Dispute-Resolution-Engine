package auditLogger

import (
	"api/db"
	"api/env"
	"api/models"
	"api/utilities"

	"gorm.io/gorm"
)

type DisputeProceedingsLogger struct {
	DB        *gorm.DB
	EnvReader env.Env
}

type LogJson struct {
	Message string
	Json    interface{}
}

func NewDisputeProceedingsLoggerDBInit() (DisputeProceedingsLogger,error) {
	DB, err := db.Init()
	if err != nil {
		return DisputeProceedingsLogger{}, err
	}
	return DisputeProceedingsLogger{DB: DB, EnvReader: env.NewEnvLoader()}, nil
}

func NewDisputeProceedingsLogger(db *gorm.DB) DisputeProceedingsLogger {
	return DisputeProceedingsLogger{DB: db, EnvReader: env.NewEnvLoader()}
}

func (d DisputeProceedingsLogger) LogDisputeProceedings(proceedingType models.EventTypes, eventData map[string]interface{}) error {
	// Initialize the logger
	logger := utilities.NewLogger()

	// Log the event data for debugging
	logger.Info(eventData)

	// Attempt to create a new event log entry in the database
	err := d.DB.Create(&models.EventLog{
		EventType: proceedingType,
		EventData: eventData,
	}).Error

	// Error handling
	if err != nil {
		logger.WithError(err).Error("Error logging dispute proceedings")
		return err
	}

	// Log success message
	logger.Info("Dispute proceedings logged successfully")
	return nil
}