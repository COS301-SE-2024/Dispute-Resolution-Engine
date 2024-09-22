package main

import (
	"orchestrator/handlers"
	"orchestrator/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Create a new controller instance
	controller := controller.NewController()
	// Start the controller
	controller.Start()

	// Create a new handler instance
	handlers := handlers.NewHandler(controller)

	// Add notify of update state machine handler
	router.POST("/restart", handlers.RestartStateMachine)
	// Add notify of START state machine handler
	router.POST("/start", handlers.StartStateMachine)


	router.Run(":8090")
}
