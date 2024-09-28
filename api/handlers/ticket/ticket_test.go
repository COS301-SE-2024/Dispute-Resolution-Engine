package ticket_test

import (
	"api/handlers/ticket"
	"api/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type mockTicketModel struct {
	throwErrors bool
}

type mockJwtModel struct {
	throwErrors bool
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
	router.POST("", handler.GetTicketList)
	router.GET("/:id", handler.GetUserTicketDetails)
	router.PATCH("/:id", handler.PatchTicketStatus)
	router.POST("/:id/messages", handler.CreateTicketMessage)
	router.POST("/create", handler.CreateTicket)

	suite.router = router
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

	var result models.Response
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &result))
	suite.NotEmpty(result.Error)
}


