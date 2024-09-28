package dispute_test

import (
	"api/handlers/dispute"
	"api/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type mockEvidence struct {
	user    int64
	dispute int64
	path    string
	data    string
}

type mockDisputeModel struct {
	throwErrors bool
	evidence    []mockEvidence
	Get_Experts []models.AdminDisputeExperts
}
type mockJwtModel struct {
	throwErrors bool
	returnUser  models.UserInfoJWT
}
type mockEmailModel struct {
	throwErrors bool
}

type mockAuditLogger struct {
}

type mockTicketModel struct {
	throwErrors bool
	Error       error
}

type DisputeErrorTestSuite struct {
	suite.Suite
	disputeMock      *mockDisputeModel
	jwtMock          *mockJwtModel
	emailMock        *mockEmailModel
	router           *gin.Engine
	auditMock        *mockAuditLogger
	mockOrchestrator *mockOrchestrator
	mockEnv          *mockEnv
	mockTicket       *mockTicketModel
}

func (suite *DisputeErrorTestSuite) SetupTest() {
	suite.disputeMock = &mockDisputeModel{}
	suite.jwtMock = &mockJwtModel{}
	suite.emailMock = &mockEmailModel{}
	suite.auditMock = &mockAuditLogger{}
	suite.mockOrchestrator = &mockOrchestrator{}
	suite.mockEnv = &mockEnv{}
	suite.mockTicket = &mockTicketModel{}

	handler := dispute.Dispute{Model: suite.disputeMock, JWT: suite.jwtMock, Email: suite.emailMock, AuditLogger: suite.auditMock, OrchestratorEntity: suite.mockOrchestrator, Env: suite.mockEnv, TicketModel: suite.mockTicket}
	gin.SetMode("release")
	router := gin.Default()
	router.Use(suite.jwtMock.JWTMiddleware)

	router.GET("/disputes", handler.GetSummaryListOfDisputes)
	router.POST("/create", handler.CreateDispute)
	router.GET("/:id", handler.GetDispute)

	router.POST("/disputes", handler.GetSummaryListOfDisputes)
	router.POST("/:id/objections", handler.ExpertObjection)
	router.PATCH("/objections/:id", handler.ExpertObjectionsReview)
	router.POST("/experts/objections", handler.ViewExpertRejections)
	router.POST("/:id/evidence", handler.UploadEvidence)
	router.POST("/:id/decision", handler.SubmitWriteup)
	router.PUT("/:id/status", handler.UpdateStatus)

	suite.router = router
}
func TestDisputeErrors(t *testing.T) {
	suite.Run(t, new(DisputeErrorTestSuite))
}

// ---------------------------------------------------------------- UTILITY FUNCTIONS

// Creates and writes to a form field with the specified value
func createStringField(w *multipart.Writer, key string, value string) {
	field, _ := w.CreateFormField(key)
	field.Write([]byte(value))
}

// Creates and writes to a form file field with the specified value
func createFileField(w *multipart.Writer, field, filename, value string) {
	file, _ := w.CreateFormFile(field, filename)
	file.Write([]byte(value))
}

// ---------------------------------------------------------------- MODEL MOCKS

//ticket mock

func (m *mockTicketModel) GetAdminTicketList(searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error) {
	if m.throwErrors {
		return nil, 0, m.Error
	}
	return nil, 0, nil
}

func (m *mockTicketModel) GetTicketsByUserID(uid int64, searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error) {
	if m.throwErrors {
		return nil, 0, m.Error
	}
	return nil, 0, nil
}

func (m *mockTicketModel) GetTicketDetails(ticketID int64, userID int64) (models.TicketsByUser, error) {
	if m.throwErrors {
		return models.TicketsByUser{}, m.Error
	}
	return models.TicketsByUser{}, nil
}

