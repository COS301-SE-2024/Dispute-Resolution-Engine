package workflow

import (
	"orchestrator/db"

	"gorm.io/gorm"
)

type DBQuery interface {
	FetchWorkflowQuery(id int) (*db.Workflowdb, error)
	FetchUserQuery(id int64) (*db.User, error)
	FetchTagsByID(id int64) (*db.Tag, error)
	CreateWorkflows(workflow *db.Workflowdb) error
	CreateLabbelledWorkdlows(labelledWorkflow *db.LabelledWorkflow) error
	SaveWorkflowInstance(activeWorkflow *db.ActiveWorkflows) error
	DeleteLabelledWorkflow(labelledWorkflow *db.LabelledWorkflow) error
	
	FetchActiveWorkflows() ([]db.ActiveWorkflows, error)
	FetchActiveWorkflow(id int) (*db.ActiveWorkflows, error)
	SaveActiveWorkflow(activeWorkflow *db.ActiveWorkflows) error
}

type WorkflowQuery struct {
	DB *gorm.DB
}

func CreateDBQuery() *WorkflowQuery {
	Database, err := db.Init()
	if err != nil {
		return nil
	}
	return &WorkflowQuery{
		DB: Database,
	}
}

func (wfq *WorkflowQuery) FetchWorkflowQuery(id int) (*db.Workflowdb, error) {
	var workflow db.Workflowdb
	result := wfq.DB.First(&workflow, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &workflow, nil
}

func (wfq *WorkflowQuery) FetchUserQuery(id int64) (*db.User, error) {
	var user db.User
	result := wfq.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (wfq *WorkflowQuery) FetchTagsByID(id int64) (*db.Tag, error) {
	var tag db.Tag
	result := wfq.DB.First(&tag, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tag, nil
}

func (wfq *WorkflowQuery) CreateWorkflows(workflow *db.Workflowdb) error {
	result := wfq.DB.Create(workflow)
	if result.Error != nil {
		return result.Error
	}
	return nil
}