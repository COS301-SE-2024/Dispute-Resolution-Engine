package ticket

import (
	"api/env"
	"api/middleware"
	"api/models"
	"api/utilities"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type TicketModel interface {
	getAdminTicketList(searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error)
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

func (t *ticketModelReal) getAdminTicketList(searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) ([]models.TicketSummaryResponse, int64, error) {
	logger := utilities.NewLogger().LogWithCaller()

	tickets := []models.TicketSummaryResponse{}
	var queryString strings.Builder
	var countString strings.Builder
	var countParams []interface{}
	var queryParams []interface{}

	queryString.WriteString("SELECT SELECT t.id, t.created_at, t.subject, t.status, u.id AS user_id, u.first_name, u.surname FROM tickets t JOIN users u ON t.created_by = u.id")
	countString.WriteString("SELECT COUNT(*) FROM disputes")
	if searchTerm != nil {
		queryString.WriteString(" WHERE disputes.title LIKE ?")
		countString.WriteString(" WHERE disputes.title LIKE ?")
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
			queryString.WriteString(filter.Attr + " = ?")
			countString.WriteString(filter.Attr + " = ?")
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

	if _, valid := validSortAttrs[sortAttr.Attr]; !valid {
		return tickets, 0, errors.New("invalid sort attribute")
	}

	if sortAttr.Order != "asc" && sortAttr.Order != "desc" {
		sortAttr.Order = "asc"
	}

	queryString.WriteString(" ORDER BY " + sortAttr.Attr + " " + sortAttr.Order)

	if limit != nil {
		queryString.WriteString(" LIMIT ?")
		queryParams = append(queryParams, *limit)
	}
	if offset != nil {
		queryString.WriteString(" OFFSET ?")
		queryParams = append(queryParams, *offset)
	}

	ticketsIntermediate := []models.TicketIntermediate{}
	err := t.db.Raw(queryString.String(), queryParams...).Scan(ticketsIntermediate).Error
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
		ticketResp.ID = ticket.Id
		ticketResp.DateCreated = ticket.CreatedAt.Format("2006-01-02")
		ticketResp.Subject = ticket.Subject
		ticketResp.Status = ticket.Status
		ticketResp.User = models.TicketUser{
			ID:       ticket.UserID,
			FullName: ticket.FirstName + " " + ticket.Surname,
		}
		tickets = append(tickets, ticketResp)
	}

	return tickets, count, err
}
