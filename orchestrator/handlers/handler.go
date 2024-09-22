package handlers

import (
	"net/http"
	"orchestrator/controller"
	"github.com/gin-gonic/gin"
)


type Handler struct {
	controller* controller.Controller
}

func NewHandler() *Handler {
	return &Handler{
		controller: controller.NewController(),
	}
}

// for when the api notifies the orchestartor that there has been a manual change to the state machine
// will stop any current state machine and start a new one
//body should contain id of the state machine in json format
func (h *Handler) RestartStateMachine(c *gin.Context) {
	// Logic for updating the state machine
	c.JSON(http.StatusOK, gin.H{
		"message": "State machine updated successfully!",
	})
}

// for when the api notifies the orchestrator to start a new state machine
// body should contain the id of the state machine json format
func (h *Handler) StartStateMachine(c *gin.Context) {
	

	c.JSON(http.StatusOK, gin.H{
		"message": "State machine started successfully!",
	})
}