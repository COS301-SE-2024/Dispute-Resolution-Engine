package auditLogger

import (
	"api/env"
	"api/models"

	"gorm.io/gorm"
)

type DisputeProceedingsLogger struct {
	DB        *gorm.DB
	EnvReader env.Env
}


func NewDisputeProceedingsLogger(db *gorm.DB) DisputeProceedingsLogger {
	return DisputeProceedingsLogger{DB: db, EnvReader: env.NewEnvLoader()}
}

func (d DisputeProceedingsLogger) LogDisputeProceedings(proceedingType models.EventTypes, jsonMessage string) {

}