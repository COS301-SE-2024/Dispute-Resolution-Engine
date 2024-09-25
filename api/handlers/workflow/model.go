package workflow

import (
	"api/auditLogger"
	"api/env"
	"api/handlers/notifications"
	"api/middleware"
	"api/models"
	"gorm.io/gorm"
)

type WorkflowDBModel interface {
	GetWorkflowsWithLimitOffset(limit, offset *int) ([]models.Workflow, error)
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


func (wfmr *workflowModelReal) GetWorkflowsWithLimitOffset(limit, offset *int) ([]models.Workflow, error) {
	var workflows []models.Workflow

    // Create a query object
    query := wfmr.DB.Model(&models.Workflow{})

    // If limit is provided, apply it
    if limit != nil {
        query = query.Limit(*limit)
    }

    // If offset is provided, apply it
    if offset != nil {
        query = query.Offset(*offset)
    }

    // Execute the query
    result := query.Find(&workflows)

    // Handle any errors
    if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
        return nil, result.Error
    }

    return workflows, nil
}
