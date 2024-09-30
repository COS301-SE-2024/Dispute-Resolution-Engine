package auditLogger

import (
	"api/env"
	"api/models"

	"gorm.io/gorm"
)

type DisputeProceedingsLogger struct {
	DB        DBDisputeProceedingsLogger
	EnvReader env.Env
}

type LogJson struct {
	Message string
	Json    interface{}
}

type DisputeProceedingsLoggerInterface interface {
	LogDisputeProceedings(proceedingType models.EventTypes, eventData map[string]interface{}) error
}


type DBDisputeProceedingsLogger interface {
	CreateLog(log models.EventLog) error
}

type DisputeProceedingsLoggerReal struct {
	DB *gorm.DB
}


func (d DisputeProceedingsLoggerReal) CreateLog(log models.EventLog) error {
	return d.DB.Create(&log).Error
}

