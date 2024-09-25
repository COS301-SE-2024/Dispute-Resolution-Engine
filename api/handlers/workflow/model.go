package workflow

import (
	"api/auditLogger"
	"api/env"
	"api/handlers/notifications"
	"api/middleware"

	"gorm.io/gorm"
)

type WorkflowDBModel interface {
}

type Workflow struct {
	DB                       WorkflowDBModel
	EnvReader                env.Env
	Emailer                  notifications.EmailSystem
	Jwt                      middleware.Jwt
	DisputeProceedingsLogger auditLogger.DisputeProceedingsLoggerInterface
}

type workflowModelReal struct {
	DB *gorm.DB
	env env.Env
}

func NewWorkflowHandler(db *gorm.DB, envReader env.Env) Workflow {
	return Workflow{
		DB:                       &workflowModelReal{DB: db, env: envReader},
		Emailer: 				 notifications.NewHandler(db),
		EnvReader:                env.NewEnvLoader(),
		Jwt:                      middleware.NewJwtMiddleware(),
		DisputeProceedingsLogger: auditLogger.NewDisputeProceedingsLogger(db, envReader),
	}
}

