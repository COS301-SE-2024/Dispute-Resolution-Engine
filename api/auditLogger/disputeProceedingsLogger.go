package auditLogger

import (
	"api/env"
	"gorm.io/gorm"
)

type DisputeProceedingsLogger struct {
	DB        *gorm.DB
	EnvReader env.Env
}

type ProceedingType string

const (
    Notification ProceedingType = "NOTIFICATION"
    Dispute      ProceedingType = "DISPUTE"
    User         ProceedingType = "USER"
    Expert       ProceedingType = "EXPERT"
    WorkFlow     ProceedingType = "WORKFLOW"
)



func NewDisputeProceedingsLogger(db *gorm.DB) DisputeProceedingsLogger {
	return DisputeProceedingsLogger{DB: db, EnvReader: env.NewEnvLoader()}
}

func (d DisputeProceedingsLogger) LogDisputeProceedings() {

}