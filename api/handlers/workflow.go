package handlers

import "github.com/gin-gonic/gin"

func SetupWorkflowRoutes(g *gin.RouterGroup, h Workflow) {
	g.GET("", h.GetWorkflows)
	g.GET("/:id", h.GetIndivualWorkflow)
	g.POST("", h.StoreWorkflow)
	g.PUT("/:id", h.UpdateWorkflow)
	g.DELETE("/:id", h.DeleteWorkflow)
}

func (w Workflow) GetWorkflows(c *gin.Context) {
}

func (w Workflow) GetIndivualWorkflow(c *gin.Context) {
}

func (w Workflow) StoreWorkflow(c *gin.Context) {
}

func (w Workflow) UpdateWorkflow(c *gin.Context) {
}

func (w Workflow) DeleteWorkflow(c *gin.Context) {
}

