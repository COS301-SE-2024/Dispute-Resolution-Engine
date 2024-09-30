package workflow_test

import (
	"api/env"
	"api/handlers/workflow"
	"api/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

//------------------------------------------------------------------------- Mocks

//Orchestrator mock

type OrchestratorMock struct {
	throwError   bool
	Error        error
	ReturnString string
}

func (o *OrchestratorMock) MakeRequestToOrchestrator(endpoint string, payload workflow.OrchestratorRequest) (string, error) {
	if o.throwError {
		return "", o.Error
	}
	return o.ReturnString, nil
}

func (o *OrchestratorMock) SendResetRequestToOrchestrator(endpoint string, payload workflow.OrchestratorResetRequest) (string, error) {
	if o.throwError {
		return "", o.Error
	}
	return o.ReturnString, nil
}
func (o *OrchestratorMock) GetTriggers() (string, error) {
	if o.throwError {
		return "", o.Error
	}
	return o.ReturnString, nil
}

//DB Model

type mockDB struct {
	throwError             bool
	Error                  error
	ReturnWorkflowArray    []models.GetWorkflowResponse
	ReturnDetailedWorkflow *models.DetailedWorkflowResponse
	ReturnWorkflow         *models.Workflow
	ReturnDispute          *models.Dispute
	ReturnTagArray         []models.Tag
	ReturnTag              *models.Tag
	ReturnActiveWorkflow   *models.ActiveWorkflows
}

func (m *mockDB) GetActiveWorkflowByDisputeID(disputeID uint64) (*models.ActiveWorkflows, error) {
	if m.throwError {
		return nil, m.Error
	}
	return m.ReturnActiveWorkflow, nil
}

func (m *mockDB) GetWorkflowsWithLimitOffset(limit, offset *int, name *string) (int64, []models.GetWorkflowResponse, error) {
	if m.throwError {
		return 0, nil, m.Error
	}
	return int64(len(m.ReturnWorkflowArray)), m.ReturnWorkflowArray, nil
}

func (m *mockDB) GetWorkflowRecordByID(id uint64) (*models.Workflow, error) {
	if m.throwError {
		return nil, m.Error
	}
	return m.ReturnWorkflow, nil
}

func (m *mockDB) GetWorkflowByID(id uint64) (*models.DetailedWorkflowResponse, error) {
	if m.throwError {
		return nil, m.Error
	}
	return m.ReturnDetailedWorkflow, nil
}

func (m *mockDB) GetActiveWorkflowByWorkflowID(workflowID uint64) (*models.ActiveWorkflows, error) {
	if m.throwError {
		return nil, m.Error
	}
	return m.ReturnActiveWorkflow, nil
}

func (m *mockDB) QueryTagsToRelatedWorkflow(workflowID uint64) ([]models.Tag, error) {
	if m.throwError {
		return nil, m.Error
	}
	return m.ReturnTagArray, nil
}

func (m *mockDB) FindDisputeByID(id uint64) (*models.Dispute, error) {
	if m.throwError {
		return nil, m.Error
	}
	return m.ReturnDispute, nil
}

func (m *mockDB) CreateWorkflow(workflow *models.Workflow) error {
	if m.throwError {
		return m.Error
	}
	return nil
}

func (m *mockDB) CreateWorkflowTag(tag *models.WorkflowTags) error {
	if m.throwError {
		return m.Error
	}
	return nil
}

func (m *mockDB) CreateActiveWorkflow(workflow *models.ActiveWorkflows) error {
	if m.throwError {
		return m.Error
	}
	return nil
}

func (m *mockDB) UpdateWorkflow(workflow *models.Workflow) error {
	if m.throwError {
		return m.Error
	}
	return nil
}

func (m *mockDB) UpdateActiveWorkflow(workflow *models.ActiveWorkflows) error {
	if m.throwError {
		return m.Error
	}
	return nil
}

func (m *mockDB) DeleteTagsByWorkflowID(workflowID uint64) error {
	if m.throwError {
		return m.Error
	}
	return nil
}

func (m *mockDB) DeleteWorkflow(wf *models.Workflow) error {
	if m.throwError {
		return m.Error
	}
	return nil
}

func (m *mockDB) DeleteActiveWorkflow(wf *models.ActiveWorkflows) error {
	if m.throwError {
		return m.Error
	}
	return nil
}

type mockJwtModel struct {
	throwErrors bool
}
type mockEmailModel struct {
	throwErrors bool
}

type mockAuditLogger struct {
}

// mock model auditlogger
func (m *mockAuditLogger) LogDisputeProceedings(proceedingType models.EventTypes, eventData map[string]interface{}) error {
	return nil
}

// mock model dispute

func (m *mockJwtModel) GenerateJWT(user models.User) (string, error) {
	if m.throwErrors {
		return "", errors.ErrUnsupported
	}
	return "mock", nil
}
func (m *mockJwtModel) StoreJWT(email string, jwt string) error {
	if m.throwErrors {
		return errors.ErrUnsupported
	}
	return nil
}
func (m *mockJwtModel) GetJWT(email string) (string, error) {
	if m.throwErrors {
		return "", errors.ErrUnsupported
	}
	return "", nil
}
func (m *mockJwtModel) JWTMiddleware(c *gin.Context) {}

func (m *mockJwtModel) GetClaims(c *gin.Context) (models.UserInfoJWT, error) {
	if m.throwErrors {
		return models.UserInfoJWT{}, errors.ErrUnsupported
	}
	return models.UserInfoJWT{
		ID:                1,
		FirstName:         "",
		Surname:           "",
		Birthdate:         time.Now(),
		Nationality:       "",
		Role:              "",
		Email:             "",
		PhoneNumber:       new(string),
		AddressID:         new(int64),
		Status:            "",
		Gender:            "",
		PreferredLanguage: new(string),
		Timezone:          new(string),
	}, nil

}

func (m *mockEmailModel) SendAdminEmail(c *gin.Context, disputeID int64, resEmail string, title string, summary string) {
}

func (m *mockEmailModel) SendDefaultUserEmail(c *gin.Context, email string, pass string, title string, summary string) {

}

func (m *mockEmailModel) NotifyDisputeStateChanged(c *gin.Context, disputeID int64, disputeStatus, description string) {
}

func (m *mockEmailModel) NotifyEvent(c *gin.Context) {
}

// mock env reader
type mockEnvReader struct {
	throwError bool
	Error      error
}

func (m *mockEnvReader) LoadFromFile(files ...string) {
}

func (m *mockEnvReader) Register(key string) {
}

func (m *mockEnvReader) RegisterDefault(key, fallback string) {
}

func (m *mockEnvReader) Get(key string) (string, error) {
	if m.throwError {
		return "", m.Error
	}
	return "mock", nil
}

//------------------------------------------------------------------------- Test Suite

type WorkflowTestSuite struct {
	suite.Suite
	mockDB           *mockDB
	mockEnvReader    *mockEnvReader
	mockJwtModel     *mockJwtModel
	mockEmailModel   *mockEmailModel
	mockAuditLogger  *mockAuditLogger
	mockOrchestrator *OrchestratorMock
	router           *gin.Engine
}

func (suite *WorkflowTestSuite) SetupTest() {
	suite.mockDB = &mockDB{Error: errors.New("Test error")}
	suite.mockJwtModel = &mockJwtModel{}
	suite.mockEmailModel = &mockEmailModel{}
	suite.mockAuditLogger = &mockAuditLogger{}
	suite.mockOrchestrator = &OrchestratorMock{}
	suite.mockEnvReader = &mockEnvReader{}

	handler := workflow.Workflow{
		DB:                       suite.mockDB,
		EnvReader:                suite.mockEnvReader,
		Jwt:                      suite.mockJwtModel,
		Emailer:                  suite.mockEmailModel,
		DisputeProceedingsLogger: suite.mockAuditLogger,
		OrchestratorEntity:       suite.mockOrchestrator,
	}

	gin.SetMode("release")
	suite.router = gin.Default()
	suite.router.Use(handler.Jwt.JWTMiddleware)
	suite.router.GET("", handler.GetWorkflows)
	suite.router.GET("/:id", handler.GetIndividualWorkflow)
	suite.router.POST("/create", handler.StoreWorkflow)
	suite.router.PATCH("/:id", handler.UpdateWorkflow)
	suite.router.DELETE("/:id", handler.DeleteWorkflow)
	suite.router.GET("/reset", handler.ResetActiveWorkflow)
}

func TestWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowTestSuite))
}

