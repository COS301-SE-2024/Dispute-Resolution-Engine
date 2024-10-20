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
	CreateLabbelledWorkflows(labelledWorkflow *db.WorkflowTags) error
	SaveWorkflowInstance(activeWorkflow *db.Workflowdb) error
	DeleteLabelledWorkflowByWorkflowId(id uint64) error

	FetchActiveWorkflows() ([]db.ActiveWorkflows, error)
	FetchActiveWorkflow(id int) (*db.ActiveWorkflows, error)
	SaveActiveWorkflowInstance(activeWorkflow *db.ActiveWorkflows) error
}

type WorkflowQuery struct {
	DB *gorm.DB
}

func CreateWorkflowQuery() *WorkflowQuery {
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

func (wfq *WorkflowQuery) CreateLabbelledWorkflows(labelledWorkflow *db.WorkflowTags) error {
	result := wfq.DB.Create(labelledWorkflow)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (wfq *WorkflowQuery) SaveWorkflowInstance(workflow *db.Workflowdb) error {
	result := wfq.DB.Save(workflow)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (wfq *WorkflowQuery) DeleteLabelledWorkflowByWorkflowId(id uint64) error {
	err := wfq.DB.Where("workflow_id = ?", id).Delete(&db.WorkflowTags{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (wfq *WorkflowQuery) FetchActiveWorkflows() ([]db.ActiveWorkflows, error) {
	var activeWorkflows []db.ActiveWorkflows
	result := wfq.DB.Table("active_workflows").
		Select("id, workflow as workflow_id, current_state, state_deadline, workflow_instance").
		Scan(&activeWorkflows)
	if result.Error != nil {
		return nil, result.Error
	}
	return activeWorkflows, nil
}

func (wfq *WorkflowQuery) FetchActiveWorkflow(id int) (*db.ActiveWorkflows, error) {
	var activeWorkflow db.ActiveWorkflows
	result := wfq.DB.
		Table("active_workflows").
		Select("id, workflow, current_state, state_deadline, workflow_instance").
		Where("id = ?", id).
		Scan(&activeWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return &activeWorkflow, nil
}

func (wfq *WorkflowQuery) SaveActiveWorkflowInstance(activeWorkflow *db.ActiveWorkflows) error {
	result := wfq.DB.Save(activeWorkflow)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