func (m *mockTicketModel) GetAdminTicketDetails(ticketID int64) (models.TicketsByUser, error) {
	if m.throwErrors {
		return models.TicketsByUser{}, m.Error
	}
	return models.TicketsByUser{}, nil
}

func (m *mockTicketModel) PatchTicketStatus(status string, ticketID int64) error {
	if m.throwErrors {
		return m.Error
	}
	return nil
}

func (m *mockTicketModel) AddUserTicketMessage(ticketID int64, userID int64, message string) (models.TicketMessage, error) {
	if m.throwErrors {
		return models.TicketMessage{}, m.Error
	}
	return models.TicketMessage{}, nil
}

func (m *mockTicketModel) AddAdminTicketMessage(ticketID int64, userID int64, message string) (models.TicketMessage, error) {
	if m.throwErrors {
		return models.TicketMessage{}, m.Error
	}
	return models.TicketMessage{}, nil
}

func (m *mockTicketModel) CreateTicket(userID int64, dispute int64, subject string, message string) (models.Ticket, error) {
	if m.throwErrors {
		return models.Ticket{}, m.Error
	}
	return models.Ticket{}, nil
}

//mock env

type mockEnv struct {
	throwErrors bool
	Error       error
}

func (m *mockEnv) LoadFromFile(files ...string) {
}

func (m *mockEnv) Register(key string) {
}

func (m *mockEnv) RegisterDefault(key, fallback string) {
}

func (m *mockEnv) Get(key string) (string, error) {
	if m.throwErrors {
		return "", m.Error
	}
	return "", nil
}

// mock orchestrator
type mockOrchestrator struct {
	throwErrors bool
	Error       error
}

func (m *mockOrchestrator) MakeRequestToOrchestrator(endpoint string, payload dispute.OrchestratorRequest) (string, error) {
	if m.throwErrors {
		return "", m.Error
	}
	return "", nil
}

// mock model auditlogger
func (m *mockAuditLogger) LogDisputeProceedings(proceedingType models.EventTypes, eventData map[string]interface{}) error {
	return nil
}

func (m *mockDisputeModel) GetWorkflowRecordByID(id uint64) (*models.Workflow, error) {
	if m.throwErrors {
		return nil, errors.ErrUnsupported
	}
	return &models.Workflow{}, nil

}

func (m *mockDisputeModel) CreateActiverWorkflow(workflow *models.ActiveWorkflows) error {
	if m.throwErrors {
		return errors.ErrUnsupported
	}
	return nil
}

func (m *mockDisputeModel) DeleteActiveWorkflow(workflow *models.ActiveWorkflows) error {
	if m.throwErrors {
		return errors.ErrUnsupported
	}
	return nil
}

// mock model dispute
func (m *mockDisputeModel) UploadEvidence(userId, disputeId int64, path string, file io.Reader) (uint, error) {
	if m.throwErrors {
		return 0, errors.ErrUnsupported
	}

	data, _ := io.ReadAll(file)
	m.evidence = append(m.evidence, mockEvidence{
		user:    userId,
		dispute: disputeId,
		path:    path,
		data:    string(data),
	})
	return 0, nil
}
func (m *mockDisputeModel) GetEvidenceByDispute(disputeId int64) ([]models.Evidence, error) {
	if m.throwErrors {
		return nil, errors.ErrUnsupported
	}
	return []models.Evidence{}, nil
}
func (m *mockDisputeModel) GetDisputeExperts(disputeId int64) ([]models.Expert, error) {
	if m.throwErrors {
		return nil, errors.ErrUnsupported
	}
	return []models.Expert{}, nil
}
func (m *mockDisputeModel) GetDisputesByUser(userId int64) ([]models.Dispute, error) {
	if m.throwErrors {
		return nil, errors.ErrUnsupported
	}
	return []models.Dispute{}, nil
}

