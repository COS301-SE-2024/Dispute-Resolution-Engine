package auditLogger

import (
	"api/env"
	"api/models"
	"api/utilities"
	"encoding/json"

	"gorm.io/gorm"
)

type DisputeProceedingsLogger struct {
	DB        *gorm.DB
	EnvReader env.Env
}


func NewDisputeProceedingsLogger(db *gorm.DB) DisputeProceedingsLogger {
	return DisputeProceedingsLogger{DB: db, EnvReader: env.NewEnvLoader()}
}

func (d DisputeProceedingsLogger) LogDisputeProceedings(proceedingType models.EventTypes, jsonMessage string) error{
	// Parse the JSON message
	logger := utilities.NewLogger()

	var eventData map[string]interface{}
	err := json.Unmarshal([]byte(jsonMessage), &eventData)
	if err != nil {
		logger.WithError(err).Error("Error parsing JSON message")
		return err
	}

	// Log the dispute proceedings
	err = d.DB.Create(&models.EventLog{
		EventType: proceedingType,
		EventData: eventData,
	}).Error
	if err != nil {
		logger.WithError(err).Error("Error logging dispute proceedings")
		return err
	}
	return nil
}