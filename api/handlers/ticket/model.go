package ticket

import (
	"api/env"
	"api/middleware"
	"api/models"
	"api/utilities"
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type TicketModel interface {
	getAdminTicketList(searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error)
	getTicketsByUserID(uid int64, searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error)
	getTicketDetails(ticketID int64, userID int64) ([]models.TicketsByUser, error)
	getAdminTicketDetails(ticketID int64) ([]models.TicketsByUser, error)
}

type Ticket struct {
	Model TicketModel
	JWT   middleware.Jwt
	Env   env.Env
}

type ticketModelReal struct {
	db  *gorm.DB
	env env.Env
}

func NewHandler(db *gorm.DB, envReader env.Env) Ticket {
	return Ticket{
		Model: &ticketModelReal{db: db, env: env.NewEnvLoader()},
		JWT:   middleware.NewJwtMiddleware(),
		Env:   envReader,
	}
}

func (t *ticketModelReal) getAdminTicketDetails(ticketID int64) ([]models.TicketsByUser, error) {
	logger := utilities.NewLogger().LogWithCaller()
	tickets := []models.TicketsByUser{}
	var IntermediateTick = models.TicketIntermediate{}
	err := t.db.Raw("SELECT t.id, t.created_at, t.subject, t.status, t.initial_message, u.id AS user_id, u.first_name, u.surname FROM tickets t JOIN users u ON t.created_by = u.id WHERE t.id = ?", ticketID).Scan(&IntermediateTick).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving ticket")
		return tickets, err
	}
	var ticketMessages = []models.TicketMessages{}
	err = t.db.Raw("SELECT tm.id, tm.content, tm.user_id, u.first_name, u.surname FROM ticket_messages tm JOIN users u ON tm.user_id = u.id WHERE tm.ticket_id = ?", ticketID).Scan(&ticketMessages).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving ticket messages")
		return tickets, err
	}
	tickets = append(tickets, models.TicketsByUser{
		TicketSummaryResponse: models.TicketSummaryResponse{
			ID:          strconv.Itoa(int(IntermediateTick.Id)),
			User:        models.TicketUser{ID: strconv.Itoa(int(IntermediateTick.UserID)), FullName: IntermediateTick.FirstName + " " + IntermediateTick.Surname},
			DateCreated: IntermediateTick.CreatedAt.Format("2006-01-02"),
			Subject:     IntermediateTick.Subject,
			Status:      IntermediateTick.Status,
		},
		Body:     *IntermediateTick.InitialMessage,
		Messages: ticketMessages,
	})
	return tickets, err
}

func (t *ticketModelReal) getTicketDetails(ticketID int64, userID int64) ([]models.TicketsByUser, error) {
	logger := utilities.NewLogger().LogWithCaller()
	tickets := []models.TicketsByUser{}
	var IntermediateTick = models.TicketIntermediate{}
	err := t.db.Raw("SELECT t.id, t.created_at, t.subject, t.status, t.initial_message, u.id AS user_id, u.first_name, u.surname FROM tickets t JOIN users u ON t.created_by = u.id WHERE t.id = ? AND u.id = ?", ticketID, userID).Scan(&IntermediateTick).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving ticket")
		return tickets, err
	}
	var ticketMessages = []models.TicketMessages{}
	err = t.db.Raw("SELECT tm.id, tm.content, tm.user_id, u.first_name, u.surname FROM ticket_messages tm JOIN users u ON tm.user_id = u.id WHERE tm.ticket_id = ? AND u.id = ?", ticketID, userID).Scan(&ticketMessages).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving ticket messages")
		return tickets, err
	}
	tickets = append(tickets, models.TicketsByUser{
		TicketSummaryResponse: models.TicketSummaryResponse{
			ID:          strconv.Itoa(int(IntermediateTick.Id)),
			User:        models.TicketUser{ID: strconv.Itoa(int(IntermediateTick.UserID)), FullName: IntermediateTick.FirstName + " " + IntermediateTick.Surname},
			DateCreated: IntermediateTick.CreatedAt.Format("2006-01-02"),
			Subject:     IntermediateTick.Subject,
			Status:      IntermediateTick.Status,
		},
		Body:     *IntermediateTick.InitialMessage,
		Messages: ticketMessages,
	})

	if len(tickets) == 0 {
		return tickets, errors.New("Unauthorized ticket access attempt")
	}

	return tickets, err
}

