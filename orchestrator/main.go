package main

import (
	"orchestrator/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	handlers := handlers.NewHandler()
	//add notify of update state machine handler
	router.POST("/restart", handlers.RestartStateMachine)
	//add notify of START state machine handler
	router.POST("/start", handlers.StartStateMachine)


	router.Run(":8090")
}
