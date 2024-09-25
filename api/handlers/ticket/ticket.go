package ticket

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupTicketRoutes(g *gin.RouterGroup, h Ticket) {
	jwt := middleware.NewJwtMiddleware()
	g.Use(jwt.JWTMiddleware)
	g.POST("", h.getTicketList)
}

func (h Ticket) getTicketList(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	claims, err := h.JWT.GetClaims(c)

	if err != nil {
		logger.Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized access"})
		return
	}

	userRole := claims.Role

	if userRole == "admin" && c.Request.Method == "POST" {
		var reqAdminTickets models.TicketsRequest
		if err := c.BindJSON(&reqAdminTickets); err != nil {
			logger.WithError(err).Error("Invalid request")
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request"})
			return
		}
		var searchTerm *string
		var limit *int
		var offset *int
		var sort *models.Sort
		var filters *models.Filter
		if reqAdminTickets.Search != nil {
			searchTerm = reqAdminTickets.Search
		}
		if reqAdminTickets.Limit != nil {
			limit = reqAdminTickets.Limit
		}
		if reqAdminTickets.Offset != nil {
			offset = reqAdminTickets.Offset
		}
		if reqAdminTickets.Sort != nil {
			sort = reqAdminTickets.Sort
		}
		if reqAdminTickets.Filter != nil {
			filters = reqAdminTickets.Filter
		}

		tickets, count

	}
}
