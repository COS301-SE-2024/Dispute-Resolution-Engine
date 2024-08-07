package handlers

import (
	"api/models"
	"api/utilities"

	"github.com/gin-gonic/gin"
)

func SetupWorkflowRoutes(g *gin.RouterGroup, h Workflow) {
	g.GET("", h.GetWorkflows)
	g.GET("/:id", h.GetIndivualWorkflow)
	g.POST("", h.StoreWorkflow)
	g.PUT("/:id", h.UpdateWorkflow)
	g.DELETE("/:id", h.DeleteWorkflow)
}

func (w Workflow) GetWorkflows(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	var workflows []models.Workflow
	result := w.DB.Find(&workflows)
	if result.Error != nil {
		logger.Error(result.Error)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, gin.H{"data": workflows})
}

func (w Workflow) GetIndivualWorkflow(c *gin.Context) {
}

func (w Workflow) StoreWorkflow(c *gin.Context) {
}

func (w Workflow) UpdateWorkflow(c *gin.Context) {
}

func (w Workflow) DeleteWorkflow(c *gin.Context) {
}