func (suite *WorkflowTestSuite) TestValidateWorkflowDefinition_Success() {
	// Arrange
	definition := models.WorkflowOrchestrator{
		Initial: "start",
		States: map[string]models.State{
			"start": {
				Label:       "Start",
				Description: "Starting state",
				Triggers: map[string]models.Trigger{
					"next": {
						Label: "Next",
						Next:  "end",
					},
				},
			},
			"end": {
				Label:       "End",
				Description: "Ending state",
			},
		},
	}

	// Act
	err := workflow.ValidateWorkflowDefinition(definition)

	// Assert
	suite.NoError(err)
}

func (suite *WorkflowTestSuite) TestValidateWorkflowDefinition_MissingInitialState() {
	// Arrange
	definition := models.WorkflowOrchestrator{
		Initial: "start",
		States: map[string]models.State{
			"end": {
				Label:       "End",
				Description: "Ending state",
			},
		},
	}

	// Act
	err := workflow.ValidateWorkflowDefinition(definition)

	// Assert
	suite.EqualError(err, "initial state 'start' does not exist in states")
}

func (suite *WorkflowTestSuite) TestValidateWorkflowDefinition_MissingStateLabel() {
	// Arrange
	definition := models.WorkflowOrchestrator{
		Initial: "start",
		States: map[string]models.State{
			"start": {
				Description: "Starting state",
			},
		},
	}

	// Act
	err := workflow.ValidateWorkflowDefinition(definition)

	// Assert
	suite.EqualError(err, "state 'start' is missing a label")
}

