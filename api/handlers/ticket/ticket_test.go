package ticket_test

import (
	"api/handlers/ticket"
	"api/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type mockTicketModel struct {
	throwErrors bool
}

type mockJwtModel struct {
	throwErrors bool
	returnUser  models.UserInfoJWT
}

type mockEnv struct {
	throwErrors bool
	Error       error
}

type TicketErrorTestSuite struct {
	suite.Suite
	ticketMock *mockTicketModel
	jwtMock    *mockJwtModel
	envMock    *mockEnv
	router     *gin.Engine
}

func (suite *TicketErrorTestSuite) SetupTest() {
	suite.ticketMock = &mockTicketModel{}
	suite.jwtMock = &mockJwtModel{}
	suite.envMock = &mockEnv{}
	router := gin.Default()

	handler := ticket.Ticket{Model: suite.ticketMock, JWT: suite.jwtMock, Env: suite.envMock}
	router.Use(suite.jwtMock.JWTMiddleware)
	router.POST("/tickets", handler.GetTicketList)
	router.GET("/:id", handler.GetUserTicketDetails)
	router.PATCH("/:id", handler.PatchTicketStatus)
	router.POST("/:id/messages", handler.CreateTicketMessage)
	router.POST("/create", handler.CreateTicket)

	suite.router = router
}

func TestTicketErrors(t *testing.T) {
	suite.Run(t, new(TicketErrorTestSuite))
}

/*-------------------------------MOCK MODELS-----------------------------------------*/
//Ticket Mocks

func (t *mockTicketModel) GetAdminTicketList(searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error) {
	if t.throwErrors {
		return nil, 0, errors.New("error")
	}
	return []models.TicketSummaryResponse{}, 0, nil
}

func (t *mockTicketModel) GetTicketsByUserID(uid int64, searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error) {
	if t.throwErrors {
		return nil, 0, errors.New("error")
	}
	return []models.TicketSummaryResponse{}, 0, nil
}

func (t *mockTicketModel) GetTicketDetails(ticketID int64, userID int64) (models.TicketsByUser, error) {
	if t.throwErrors {
		return models.TicketsByUser{}, errors.New("error")
	}
	return models.TicketsByUser{}, nil
}

func (t *mockTicketModel) GetAdminTicketDetails(ticketID int64) (models.TicketsByUser, error) {
	if t.throwErrors {
		return models.TicketsByUser{}, errors.New("error")
	}
	return models.TicketsByUser{}, nil
}

func (t *mockTicketModel) PatchTicketStatus(status string, ticketID int64) error {
	if t.throwErrors {
		return errors.New("error")
	}
	return nil
}

func (t *mockTicketModel) AddUserTicketMessage(ticketID int64, userID int64, message string) (models.TicketMessage, error) {
	if t.throwErrors {
		return models.TicketMessage{}, errors.New("error")
	}
	return models.TicketMessage{}, nil
}

func (t *mockTicketModel) AddAdminTicketMessage(ticketID int64, userID int64, message string) (models.TicketMessage, error) {
	if t.throwErrors {
		return models.TicketMessage{}, errors.New("error")
	}
	return models.TicketMessage{}, nil
}

func (t *mockTicketModel) CreateTicket(userID int64, dispute int64, subject string, message string) (models.Ticket, error) {
	if t.throwErrors {
		return models.Ticket{}, errors.New("error")
	}
	return models.Ticket{}, nil
}

// JWT Mocks
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

//mock env

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

// ---------------------------------------------------------------- CREATE TICKET

