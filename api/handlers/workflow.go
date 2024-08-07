package handlers

import (
	"api/models"
	"api/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, models.Response{Data: workflows})
}


func (w Workflow) GetIndivualWorkflow(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	id := c.Param("id")

	var workflow models.Workflow
	result := w.DB.First(&workflow, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Warnf("Workflow with ID %s not found", id)
			c.JSON(http.StatusOK, models.Response{Data: nil})
		} else {
			logger.Error(result.Error)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		}
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: workflow})
}

func (w Workflow) StoreWorkflow(c *gin.Context) {
}

func (w Workflow) UpdateWorkflow(c *gin.Context) {
}

func (w Workflow) DeleteWorkflow(c *gin.Context) {
}

