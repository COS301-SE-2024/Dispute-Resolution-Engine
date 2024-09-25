package ticket

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"encoding/json"
	"io"
	"net/http"
	"strings"

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
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.WithError(err).Error("Error reading request body")
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request body"})
			return
		}

		// Reset the body so it can be read again by BindJSON
		c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

		// Check if the body is valid JSON and not empty
		var bodyMap map[string]interface{}
		if err := json.Unmarshal(body, &bodyMap); err != nil {
			logger.WithError(err).Error("Invalid JSON format")
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request body"})
			return
		}

		// If the body contains no key-value pairs, consider it empty
		if len(bodyMap) == 0 {
			logger.Info("Empty request body")
			tickets, count, err := h.Model.getAdminTicketList(nil, nil, nil, nil, nil)
			if err != nil {
				logger.WithError(err).Error("error retrieving tickets")
				c.JSON(http.StatusInternalServerError, models.Response{Error: "Error while retrieving tickets"})
				return
			}
			if count == 0 {
				logger.Info("No disputes found")
				c.JSON(http.StatusOK, models.Response{Data: gin.H{
					"tickets": tickets,
					"total":   count,
				}})
				return
			}
			c.JSON(http.StatusOK, models.Response{Data: gin.H{
				"tickets": tickets,
				"total":   count,
			}})
			return
		}
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
		tickets, count, err := h.Model.getAdminTicketList(searchTerm, limit, offset, sort, filters)
		if err != nil {
			logger.WithError(err).Error("error retrieving tickets")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error while retrieving tickets"})
			return
		}
		if count == 0 {
			logger.Info("No tickets found")
			c.JSON(http.StatusOK, models.Response{Data: gin.H{
				"tickets": tickets,
				"total":   count,
			}})
			return
		}
		c.JSON(http.StatusOK, models.Response{Data: gin.H{
			"tickets": tickets,
			"total":   count,
		}})
		return
	}
}