func (suite *TicketErrorTestSuite) TestCreateUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("POST", "/create", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestCreateUnauthorizedAdmin() {
	suite.jwtMock.returnUser.Role = "admin"
	body := `{"dispute": 1, "subject": "test", "message": "test"}`
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestCreateBadRequest() {
	req, _ := http.NewRequest("POST", "/create", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestCreateBadRequest2() {
	body := `{"dispute": kjhdsak, "subject": "test","message": "test"}`
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestCreateError() {
	suite.ticketMock.throwErrors = true

	body := `{"dispute": 1, "subject": "test", "message": "test"}`
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestCreateSuccess() {
	body := `{"dispute": 1, "subject": "test", "message": "test"}`
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusCreated, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result)
}

// ---------------------------------------------------------------- CREATE TICKET MESSAGE
func (suite *TicketErrorTestSuite) TestCreateMessageUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("POST", "/1/messages", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestCreateMessageBadRequest() {
	req, _ := http.NewRequest("POST", "/1/messages", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestCreateTicketMessageBadID() {
	body := `{"message": "test"}`
	req, _ := http.NewRequest("POST", "/$/messages", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestCreateTicketMessageError() {
	suite.ticketMock.throwErrors = true
	body := `{"message": "test"}`
	req, _ := http.NewRequest("POST", "/1/messages", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestCreateTicketSuccess() {
	body := `{"message": "test"}`
	req, _ := http.NewRequest("POST", "/1/messages", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusCreated, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

// ---------------------------------------------------------------- PATCH TICKET STATUS

func (suite *TicketErrorTestSuite) TestPatchUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("PATCH", "/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestPatchUnauthorizedUser() {
	suite.jwtMock.returnUser.Role = "user"
	req, _ := http.NewRequest("PATCH", "/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestPatchError() {
	suite.ticketMock.throwErrors = true
	suite.jwtMock.returnUser.Role = "admin"
	body := `{"status": "Open"}`
	req, _ := http.NewRequest("PATCH", "/1", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestPatchBadRequest() {
	suite.jwtMock.returnUser.Role = "admin"
	req, _ := http.NewRequest("PATCH", "/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestPatchSuccess() {
	suite.jwtMock.returnUser.Role = "admin"
	body := `{"status": "Open"}`
	req, _ := http.NewRequest("PATCH", "/1", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusNoContent, w.Code)
}

// ---------------------------------------------------------------- GET TICKET DETAILS

func (suite *TicketErrorTestSuite) TestGetTicketDetailsUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestGetTicketDetailsError() {
	suite.ticketMock.throwErrors = true
	suite.jwtMock.returnUser.Role = "admin"
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestGetTicketDetailsSuccess() {
	suite.jwtMock.returnUser.Role = "admin"
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

// ---------------------------------------------------------------- GET USER TICKET DETAILS

func (suite *TicketErrorTestSuite) TestGetUserTicketDetailsBadRequest() {
	req, _ := http.NewRequest("GET", "/$", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestGetUserTicketDetailsUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestGetUserTicketDetailsError() {
	suite.ticketMock.throwErrors = true
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Error string `json:"error"`
	}
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestGetUserTicketDetailsSuccess() {
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

// ---------------------------------------------------------------- GET TICKET LIST

func (suite *TicketErrorTestSuite) TestGetTicketListUnauthorized() {
	suite.jwtMock.throwErrors = true
	req, _ := http.NewRequest("POST", "/tickets", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestGetTicketListError() {
	suite.ticketMock.throwErrors = true
	body := `{}`
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result models.Response

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}

func (suite *TicketErrorTestSuite) TestGetTicketListSuccess() {
	body := `{}`
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

func (suite *TicketErrorTestSuite) TestGetTicketListAdminSuccess() {
	suite.jwtMock.returnUser.Role = "admin"
	body := `{}`
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

func (suite *TicketErrorTestSuite) TestGetTicketListSuccess2() {
	body := `{"limit": 10, "offset": 0}`
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}

func (suite *TicketErrorTestSuite) TestGetTicketListAdminSuccess2() {
	suite.jwtMock.returnUser.Role = "admin"
	body := `{"limit": 10, "offset": 0}`
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var result struct {
		Data map[string]interface{} `json:"data"`
	}
	suite.Equal(http.StatusOK, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Data)
}
