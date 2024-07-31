package dispute_test

import (
	"api/handlers/dispute"
	"api/models"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

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

type DisputeErrorTestSuite struct {
	suite.Suite
	mock   *mockDisputeModel
	router *gin.Engine
}

func (suite *DisputeErrorTestSuite) SetupTest() {
	suite.mock = &mockDisputeModel{}

	handler := dispute.Dispute{Model: suite.mock}
	gin.SetMode("release")
	router := gin.Default()
	router.POST("/:id/evidence", handler.UploadEvidence)

	suite.router = router
}
func TestDisputeErrors(t *testing.T) {
	suite.Run(t, new(DisputeErrorTestSuite))
}

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
	return nil, nil
}
func (m *mockDisputeModel) GetDisputeExperts(disputeId int64) ([]models.Expert, error) {
	if m.throwErrors {
		return nil, errors.ErrUnsupported
	}
	return nil, nil
}
func (m *mockDisputeModel) GetDisputesByUser(userId int64) ([]models.Dispute, error) {
	if m.throwErrors {
		return nil, errors.ErrUnsupported
	}
	return nil, nil
}
func (m *mockDisputeModel) GetDispute(disputeId int64) (models.Dispute, error) {
	if m.throwErrors {
		return models.Dispute{}, errors.ErrUnsupported
	}
	return models.Dispute{}, nil
}
func (m *mockDisputeModel) GetUserByEmail(email string) (models.User, error) {
	if m.throwErrors {
		return models.User{}, errors.ErrUnsupported
	}
	return models.User{}, nil
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

func (suite *DisputeErrorTestSuite) TestEvidenceUnauthorized() {
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
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
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

func (suite *DisputeErrorTestSuite) TestEvidenceUploadFail() {
	suite.mock.throwErrors = true

	req, _ := http.NewRequest("POST", "/1/evidence", nil)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *DisputeErrorTestSuite) TestEvidenceUploadSingle() {
	data := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(data)

	file, _ := form.CreateFormFile("files", "file1.txt")
	file.Write([]byte("file contents"))

	req, _ := http.NewRequest("POST", "/1/evidence", data)
	req.Header.Add("Authorization", "Bearer mock")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data string `json:"data"`
	}
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
	suite.Equal([]mockEvidence{
		{
			user:    0, // TODO: Replace this mocked data from the JWT
			dispute: 1,
			path:    filepath.Join("1/file1.txt"),
			data:    "file contents",
		},
	}, suite.mock.evidence)
}