func (suite *WorkflowTestSuite) TestValidateWorkflowDefinition_MissingStateDescription() {
	// Arrange
	definition := models.WorkflowOrchestrator{
		Initial: "start",
		States: map[string]models.State{
			"start": {
				Label: "Start",
			},
		},
	}

	// Act
	err := workflow.ValidateWorkflowDefinition(definition)

	// Assert
	suite.EqualError(err, "state 'start' is missing a description")
}

func (suite *WorkflowTestSuite) TestValidateWorkflowDefinition_MissingTriggerLabel() {
	// Arrange
	definition := models.WorkflowOrchestrator{
		Initial: "start",
		States: map[string]models.State{
			"start": {
				Label:       "Start",
				Description: "Starting state",
				Triggers: map[string]models.Trigger{
					"next": {
						Next: "end",
					},
				},
			},
			"end": {
				Label:       "End",
				Description: "Ending state",
			},
		},
	}

	// Act
	err := workflow.ValidateWorkflowDefinition(definition)

	// Assert
	suite.EqualError(err, "trigger 'next' in state 'start' is missing a label")
}

func (suite *WorkflowTestSuite) TestValidateWorkflowDefinition_NonExistentNextState() {
	// Arrange
	definition := models.WorkflowOrchestrator{
		Initial: "start",
		States: map[string]models.State{
			"start": {
				Label:       "Start",
				Description: "Starting state",
				Triggers: map[string]models.Trigger{
					"next": {
						Label: "Next",
						Next:  "nonexistent",
					},
				},
			},
		},
	}

	// Act
	err := workflow.ValidateWorkflowDefinition(definition)

	// Assert
	suite.EqualError(err, "trigger 'next' in state 'start' points to a non-existent state 'nonexistent'")
}

