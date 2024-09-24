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
	getAdminTicketList(searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter)
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

func (t *ticketModelReal) getAdminTicketList(searchTerm *string, limit *int, offset *int, sortAttr *models.Sort, filters *[]models.Filter) {
	logger := utilities.NewLogger().LogWithCaller()
	
	var disputes []models.AdminDisputeSummariesResponse = []models.AdminDisputeSummariesResponse{}
	var queryString strings.Builder
	var countString strings.Builder
	var countParams []interface{}
	var queryParams []interface{}

	queryString.WriteString("SELECT id, title, status, case_date, date_resolved FROM disputes")
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
		"status":        true,
	}

	if _, valid := validSortAttrs[sortAttr.Attr]; !valid {
		return disputes, 0, errors.New("invalid sort attribute")
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

}
