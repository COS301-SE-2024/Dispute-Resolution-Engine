package handlers

import (
	"api/models"
	"api/utilities"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupWorkflowRoutes(g *gin.RouterGroup, h Workflow) {

	g.GET("", h.GetWorkflows)
	g.GET("/:id", h.GetIndividualWorkflow)
	g.POST("", h.StoreWorkflow)
	g.PUT("/:id", h.UpdateWorkflow)
	g.DELETE("/:id", h.DeleteWorkflow)

	//manage active workflows
	g.POST("/activate", h.NewActiveWorkflow)
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
	var workflows []models.Workflow

	// Fetch workflows with tags, limiting to 10 results
	result := w.DB.Limit(10).Find(&workflows)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		logger.Error(result.Error)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	// Prepare the response structure
	var response []struct {
		models.Workflow
		Tags []models.Tag `json:"tags"`
	}

	for _, workflow := range workflows {
		var taggedWorkflow struct {
			models.Workflow
			Tags []models.Tag `json:"tags"`
		}

		taggedWorkflow.Workflow = workflow
		// Query for tags related to each workflow, explicitly selecting the fields
		err := w.DB.Table("workflow_tags").
			Select("tags.id, tags.tag_name").
			Joins("join tags on workflow_tags.tag_id = tags.id").
			Where("workflow_tags.workflow_id = ?", workflow.ID).
			Scan(&taggedWorkflow.Tags).Error

		if err != nil {
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
			return
		}

		response = append(response, taggedWorkflow)
	}

	c.JSON(http.StatusOK, models.Response{Data: response})
}

func (w Workflow) GetIndividualWorkflow(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	id := c.Param("id")

	var workflow models.Workflow
	result := w.DB.First(&workflow, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Warnf("Workflow with ID %s not found", id)
			c.JSON(http.StatusOK, models.Response{Error: "Workflow not found"})
		} else {
			logger.Error(result.Error)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		}
		return
	}

	// Fetch the tags associated with this workflow
	var tags []models.Tag
	err := w.DB.Table("workflow_tags").
		Select("tags.id, tags.tag_name").
		Joins("join tags on workflow_tags.tag_id = tags.id").
		Where("workflow_tags.workflow_id = ?", workflow.ID).
		Scan(&tags).Error

	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	// Create the response structure
	response := struct {
		models.Workflow
		Tags []models.Tag `json:"tags"`
	}{
		Workflow: workflow,
		Tags:     tags,
	}

	c.JSON(http.StatusOK, models.Response{Data: response})
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

	//comvert map[string] to raw json
	workflowDefinition, err := json.Marshal(workflow.Definition)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to process workflow definition"})
		return
	}

	//check all fields present
	if workflow.Name == "" || workflow.Author == nil || (workflow.Definition == nil || len(workflow.Definition) == 0) {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Missing required fields"})
		return
	}
	//put into struct
	res := &models.Workflow{
		Name:       workflow.Name,
		Definition: workflowDefinition,
		AuthorID:   *workflow.Author,
	}

	// Store the workflow in the database
	result := w.DB.Create(res)
	if result.Error != nil {
		logger.Error(result.Error)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	for _, tagID := range workflow.Category {
		labelledWorkflow := models.WorkflowTags{
			WorkflowID: res.ID,
			TagID:      uint64(tagID),
		}
		if err := w.DB.Create(&labelledWorkflow).Error; err != nil {
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to link workflow with tags"})
			return
		}
	}

	c.JSON(http.StatusOK, models.Response{Data: res})
}

func (w Workflow) UpdateWorkflow(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	id := c.Param("id")

	// Find the existing workflow
	var existingWorkflow models.Workflow
	result := w.DB.First(&existingWorkflow, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Warnf("Workflow with ID %s not found", id)
			c.JSON(http.StatusNotFound, models.Response{Error: "Workflow not found"})
		} else {
			logger.Error(result.Error)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		}
		return
	}

	// Bind the request payload to the UpdateWorkflow struct
	var updateData models.UpdateWorkflow
	err := c.BindJSON(&updateData)
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
		workflowDefinition, err := json.Marshal(*updateData.WorkflowDefinition)
		if err != nil {
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to process workflow definition"})
			return
		}
		existingWorkflow.Definition = workflowDefinition
	}

	// Update the AuthorID if provided
	if updateData.Author != nil {
		existingWorkflow.AuthorID = *updateData.Author
	}

	// Save the updated workflow
	result = w.DB.Save(&existingWorkflow)
	if result.Error != nil {
		logger.Error(result.Error)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	// Manage categories (tags) in labelled_workflow if provided
	if updateData.Category != nil {
		// Remove existing tags
		err = w.DB.Where("workflow_id = ?", existingWorkflow.ID).Delete(&models.WorkflowTags{}).Error
		if err != nil {
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to update categories"})
			return
		}

		// Insert new tags
		for _, categoryID := range *updateData.Category {
			labelledWorkflow := models.WorkflowTags{
				WorkflowID: existingWorkflow.ID,
				TagID:      uint64(categoryID),
			}
			err = w.DB.Create(&labelledWorkflow).Error
			if err != nil {
				logger.Error(err)
				c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to update categories"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, models.Response{Data: "Workflow updated"})
}

func (w Workflow) DeleteWorkflow(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	id := c.Param("id")

	// Find the workflow record
	var workflow models.Workflow
	result := w.DB.First(&workflow, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Warnf("Workflow with ID %s not found", id)
			c.JSON(http.StatusNotFound, models.Response{Error: "Workflow not found"})
		} else {
			logger.Error(result.Error)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		}
		return
	}

	// Delete the tags associated with the workflow
	err := w.DB.Where("workflow_id = ?", workflow.ID).Delete(&models.WorkflowTags{}).Error
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to delete associated tags"})
		return
	}

	// Delete the workflow record itself
	result = w.DB.Delete(&workflow)
	if result.Error != nil {
		logger.Error(result.Error)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Failed to delete workflow"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Workflow and associated tags deleted"})
}

