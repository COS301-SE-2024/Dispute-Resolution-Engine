package workflow

import (
	"encoding/json"
	"fmt"
	"orchestrator/db"
	"time"

	"gorm.io/gorm"
)

// Interface for the workflow API
type API interface {
	FetchWorkflow(id int) (*db.Workflowdb, error)
	StoreWorkflow(name string, workflow Workflow, categories []int64, Author int64) error
	UpdateWorkflow(id int, name *string, workflow *Workflow, categories *[]int64, author *int64) error
	FetchActiveWorkflows() ([]db.ActiveWorkflows, error)
	FetchActiveWorkflow(id int) (*db.ActiveWorkflows, error)
	UpdateActiveWorkflow(id int, workflowID *int, currentState *string, dateSubmitted *time.Time, stateDeadline *time.Time, workflowInstance *json.RawMessage) error
}

// APIWorkflow is the implementation of the API interface
type APIWorkflow struct {
	DB      *gorm.DB
	WfQuery DBQuery
}

// Workflow is the struct that represents a workflow
func CreateAPIWorkflow() *APIWorkflow {
	Database, err := db.Init()
	if err != nil {
		return nil
	}
	return &APIWorkflow{
		DB:      Database,
		WfQuery: CreateDBQuery(),
	}
}

func (api *APIWorkflow) FetchWorkflow(id int) (*db.Workflowdb, error) {
	// Fetch the workflow from the database
	workflow, err := api.WfQuery.FetchWorkflowQuery(id)
	if err != nil {
		return nil, err
	}

	return workflow, nil
}

func (api *APIWorkflow) StoreWorkflow(name string, workflow Workflow, categories []int64, Author int64) error {
	marshal, err := json.Marshal(workflow)
	if err != nil {
		fmt.Println("Error marshalling workflow: ")
		return err
	}

	//check if use exist in users table

	_, err = api.WfQuery.FetchUserQuery(Author)
	if err != nil {
		return err
	}

	//check if category exist in tags table
	for _, category := range categories {
		_, err := api.WfQuery.FetchTagsByID(category)
		if err != nil {
			return err
		}
	}

	//add entry in the db
	workflowDbEntry := &db.Workflowdb{
		Name:       name,
		Definition: marshal,
		AuthorID:   Author,
	}

	err = api.WfQuery.CreateWorkflows(workflowDbEntry)
	if err != nil {
		return err
	}

	//add associated tags
	for _, category := range categories {
		labelledWorkflow := &db.LabelledWorkflow{
			WorkflowID: workflowDbEntry.ID,
			TagID:      uint64(category),
		}
		err = api.WfQuery.CreateLabbelledWorkdlows(labelledWorkflow)
		if err != nil {
			return err
		}
	}

	return nil
}

func (api *APIWorkflow) UpdateWorkflow(id int, name *string, workflow *Workflow, categories *[]int64, author *int64) error {
	//get existign workflow
	existingWorkflow, err := api.WfQuery.FetchWorkflowQuery(id)
	if err != nil {
		return err
	}

	// Update the name if provided
	if name != nil {
		existingWorkflow.Name = *name
	}

	// Update the WorkflowDefinition if provided
	if workflow != nil {
		workflowDefinition, err := json.Marshal(*workflow)
		if err != nil {
			return err
		}
		existingWorkflow.Definition = workflowDefinition
	}

	// Update the AuthorID if provided
	if author != nil {
		existingWorkflow.AuthorID = *author
	}
	// Save the updated workflow
	err = api.WfQuery.SaveWorkflowInstance(existingWorkflow)
	if err != nil {
		return err
	}


	// Manage categories (tags) in labelled_workflow if provided
	if categories != nil {
		// Remove existing tags
		err = api.WfQuery.DeleteLabelledWorkfloByWorkflowId(existingWorkflow.ID)
		if err != nil {
			return err
		}

		// Insert new tags
		for _, categoryID := range *categories {
			labelledWorkflow := &db.LabelledWorkflow{
				WorkflowID: existingWorkflow.ID,
				TagID:      uint64(categoryID),
			}
			err = api.DB.Create(&labelledWorkflow).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (api *APIWorkflow) FetchActiveWorkflows() ([]db.ActiveWorkflows, error) {
	// Define a slice to hold the result
	var activeWorkflows []db.ActiveWorkflows

	// Use a join to fetch active workflows and their related workflow definitions
	result := api.DB.
		Table("active_workflows").
		Select("id, workflow as workflow_id, current_state, state_deadline, workflow_instance").
		Scan(&activeWorkflows)
	// Check for errors in the result
	if result.Error != nil {
		return nil, result.Error
	}

	return activeWorkflows, nil
}

func (api *APIWorkflow) FetchActiveWorkflow(id int) (*db.ActiveWorkflows, error) {
	var activeWorkflow db.ActiveWorkflows

	result := api.DB.
		Table("active_workflows").
		Select("id, workflow as workflow_id, current_state, state_deadline, workflow_instance").
		Where("id = ?", id).
		Scan(&activeWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return &activeWorkflow, nil
}

func (api *APIWorkflow) UpdateActiveWorkflow(id int, workflowID *int, currentState *string, dateSubmitted *time.Time, stateDeadline *time.Time, workflowInstance *json.RawMessage) error {
	// Fetch the active workflow
	var activeWorkflow db.ActiveWorkflows
	result := api.DB.First(&activeWorkflow, id)
	if result.Error != nil {
		return result.Error
	}

	// Update the workflowID if provided
	if workflowID != nil {
		activeWorkflow.WorkflowID = int64(*workflowID)
	}

	// Update the currentState if provided
	if currentState != nil {
		activeWorkflow.CurrentState = *currentState
	}

	// Update the dateSubmitted if provided
	if dateSubmitted != nil {
		activeWorkflow.DateSubmitted = *dateSubmitted
	}

	// Update the stateDeadline if provided
	if stateDeadline != nil {
		activeWorkflow.StateDeadline = *stateDeadline
	}

	// Update the workflowInstance if provided
	if workflowInstance != nil {
		activeWorkflow.WorkflowInstance = *workflowInstance
	}

	// Save the updated active workflow
	result = api.DB.Save(&activeWorkflow)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