func (m *mockDisputeModel) GetAdminDisputes(searchTerm *string, limit *int, offset *int, sort *models.Sort, filters *[]models.Filter, dateFilter *models.DateFilter) ([]models.AdminDisputeSummariesResponse, int64, error) {
	if m.throwErrors {
		return nil, 0, errors.ErrUnsupported
	}
	return []models.AdminDisputeSummariesResponse{}, 0, nil
}

func (m *mockDisputeModel) GetDispute(disputeId int64) (models.Dispute, error) {
	if m.throwErrors {
		return models.Dispute{}, errors.ErrUnsupported
	}
	return models.Dispute{
		ID:         new(int64),
		Workflow:   new(int64),
		Respondant: new(int64),
	}, nil
}

func (m *mockDisputeModel) GetAdminDisputeDetails(disputeId int64) (models.AdminDisputeDetailsResponse, error) {
	if m.throwErrors {
		return models.AdminDisputeDetailsResponse{}, errors.ErrUnsupported
	}
	return models.AdminDisputeDetailsResponse{}, nil
}

func (m *mockDisputeModel) GetUser(userID int64) (models.UserDetails, error) {
	if m.throwErrors {
		return models.UserDetails{}, errors.ErrUnsupported
	}
	return models.UserDetails{
		FullName: "name",
		Email:    "email",
		Address:  "address",
	}, nil
}

func (m *mockDisputeModel) GetUserByEmail(email string) (models.User, error) {
	if m.throwErrors {
		return models.User{}, errors.ErrUnsupported
	}
	return models.User{
		PhoneNumber:       new(string),
		AddressID:         new(int64),
		LastLogin:         new(time.Time),
		PreferredLanguage: new(string),
		Timezone:          new(string),
	}, nil
}

func (m *mockDisputeModel) GetExpertRejections(expertID, disputeID *int64, limit, offset *int) ([]models.ExpertObjectionsView, error) {
	if m.throwErrors {
		return nil, errors.ErrUnsupported
	}
	return []models.ExpertObjectionsView{
		{
			ObjectionID:     1,
			ExpertID:        1,
			ExpertFullName:  "name",
			DisputeID:       1,
			DisputeTitle:    "title",
			UserID:          1,
			UserFullName:    "name",
			ObjectionStatus: "status",
		},
	}, nil
}

func (m *mockDisputeModel) CreateDispute(dispute models.Dispute) (int64, error) {
	if m.throwErrors {
		return 0, errors.ErrUnsupported
	}
	return 0, nil
}
func (m *mockDisputeModel) UpdateDisputeStatus(disputeId int64, status string) error {
	if m.throwErrors {
		return errors.ErrUnsupported
	}
	return nil
}
func (m *mockDisputeModel) ObjectExpert(disputeId, expertId, ticketId int64) error {
	if m.throwErrors {
		return errors.ErrUnsupported
	}
	return nil
}
func (m *mockDisputeModel) ReviewExpertObjection(expertId int64, approved models.ExpObjStatus) error {
	if m.throwErrors {
		return errors.ErrUnsupported
	}
	return nil
}

