package workflow

import (
	"encoding/json"
	"fmt"
	"orchestrator/db"
	"time"

	"gorm.io/gorm"
)

// ----------------------------API--------------------------------
type API interface {
	Fetch(id int) (*db.Workflowdb, error)
	Store(workflow Workflow, categories []int64, Author int64) error
	Update(id int, name *string, workflow *Workflow, categories *[]int64, author *int64) error
}

type APIWorkflow struct {
	DB *gorm.DB
}

func CreateAPIWorkflow() *APIWorkflow {
	Database, err := db.Init()
	if err != nil {
		return nil
	}
	return &APIWorkflow{
		DB: Database,
	}
}

func (api *APIWorkflow) FetchWorkflow(id int) (*db.Workflowdb, error) {
	var workflow db.Workflowdb
	result := api.DB.First(&workflow, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &workflow, nil
}

func (api *APIWorkflow) StoreWorkflow(name string, workflow Workflow, categories []int64, Author int64) error {
	marshal, err := json.Marshal(workflow)
	if err != nil {
		fmt.Println("Error marshalling workflow: ")
		return err
	}

	//check if use exist in users table

	var user db.User
	result := api.DB.First(&user, Author)
	if result.Error != nil {
		return result.Error
	}

	//check if category exist in tags table
	for _, category := range categories {
		var tag db.Tag
		result := api.DB.First(&tag, category)
		if result.Error != nil {
			return result.Error
		}
	}

	//add entry in the db
	workflowDbEntry := &db.Workflowdb{
		Name:       name,
		Definition: marshal,
		AuthorID:   Author,
	}

	result = api.DB.Create(&workflowDbEntry)
	if result.Error != nil {
		return result.Error
	}

	//add associated tags
	for _, category := range categories {
		labelledWorkflow := &db.LabelledWorkflow{
			WorkflowID: workflowDbEntry.ID,
			TagID:      uint64(category),
		}
		result = api.DB.Create(&labelledWorkflow)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (api *APIWorkflow) UpdateWorkflow(id int, name *string, workflow *Workflow, categories *[]int64, author *int64) error {
	//get existign workflow
	var existingWorkflow db.Workflowdb
	result := api.DB.First(&existingWorkflow, id)
	if result.Error != nil {
		return result.Error
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
	result = api.DB.Save(&existingWorkflow)
	if result.Error != nil {
		return result.Error
	}

	// Manage categories (tags) in labelled_workflow if provided
	if categories != nil {
		// Remove existing tags
		err := api.DB.Where("workflow_id = ?", existingWorkflow.ID).Delete(&db.LabelledWorkflow{}).Error
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

func (api *APIWorkflow) FetchActiveWorkflows() ([]ActiveWorkflowsResponse, error) {
	// Define a slice to hold the result
	var activeWorkflows []ActiveWorkflowsResponse

	// Use a join to fetch active workflows and their related workflow definitions
	result := api.DB.
		Table("active_workflows").
		Select("active_workflows.id, active_workflows.workflow as workflow_id, workflows.definition as workflow_definition, active_workflows.current_state, active_workflows.state_deadline").
		Joins("join workflows on workflows.id = active_workflows.workflow").
		Scan(&activeWorkflows)
	// Check for errors in the result
	if result.Error != nil {
		return nil, result.Error
	}

	return activeWorkflows, nil
}

func (api *APIWorkflow) FetchActiveWorkflow(id int) (*ActiveWorkflowsResponse, error) {
	var activeWorkflow ActiveWorkflowsResponse

	result := api.DB.
		Table("active_workflows").
		Select("active_workflows.id, active_workflows.workflow as workflow_id, workflows.name as workflow_definition, active_workflows.current_state, active_workflows.state_deadline").
		Joins("join workflows on workflows.id = active_workflows.workflow").
		Where("active_workflows.id = ?", id).
		Scan(&activeWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return &activeWorkflow, nil
}

func (api *APIWorkflow) UpdateActiveWorkflow(id int, workflowID *int, currentState *string,  dateSubmitted *time.Time, stateDeadline *time.Time) error {
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

	// Save the updated active workflow
	result = api.DB.Save(&activeWorkflow)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