func (suite *WorkflowTestSuite) TestValidateWorkflowDefinition_InvalidTimerDuration() {
	// Arrange
	definition := models.WorkflowOrchestrator{
		Initial: "start",
		States: map[string]models.State{
			"start": {
				Label:       "Start",
				Description: "Starting state",
				Timer: &models.Timer{
					Duration: models.DurationWrapper{Duration: 0},
					OnExpire: "expire",
				},
			},
		},
	}

	// Act
	err := workflow.ValidateWorkflowDefinition(definition)

	// Assert
	suite.EqualError(err, "timer in state 'start' must have a non-zero duration")
}

func (suite *WorkflowTestSuite) TestValidateWorkflowDefinition_NonExistentTimerTrigger() {
	// Arrange
	definition := models.WorkflowOrchestrator{
		Initial: "start",
		States: map[string]models.State{
			"start": {
				Label:       "Start",
				Description: "Starting state",
				Timer: &models.Timer{
					Duration: models.DurationWrapper{Duration: 10},
					OnExpire: "nonexistent",
				},
			},
		},
	}

	// Act
	err := workflow.ValidateWorkflowDefinition(definition)

	// Assert
	suite.EqualError(err, "timer in state 'start' points to a non-existent trigger 'nonexistent'")
}

func (suite *WorkflowTestSuite) TestResetActiveWorkflow_Success() {

	// Arrange
	suite.mockDB.throwError = false
	suite.mockOrchestrator.throwError = false
	suite.mockDB.ReturnActiveWorkflow = &models.ActiveWorkflows{
		ID:           1,
		CurrentState: "initial",
	}

	resetRequest := models.ResetActiveWorkflow{
		DisputeID:    new(int64),
		CurrentState: new(string),
		Deadline:     new(time.Time),
	}
	*resetRequest.DisputeID = 1
	*resetRequest.CurrentState = "new_state"
	*resetRequest.Deadline = time.Now().Add(24 * time.Hour)

	body, _ := json.Marshal(resetRequest)
	req, _ := http.NewRequest("POST", "/reset", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	w := workflow.Workflow{
		DB:                 suite.mockDB,
		EnvReader:          suite.mockEnvReader,
		OrchestratorEntity: suite.mockOrchestrator,
	}

	// Act
	w.ResetActiveWorkflow(c)
	//print body of request

	// Assert
	suite.Equal(http.StatusOK, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Database updated and request to Reset workflow sent", response.Data)
}

func (suite *WorkflowTestSuite) TestResetActiveWorkflow_InvalidPayload() {
	// Arrange
	req, _ := http.NewRequest("POST", "/reset", bytes.NewBuffer([]byte("invalid payload")))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	// Act
	w.ResetActiveWorkflow(c)

	// Assert
	suite.Equal(http.StatusBadRequest, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Invalid request payload", response.Error)
}

func (suite *WorkflowTestSuite) TestResetActiveWorkflow_MissingFields() {
	// Arrange
	resetRequest := models.ResetActiveWorkflow{}
	body, _ := json.Marshal(resetRequest)
	req, _ := http.NewRequest("POST", "/reset", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	// Act
	w.ResetActiveWorkflow(c)

	// Assert
	suite.Equal(http.StatusBadRequest, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Missing required fields", response.Error)
}

func (suite *WorkflowTestSuite) TestResetActiveWorkflow_NotFound() {
	// Arrange
	suite.mockDB.throwError = true
	suite.mockDB.Error = gorm.ErrRecordNotFound

	resetRequest := models.ResetActiveWorkflow{
		DisputeID:    new(int64),
		CurrentState: new(string),
		Deadline:     new(time.Time),
	}
	*resetRequest.DisputeID = 1
	*resetRequest.CurrentState = "new_state"
	*resetRequest.Deadline = time.Now().Add(24 * time.Hour)

	body, _ := json.Marshal(resetRequest)
	req, _ := http.NewRequest("POST", "/reset", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	w := workflow.Workflow{
		DB: suite.mockDB,
		EnvReader: 		suite.mockEnvReader,
		OrchestratorEntity: suite.mockOrchestrator,
	}

	// Act
	w.ResetActiveWorkflow(c)

	// Assert
	suite.Equal(http.StatusNotFound, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Active workflow not found", response.Error)
}

func (suite *WorkflowTestSuite) TestResetActiveWorkflow_DBError() {
	// Arrange
	suite.mockDB.throwError = true
	suite.mockOrchestrator.throwError = true
	suite.mockEnvReader.throwError = true
	suite.mockEnvReader.Error = errors.New("Internal Server Error")

	resetRequest := models.ResetActiveWorkflow{
		DisputeID:    new(int64),
		CurrentState: new(string),
		Deadline:     new(time.Time),
	}
	*resetRequest.DisputeID = 1
	*resetRequest.CurrentState = "new_state"
	*resetRequest.Deadline = time.Now().Add(24 * time.Hour)

	body, _ := json.Marshal(resetRequest)
	req, _ := http.NewRequest("POST", "/reset", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	w := workflow.Workflow{
		DB: suite.mockDB,
		OrchestratorEntity: suite.mockOrchestrator,
		EnvReader: 			suite.mockEnvReader,
	}

	// Act
	w.ResetActiveWorkflow(c)

	// Assert
	suite.Equal(http.StatusInternalServerError, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Internal Server Error", response.Error)
}

func (suite *WorkflowTestSuite) TestResetActiveWorkflow_OrchestratorError() {
	// Arrange
	suite.mockDB.throwError = false
	suite.mockDB.ReturnActiveWorkflow = &models.ActiveWorkflows{
		ID:           1,
		CurrentState: "initial",
	}
	suite.mockOrchestrator.throwError = true
	suite.mockOrchestrator.Error = errors.New("Orchestrator error")

	resetRequest := models.ResetActiveWorkflow{
		DisputeID:    new(int64),
		CurrentState: new(string),
		Deadline:     new(time.Time),
	}
	*resetRequest.DisputeID = 1
	*resetRequest.CurrentState = "new_state"
	*resetRequest.Deadline = time.Now().Add(24 * time.Hour)

	body, _ := json.Marshal(resetRequest)
	req, _ := http.NewRequest("POST", "/reset", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	w := workflow.Workflow{
		DB:                 suite.mockDB,
		EnvReader:          env.NewEnvLoader(),
		OrchestratorEntity: suite.mockOrchestrator,
	}

	// Act
	w.ResetActiveWorkflow(c)

	// Assert
	suite.Equal(http.StatusInternalServerError, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Internal Server Error", response.Error)
}
func (suite *WorkflowTestSuite) TestDeleteWorkflow_Success() {
	// Arrange
	suite.mockDB.throwError = false
	suite.mockDB.ReturnWorkflow = &models.Workflow{ID: 1, Name: "Workflow 1"}

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("DELETE", "/1", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.DeleteWorkflow(c)

	// Assert
	suite.Equal(http.StatusOK, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Workflow and associated tags deleted", response.Data)
}

func (suite *WorkflowTestSuite) TestDeleteWorkflow_InvalidID() {
	// Arrange
	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("DELETE", "/invalid", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	// Act
	w.DeleteWorkflow(c)

	// Assert
	suite.Equal(http.StatusBadRequest, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Invalid ID parameter", response.Error)
}

func (suite *WorkflowTestSuite) TestDeleteWorkflow_NotFound() {
	// Arrange
	suite.mockDB.throwError = true
	suite.mockDB.Error = gorm.ErrRecordNotFound

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("DELETE", "/1", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.DeleteWorkflow(c)

	// Assert
	suite.Equal(http.StatusNotFound, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Workflow not found", response.Error)
}

func (suite *WorkflowTestSuite) TestDeleteWorkflow_DBError() {
	// Arrange
	suite.mockDB.throwError = true

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("DELETE", "/1", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.DeleteWorkflow(c)

	// Assert
	suite.Equal(http.StatusInternalServerError, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Internal Server Error", response.Error)
}

func (suite *WorkflowTestSuite) TestDeleteWorkflow_FailedToDeleteTags() {
	// Arrange
	suite.mockDB.throwError = false
	suite.mockDB.ReturnWorkflow = &models.Workflow{ID: 1, Name: "Workflow 1"}

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	// Simulate error when deleting tags
	suite.mockDB.throwError = true

	req, _ := http.NewRequest("DELETE", "/1", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.DeleteWorkflow(c)

	// Assert
	suite.Equal(http.StatusInternalServerError, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Internal Server Error", response.Error)
}

func (suite *WorkflowTestSuite) TestDeleteWorkflow_FailedToDeleteWorkflow() {
	// Arrange
	suite.mockDB.throwError = false
	suite.mockDB.ReturnWorkflow = &models.Workflow{ID: 1, Name: "Workflow 1"}

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	// Simulate error when deleting workflow
	suite.mockDB.throwError = true

	req, _ := http.NewRequest("DELETE", "/1", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.DeleteWorkflow(c)

	// Assert
	suite.Equal(http.StatusInternalServerError, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Internal Server Error", response.Error)
}

func (suite *WorkflowTestSuite) TestStoreWorkflow_Success() {
	// Arrange
	suite.mockDB.throwError = false
	suite.mockJwtModel.throwErrors = false
	// authorID := int64(1)
	workflows := models.CreateWorkflow{
		Name: "New Workflow",
		// Author:     &authorID,
		Definition: models.WorkflowOrchestrator{
			Initial: "initial",
			States: map[string]models.State{
				"initial": {
					Label:       "State 1",
					Description: "Description 1",
				},
			},
		},
		// Category:   []int64{1, 2},
	}

	body, _ := json.Marshal(workflows)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	w := workflow.Workflow{
		DB:  suite.mockDB,
		Jwt: suite.mockJwtModel,
	}

	// Act
	w.StoreWorkflow(c)

	// Assert
	suite.Equal(http.StatusOK, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.NotNil(response.Data)
}

func (suite *WorkflowTestSuite) TestStoreWorkflow_InvalidPayload() {
	// Arrange
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer([]byte("invalid payload")))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	// Act
	w.StoreWorkflow(c)

	// Assert
	suite.Equal(http.StatusBadRequest, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Invalid request payload", response.Error)
}

func (suite *WorkflowTestSuite) TestStoreWorkflow_MissingFields() {
	// Arrange
	workflows := models.CreateWorkflow{}

	body, _ := json.Marshal(workflows)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	w := workflow.Workflow{
		DB:  suite.mockDB,
		Jwt: suite.mockJwtModel,
	}

	// Act
	w.StoreWorkflow(c)

	// Assert
	suite.Equal(http.StatusBadRequest, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("initial state '' does not exist in states", response.Error)
}

func (suite *WorkflowTestSuite) TestStoreWorkflow_DBError() {
	// Arrange
	suite.mockDB.throwError = true
	// authorID := int64(1)
	workflows := models.CreateWorkflow{
		Name: "New Workflow",
		// Author:     &authorID,
		Definition: models.WorkflowOrchestrator{
			Initial: "initial",
			States: map[string]models.State{
				"initial": {
					Label:       "State 1",
					Description: "Description 1",
				},
			},
		},
		// Category:   []int64{1, 2},
	}

	body, _ := json.Marshal(workflows)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	w := workflow.Workflow{
		DB:  suite.mockDB,
		Jwt: suite.mockJwtModel,
	}

	// Act
	w.StoreWorkflow(c)

	// Assert
	suite.Equal(http.StatusInternalServerError, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Internal Server Error", response.Error)
}

func (suite *WorkflowTestSuite) TestUpdateWorkflow_Success() {
	// Arrange
	suite.mockDB.throwError = false
	suite.mockDB.ReturnWorkflow = &models.Workflow{ID: 1, Name: "Old Workflow"}

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	// authorID := int64(1)
	updated := models.WorkflowOrchestrator{
		Initial: "initial",
		States: map[string]models.State{
			"initial": {
				Label:       "State 1",
				Description: "Description 1",
			},
		},
	}

	updateData := models.UpdateWorkflow{
		Name:               new(string),
		WorkflowDefinition: &updated,
		// Author:   &authorID,
		// Category: &[]int64{1, 2},
	}
	*updateData.Name = "Updated Workflow"
	// *updateData.Author = 2

	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PATCH", "/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.UpdateWorkflow(c)

	// Assert
	suite.Equal(http.StatusOK, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Workflow updated", response.Data)
}

func (suite *WorkflowTestSuite) TestUpdateWorkflow_InvalidID() {
	// Arrange
	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("PUT", "/invalid", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	// Act
	w.UpdateWorkflow(c)

	// Assert
	suite.Equal(http.StatusBadRequest, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Invalid ID parameter", response.Error)
}

func (suite *WorkflowTestSuite) TestUpdateWorkflow_NotFound() {
	// Arrange
	suite.mockDB.throwError = true
	suite.mockDB.Error = gorm.ErrRecordNotFound

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	updateData := models.UpdateWorkflow{
		Name: new(string),
	}
	*updateData.Name = "Updated Workflow"

	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.UpdateWorkflow(c)

	// Assert
	suite.Equal(http.StatusNotFound, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Workflow not found", response.Error)
}

func (suite *WorkflowTestSuite) TestUpdateWorkflow_InvalidPayload() {
	// Arrange
	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("PUT", "/1", bytes.NewBuffer([]byte("invalid payload")))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.UpdateWorkflow(c)

	// Assert
	suite.Equal(http.StatusBadRequest, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Invalid request payload", response.Error)
}

func (suite *WorkflowTestSuite) TestUpdateWorkflow_DBError() {
	// Arrange
	suite.mockDB.throwError = true

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	updateData := models.UpdateWorkflow{
		Name: new(string),
	}
	*updateData.Name = "Updated Workflow"

	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.UpdateWorkflow(c)

	// Assert
	suite.Equal(http.StatusInternalServerError, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Internal Server Error", response.Error)
}

func (suite *WorkflowTestSuite) TestUpdateWorkflow_FailedToUpdateCategories() {
	// Arrange
	suite.mockDB.throwError = false
	suite.mockDB.ReturnWorkflow = &models.Workflow{ID: 1, Name: "Old Workflow"}
	suite.mockDB.ReturnTagArray = []models.Tag{
		{ID: 1, TagName: "Tag 1"},
		{ID: 2, TagName: "Tag 2"},
	}

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	updateData := models.UpdateWorkflow{
		Name: new(string),
		// Category: &[]int64{1, 2},
	}
	*updateData.Name = "Updated Workflow"

	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Simulate error when updating categories
	suite.mockDB.throwError = true

	// Act
	w.UpdateWorkflow(c)

	// Assert
	suite.Equal(http.StatusInternalServerError, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Internal Server Error", response.Error)
}

func (suite *WorkflowTestSuite) TestGetIndividualWorkflow_Success() {
	// Arrange
	suite.mockDB.throwError = false
	suite.mockDB.ReturnDetailedWorkflow = &models.DetailedWorkflowResponse{GetWorkflowResponse: models.GetWorkflowResponse{ID: 1, Name: "Workflow 1"}}

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("GET", "/1", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.GetIndividualWorkflow(c)

	// Assert
	suite.Equal(http.StatusOK, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.NotNil(response.Data)
}

func (suite *WorkflowTestSuite) TestGetIndividualWorkflow_InvalidID() {
	// Arrange
	w := workflow.Workflow{
		DB: suite.mockDB,
	}
	suite.mockDB.ReturnDetailedWorkflow = &models.DetailedWorkflowResponse{GetWorkflowResponse: models.GetWorkflowResponse{ID: 1, Name: "Workflow 1"}}
	req, _ := http.NewRequest("GET", "/invalid", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	// Act
	w.GetIndividualWorkflow(c)

	// Assert
	suite.Equal(http.StatusBadRequest, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Invalid ID parameter", response.Error)
}

func (suite *WorkflowTestSuite) TestGetIndividualWorkflow_NotFound() {
	// Arrange
	suite.mockDB.throwError = true
	suite.mockDB.Error = gorm.ErrRecordNotFound

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("GET", "/1", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.GetIndividualWorkflow(c)

	// Assert
	suite.Equal(http.StatusOK, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Record not found", response.Error)
}

func (suite *WorkflowTestSuite) TestGetIndividualWorkflow_DBError() {
	// Arrange
	suite.mockDB.throwError = true

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("GET", "/1", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.GetIndividualWorkflow(c)

	// Assert
	suite.Equal(http.StatusInternalServerError, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Internal Server Error", response.Error)
}

func (suite *WorkflowTestSuite) TestGetIndividualWorkflow_QueryTagsError() {
	// Arrange
	suite.mockDB.throwError = false
	suite.mockDB.ReturnWorkflow = &models.Workflow{ID: 1, Name: "Workflow 1"}
	suite.mockDB.throwError = true

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("GET", "/1", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Act
	w.GetIndividualWorkflow(c)

	// Assert
	suite.Equal(http.StatusInternalServerError, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Internal Server Error", response.Error)
}

func (suite *WorkflowTestSuite) TestGetWorkflows_Success() {
	// Arrange
	suite.mockDB.throwError = false
	suite.mockDB.ReturnWorkflowArray = []models.GetWorkflowResponse{
		{ID: 1, Name: "Workflow 1"},
		{ID: 2, Name: "Workflow 2"},
	}

	w := workflow.Workflow{
		DB: suite.mockDB,
	}
	//empty json body
	body, _ := json.Marshal(models.GetWorkflow{})

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	// Act
	w.GetWorkflows(c)

	// Assert
	suite.Equal(http.StatusOK, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.NotNil(response.Data)
}

func (suite *WorkflowTestSuite) TestGetWorkflows_InvalidLimit() {
	// Arrange
	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("GET", "/?limit=invalid&offset=0", nil)
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	// Act
	w.GetWorkflows(c)

	// Assert
	suite.Equal(http.StatusBadRequest, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Invalid request payload", response.Error)
}

func (suite *WorkflowTestSuite) TestGetWorkflows_InvalidOffset() {
	// Arrange
	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	//body with invalid offset
	body, _ := json.Marshal(map[string]interface{}{"offset": "invalid"})

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	// Act
	w.GetWorkflows(c)

	// Assert
	suite.Equal(http.StatusBadRequest, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Invalid request payload", response.Error)
}

func (suite *WorkflowTestSuite) TestGetWorkflows_DBError() {
	// Arrange
	suite.mockDB.throwError = true

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	//empty json body
	body, _ := json.Marshal(models.GetWorkflow{})

	req, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	wr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	// Act
	w.GetWorkflows(c)

	// Assert
	suite.Equal(http.StatusInternalServerError, wr.Code)
	var response models.Response
	err := json.Unmarshal(wr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Internal Server Error", response.Error)
}
