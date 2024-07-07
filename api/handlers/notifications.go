package handlers

import (
	"api/middleware"
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupNotificationRoutes(g *gin.RouterGroup, h Notification) {
	g.POST("/invite", h.InviteParty)
}

func (h Notification) InviteParty(c *gin.Context) {
	jwtClaims := middleware.GetClaims(c)
	if jwtClaims == nil {
		c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	var reqInv models.DisputeNotifyInvite
	if err := c.BindJSON(&reqInv); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request"})
		return
	}

	var dbDispute models.Dispute
	h.DB.Where("id = ?", reqInv.DisputeID).First(&dbDispute)


}
