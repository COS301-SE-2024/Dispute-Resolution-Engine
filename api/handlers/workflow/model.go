//coverage:ignore file
package workflow

import (
	"api/auditLogger"
	"api/env"
	"api/handlers/notifications"
	"api/middleware"
	"api/models"
	"api/utilities"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type WorkflowDBModel interface {
	GetWorkflowRecordByID(id uint64) (*models.Workflow, error)
	GetWorkflowsWithLimitOffset(limit, offset *int, search *string) (int64, []models.GetWorkflowResponse, error)
	GetWorkflowByID(id uint64) (*models.DetailedWorkflowResponse, error)
	GetActiveWorkflowByWorkflowID(workflowID uint64) (*models.ActiveWorkflows, error)
	QueryTagsToRelatedWorkflow(workflowID uint64) ([]models.Tag, error)
	FindDisputeByID(id uint64) (*models.Dispute, error)

	CreateWorkflow(workflow *models.Workflow) error
	CreateWorkflowTag(tag *models.WorkflowTags) error
	CreateActiveWorkflow(workflow *models.ActiveWorkflows) error

	UpdateWorkflow(workflow *models.Workflow) error
	UpdateActiveWorkflow(workflow *models.ActiveWorkflows) error
	GetActiveWorkflowByDisputeID(disputeID uint64) (*models.ActiveWorkflows, error)

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
	OrchestratorEntity       WorkflowOrchestrator
}

type WorkflowModelReal struct {
	DB  *gorm.DB
	env env.Env
}

func NewWorkflowHandler(db *gorm.DB, envReader env.Env) Workflow {
	return Workflow{
		DB:                       &WorkflowModelReal{DB: db, env: envReader},
		Emailer:                  notifications.NewHandler(db),
		EnvReader:                env.NewEnvLoader(),
		Jwt:                      middleware.NewJwtMiddleware(),
		DisputeProceedingsLogger: auditLogger.NewDisputeProceedingsLogger(db, envReader),
		OrchestratorEntity:       OrchestratorReal{},
	}
}

func (wfmr *WorkflowModelReal) GetWorkflowsWithLimitOffset(limit, offset *int, search *string) (int64, []models.GetWorkflowResponse, error) {
	var workflows []models.Workflow
	var total int64

	// Create a query object
	query := wfmr.DB.Model(&models.Workflow{})
	countQuery := wfmr.DB.Model(&models.Workflow{})

	// If search is provided, apply it (search by name)
	if search != nil {
		query = query.Where("name LIKE ?", "%"+*search+"%")
		countQuery = countQuery.Where("name LIKE ?", "%"+*search+"%")
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
	if result.Error != nil {
		return 0, nil, result.Error
	}

	resultCount := countQuery.Count(&total)

	// Handle any errors
	if resultCount.Error != nil {
		return 0, nil, result.Error
	}

	response := make([]models.GetWorkflowResponse, len(workflows))

	//read into response struct
	for i, workflow := range workflows {
		var author models.User
		result := wfmr.DB.First(&author, workflow.AuthorID)
		if result.Error != nil {
			return 0, nil, result.Error
		}

		response[i] = models.GetWorkflowResponse{
			ID:          int64(workflow.ID),
			Name:        workflow.Name,
			DateCreated: workflow.CreatedAt,
			LastUpdated: workflow.LastUpdated,
			Author: models.AuthorSum{
				ID:       author.ID,
				FullName: (author.FirstName + " " + author.Surname),
			},
		}

	}
	return total, response, nil
}

func (wfmr *WorkflowModelReal) GetWorkflowRecordByID(id uint64) (*models.Workflow, error) {
	var workflow models.Workflow

	// Create a query object
	query := wfmr.DB.Model(&models.Workflow{})
	query = query.Where("id = ?", id)

	// Execute the query
	result := query.First(&workflow)

	// Handle any errors
	if result.Error != nil {
		return nil, result.Error
	}

	return &workflow, nil
}

func (wfmr *WorkflowModelReal) GetWorkflowByID(id uint64) (*models.DetailedWorkflowResponse, error) {
	var workflow models.Workflow

	// Create a query object
	query := wfmr.DB.Model(&models.Workflow{})
	query = query.Where("id = ?", id)

	// Execute the query
	result := query.First(&workflow)

	// Handle any errors
	if result.Error != nil {
		return nil, result.Error
	}

	var author models.User
	result = wfmr.DB.First(&author, workflow.AuthorID)
	if result.Error != nil {
		return nil, result.Error
	}

	var orchestrator models.WorkflowOrchestrator
	err := json.Unmarshal([]byte(workflow.Definition), &orchestrator)
	if err != nil {
		return nil, err
	}

	//map to response model
	response := models.DetailedWorkflowResponse{
		GetWorkflowResponse: models.GetWorkflowResponse{
			ID:          int64(workflow.ID),
			Name:        workflow.Name,
			DateCreated: workflow.CreatedAt,
			LastUpdated: workflow.LastUpdated,
			Author: models.AuthorSum{
				ID:       author.ID,
				FullName: (author.FirstName + " " + author.Surname),
			},
		},
		Definition: orchestrator,
	}

	return &response, result.Error
}

func (wfmr *WorkflowModelReal) FindDisputeByID(id uint64) (*models.Dispute, error) {
	var dispute models.Dispute

	// Create a query object
	results := wfmr.DB.First(&dispute, id)
	return &dispute, results.Error
}

func (wfmr *WorkflowModelReal) QueryTagsToRelatedWorkflow(workflowID uint64) ([]models.Tag, error) {
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

func (wfmr *WorkflowModelReal) CreateWorkflow(workflow *models.Workflow) error {
	result := wfmr.DB.Create(workflow)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *WorkflowModelReal) CreateWorkflowTag(tag *models.WorkflowTags) error {
	result := wfmr.DB.Create(tag)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *WorkflowModelReal) UpdateWorkflow(workflow *models.Workflow) error {
	result := wfmr.DB.Save(workflow)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *WorkflowModelReal) DeleteTagsByWorkflowID(workflowID uint64) error {
	result := wfmr.DB.Where("workflow_id = ?", workflowID).Delete(&models.WorkflowTags{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *WorkflowModelReal) DeleteWorkflow(wf *models.Workflow) error {
	result := wfmr.DB.Delete(wf)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *WorkflowModelReal) CreateActiveWorkflow(workflow *models.ActiveWorkflows) error {
	result := wfmr.DB.Create(workflow)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *WorkflowModelReal) DeleteActiveWorkflow(workflow *models.ActiveWorkflows) error {
	result := wfmr.DB.Delete(workflow)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (wfmr *WorkflowModelReal) GetActiveWorkflowByWorkflowID(workflowID uint64) (*models.ActiveWorkflows, error) {
	var activeWorkflow models.ActiveWorkflows

	query := wfmr.DB.First(&activeWorkflow, workflowID)

	return &activeWorkflow, query.Error
}

func (wfmr *WorkflowModelReal) UpdateActiveWorkflow(workflow *models.ActiveWorkflows) error {
	result := wfmr.DB.Save(workflow)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Orchestrator for the workflow
type OrchestratorRequest struct {
	ID int64 `json:"id"`
}

type OrchestratorResetRequest struct {
	ID           int64      `json:"id"`
	CurrentState *string    `json:"current_state"`
	Deadline     *time.Time `json:"deadline"`
}

type WorkflowOrchestrator interface {
	MakeRequestToOrchestrator(endpoint string, payload OrchestratorRequest) (string, error)
	SendResetRequestToOrchestrator(endpoint string, payload OrchestratorResetRequest) (string, error)
	GetTriggers() (string, error)
}

type OrchestratorReal struct {
}

func (w OrchestratorReal) MakeRequestToOrchestrator(endpoint string, payload OrchestratorRequest) (string, error) {
	logger := utilities.NewLogger().LogWithCaller()

	// Marshal the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("marshal error: ", err)
		return "", fmt.Errorf("internal server error")
	}
	logger.Info("Payload: ", string(payloadBytes))

	// Send the POST request to the orchestrator
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Error("post error: ", err)
		return "", fmt.Errorf("internal server error")
	}
	defer resp.Body.Close()

	// Check for a successful status code (200 OK)

	if resp.StatusCode == http.StatusInternalServerError {
		logger.Error("status code error: ", resp.StatusCode)
		return "", fmt.Errorf("Check that you gave the correct state name if resetting")
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("status code error: ", resp.StatusCode)
		rsponseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("read body error: ", err)
			return "", fmt.Errorf("internal server error")
		}

		return string(rsponseBody), fmt.Errorf("internal server error")
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read body error: ", err)
		return "", fmt.Errorf("internal server error")
	}

	// Convert the response body to a string
	responseBody := string(bodyBytes)

	// Log the response body for debugging
	logger.Info("Response Body: ", responseBody)

	return responseBody, nil
}

func (w OrchestratorReal) SendResetRequestToOrchestrator(endpoint string, payload OrchestratorResetRequest) (string, error) {
	logger := utilities.NewLogger().LogWithCaller()

	// Marshal the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("marshal error: ", err)
		return "", errors.New("internal server error")
	}
	logger.Info("Payload: ", string(payloadBytes))

	// Send the POST request to the orchestrator
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Error("post error: ", err)
		return "", errors.New("failed to send request to orchestrator")
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error reading response body: ", err)
		return "", errors.New("failed to read response from orchestrator")
	}

	// If the status code is not 200, treat it as an error
	if resp.StatusCode != http.StatusOK {
		// Log the full response body for debugging purposes
		logger.Error("Received non-200 status code: ", resp.Status, "Response Body: ", string(body))

		// Return the response body as the error message

		//unmarshal the response body
		var marshalledResponse map[string]interface{}
		err := json.Unmarshal(body, &marshalledResponse)
		if err != nil {
			logger.Error("error unmarshalling response body: ", err)
			return "", fmt.Errorf("received error from orchestrator: %s", string(body))
		}
		return "", fmt.Errorf(marshalledResponse["error"].(string))
	}

	// Log the success and return the response body (if any)
	logger.Info("Orchestrator response: ", string(body))
	return string(body), nil
}


func (w OrchestratorReal) GetTriggers() (string, error) {
	logger := utilities.NewLogger().LogWithCaller()

	// Send the GET request to the orchestrator
	resp, err := http.Get("http://orchestrator:8090/triggers")
	if err != nil {
		logger.Error("get error: ", err)
		return "", fmt.Errorf("internal server error")
	}
	defer resp.Body.Close()

	// Check for a successful status code (200 OK)
	if resp.StatusCode != http.StatusOK {
		logger.Error("status code error: ", resp.StatusCode)
		return "", fmt.Errorf("internal server error")
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read body error: ", err)
		return "", fmt.Errorf("internal server error")
	}

	// Convert the response body to a string
	responseBody := string(bodyBytes)

	// Log the response body for debugging
	logger.Info("Response Body: ", responseBody)

	return responseBody, nil
}

func (wfmr *WorkflowModelReal) GetActiveWorkflowByDisputeID(disputeID uint64) (*models.ActiveWorkflows, error) {
	var activeWorkflow models.ActiveWorkflows

	// Perform a join query to get the active workflow using the workflow_id in the Dispute table
	query := wfmr.DB.
		Table("active_workflows").
		Joins("JOIN disputes ON disputes.workflow = active_workflows.id").
		Where("disputes.id = ?", disputeID).
		Select("active_workflows.*"). // Select only fields from active_workflows table
		First(&activeWorkflow)

	return &activeWorkflow, query.Error
}


