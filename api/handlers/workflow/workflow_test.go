package workflow_test

import (
	"api/env"
	"api/handlers/workflow"
	"api/models"
	"errors"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

//------------------------------------------------------------------------- Mocks

//DB Model

type mockDB struct {
	throwError bool
	ReturnWorkflowArray []models.Workflow
	ReturnWorkflow *models.Workflow
	ReturnDispute *models.Dispute
	ReturnTagArray []models.Tag
	ReturnTag *models.Tag
	ReturnActiveWorkflow *models.ActiveWorkflows
}

func (m *mockDB) GetWorkflowsWithLimitOffset(limit, offset *int) ([]models.Workflow, error) {
	if m.throwError {
		return nil, errors.New("error")
	}
	return []models.Workflow{}, nil
}

func (m *mockDB) GetWorkflowByID(id uint64) (*models.Workflow, error) {
	if m.throwError {
		return nil, errors.New("error")
	}
	return &models.Workflow{}, nil
}

func (m *mockDB) GetActiveWorkflowByWorkflowID(workflowID uint64) (*models.ActiveWorkflows, error) {
	if m.throwError {
		return nil, errors.New("error")
	}
	return &models.ActiveWorkflows{}, nil
}

func (m *mockDB) QueryTagsToRelatedWorkflow(workflowID uint64) ([]models.Tag, error) {
	if m.throwError {
		return nil, errors.New("error")
	}
	return []models.Tag{}, nil
}

func (m *mockDB) FindDisputeByID(id uint64) (*models.Dispute, error) {  // Fixed typo: FindDipsuteByID
	if m.throwError {
		return nil, errors.New("error")
	}
	return &models.Dispute{}, nil
}

func (m *mockDB) CreateWorkflow(workflow *models.Workflow) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) CreateWorkflowTag(tag *models.WorkflowTags) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) CreateActiveWorkflow(workflow *models.ActiveWorkflows) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) UpdateWorkflow(workflow *models.Workflow) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) UpdateActiveWorkflow(workflow *models.ActiveWorkflows) error {  // Fixed signature
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) DeleteTagsByWorkflowID(workflowID uint64) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) DeleteWorkflow(wf *models.Workflow) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) DeleteActiveWorkflow(wf *models.ActiveWorkflows) error {  // Fixed signature
	if m.throwError {
		return errors.New("error")
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
	mockDB *mockDB
	mockJwtModel *mockJwtModel
	mockEmailModel *mockEmailModel
	mockAuditLogger *mockAuditLogger
	router *gin.Engine
}


func (suite *WorkflowTestSuite) SetupTest() {
	suite.mockDB = &mockDB{}
	suite.mockJwtModel = &mockJwtModel{}
	suite.mockEmailModel = &mockEmailModel{}
	suite.mockAuditLogger = &mockAuditLogger{}

	handler := workflow.Workflow{
		DB: suite.mockDB,
		EnvReader: env.NewEnvLoader(),
		Jwt: suite.mockJwtModel,
		Emailer: suite.mockEmailModel,
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