func (t *ticketModelReal) getTicketsByUserID(uid int64, searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error) {
	logger := utilities.NewLogger().LogWithCaller()
	tickets := []models.TicketSummaryResponse{}
	var queryString strings.Builder
	var countString strings.Builder
	var countParams []interface{}
	var queryParams []interface{}

	queryString.WriteString("SELECT t.id , t.created_at, t.subject, t.status AS status, u.id AS user_id, u.first_name, u.surname FROM tickets t JOIN users u ON t.created_by = u.id WHERE u.id = ?")

	countString.WriteString("SELECT COUNT(*) FROM tickets t JOIN users u ON t.created_by = u.id WHERE u.id = ?")

	queryParams = append(queryParams, uid)
	countParams = append(countParams, uid)
	if searchTerm != nil {
		queryString.WriteString(" AND WHERE t.subject LIKE ?")
		countString.WriteString(" AND WHERE t.subject LIKE ?")
		queryParams = append(queryParams, "%"+*searchTerm+"%")
		countParams = append(countParams, "%"+*searchTerm+"%")
	}

	if filters != nil && len(*filters) > 0 {
		if searchTerm != nil {
			queryString.WriteString(" AND ")
			countString.WriteString(" AND ")
		} else {
			queryString.WriteString(" AND WHERE ")
			countString.WriteString(" AND WHERE ")
		}
		for i, filter := range *filters {
			queryString.WriteString("t." + filter.Attr + " = ?")
			countString.WriteString("t." + filter.Attr + " = ?")
			queryParams = append(queryParams, filter.Value)
			countParams = append(countParams, filter.Value)
			if i < len(*filters)-1 {
				queryString.WriteString(" AND ")
				countString.WriteString(" AND ")
			}
		}
	}

	validSortAttrs := map[string]bool{
		"date_created": true,
		"subject":      true,
		"user":         true,
		"status":       true,
	}

	if sortAttr != nil {
		if _, valid := validSortAttrs[sortAttr.Attr]; !valid {
			return tickets, 0, errors.New("invalid sort attribute")
		}

		if sortAttr.Order != "asc" && sortAttr.Order != "desc" {
			sortAttr.Order = "asc"
		}

		queryString.WriteString(" ORDER BY " + sortAttr.Attr + " " + sortAttr.Order)
	}

	if limit != nil {
		queryString.WriteString(" LIMIT ?")
		queryParams = append(queryParams, *limit)
	}
	if offset != nil {
		queryString.WriteString(" OFFSET ?")
		queryParams = append(queryParams, *offset)
	}

	ticketsIntermediate := []models.TicketIntermediate{}
	err := t.db.Raw(queryString.String(), queryParams...).Scan(&ticketsIntermediate).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving tickets")
		return tickets, 0, err
	}

	var count int64 = 0
	err = t.db.Raw(countString.String(), countParams...).Scan(&count).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving ticket count")
		return tickets, 0, err
	}

	for _, ticket := range ticketsIntermediate {
		var ticketResp models.TicketSummaryResponse
		ticketResp.ID = strconv.Itoa(int(ticket.Id))
		ticketResp.DateCreated = ticket.CreatedAt.Format("2006-01-02")
		ticketResp.Subject = ticket.Subject
		ticketResp.Status = ticket.Status
		ticketResp.User = models.TicketUser{
			ID:       strconv.Itoa(int(ticket.UserID)),
			FullName: ticket.FirstName + " " + ticket.Surname,
		}
		tickets = append(tickets, ticketResp)
	}

	return tickets, count, err

}

func (t *ticketModelReal) getAdminTicketList(searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error) {
	logger := utilities.NewLogger().LogWithCaller()
	tickets := []models.TicketSummaryResponse{}
	var queryString strings.Builder
	var countString strings.Builder
	var countParams []interface{}
	var queryParams []interface{}

	queryString.WriteString("SELECT t.id , t.created_at, t.subject, t.status AS status, u.id AS user_id, u.first_name, u.surname FROM tickets t JOIN users u ON t.created_by = u.id")
	countString.WriteString("SELECT COUNT(*) FROM tickets t JOIN users u ON t.created_by = u.id")
	if searchTerm != nil {
		queryString.WriteString(" WHERE t.subject LIKE ?")
		countString.WriteString(" WHERE t.subject LIKE ?")
		queryParams = append(queryParams, "%"+*searchTerm+"%")
		countParams = append(countParams, "%"+*searchTerm+"%")
	}

	if filters != nil && len(*filters) > 0 {
		if searchTerm != nil {
			queryString.WriteString(" AND ")
			countString.WriteString(" AND ")
		} else {
			queryString.WriteString(" WHERE ")
			countString.WriteString(" WHERE ")
		}
		for i, filter := range *filters {
			queryString.WriteString("t." + filter.Attr + " = ?")
			countString.WriteString("t." + filter.Attr + " = ?")
			queryParams = append(queryParams, filter.Value)
			countParams = append(countParams, filter.Value)
			if i < len(*filters)-1 {
				queryString.WriteString(" AND ")
				countString.WriteString(" AND ")
			}
		}
	}

	validSortAttrs := map[string]bool{
		"date_created": true,
		"subject":      true,
		"user":         true,
		"status":       true,
	}

	if sortAttr != nil {
		if _, valid := validSortAttrs[sortAttr.Attr]; !valid {
			return tickets, 0, errors.New("invalid sort attribute")
		}

		if sortAttr.Order != "asc" && sortAttr.Order != "desc" {
			sortAttr.Order = "asc"
		}

		queryString.WriteString(" ORDER BY " + sortAttr.Attr + " " + sortAttr.Order)
	}

	if limit != nil {
		queryString.WriteString(" LIMIT ?")
		queryParams = append(queryParams, *limit)
	}
	if offset != nil {
		queryString.WriteString(" OFFSET ?")
		queryParams = append(queryParams, *offset)
	}

	ticketsIntermediate := []models.TicketIntermediate{}
	err := t.db.Raw(queryString.String(), queryParams...).Scan(&ticketsIntermediate).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving tickets")
		return tickets, 0, err
	}

	var count int64 = 0
	err = t.db.Raw(countString.String(), countParams...).Scan(&count).Error
	if err != nil {
		logger.WithError(err).Error("Error retrieving ticket count")
		return tickets, 0, err
	}

	for _, ticket := range ticketsIntermediate {
		var ticketResp models.TicketSummaryResponse
		ticketResp.ID = strconv.Itoa(int(ticket.Id))
		ticketResp.DateCreated = ticket.CreatedAt.Format("2006-01-02")
		ticketResp.Subject = ticket.Subject
		ticketResp.Status = ticket.Status
		ticketResp.User = models.TicketUser{
			ID:       strconv.Itoa(int(ticket.UserID)),
			FullName: ticket.FirstName + " " + ticket.Surname,
		}
		tickets = append(tickets, ticketResp)
	}

	return tickets, count, err
}
