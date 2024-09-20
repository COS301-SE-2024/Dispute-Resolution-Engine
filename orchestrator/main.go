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
	// "orchestrator/controller"
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

	wf2, shouldReturn := readWorkflowFromFile("templates/v2.json")
	if shouldReturn {
		return
	}

	fmt.Println(wf.GetWorkflowString())

	// manualTestStoreWorkflow(wf)
	// manualTestFetchWorkflow(1)

	manualTestUpdateWorkflow(1, wf2)

	// // ======== Statemachine Tests =========
	// // Create a new controller
	// c := controller.NewController()
	// c.Start()
	// // Register the workflow with the controller
	// c.RegisterStateMachine("wf1", wf)
	// // Wait for a signal to shutdown
	// c.WaitForSignal()
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


func manualTestStoreWorkflow(wf workflow.Workflow) {
	api := workflow.CreateAPIWorkflow()
	fmt.Println("storing workflow")
	// Store the workflow to the API
	err := api.Store("test",wf, []int64{}, 1)
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

	var unmarshalledWorkflow workflow.Workflow
	err = json.Unmarshal(wf.Definition, &unmarshalledWorkflow)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(unmarshalledWorkflow.GetWorkflowString())
	fmt.Println("Workflow fetched successfully")
}


func manualTestUpdateWorkflow(id int, wf workflow.Workflow) {
	api := workflow.CreateAPIWorkflow()
	fmt.Println("updating workflow")
	// Update the workflow in the API
	bob:= "bob"
	err := api.Update(id, &bob, &wf, nil, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Workflow updated successfully")
}