package workflow

import (
	"api/models"
	"api/utilities"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	// "time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupWorkflowRoutes(g *gin.RouterGroup, h Workflow) {

	g.POST("", h.GetWorkflows)
	g.GET("/:id", h.GetIndividualWorkflow)
	g.POST("/create", h.StoreWorkflow)
	g.PATCH("/:id", h.UpdateWorkflow)
	g.DELETE("/:id", h.DeleteWorkflow)

	//manage active workflows
	// g.POST("/activate", h.NewActiveWorkflow)
	g.POST("/reset", h.ResetActiveWorkflow)
	// g.POST("/complete", h.CompleteActiveWorkflow)
}

type WorkflowResult struct {
	Workflow models.Workflow `json:"workflow"`
	Author   models.User     `json:"author,omitempty"`
	Category models.Tag      `json:"category,omitempty"`
}

func (w Workflow) GetWorkflows(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	//get parameters
	body := models.GetWorkflow{}
	err := c.BindJSON(&body)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request payload"})
		return
	}

	total, response, err := w.DB.GetWorkflowsWithLimitOffset(body.Limit, body.Offset, body.Search)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, models.Response{Data: []models.Workflow{}})
			return
		}

		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	// // Prepare the response structure
	// var response []struct {
	// 	models.Workflow
	// 	Tags []models.Tag `json:"tags"`
	// }

	// for _, workflow := range workflows {
	// 	var taggedWorkflow struct {
	// 		models.Workflow
	// 		Tags []models.Tag `json:"tags"`
	// 	}

	// 	taggedWorkflow.Workflow = workflow
	// 	// Query for tags related to each workflow, explicitly selecting the fields
	// 	taggedWorkflow.Tags, err = w.DB.QueryTagsToRelatedWorkflow(workflow.ID)
	// 	if err != nil {
	// 		logger.Error(err)
	// 		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
	// 		return
	// 	}

	// 	response = append(response, taggedWorkflow)
	// }

	if len(response) == 0 {

		c.JSON(http.StatusOK, models.Response{Data: models.WorkflowResult{Total: 0, Workflows: []models.GetWorkflowResponse{}}})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: models.WorkflowResult{Total: int(total), Workflows: response}})
}

func (w Workflow) GetIndividualWorkflow(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid ID parameter"})
		return
	}

	workflow, err := w.DB.GetWorkflowByID(uint64(idInt))
	if err != nil && err == gorm.ErrRecordNotFound {
		logger.Error(err)
		c.JSON(http.StatusOK, models.Response{Error: "Record not found"})
		return
	}
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	// // Fetch the tags associated with this workflow
	// var tags []models.Tag
	// tags, err = w.DB.QueryTagsToRelatedWorkflow(uint64(idInt))
	// if err != nil {
	// 	logger.Error(err)
	// 	c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
	// 	return
	// }

	// response := struct {
	// 	models.Workflow
	// 	Tags []models.Tag `json:"tags"`
	// }{
	// 	Workflow: *workflow,
	// 	Tags:     tags,
	// }
	c.JSON(http.StatusOK, models.Response{Data: workflow})
}
func (w Workflow) StoreWorkflow(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	var workflow models.CreateWorkflow

	// Bind incoming JSON to the struct
	err := c.BindJSON(&workflow)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request payload"})
		return
	}

	// Get details from JWT
	claims, err := w.Jwt.GetClaims(c)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to get user details"})
		return
	}

	// Validate the workflow definition
	if err := ValidateWorkflowDefinition(workflow.Definition); err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	// Convert map[string] to raw JSON
	workflowDefinition, err := json.Marshal(workflow.Definition)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to process workflow definition"})
		return
	}

	// Check all fields are present
	if workflow.Name == "" || workflowDefinition == nil || workflow.Definition.Initial == "" || workflow.Definition.States == nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Missing required fields"})
		return
	}

	// Put into struct
	res := &models.Workflow{
		Name:       workflow.Name,
		Definition: workflowDefinition,
		AuthorID:   claims.ID,
	}

	// Store the workflow in the database
	result := w.DB.CreateWorkflow(res)
	if result != nil {
		logger.Error(result)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	// Link workflow with tags (if necessary)
	// for _, tagID := range workflow.Category {
	// 	labelledWorkflow := models.WorkflowTags{
	// 		WorkflowID: res.ID,
	// 		TagID:      uint64(tagID),
	// 	}
	// 	if err := w.DB.CreateWorkflowTag(&labelledWorkflow); err != nil {
	// 		logger.Error(err)
	// 		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to link workflow with tags"})
	// 		return
	// 	}
	// }

	c.JSON(http.StatusOK, models.Response{Data: res})
}

