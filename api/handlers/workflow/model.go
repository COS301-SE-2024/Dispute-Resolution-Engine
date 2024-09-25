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
	GetWorkflowByID(id uint64) (*models.Workflow, error)
	GetActiveWorkflowByWorkflowID(workflowID uint64) (*models.ActiveWorkflows, error)
	QueryTagsToRelatedWorkflow(workflowID uint64) ([]models.Tag, error)
	FindDisputeByID(id uint64) (*models.Dispute, error)

	CreateWorkflow(workflow *models.Workflow) error
	CreateWorkflowTag(tag *models.WorkflowTags) error
	CreateActiveWorkflow(workflow *models.ActiveWorkflows) error

	UpdateWorkflow(workflow *models.Workflow) error
	UpdateActiveWorkflow(workflow *models.ActiveWorkflows) error

	DeleteTagsByWorkflowID(workflowID uint64) error
	DeleteWorkflow(wf *models.Workflow) error
	DeleteActiveWorkflow(workflow *models.ActiveWorkflows) error
}

type Workflow struct {
	DB                       WorkflowDBModel
	EnvReader                env.Env
	Emailer                  notifications.EmailSystem
	Jwt                      middleware.Jwt
	DisputeProceedingsLogger auditLogger.DisputeProceedingsLoggerInterface
}

type workflowModelReal struct {
	DB  *gorm.DB
	env env.Env
}

func NewWorkflowHandler(db *gorm.DB, envReader env.Env) Workflow {
	return Workflow{
		DB:                       &workflowModelReal{DB: db, env: envReader},
		Emailer:                  notifications.NewHandler(db),
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

func (wfmr *workflowModelReal) GetWorkflowByID(id uint64) (*models.Workflow, error) {
	var workflow models.Workflow

	// Create a query object
	query := wfmr.DB.Model(&models.Workflow{})
	query = query.Where("id = ?", id)

	// Execute the query
	result := query.First(&workflow)

	return &workflow, result.Error
}

func (wfmr *workflowModelReal) FindDisputeByID(id uint64) (*models.Dispute, error) {
	var dispute models.Dispute

	// Create a query object
	results := wfmr.DB.First(&dispute, id)
	return &dispute, results.Error
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

func (wfmr *workflowModelReal) CreateWorkflowTag(tag *models.WorkflowTags) error {
	result := wfmr.DB.Create(tag)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *workflowModelReal) UpdateWorkflow(workflow *models.Workflow) error {
	result := wfmr.DB.Save(workflow)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *workflowModelReal) DeleteTagsByWorkflowID(workflowID uint64) error {
	result := wfmr.DB.Where("workflow_id = ?", workflowID).Delete(&models.WorkflowTags{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *workflowModelReal) DeleteWorkflow(wf *models.Workflow) error {
	result := wfmr.DB.Delete(wf)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *workflowModelReal) CreateActiveWorkflow(workflow *models.ActiveWorkflows) error {
	result := wfmr.DB.Create(workflow)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *workflowModelReal) DeleteActiveWorkflow(workflow *models.ActiveWorkflows) error {
	result := wfmr.DB.Delete(workflow)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *workflowModelReal) GetActiveWorkflowByWorkflowID(workflowID uint64) (*models.ActiveWorkflows, error) {
	var activeWorkflow models.ActiveWorkflows

	query := wfmr.DB.First(&activeWorkflow, workflowID)

	return &activeWorkflow, query.Error
}

func (wfmr *workflowModelReal) UpdateActiveWorkflow(workflow *models.ActiveWorkflows) error {
	result := wfmr.DB.Save(workflow)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
