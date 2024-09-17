package workflow

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"orchestrator/db"
	"orchestrator/env"
	"time"

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
	Fetch(id int) (*Workflow, error)
	Store(apiURL string, workflow Workflow, categories []int64, Author *int64) error
	Update(apiURL string, workflow *Workflow, categories *[]int64, author *int64) error
}

type APIWorkflow struct{
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

func (api *APIWorkflow) Fetch(id int) (*Workflow, error) {
	var workflow Workflow
	result := api.DB.First(&workflow, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &workflow, nil
}

func FetchWorkflowFromAPI(apiURL string, secretKey string) (*Workflow, error) {
	// Create a new GET request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Set the X-Orchestrator-Key header
	req.Header.Set("X-Orchestrator-Key", secretKey)

	// Perform the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch workflow: " + resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Define a temporary structure to extract the data field and ID
	var responseData struct {
		Data struct {
			ID                 int      `json:"ID"`
			WorkflowDefinition Workflow `json:"WorkflowDefinition"`
		} `json:"data"`
	}

	// Unmarshal the JSON response to extract the "data" field
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData.Data.WorkflowDefinition, nil
}

func StoreWorkflowToAPI(apiURL string, workflow Workflow, categories []int64, Author *int64) error {
	store := StoreWorkflowRequest{
		WorkflowDefinition: workflow,
		Category:           categories,
		Author:             Author,
	}
	storeJSON, err := json.Marshal(store)
	if err != nil {
		return err
	}

	// Create a new POST request with the workflow JSON as the body
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(storeJSON))
	if err != nil {
		return err
	}

	// Set the appropriate content-type header
	key, err := env.Get("ORCHESTRATOR_KEY")
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Orchestrator-Key", key)

	// Perform the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("failed to store workflow: " + resp.Status)
	}

	return nil
}

func UpdateWorkflowToAPI(apiURL string, workflow *Workflow, categories *[]int64, author *int64) error {
	// Prepare the update request structure
	var update UpdateWorkflowRequest
	update.WorkflowDefinition = workflow

	// Add categories if provided
	if categories != nil {
		update.Category = categories
	}

	// Add author if provided
	if author != nil {
		update.Author = author
	}

	// Marshal the update request object to JSON
	updateJSON, err := json.Marshal(update)
	if err != nil {
		return err
	}

	// Create a new PUT request with the update JSON as the body
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(updateJSON))
	if err != nil {
		return err
	}

	// Set the appropriate headers
	key, err := env.Get("ORCHESTRATOR_KEY")
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Orchestrator-Key", key)

	// Perform the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errors.New("failed to update workflow: " + resp.Status)
	}

	return nil
}