// ValidateWorkflowDefinition checks if the workflow definition is valid
func ValidateWorkflowDefinition(definition models.WorkflowOrchestrator) error {
	// Check if Initial state exists in States
	if _, exists := definition.States[definition.Initial]; !exists {
		return fmt.Errorf("initial state '%s' does not exist in states", definition.Initial)
	}

	for stateID, state := range definition.States {
		// Check if state has a valid label
		if state.Label == "" {
			return fmt.Errorf("state '%s' is missing a label", stateID)
		}

		// Check if state has a description
		if state.Description == "" {
			return fmt.Errorf("state '%s' is missing a description", stateID)
		}

		// Validate each trigger in the state
		for triggerID, trigger := range state.Triggers {
			if trigger.Label == "" {
				return fmt.Errorf("trigger '%s' in state '%s' is missing a label", triggerID, stateID)
			}

			// Check if the next state exists in the workflow
			if _, exists := definition.States[trigger.Next]; !exists {
				return fmt.Errorf("trigger '%s' in state '%s' points to a non-existent state '%s'", triggerID, stateID, trigger.Next)
			}
		}

		// Validate the timer if present
		if state.Timer != nil {
			if state.Timer.Duration.Duration == 0 {
				return fmt.Errorf("timer in state '%s' must have a non-zero duration", stateID)
			}

			if _, exists := state.Triggers[state.Timer.OnExpire]; !exists {
				return fmt.Errorf("timer in state '%s' points to a non-existent trigger '%s'", stateID, state.Timer.OnExpire)
			}
		}
	}

	return nil
}

func (w Workflow) UpdateWorkflow(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid ID parameter"})
		return
	}

	// Find the existing workflow
	var existingWorkflow *models.Workflow
	existingWorkflow, err = w.DB.GetWorkflowRecordByID(uint64(idInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warnf("Workflow with ID %s not found", id)
			c.JSON(http.StatusNotFound, models.Response{Error: "Workflow not found"})
		} else {
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		}
		return
	}

	// Bind the request payload to the UpdateWorkflow struct
	var updateData models.UpdateWorkflow
	err = c.BindJSON(&updateData)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request payload"})
		return
	}

	// Update the Name if provided
	if updateData.Name != nil {
		existingWorkflow.Name = *updateData.Name
	}

	// Update the WorkflowDefinition if provided
	if updateData.WorkflowDefinition != nil {
		//validate the workflow definition
		if err := ValidateWorkflowDefinition(*updateData.WorkflowDefinition); err != nil {
			logger.Error(err)
			c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
			return
		}

		workflowDefinition, err := json.Marshal(*updateData.WorkflowDefinition)
		if err != nil {

			logger.Error(err)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to process workflow definition"})
			return
		}
		existingWorkflow.Definition = workflowDefinition
	}

	// // Update the AuthorID if provided
	// if updateData.Author != nil {
	// 	existingWorkflow.AuthorID = *updateData.Author
	// }

	// Save the updated workflow
	err = w.DB.UpdateWorkflow(existingWorkflow)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	// // Manage categories (tags) in labelled_workflow if provided
	// if updateData.Category != nil {
	// 	// Remove existing tags
	// 	err = w.DB.DeleteTagsByWorkflowID(existingWorkflow.ID)
	// 	if err != nil {
	// 		fmt.Println("here")

	// 		logger.Error(err)
	// 		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to update categories"})
	// 		return
	// 	}

	// 	// Insert new tags
	// 	for _, categoryID := range *updateData.Category {
	// 		labelledWorkflow := models.WorkflowTags{
	// 			WorkflowID: existingWorkflow.ID,
	// 			TagID:      uint64(categoryID),
	// 		}
	// 		err = w.DB.CreateWorkflowTag(&labelledWorkflow)
	// 		if err != nil {
	// 			logger.Error(err)
	// 			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to update categories"})
	// 			return
	// 		}
	// 	}
	// }

	c.JSON(http.StatusOK, models.Response{Data: "Workflow updated"})
}

func (w Workflow) DeleteWorkflow(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid ID parameter"})
		return
	}

	// Find the workflow record
	var workflow *models.Workflow
	workflow, result := w.DB.GetWorkflowRecordByID(uint64(idInt))
	if result != nil {
		if result == gorm.ErrRecordNotFound {
			logger.Warnf("Workflow with ID %s not found", id)
			c.JSON(http.StatusNotFound, models.Response{Error: "Workflow not found"})
		} else {
			logger.Error(result)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		}
		return
	}

	// Delete the tags associated with the workflow
	err = w.DB.DeleteTagsByWorkflowID(workflow.ID)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to delete associated tags"})
		return
	}

	// Delete the workflow record itself
	result = w.DB.DeleteWorkflow(workflow)
	if result != nil {
		logger.Error(result)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to delete workflow"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Workflow and associated tags deleted"})
}

// func (w Workflow) NewActiveWorkflow(c *gin.Context) {
// 	logger := utilities.NewLogger().LogWithCaller()

// 	var newActiveWorkflow models.NewActiveWorkflow
// 	err := c.BindJSON(&newActiveWorkflow)
// 	if err != nil {
// 		logger.Error(err)
// 		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request payload"})
// 		return
// 	}

