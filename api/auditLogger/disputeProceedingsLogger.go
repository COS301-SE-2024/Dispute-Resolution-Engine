package auditLogger

import (
	"api/env"
	"api/models"
	"api/utilities"

	"gorm.io/gorm"
)

// func NewDisputeProceedingsLoggerDBInit() (DisputeProceedingsLoggerInterface,error) {
// 	DB, err := db.Init()
// 	if err != nil {
// 		return DisputeProceedingsLogger{}, err
// 	}
// 	return DisputeProceedingsLogger{DB: DB, EnvReader: env.NewEnvLoader()}, nil
// }

func NewDisputeProceedingsLogger(db *gorm.DB, envLoader env.Env) DisputeProceedingsLoggerInterface {
	return DisputeProceedingsLogger{DB: DisputeProceedingsLoggerReal{DB: db}, EnvReader: envLoader}
}

func (d DisputeProceedingsLogger) LogDisputeProceedings(proceedingType models.EventTypes, eventData map[string]interface{}) error {
	// Initialize the logger
	logger := utilities.NewLogger().LogWithCaller()

	// Log the event data for debugging
	logger.Info(eventData)

	// Attempt to create a new event log entry in the database
	err := d.DB.CreateLog(models.EventLog{
		EventType: proceedingType,
		EventData: eventData,
	})

	// Error handling
	if err != nil {
		logger.WithError(err).Error("Error logging dispute proceedings")
		return err
	}

	// Log success message
	logger.Info("Dispute proceedings logged successfully")
	return nil
}
