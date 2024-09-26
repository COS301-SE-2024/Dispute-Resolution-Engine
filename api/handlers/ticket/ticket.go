package ticket

import (
	"api/middleware"
	"api/models"
	"api/utilities"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupTicketRoutes(g *gin.RouterGroup, h Ticket) {
	jwt := middleware.NewJwtMiddleware()
	g.Use(jwt.JWTMiddleware)
	g.POST("", h.getTicketList)
	g.GET("/:id", h.getUserTicketDetails)
	g.PATCH("/:id", h.patchTicketStatus)
	g.POST("/:id/messages", h.createTicketMessage)
	g.POST("/disputes/:id/ticket", h.createTicket)
}

func (h Ticket) createTicket(c *gin.Context) {
	createReq := models.TicketCreate{}
	logger := utilities.NewLogger().LogWithCaller()
	if err := c.BindJSON(&createReq); err != nil {
		logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request"})
		return
	}

	claims, err := h.JWT.GetClaims(c)
	if err != nil {
		logger.Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized access"})
		return
	}

	userRole := claims.Role
	//admins may not create tickets
	if userRole == "admin" {
		logger.Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized access"})
		return
	}

	disputeID := c.Param("id")
	if disputeID == "" {
		logger.Error("No dispute ID provided")
		c.JSON(http.StatusBadRequest, models.Response{Error: "No dispute ID provided"})
		return
	}

	disputeIntID, err := strconv.Atoi(disputeID)
	if err != nil {
		logger.WithError(err).Error("Invalid dispute ID")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid dispute ID"})
		return
	}

	err = h.Model.createTicket(claims.ID, int64(disputeIntID), createReq.Subject, createReq.Body)
	if err != nil {
		logger.WithError(err).Error("Error creating ticket")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error creating ticket"})
		return
	}
	

}

func (h Ticket) createTicketMessage(c *gin.Context) {
	tickReq := models.TicketMessageCreate{}
	logger := utilities.NewLogger().LogWithCaller()
	if err := c.BindJSON(&tickReq); err != nil {
		logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request"})
		return
	}
	claims, err := h.JWT.GetClaims(c)
	if err != nil {
		logger.Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized access"})
		return
	}
	ticketID := c.Param("id")
	if ticketID == "" {
		logger.Error("No ticket ID provided")
		c.JSON(http.StatusBadRequest, models.Response{Error: "No ticket ID provided"})
		return
	}
	ticketIDInt, err := strconv.Atoi(ticketID)
	if err != nil {
		logger.WithError(err).Error("Invalid ticket ID")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid ticket ID"})
		return
	}
	userRole := claims.Role
	if userRole == "admin" {
		err = h.Model.addAdminTicketMessage(int64(ticketIDInt), claims.ID, tickReq.Message)
		if err != nil {
			logger.WithError(err).Error("Error adding ticket message")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error adding ticket message"})
			return
		}
		c.JSON(http.StatusNoContent, nil)
		return
	}
	err = h.Model.addUserTicketMessage(int64(ticketIDInt), claims.ID, tickReq.Message)
	if err != nil {
		logger.WithError(err).Error("Error adding ticket message")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error adding ticket message"})
		return
	}
	c.JSON(http.StatusNoContent, nil)

}

func (h Ticket) patchTicketStatus(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	claims, err := h.JWT.GetClaims(c)
	if err != nil || claims.Role != "admin" {
		logger.Error("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized access"})
		return
	}

	ticketID := c.Param("id")
	if ticketID == "" {
		logger.Error("No ticket ID provided")
		c.JSON(http.StatusBadRequest, models.Response{Error: "No ticket ID provided"})
		return
	}

	tickReq := models.PatchTicketStatus{}
	if err := c.BindJSON(&tickReq); err != nil {
		logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request"})
		return
	}

	ticketIDInt, err := strconv.Atoi(ticketID)
	if err != nil {
		logger.WithError(err).Error("Invalid ticket ID")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid ticket ID"})
		return
	}

	err = h.Model.patchTicketStatus(tickReq.Status, int64(ticketIDInt))
	if err != nil {
		logger.WithError(err).Error("Error updating ticket status")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error updating ticket status"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h Ticket) getUserTicketDetails(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	ticketID := c.Param("id")
	if ticketID == "" {
		logger.Error("No ticket ID provided")
		c.JSON(http.StatusBadRequest, models.Response{Error: "No ticket ID provided"})
		return
	}

	ticketIDInt, err := strconv.Atoi(ticketID)
	if err != nil {
		logger.WithError(err).Error("Invalid ticket ID")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid ticket ID"})
		return
	}

	claims, err := h.JWT.GetClaims(c)
	if err != nil {
		logger.Error("Unauthorized access attempt, failed to fetch claims")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized access"})
		return
	}

	userRole := claims.Role

	if userRole == "admin" {
		ticketDetails, err := h.Model.getAdminTicketDetails(int64(ticketIDInt))
		if err != nil {
			logger.WithError(err).Error("Error retrieving ticket details")
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving ticket details"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Data: ticketDetails})
	}

	ticketDetails, err := h.Model.getTicketDetails(int64(ticketIDInt), claims.ID)

	if err.Error() == "Unauthorized ticket access attempt" {
		logger.Error("Unauthorized ticket access attempt")
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized ticket access attempt"})
		return
	}

	if err != nil {
		logger.WithError(err).Error("Error retrieving ticket details")
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Error retrieving ticket details"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: ticketDetails})

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
	userID := claims.ID

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
		var filters *[]models.Filter
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
			filters = &reqAdminTickets.Filter
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

	var reqTickets models.TicketsRequest
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
		tickets, count, err := h.Model.getTicketsByUserID(userID, nil, nil, nil, nil, nil)
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
	if err := c.BindJSON(&reqTickets); err != nil {
		logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request"})
		return
	}
	var searchTerm *string
	var limit *int
	var offset *int
	var sort *models.Sort
	var filters *[]models.Filter
	if reqTickets.Search != nil {
		searchTerm = reqTickets.Search
	}
	if reqTickets.Limit != nil {
		limit = reqTickets.Limit
	}
	if reqTickets.Offset != nil {
		offset = reqTickets.Offset
	}
	if reqTickets.Sort != nil {
		sort = reqTickets.Sort
	}
	if reqTickets.Filter != nil {
		filters = &reqTickets.Filter
	}
	tickets, count, err := h.Model.getTicketsByUserID(userID, searchTerm, limit, offset, sort, filters)
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

}