type OrchestratorRequest struct {
	ID int64 `json:"id"`
}

func (w Workflow) NewActiveWorkflow(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	var newActiveWorkflow models.NewActiveWorkflow
	err := c.BindJSON(&newActiveWorkflow)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request payload"})
		return
	}

	// Check if all fields are present
	if newActiveWorkflow.DisputeID == nil || newActiveWorkflow.Workflow == nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Missing required fields"})
		return
	}

	// Find the dispute
	var dispute models.Dispute
	result := w.DB.First(&dispute, newActiveWorkflow.DisputeID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Warnf("Dispute with ID %d not found", newActiveWorkflow.DisputeID)
			c.JSON(http.StatusNotFound, models.Response{Error: "Dispute not found"})
		} else {
			logger.Error(result.Error)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		}
		return
	}

	// Find the workflow
	var workflow models.Workflow
	result = w.DB.First(&workflow, newActiveWorkflow.Workflow)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Warnf("Workflow with ID %d not found", newActiveWorkflow.Workflow)
			c.JSON(http.StatusNotFound, models.Response{Error: "Workflow not found"})
		} else {
			logger.Error(result.Error)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		}
		return
	}

	//add entry to active Workflows
	timeNow := time.Now()
	workflowID := int64(workflow.ID)
	activeWorkflow := models.ActiveWorkflows{
		ID:           *dispute.ID,
		Workflow:     workflowID,
		WorkflowInstance: workflow.Definition,
		DateSubmitted: timeNow,
	}

	result = w.DB.Create(&activeWorkflow)
	if result.Error != nil && result.Error != gorm.ErrDuplicatedKey {
		logger.Error(result.Error)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Request already exists"})
		return
	}else if result.Error != nil {
		logger.Error(result.Error)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	//send request to orchestrator to activate workflow

	// Get the environment variables
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

	startEndpoint, err := w.EnvReader.Get("ORCH_START")
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		return
	}

	// Send the request to the orchestrator
	payload := OrchestratorRequest{ID: activeWorkflow.ID}

	_, err = w.MakeRequestToOrchestrator(fmt.Sprintf("http://%s:%s%s", url, port, startEndpoint), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		//delete the active workflow from table
		w.DB.Delete(&activeWorkflow)
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: "Databse updated and request to Activate workflow sent"})

}



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
	var activeWorkflow models.ActiveWorkflows
	result := w.DB.First(&activeWorkflow, resetActiveWorkflow.DisputeID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Warnf("Active workflow with ID %d not found", resetActiveWorkflow.DisputeID)
			c.JSON(http.StatusNotFound, models.Response{Error: "Active workflow not found"})
		} else {
			logger.Error(result.Error)
			c.JSON(http.StatusInternalServerError, models.Response{Error: "Internal Server Error"})
		}
		return
	}

	// Update the active workflow
	activeWorkflow.CurrentState = *resetActiveWorkflow.CurrentState
	activeWorkflow.StateDeadline = *resetActiveWorkflow.Deadline
	result = w.DB.Save(&activeWorkflow)
	if result.Error != nil {
		logger.Error(result.Error)
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

	_, err = w.MakeRequestToOrchestrator(fmt.Sprintf("http://%s:%s%s", url, port, resetEndpoint), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Data: "Database updated and request to Reset workflow sent"})
}

//helper funciton to complete active workflow

func (w Workflow) MakeRequestToOrchestrator(endpoint string, payload OrchestratorRequest) (string, error) {
	logger := utilities.NewLogger().LogWithCaller()

	// Marshal the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("marshal error: ", err)
		return "", fmt.Errorf("internal server error")
	}
	logger.Info("Payload: ", string(payloadBytes))

	// Send the POST request to the orchestrator
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Error("post error: ", err)
		return "", fmt.Errorf("internal server error")
	}
	defer resp.Body.Close()

	// Check for a successful status code (200 OK)

	if resp.StatusCode  == http.StatusInternalServerError {
		logger.Error("status code error: ", resp.StatusCode)
		return "", fmt.Errorf("Check theat you gave the correct state name if resetting")
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("status code error: ", resp.StatusCode)
		return "", fmt.Errorf("internal server error")
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read body error: ", err)
		return "", fmt.Errorf("internal server error")
	}

	// Convert the response body to a string
	responseBody := string(bodyBytes)

	// Log the response body for debugging
	logger.Info("Response Body: ", responseBody)

	return responseBody, nil
}
	
	
