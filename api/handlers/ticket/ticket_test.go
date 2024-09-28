package ticket_test

import (
	"api/handlers/ticket"
	"api/models"
	"errors"

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
	suite.router = gin.Default()

	handler := ticket.Ticket{Model: suite.ticketMock, JWT: suite.jwtMock, Env: suite.envMock}
}

/*-------------------------------MOCK MODELS-----------------------------------------*/
/*------------------------------MOCK TICKET MODEL-----------------------------------------*/

func (t *mockTicketModel) getAdminTicketList(searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error) {
	if t.throwErrors {
		return nil, 0, errors.New("error")
	}
	return []models.TicketSummaryResponse{}, 0, nil
}

func (t *mockTicketModel) getTicketsByUserID(uid int64, searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error) {
	if t.throwErrors {
		return nil, 0, errors.New("error")
	}
	return []models.TicketSummaryResponse{}, 0, nil
}

func (t *mockTicketModel) getTicketDetails(ticketID int64, userID int64) (models.TicketsByUser, error) {
	if t.throwErrors {
		return models.TicketsByUser{}, errors.New("error")
	}
	return models.TicketsByUser{}, nil
}

func (t *mockTicketModel) getAdminTicketDetails(ticketID int64) (models.TicketsByUser, error) {
	if t.throwErrors {
		return models.TicketsByUser{}, errors.New("error")
	}
	return models.TicketsByUser{}, nil
}

func (t *mockTicketModel) patchTicketStatus(status string, ticketID int64) error {
	if t.throwErrors {
		return errors.New("error")
	}
	return nil
}

func (t *mockTicketModel) addUserTicketMessage(ticketID int64, userID int64, message string) (models.TicketMessage, error) {
	if t.throwErrors {
		return models.TicketMessage{}, errors.New("error")
	}
	return models.TicketMessage{}, nil
}

func (t *mockTicketModel) addAdminTicketMessage(ticketID int64, userID int64, message string) (models.TicketMessage, error) {
	if t.throwErrors {
		return models.TicketMessage{}, errors.New("error")
	}
	return models.TicketMessage{}, nil
}

func (t *mockTicketModel) createTicket(userID int64, dispute int64, subject string, message string) (models.Ticket, error) {
	if t.throwErrors {
		return models.Ticket{}, errors.New("error")
	}
	return models.Ticket{}, nil
}
