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

//DB Model

type mockDB struct {
	throwError           bool
	Error                error
	ReturnWorkflowArray  []models.Workflow
	ReturnWorkflow       *models.Workflow
	ReturnDispute        *models.Dispute
	ReturnTagArray       []models.Tag
	ReturnTag            *models.Tag
	ReturnActiveWorkflow *models.ActiveWorkflows
}

func (m *mockDB) GetWorkflowsWithLimitOffset(limit, offset *int) ([]models.Workflow, error) {
	if m.throwError {
		return nil, m.Error
	}
	return []models.Workflow{}, nil
}

func (m *mockDB) GetWorkflowByID(id uint64) (*models.Workflow, error) {
	if m.throwError {
		return nil, m.Error
	}
	return &models.Workflow{}, nil
}

func (m *mockDB) GetActiveWorkflowByWorkflowID(workflowID uint64) (*models.ActiveWorkflows, error) {
	if m.throwError {
		return nil, m.Error
	}
	return &models.ActiveWorkflows{}, nil
}

func (m *mockDB) QueryTagsToRelatedWorkflow(workflowID uint64) ([]models.Tag, error) {
	if m.throwError {
		return nil, m.Error
	}
	return []models.Tag{}, nil
}

func (m *mockDB) FindDisputeByID(id uint64) (*models.Dispute, error) {
	if m.throwError {
		return nil, m.Error
	}
	return &models.Dispute{}, nil
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
		ID:                0,
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

func (m *mockEmailModel) NotifyDisputeStateChanged(c *gin.Context, disputeID int64, disputeStatus string) {
}

//------------------------------------------------------------------------- Test Suite

type WorkflowTestSuite struct {
	suite.Suite
	mockDB          *mockDB
	mockJwtModel    *mockJwtModel
	mockEmailModel  *mockEmailModel
	mockAuditLogger *mockAuditLogger
	router          *gin.Engine
}

func (suite *WorkflowTestSuite) SetupTest() {
	suite.mockDB = &mockDB{Error: errors.New("Test error")}
	suite.mockJwtModel = &mockJwtModel{}
	suite.mockEmailModel = &mockEmailModel{}
	suite.mockAuditLogger = &mockAuditLogger{}

	handler := workflow.Workflow{
		DB:                       suite.mockDB,
		EnvReader:                env.NewEnvLoader(),
		Jwt:                      suite.mockJwtModel,
		Emailer:                  suite.mockEmailModel,
		DisputeProceedingsLogger: suite.mockAuditLogger,
	}

	gin.SetMode("release")
	suite.router = gin.Default()
	suite.router.GET("", handler.GetWorkflows)
	suite.router.GET("/:id", handler.GetIndividualWorkflow)
	suite.router.POST("", handler.StoreWorkflow)
	suite.router.PUT("/:id", handler.UpdateWorkflow)
	suite.router.DELETE("/:id", handler.DeleteWorkflow)
	suite.router.POST("/active", handler.NewActiveWorkflow)
	suite.router.GET("/reset", handler.ResetActiveWorkflow)
}

func TestWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowTestSuite))
}
func (suite *WorkflowTestSuite) TestStoreWorkflow_Success() {
	// Arrange
	suite.mockDB.throwError = false
	authorID := int64(1)
	workflows := models.CreateWorkflow{
		Name:       "New Workflow",
		Author:     &authorID,
		Definition: map[string]interface{}{"key": "value"},
		Category:   []int64{1, 2},
	}

	body, _ := json.Marshal(workflows)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
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


func (suite *WorkflowTestSuite) TestUpdateWorkflow_Success() {
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

	authorID := int64(1)

	updateData := models.UpdateWorkflow{
		Name: new(string),
		WorkflowDefinition: &map[string]interface{}{
			"key": "new value",
		},
		Author:   &authorID,
		Category: &[]int64{1, 2},
	}
	*updateData.Name = "Updated Workflow"
	*updateData.Author = 2

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
		Name:     new(string),
		Category: &[]int64{1, 2},
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
	suite.mockDB.ReturnWorkflow = &models.Workflow{ID: 1, Name: "Workflow 1"}
	suite.mockDB.ReturnTagArray = []models.Tag{
		{ID: 1, TagName: "Tag 1"},
		{ID: 2, TagName: "Tag 2"},
	}

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
	suite.mockDB.ReturnWorkflowArray = []models.Workflow{
		{ID: 1, Name: "Workflow 1"},
		{ID: 2, Name: "Workflow 2"},
	}
	suite.mockDB.ReturnTagArray = []models.Tag{
		{ID: 1, TagName: "Tag 1"},
		{ID: 2, TagName: "Tag 2"},
	}

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("GET", "/?limit=10&offset=0", nil)
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
	suite.Equal("Invalid limit parameter", response.Error)
}

func (suite *WorkflowTestSuite) TestGetWorkflows_InvalidOffset() {
	// Arrange
	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("GET", "/?limit=10&offset=invalid", nil)
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
	suite.Equal("Invalid offset parameter", response.Error)
}

func (suite *WorkflowTestSuite) TestGetWorkflows_DBError() {
	// Arrange
	suite.mockDB.throwError = true

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("GET", "/?limit=10&offset=0", nil)
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

func (suite *WorkflowTestSuite) TestGetWorkflows_QueryTagsError() {
	// Arrange
	suite.mockDB.throwError = false
	suite.mockDB.ReturnWorkflowArray = []models.Workflow{
		{ID: 1, Name: "Workflow 1"},
	}
	suite.mockDB.throwError = true

	w := workflow.Workflow{
		DB: suite.mockDB,
	}

	req, _ := http.NewRequest("GET", "/?limit=10&offset=0", nil)
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