func (m *mockDisputeModel) GetUserById(userId int64) (models.User, error) {
	if m.throwErrors {
		return models.User{}, errors.ErrUnsupported
	}
	return models.User{
		PhoneNumber:       new(string),
		AddressID:         new(int64),
		LastLogin:         new(time.Time),
		PreferredLanguage: new(string),
		Timezone:          new(string),
	}, nil
}

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
func (m *mockJwtModel) JWTMiddleware(c *gin.Context) {
	if m.throwErrors {
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func (m *mockJwtModel) GetClaims(c *gin.Context) (models.UserInfoJWT, error) {
	if m.throwErrors {
		return models.UserInfoJWT{}, errors.ErrUnsupported
	}
	return m.returnUser, nil
}

func (m *mockDisputeModel) UploadWriteup(userId, disputeId int64, path string, file io.Reader) error {
	return nil
}
func (m *mockDisputeModel) GenerateAISummary(disputeID int64, disputeDesc string, apiKey string) {

}

func (m *mockDisputeModel) GetExperts(disputeID int64) ([]models.AdminDisputeExperts, error) {
	if m.throwErrors {
		return nil, errors.ErrUnsupported
	}
	return m.Get_Experts, nil
}

func (m *mockDisputeModel) AssignExpertsToDispute(disputeID int64) ([]models.User, error) {
	return nil, nil
}

func (m *mockDisputeModel) CreateDefaultUser(email string, fullName string, pass string) error {
	return nil
}

func (m *mockEmailModel) SendAdminEmail(c *gin.Context, disputeID int64, resEmail string, title string, summary string) {
}

func (m *mockEmailModel) SendDefaultUserEmail(c *gin.Context, email string, pass string, title string, summary string) {

}

func (m *mockEmailModel) NotifyDisputeStateChanged(c *gin.Context, disputeID int64, disputeStatus string) {
}

func (m *mockEmailModel) NotifyEvent(c *gin.Context) {
}

// ---------------------------------------------------------------- Get Summary List Tests

func (suite *DisputeErrorTestSuite) TestGetSummaryListUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("GET", "/disputes", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Unauthorized", result.Error)
}

func (suite *DisputeErrorTestSuite) TestGetSummaryListNoDisputes() {
	req, _ := http.NewRequest("GET", "/disputes", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data []models.DisputeSummaryResponse `json:"data"`
	}
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Empty(result.Data)
}

func (suite *DisputeErrorTestSuite) TestGetSummaryListWithDisputes() {
	// Mock disputes
	suite.disputeMock.evidence = []mockEvidence{
		{
			user:    1,
			dispute: 1,
			path:    "path/to/file",
			data:    "evidence data",
		},
	}
}

// 	req, _ := http.NewRequest("GET", "/summary", nil)
// 	req.Header.Add("Authorization", "Bearer mock")

// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)

// 	var result models.Response
// 	suite.Equal(http.StatusOK, w.Code)
// 	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
// 	suite.Empty(result.Error)

// 	// Check if the response has the correct number of disputes
// 	suite.NotEmpty(result.Data)
// 	suite.Equal(len(suite.disputeMock.evidence), len(result.Data.([]models.DisputeSummaryResponse)))
// }

func (suite *DisputeErrorTestSuite) TestGetSummaryListAdminSuccess() {
	suite.jwtMock.returnUser.Role = "admin"
	body := `{}`
	req, _ := http.NewRequest("POST", "/disputes", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

func (suite *DisputeErrorTestSuite) TestGetSummaryListAdminSuccess2() {
	suite.jwtMock.returnUser.Role = "admin"
	body := `{
  "limit": 50,
  "offset": 0,
  "sort": {
    "attr": "case_date",
    "order": "desc"
  },
  "filter": [
    {
      "attr": "workflow",
      "value": "1"
    }
  ],
  "dateFilter": {
    "filed": {
      "before": "2024-09-20",
      "after": "2024-01-01"
    },
    "resolved": {
      "before": "2024-09-20"
    }
  }
}`
	req, _ := http.NewRequest("POST", "/disputes", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

func (suite *DisputeErrorTestSuite) TestGetSummaryListErrorRetrievingDisputes() {
	suite.disputeMock.throwErrors = true
	req, _ := http.NewRequest("GET", "/disputes", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Error while retrieving disputes", result.Error)
}

func (suite *DisputeErrorTestSuite) TestGetSummaryListAdminBadRequest() {
	suite.jwtMock.returnUser.Role = "admin"
	body := `{asdasdas}`
	req, _ := http.NewRequest("POST", "/disputes", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestGetSummaryListAdminError() {
	suite.jwtMock.returnUser.Role = "admin"
	suite.disputeMock.throwErrors = true
	body := `{}`
	req, _ := http.NewRequest("POST", "/disputes", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

// ---------------------------------------------------------------- EVIDENCE UPLOAD
func (suite *DisputeErrorTestSuite) TestEvidenceUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("POST", "/1/evidence", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestEvidenceBadID() {
	req, _ := http.NewRequest("POST", "/asdasd/evidence", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestEvidenceEmptyMultipartForm() {
	req, _ := http.NewRequest("POST", "/1/evidence", bytes.NewBuffer([]byte{}))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "multipart/form-data")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

// func (suite *DisputeErrorTestSuite) TestEvidenceErrorOpeningMultipartFile() {
// 	data := bytes.NewBuffer([]byte{})
// 	form := multipart.NewWriter(data)
// 	form.CreateFormFile("files", "file1.txt")
// 	form.Close()

// 	req, _ := http.NewRequest("POST", "/1/evidence", data)
// 	req.Header.Add("Authorization", "Bearer mock")
// 	req.Header.Add("Content-Type", form.FormDataContentType())

// 	// Override the mock to simulate an error when opening the file
// 	mockFile := &mockDisputeModel{}
// 	mockFile.throwErrors = true

// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)

// 	suite.Equal(http.StatusInternalServerError, w.Code)
// }

func (suite *DisputeErrorTestSuite) TestEvidenceErrorDuringUpload() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)
	createFileField(form, "files", "file1.txt", "file contents")
	form.Close()

	req, _ := http.NewRequest("POST", "/1/evidence", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", form.FormDataContentType())

	suite.disputeMock.throwErrors = true

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
}

func (suite *DisputeErrorTestSuite) TestEvidenceMultipleFilesUpload() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)
	createFileField(form, "files", "file1.txt", "file contents 1")
	createFileField(form, "files", "file2.txt", "file contents 2")
	form.Close()

	req, _ := http.NewRequest("POST", "/1/evidence", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", form.FormDataContentType())

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusCreated, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Equal("Files uploaded", result.Data)
}

// func (suite *DisputeErrorTestSuite) TestEvidenceInvalidAuthorizationHeader() {
// 	data := bytes.NewBuffer([]byte{})
// 	form := multipart.NewWriter(data)
// 	createFileField(form, "files", "file1.txt", "file contents")
// 	form.Close()

// 	req, _ := http.NewRequest("POST", "/1/evidence", data)
// 	req.Header.Add("Authorization", "Bearer invalid-token")
// 	req.Header.Add("Content-Type", form.FormDataContentType())

// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)

//		var result models.Response
//		suite.Equal(http.StatusUnauthorized, w.Code)
//		suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
//		suite.NotEmpty(result.Error)
//	}
func (suite *DisputeErrorTestSuite) TestEvidenceValidUpload() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)
	createFileField(form, "files", "file1.txt", "file contents")
	form.Close()

	req, _ := http.NewRequest("POST", "/1/evidence", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", form.FormDataContentType())

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusCreated, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Equal("Files uploaded", result.Data)
}

func (suite *DisputeErrorTestSuite) TestEvidenceBadBody() {
	req, _ := http.NewRequest("POST", "/1/evidence", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestEvidenceBadHeaders() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)

	createFileField(form, "files", "file1.txt", "file contents")

	req, _ := http.NewRequest("POST", "/1/evidence", data)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestEvidenceUploadFail() {
	suite.disputeMock.throwErrors = true

	req, _ := http.NewRequest("POST", "/1/evidence", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestEvidenceUploadSingle() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)

	createFileField(form, "files", "file1.txt", "file contents")
	form.Close()

	req, _ := http.NewRequest("POST", "/1/evidence", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", form.Boundary()))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data string `json:"data"`
	}
	suite.Equal(http.StatusCreated, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
	suite.Equal([]mockEvidence{
		{
			user:    0, // TODO: Replace this mocked data from the JWT
			dispute: 1,
			path:    filepath.Join("1/file1.txt"),
			data:    "file contents",
		},
	}, suite.disputeMock.evidence)
}

// ---------------------------------------------------------------- CREATE DISPUTE

func (suite *DisputeErrorTestSuite) TestCreateUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("POST", "/create", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestCreateNoBody() {
	req, _ := http.NewRequest("POST", "/create", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestCreateMissingFields() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)

	title, _ := form.CreateFormField("title")
	title.Write([]byte("Title"))
	form.Close()

	req, _ := http.NewRequest("POST", "/create", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", form.Boundary()))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestCreateFailed() {
	suite.disputeMock.throwErrors = true

	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)

	createStringField(form, "title", "Title")
	createStringField(form, "description", "Desc")
	createStringField(form, "respondent[full_name]", "First Last")
	createStringField(form, "respondent[email]", "Email")
	createStringField(form, "respondent[workflow]", "1")
	form.Close()

	req, _ := http.NewRequest("POST", "/create", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", form.Boundary()))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestCreateSuccess() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)

	createStringField(form, "title", "Title")
	createStringField(form, "description", "Desc")
	createStringField(form, "respondent[full_name]", "First Last")
	createStringField(form, "respondent[email]", "Email")
	createStringField(form, "respondent[workflow]", "1")
	form.Close()

	req, _ := http.NewRequest("POST", "/create", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", form.Boundary()))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusCreated, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

func (suite *DisputeErrorTestSuite) TestCreateFileUploads() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)

	createStringField(form, "title", "Title")
	createStringField(form, "description", "Desc")
	createStringField(form, "respondent[full_name]", "First Last")
	createStringField(form, "respondent[email]", "Email")
	createStringField(form, "respondent[workflow]", "1")

	createFileField(form, "files", "file1.txt", "contents 1")
	createFileField(form, "files", "file2.txt", "contents 2")
	createFileField(form, "files", "file3.txt", "contents 3")
	form.Close()

	req, _ := http.NewRequest("POST", "/create", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", form.Boundary()))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusCreated, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
	suite.Equal([]mockEvidence{
		{
			user:    0, // TODO: Replace this mocked data from the JWT
			dispute: 0,
			path:    filepath.Join("0/file1.txt"),
			data:    "contents 1",
		},
		{
			user:    0, // TODO: Replace this mocked data from the JWT
			dispute: 0,
			path:    filepath.Join("0/file2.txt"),
			data:    "contents 2",
		},
		{
			user:    0, // TODO: Replace this mocked data from the JWT
			dispute: 0,
			path:    filepath.Join("0/file3.txt"),
			data:    "contents 3",
		},
	}, suite.disputeMock.evidence)
}

// ---------------------------------------------------------------- GET DISPUTE

func (suite *DisputeErrorTestSuite) TestGetUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestGetBadID() {
	req, _ := http.NewRequest("GET", "/asdasd", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestGetFail() {
	suite.disputeMock.throwErrors = true

	req, _ := http.NewRequest("GET", "/1", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestGetSuccess() {
	req, _ := http.NewRequest("GET", "/1", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

// func (suite *DisputeErrorTestSuite) TestGetNoEvidence() {
// 	req, _ := http.NewRequest("GET", "/1", nil)
// 	req.Header.Add("Authorization", "Bearer mock")

// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)

// 	var result models.Response
// 	suite.Equal(http.StatusOK, w.Code)
// 	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
// 	suite.NotEmpty(result.Data)

//		disputeDetails := result.Data.(models.DisputeDetailsResponse)
//		suite.Empty(disputeDetails.Evidence)
//	}
func (suite *DisputeErrorTestSuite) TestGetNoExperts() {
	req, _ := http.NewRequest("GET", "/1", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data models.DisputeDetailsResponse `json:"data"`
	}
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)

	suite.Empty(result.Data.Experts)
}

func (suite *DisputeErrorTestSuite) TestGetLoggerInitializationError() {
	req, _ := http.NewRequest("GET", "/1", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

func (suite *DisputeErrorTestSuite) TestViewExpertRejectionsInvalidBody() {
	req, _ := http.NewRequest("POST", "/experts/objections", bytes.NewBuffer([]byte("invalid body")))

	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)

	suite.Equal("Invalid Body", result.Error)
}

//---------------------------------------------------------------- Expert Objection Review

func (suite *DisputeErrorTestSuite) TestExpertObjectionsReviewUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("POST", "/1/experts/review-rejection", nil)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Unauthorized", result.Error)
}

func (suite *DisputeErrorTestSuite) TestExpertObjectionsReviewInvalidDisputeID() {
	req, _ := http.NewRequest("PATCH", "/objections/invalid", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Invalid Dispute ID", result.Error)
}

func (suite *DisputeErrorTestSuite) TestExpertObjectionsReviewInvalidRequestBody() {
	req, _ := http.NewRequest("PATCH", "/objections/1", bytes.NewBuffer([]byte("invalid body")))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Invalid body", result.Error)
}

func (suite *DisputeErrorTestSuite) TestViewExpertRejectionsErrorRetrieving() {
	suite.disputeMock.throwErrors = true

	body := `{"Expert_id": 1, "Dispute_id": 1, "Limits": 10, "Offset": 0}`
	req, _ := http.NewRequest("POST", "/experts/objections", bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Internal Server Error", result.Error)
}

func (suite *DisputeErrorTestSuite) TestExpertObjectionsReviewErrorReviewingObjection() {
	reqBody := `{"expert_id": 1, "accepted": true}`
	req, _ := http.NewRequest("PATCH", "/objections/1", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	suite.disputeMock.throwErrors = true

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Missing fields in request", result.Error)
}

func (suite *DisputeErrorTestSuite) TestExpertObjectionsReviewSuccess() {
	suite.jwtMock.throwErrors = false
	reqBody := `{"status": "Overruled"}`
	req, _ := http.NewRequest("PATCH", "/objections/1", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(result.Error)
	fmt.Println("BODY: ", w.Body.String())
	suite.Equal(nil, result.Data)
}

//---------------------------------------------------------------- Expert Objection

func (suite *DisputeErrorTestSuite) TestExpertObjectionErrorDuringObjection() {
	reqBody := `{"expert_id": 1, "reason": "Conflict of interest"}`
	req, _ := http.NewRequest("POST", "/1/objections", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	suite.disputeMock.throwErrors = true

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)

	suite.Equal("Failed to get Expert ID", result.Error)
}

func (suite *DisputeErrorTestSuite) TestViewExpertRejectionsSuccess() {
	suite.disputeMock.throwErrors = false
	body := `{"Expert_id": 1, "Dispute_id": 1, "Limits": 10, "Offset": 0}`
	req, _ := http.NewRequest("POST", "/experts/objections", bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Parse the response into a structured result object
	var result models.Response
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Empty(result.Error)
	suite.NotEmpty(result.Data)

	var resultData []models.ExpertObjectionsView
	dataBytes, _ := json.Marshal(result.Data)
	json.Unmarshal(dataBytes, &resultData)

	fmt.Println("BODY: ", w.Body.String())

	// Expected result
	expected := []models.ExpertObjectionsView{
		{
			ObjectionID:     1,
			ExpertID:        1,
			ExpertFullName:  "name",
			DisputeID:       1,
			DisputeTitle:    "title",
			UserID:          1,
			UserFullName:    "name",
			ObjectionStatus: "status",
		},
	}

	// Assert the result matches the expected value
	suite.Equal(expected, resultData)
}

func (suite *DisputeErrorTestSuite) TestExpertObjectionSuccess() {
	reqBody := `{"expert_id": 1, "reason": "Conflict of interest"}`
	req, _ := http.NewRequest("POST", "/1/objections", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	//inject the mock
	suite.disputeMock.throwErrors = false
	suite.disputeMock.Get_Experts = []models.AdminDisputeExperts{
		{
			ExpertID: 1,
			FullName: "name",
			Status:   string(models.ObjectionOverruled),
		},
	}

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Empty(result.Error)
	suite.Equal(float64(0), result.Data)
}

func (suite *DisputeErrorTestSuite) TestExpertObjectionUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("POST", "/1/experts/reject", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Unauthorized", result.Error)
}

func (suite *DisputeErrorTestSuite) TestExpertObjectionInvalidDisputeID() {
	req, _ := http.NewRequest("POST", "/invalid/objections", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Invalid Dispute ID", result.Error)
}

func (suite *DisputeErrorTestSuite) TestExpertObjectionInvalidRequestBody() {
	req, _ := http.NewRequest("POST", "/1/objections", bytes.NewBuffer([]byte("invalid body")))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

//---------------------------------------------------------------- Update Dispute Status

func (suite *DisputeErrorTestSuite) TestUpdateStatusInvalidRequestBody() {
	req, _ := http.NewRequest("PUT", "/dispute/status", bytes.NewBuffer([]byte("invalid body")))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	var result models.Response
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Equal("Invalid request body", result.Error)
}

func (suite *DisputeErrorTestSuite) TestUpdateStatusUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("PUT", "/dispute/status", bytes.NewBuffer([]byte(`{"dispute_id": 1, "status": "Resolved"}`)))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	var result models.Response
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Equal("Unauthorized", result.Error)
}

func (suite *DisputeErrorTestSuite) TestUpdateStatusInternalError() {
	suite.jwtMock.throwErrors = false
	suite.disputeMock.throwErrors = true
	req, _ := http.NewRequest("PUT", "/1/status", bytes.NewBuffer([]byte(`{"status": "Resolved"}`)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer mock")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	var result models.Response
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Equal("Something went wrong", result.Error)
}

func (suite *DisputeErrorTestSuite) TestUpdateStatusSuccess() {
	suite.jwtMock.throwErrors = false
	suite.disputeMock.throwErrors = false
	req, _ := http.NewRequest("PUT", "/1/status", bytes.NewBuffer([]byte(`{"status": "Resolved"}`)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer mock")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	var result models.Response
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Equal("Dispute status update successful", result.Data)
}

// ---------------------------------------------------------------- Create Dispute

func (suite *DisputeErrorTestSuite) TestCreateDisputeUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("POST", "/create", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Unauthorized", result.Error)
}

func (suite *DisputeErrorTestSuite) TestCreateDisputeMissingTitle() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)
	form.CreateFormField("description")
	form.CreateFormField("respondent[full_name]")
	form.CreateFormField("respondent[email]")
	form.Close()

	req, _ := http.NewRequest("POST", "/create", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", form.FormDataContentType())

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("missing field in form: title", result.Error)
}

func (suite *DisputeErrorTestSuite) TestCreateDisputeMissingDescription() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)
	form.CreateFormField("title")
	form.CreateFormField("respondent[full_name]")
	form.CreateFormField("respondent[email]")
	form.Close()

	req, _ := http.NewRequest("POST", "/create", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", form.FormDataContentType())

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("missing field in form: description", result.Error)
}

func (suite *DisputeErrorTestSuite) TestCreateDisputeMissingRespondentFullName() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)
	form.CreateFormField("title")
	form.CreateFormField("description")
	form.CreateFormField("respondent[email]")
	form.Close()

	req, _ := http.NewRequest("POST", "/create", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", form.FormDataContentType())

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("missing field in form: respondent[full_name]", result.Error)
}

func (suite *DisputeErrorTestSuite) TestCreateDisputeMissingRespondentEmail() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)
	form.CreateFormField("title")
	form.CreateFormField("description")
	form.CreateFormField("respondent[full_name]")
	form.Close()

	req, _ := http.NewRequest("POST", "/create", data)
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", form.FormDataContentType())

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("missing field in form: respondent[email]", result.Error) // This should match now
}
