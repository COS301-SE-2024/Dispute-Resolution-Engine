package handlers

import (
	"net/http"
	"strconv"
	"encoding/json"
	"orchestrator/controller"
	"orchestrator/utilities"
	"orchestrator/workflow"

	"github.com/gin-gonic/gin"
)


type Handler struct {
	controller* controller.Controller // Pointer to the controller
	logger utilities.Logger
	api workflow.API
}

func NewHandler(ctrlr *controller.Controller) *Handler {
	return &Handler{
		controller: ctrlr,
		logger: *utilities.NewLogger(),
		api: workflow.CreateAPIWorkflow(),
	}
}

// For when the api notifies the orchestrator to start a new state machine.
// Body should contain the id of the state machine json format.
func (h *Handler) StartStateMachine(c *gin.Context) {
	h.logger.Info("Starting state machine...")
	// Get the workflow ID from the request
	workflowIDStr := c.PostForm("id")
	// Convert the workflow ID to an integer
	workflowID, err := strconv.Atoi(workflowIDStr)
	if err != nil {
		h.logger.Error("Invalid workflow ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid workflow ID",
		})
		return
	}
	// Fetch the workflow definition from the database
	wf_record, err := h.api.FetchWorkflow(workflowID)
	if err != nil {
		h.logger.Error("Error fetching workflow definition")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching workflow definition",
		})
		return
	}

	// Unmarshal the workflow definition
	var wf workflow.Workflow
	err = json.Unmarshal([]byte(wf_record.Definition), &wf)
	if err != nil {
		h.logger.Error("Error unmarshalling workflow definition")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error unmarshalling workflow definition",
		})
		return
	}

	// Register the state machine with the controller
	h.controller.RegisterStateMachine(workflowIDStr, wf)

	// Return a success response
	h.logger.Info("State machine started successfully!")
	c.JSON(http.StatusOK, gin.H{
		"message": "State machine started successfully!",
	})
}

// for when the api notifies the orchestartor that there has been a manual change to the state machine
// will stop any current state machine and start a new one
//body should contain id of the state machine in json format
func (h *Handler) RestartStateMachine(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "State machine updated successfully!",
	})
}
