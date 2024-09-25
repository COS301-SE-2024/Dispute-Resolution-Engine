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
	GetWorkflowsWithLimitOffset(id, limit, offset *int) ([]models.Workflow, error)
	QueryTagsToRelatedWorkflow(workflowID uint64) ([]models.Tag, error)
	CreateWorkflow(workflow *models.Workflow) error
	CreateWokrflowTag(tag *models.WorkflowTags) error
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


func (wfmr *workflowModelReal) GetWorkflowsWithLimitOffset(id, limit, offset *int) ([]models.Workflow, error) {
	var workflows []models.Workflow

    // Create a query object
    query := wfmr.DB.Model(&models.Workflow{})

	// If id is provided, apply it
	if id != nil {
		query = query.Where("id = ?", *id)
	}

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


func (wfmr *workflowModelReal) QueryTagsToRelatedWorkflow(workflowID uint64) ([]models.Tag, error) {
	var tags []models.Tag

	// Create a query object
	query := wfmr.DB.Model(&models.Tag{})

	// If limit is provided, apply it
	query = query.Joins("JOIN workflow_tags ON tags.id = workflow_tags.tag_id")
	query = query.Where("workflow_tags.workflow_id = ?", workflowID)

	// Execute the query
	result := query.Find(&tags)

	// Handle any errors
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return tags, nil
}

func (wfmr *workflowModelReal) CreateWorkflow(workflow *models.Workflow) error {
	result := wfmr.DB.Create(workflow)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *workflowModelReal) CreateWokrflowTag(tag *models.WorkflowTags) error {
	result := wfmr.DB.Create(tag)

	if result.Error != nil {
		return result.Error
	}

	return nil
}