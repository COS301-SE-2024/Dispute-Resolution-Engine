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
		var reqAdminTickets models.Tick
	}

}
