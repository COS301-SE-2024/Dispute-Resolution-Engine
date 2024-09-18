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
	// "orchestrator/wf"
	"orchestrator/workflow"
)

var RequiredEnvVariables = []string{
	// PostGres-related variables
	"DATABASE_URL",
	"DATABASE_PORT",
	"DATABASE_USER",
	"DATABASE_PASSWORD",
	"DATABASE_NAME",
	"ORCHESTRATOR_KEY",
}

func main() {
	// ======== Json Tests =========
	//read template workflow form json file
	file, err := os.Open("templates/v2.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Read the JSON file
	jsonData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var wf workflow.Workflow
	err = json.Unmarshal(jsonData, &wf)
	if err != nil {
		fmt.Println("Error:", err)
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


func manualTestStoreWorkflow(wf workflow.Workflow) {
	api := workflow.CreateAPIWorkflow()
	fmt.Println("storing workflow")
	// Store the workflow to the API
	err := api.Store(wf, []int64{}, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Workflow stored successfully")
}

func manualTestFetchWorkflow(id int) {
	api := workflow.CreateAPIWorkflow()
	fmt.Println("fetching workflow")
	// Fetch the workflow from the API
	wf, err := api.Fetch(id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(wf.GetWorkflowString())
	fmt.Println("Workflow fetched successfully")
}
