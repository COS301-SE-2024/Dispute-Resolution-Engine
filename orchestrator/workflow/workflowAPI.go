package workflow

import (
	"encoding/json"
	"fmt"
	"orchestrator/db"

	"gorm.io/gorm"
)

// ----------------------------API--------------------------------

type StoreWorkflowRequest struct {
	WorkflowDefinition Workflow `json:"workflow_definition,omitempty"`
	Category           []int64  `json:"category,omitempty"`
	Author             *int64   `json:"author,omitempty"`
}

type UpdateWorkflowRequest struct {
	WorkflowDefinition *Workflow `json:"workflow_definition,omitempty"`
	Category           *[]int64  `json:"category,omitempty"`
	Author             *int64    `json:"author,omitempty"`
}

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

func (api *APIWorkflow) Fetch(id int) (*db.Workflowdb, error) {
	var workflow db.Workflowdb
	result := api.DB.First(&workflow, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &workflow, nil
}

func (api *APIWorkflow) Store(name string, workflow Workflow, categories []int64, Author int64) error {
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

func (api *APIWorkflow) Update(id int, name *string, workflow *Workflow, categories *[]int64, author *int64) error {
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
