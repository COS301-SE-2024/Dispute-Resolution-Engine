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
}
type mockJwtModel struct {
	throwErrors bool
}
type mockEmailModel struct {
	throwErrors bool
}

type mockAuditLogger struct {
}

type DisputeErrorTestSuite struct {
	suite.Suite
	disputeMock *mockDisputeModel
	jwtMock     *mockJwtModel
	emailMock   *mockEmailModel
	router      *gin.Engine
	auditMock   *mockAuditLogger
}

func (suite *DisputeErrorTestSuite) SetupTest() {
	suite.disputeMock = &mockDisputeModel{}
	suite.jwtMock = &mockJwtModel{}
	suite.emailMock = &mockEmailModel{}
	suite.auditMock = &mockAuditLogger{}

	handler := dispute.Dispute{Model: suite.disputeMock, JWT: suite.jwtMock, Email: suite.emailMock, AuditLogger: suite.auditMock}
	gin.SetMode("release")
	router := gin.Default()
	router.Use(suite.jwtMock.JWTMiddleware)
	router.POST("/:id/evidence", handler.UploadEvidence)
	router.POST("/create", handler.CreateDispute)
	router.GET("/:id", handler.GetDispute)
	router.POST("", handler.GetSummaryListOfDisputes)
	router.POST("/:id/experts/reject", handler.ExpertObjection)
	router.POST("/:id/experts/review-rejection", handler.ExpertObjectionsReview)
	router.PUT("/dispute/status", handler.UpdateStatus)


	suite.router = router
}
func TestDisputeErrors(t *testing.T) {
	suite.Run(t, new(DisputeErrorTestSuite))
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
	req, _ := http.NewRequest("POST", "/invalid/experts/reject", nil)
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
	req, _ := http.NewRequest("POST", "/1/experts/reject", bytes.NewBuffer([]byte("invalid body")))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
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
// mock model auditlogger
func (m *mockAuditLogger) LogDisputeProceedings(proceedingType models.EventTypes, eventData map[string]interface{}) error {
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
func (m *mockDisputeModel) GetUserByEmail(email string) (models.User, error) {
	if m.throwErrors {
		return models.User{}, errors.ErrUnsupported
	}
	return models.User{
		PhoneNumber:       new(string),
		AddressID:         new(int64),
		LastUpdate:        new(time.Time),
		LastLogin:         new(time.Time),
		PreferredLanguage: new(string),
		Timezone:          new(string),
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
func (m *mockDisputeModel) ObjectExpert(userId, disputeId, expertId int64, reason string) error {
	if m.throwErrors {
		return errors.ErrUnsupported
	}
	return nil
}
func (m *mockDisputeModel) ReviewExpertObjection(userId, disputeId, expertId int64, approved bool) error {
	if m.throwErrors {
		return errors.ErrUnsupported
	}
	return nil
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

// ---------------------------------------------------------------- Get Summary List Tests

func (suite *DisputeErrorTestSuite) TestGetSummaryListUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("GET", "/summary", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusUnauthorized, w.Code)
	fmt.Println("THIS BODY:", w.Body.String())
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Unauthorized", result.Error)
}

// func (suite *DisputeErrorTestSuite) TestGetSummaryListNoDisputes() {
// 	req, _ := http.NewRequest("GET", "/summary", nil)
// 	req.Header.Add("Authorization", "Bearer mock")

// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)

// 	var result models.Response
// 	suite.Equal(http.StatusOK, w.Code)
// 	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
// 	suite.Empty(result.Error)
// 	suite.Empty(result.Data)
// }

// func (suite *DisputeErrorTestSuite) TestGetSummaryListWithDisputes() {
// 	// Mock disputes
// 	suite.disputeMock.evidence = []mockEvidence{
// 		{
// 			user:    1,
// 			dispute: 1,
// 			path:    "path/to/file",
// 			data:    "evidence data",
// 		},
// 	}

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

func (suite *DisputeErrorTestSuite) TestGetSummaryListErrorRetrievingDisputes() {
	suite.disputeMock.throwErrors = true
	req, _ := http.NewRequest("GET", "/summary", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Invalid dispute id 'summary'", result.Error)
}

func (m *mockDisputeModel) GenerateAISummary(disputeID int64, disputeDesc string, apiKey string) {

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
	req, _ := http.NewRequest("POST", "/invalid/experts/review-rejection", nil)
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
	req, _ := http.NewRequest("POST", "/1/experts/review-rejection", bytes.NewBuffer([]byte("invalid body")))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("Request failed", result.Error)
}

func (suite *DisputeErrorTestSuite) TestExpertObjectionsReviewErrorReviewingObjection() {
	reqBody := `{"expert_id": 1, "accepted": true}`
	req, _ := http.NewRequest("POST", "/1/experts/review-rejection", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	suite.disputeMock.throwErrors = true

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
	suite.Equal("failed to review objection", result.Error)
}

func (suite *DisputeErrorTestSuite) TestExpertObjectionsReviewSuccess() {
	suite.jwtMock.throwErrors = false
	reqBody := `{"expert_id": 1, "accepted": true}`
	req, _ := http.NewRequest("POST", "/1/experts/review-rejection", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Add("Authorization", "Bearer mock")
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	fmt.Println("BODY:", w.Body.String())
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.Empty(result.Error)
	suite.Equal("Expert objections reviewed successfully", result.Data)
}
