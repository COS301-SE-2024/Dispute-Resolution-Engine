package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	// "orchestrator/db"
	// "orchestrator/env"
	// "orchestrator/utilities"
	// "orchestrator/env"
	// "orchestrator/scheduler"
	// "orchestrator/statemachine"
	"orchestrator/controller"
	"orchestrator/handlers"
	// "orchestrator/wf"
	"orchestrator/workflow"

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


func testInitalWorkflow(){
	// ======== Json Tests =========
	//read template workflow form json file
	// Read the JSON file

	wf, shouldReturn := readWorkflowFromFile("templates/v2.json")
	if shouldReturn {
		return
	}
	fmt.Println(wf.GetWorkflowString())

	

	// ======== Statemachine Tests =========
	// Create a new controller
	c := controller.NewController()
	c.Start()
	// Register the workflow with the controller
	c.RegisterStateMachine("wf1", wf)
	// Wait for a signal to shutdown
	c.WaitForSignal()
}

func readWorkflowFromFile(fileName string) (workflow.Workflow, bool) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return workflow.Workflow{}, true
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error:", err)
		return workflow.Workflow{}, true
	}

	var wf workflow.Workflow
	err = json.Unmarshal(jsonData, &wf)
	if err != nil {
		fmt.Println("Error:", err)
		return workflow.Workflow{}, true
	}
	return wf, false
}