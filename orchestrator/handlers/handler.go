package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"orchestrator/controller"
	"orchestrator/utilities"
	"orchestrator/workflow"

	"github.com/gin-gonic/gin"
)

type Response struct {
	ID int64 `json:"id"`
}

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
	var Res Response
	if err := c.ShouldBindJSON(&Res); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}
	active_wf_id := int(Res.ID)
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
	active_wf_id_str := strconv.Itoa(active_wf_id)
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
	var Res Response
	if err := c.ShouldBindJSON(&Res); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}
	active_wf_id := int(Res.ID)

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
	fmt.Println("State deadline: ", state_deadline)
	fmt.Println("Current State: ", current_state)
	fmt.Println("Current State Deadline: ", wf.States[current_state])
	fmt.Println("Workflow: ", wf.GetWorkflowString())

	//check if state deadline is in future
	if time.Now().Before(state_deadline) {
		// If the state deadline is in the future, update the timer duration
		wf.States[current_state].Timer.Duration.Duration = time.Until(state_deadline)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "State deadline is in the past",
		})
		return
	}

	// Register the state machine with the controller
	active_wf_id_str := strconv.Itoa(active_wf_id)
	h.controller.RegisterStateMachine(active_wf_id_str, wf)

	c.JSON(http.StatusOK, gin.H{
		"message": "State machine updated successfully!",
	})
}
