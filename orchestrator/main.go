package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	// "orchestrator/db"
	// "orchestrator/env"
	// "orchestrator/utilities"
	// "orchestrator/env"
	// "orchestrator/scheduler"
	// "orchestrator/statemachine"
	"orchestrator/controller"
	// "orchestrator/wf"
	"orchestrator/workflow"
)

func main() {
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

func manualTestUpdateActiveWorkflow(){
	api := workflow.CreateAPIWorkflow()
	fmt.Println("updating active workflow")
	// Update the active workflow in the API
	id := 1
	currentState := "new state"
	dateSubmitted := time.Now()
	stateDeadline := time.Now().Add(24 * time.Hour)

	err := api.UpdateActiveWorkflow(id, nil, &currentState, &dateSubmitted, &stateDeadline)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}