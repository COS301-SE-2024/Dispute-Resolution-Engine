package handlers

import (
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupExpertRoutes(g *gin.RouterGroup, h Expert) {
	g.Use(middleware.JWTMiddleware)
}

func (h Expert) getExpert(c *gin.Context) {
	
}
