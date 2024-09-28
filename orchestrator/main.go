package main

import (
	"orchestrator/controller"
	"orchestrator/handlers"
	"orchestrator/workflow"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Create a new query engine instance
	queryEngine := workflow.CreateWorkflowQuery()

	// Create a new controller instance
	controller := controller.NewController()

	// Start the controller
	controller.Start()

	// Create a new handler instance
	apiHandler := workflow.CreateAPIWorkflow(queryEngine)
	handlers := handlers.NewHandler(controller, &apiHandler)

	// Add notify of update state machine handler
	router.POST("/restart", handlers.RestartStateMachine)

	// Add notify of START state machine handler
	router.POST("/start", handlers.StartStateMachine)

	// Add notify of EVENT state machine handler
	router.POST("/event", handlers.TransitionStateMachine)

	// Add triggers endpoint
	router.GET("/triggers", handlers.GetTriggers)

	router.Run(":8090")

	// Wait for a signal to shutdown
	controller.WaitForSignal()
}
