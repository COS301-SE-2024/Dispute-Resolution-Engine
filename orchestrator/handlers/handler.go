package handlers

import (
	"encoding/json"
	"net/http"

	"orchestrator/controller"
	"orchestrator/utilities"
	"orchestrator/workflow"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)


type Handler struct {
	controller* controller.Controller // Pointer to the controller
	logger utilities.Logger
	api workflow.API
}

func NewHandler(ctrlr *controller.Controller, apiHandler workflow.API) *Handler {
	return &Handler{
		controller: ctrlr,
		logger: *utilities.NewLogger(),
		api: apiHandler,
	}
}

// For when the api notifies the orchestrator to start a new state machine.
// Body should contain the id of the state machine to be started in json.
func (h *Handler) StartStateMachine(c *gin.Context) {
	h.logger.Info("Starting state machine...")
	// Get the workflow ID from the request
	active_wf_id_str := c.PostForm("id")
	// Convert the workflow ID to an integer
	active_wf_id, err := strconv.Atoi(active_wf_id_str)
	if err != nil {
		h.logger.Error("Invalid workflow ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid workflow ID",
		})
		return
	}
	// Fetch the workflow definition from the database
	active_wf_record, err := h.api.FetchActiveWorkflow(active_wf_id)
	if err != nil {
		h.logger.Error("Error fetching workflow definition")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching workflow definition",
		})
		return
	}
	// Unmarshal the workflow definition
	var wf workflow.Workflow
	err = json.Unmarshal([]byte(active_wf_record.WorkflowInstance), &wf)
	if err != nil {
		h.logger.Error("Error unmarshalling workflow definition")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error unmarshalling workflow definition",
		})
		return
	}

	// Register the state machine with the controller
	h.controller.RegisterStateMachine(active_wf_id_str, wf)

	// Update the active workflow in the database
	dateSubmitted := time.Now()
	stateDeadline := time.Now().Add(wf.States[wf.Initial].Timer.Duration.Duration)
	err = h.api.UpdateActiveWorkflow(active_wf_id, nil, &wf.Initial, &dateSubmitted, &stateDeadline, nil)
	if err != nil {
		h.logger.Error("Error updating active workflow")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error updating active workflow",
		})
		return
	}

	// Return a success response
	h.logger.Info("State machine started successfully!")
	c.JSON(http.StatusOK, gin.H{
		"message": "State machine started successfully!",
	})
}

// For when the API notifies the orchestrator that there has been a manual change to the state machine.
// Will overwrite any current state machine with a new one.
// Body should contain the ID of the state machine in JSON format.
func (h *Handler) RestartStateMachine(c *gin.Context) {
	h.logger.Info("Restarting state machine...")

	// Get the workflow ID from the request
	active_wf_id_str := c.PostForm("id")
	
	// Convert the workflow ID to an integer
	active_wf_id, err := strconv.Atoi(active_wf_id_str)
	if err != nil {
		h.logger.Error("Invalid workflow ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid workflow ID",
		})
		return
	}

	// Fetch the workflow definition from the database
	active_wf_record, err := h.api.FetchActiveWorkflow(active_wf_id)
	if err != nil {
		h.logger.Error("Error fetching workflow definition")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching workflow definition",
		})
		return
	}

	// Unmarshal the workflow definition
	var wf workflow.Workflow
	err = json.Unmarshal([]byte(active_wf_record.WorkflowInstance), &wf)
	if err != nil {
		h.logger.Error("Error unmarshalling workflow definition")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error unmarshalling workflow definition",
		})
		return
	}

	// Get the current state and state deadline
	state_deadline := active_wf_record.StateDeadline
	current_state := active_wf_record.CurrentState

	// Make the current workflow's initial state the current state
	wf.Initial = current_state

	// Update the state deadline
	wf.States[current_state].Timer.Duration.Duration = time.Until(state_deadline)

	// Register the state machine with the controller
	h.controller.RegisterStateMachine(active_wf_id_str, wf)

	c.JSON(http.StatusOK, gin.H{
		"message": "State machine updated successfully!",
	})
}