// 	// Check if all fields are present
// 	if newActiveWorkflow.DisputeID == nil || newActiveWorkflow.Workflow == nil {
// 		c.JSON(http.StatusBadRequest, models.Response{Error: "Missing required fields"})
// 		return
// 	}

// 	// Find the dispute
// 	disputeID := uint64(*(newActiveWorkflow).DisputeID)
// 	dispute, result := w.DB.FindDisputeByID(disputeID)
// 	if result != nil {
// 		if result == gorm.ErrRecordNotFound {
// 			logger.Warnf("Dispute with ID %d not found", newActiveWorkflow.DisputeID)
// 			c.JSON(http.StatusNotFound, models.Response{Error: "Dispute not found"})
// 		} else {
// 			logger.Error(result)
// 			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
// 		}
// 		return
// 	}

// 	// Find the workflow
// 	workflow, result := w.DB.GetWorkflowRecordByID(uint64(*newActiveWorkflow.Workflow))
// 	if result != nil {
// 		if result == gorm.ErrRecordNotFound {
// 			logger.Warnf("Workflow with ID %d not found", newActiveWorkflow.Workflow)
// 			c.JSON(http.StatusNotFound, models.Response{Error: "Workflow not found"})
// 		} else {
// 			logger.Error(result)
// 			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
// 		}
// 		return
// 	}

// 	//add entry to active Workflows
// 	timeNow := time.Now()
// 	workflowID := int64(workflow.ID)
// 	activeWorkflow := models.ActiveWorkflows{
// 		ID:               *dispute.ID,
// 		Workflow:         workflowID,
// 		WorkflowInstance: workflow.Definition,
// 		DateSubmitted:    timeNow,
// 	}

// 	result = w.DB.CreateActiveWorkflow(&activeWorkflow)
// 	if result != nil && result != gorm.ErrDuplicatedKey {
// 		logger.Error(result)
// 		c.JSON(http.StatusInternalServerError, models.Response{Error: "Request already exists"})
// 		return
// 	} else if result != nil {
// 		logger.Error(result)
// 		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
// 		return
// 	}

// 	//send request to orchestrator to activate workflow

// 	// Get the environment variables
// 	url, err := w.EnvReader.Get("ORCH_URL")
// 	if err != nil {
// 		logger.Error(err)
// 		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
// 		return
// 	}

// 	port, err := w.EnvReader.Get("ORCH_PORT")
// 	if err != nil {
// 		logger.Error(err)
// 		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
// 		return
// 	}

// 	startEndpoint, err := w.EnvReader.Get("ORCH_START")
// 	if err != nil {
// 		logger.Error(err)
// 		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
// 		return
// 	}

// 	// Send the request to the orchestrator
// 	payload := OrchestratorRequest{ID: activeWorkflow.ID}

// 	_, err = w.OrchestratorEntity.MakeRequestToOrchestrator(fmt.Sprintf("http://%s:%s%s", url, port, startEndpoint), payload)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
// 		//delete the active workflow from table
// 		w.DB.DeleteActiveWorkflow(&activeWorkflow)
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.Response{Data: "Databse updated and request to Activate workflow sent"})

// }

func (w Workflow) ResetActiveWorkflow(c *gin.Context) {

	logger := utilities.NewLogger().LogWithCaller()

	var resetActiveWorkflow models.ResetActiveWorkflow
	err := c.BindJSON(&resetActiveWorkflow)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request payload"})
		return
	}

	// check if all fields are present
	if resetActiveWorkflow.DisputeID == nil || resetActiveWorkflow.CurrentState == nil || resetActiveWorkflow.Deadline == nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Missing required fields"})
		return
	}

	// find active workflow
	activeWorkflow, result := w.DB.GetActiveWorkflowByWorkflowID(uint64(*resetActiveWorkflow.DisputeID))
	if result != nil {
		if result == gorm.ErrRecordNotFound {
			logger.Warnf("Active workflow with ID %d not found", resetActiveWorkflow.DisputeID)
			c.JSON(http.StatusNotFound, models.Response{Error: "Active workflow not found"})
		} else {
			logger.Error(result)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		}
		return
	}

	// Update the active workflow
	activeWorkflow.CurrentState = *resetActiveWorkflow.CurrentState
	activeWorkflow.StateDeadline = *resetActiveWorkflow.Deadline
	result = w.DB.UpdateActiveWorkflow(activeWorkflow)
	if result != nil {
		logger.Error(result)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	//get Environment variables
	url, err := w.EnvReader.Get("ORCH_URL")
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	port, err := w.EnvReader.Get("ORCH_PORT")
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	resetEndpoint, err := w.EnvReader.Get("ORCH_RESET")
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	// Send the request to the orchestrator
	payload := OrchestratorRequest{ID: activeWorkflow.ID}

	body, err := w.OrchestratorEntity.MakeRequestToOrchestrator(fmt.Sprintf("http://%s:%s%s", url, port, resetEndpoint), payload)
	if err != nil {
		if strings.Contains(body, "deadline") {
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid deadline specified"})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, models.Response{Data: "Database updated and request to Reset workflow sent"})
}
