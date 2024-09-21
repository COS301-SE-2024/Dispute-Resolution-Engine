package workflow_test

import (
	// "encoding/json"
	// "fmt"
	// "net/http"
	// "net/http/httptest"
	"encoding/json"
	"fmt"
	"orchestrator/workflow"
	// "testing"
	// "github.com/stretchr/testify/assert"
)

type testData struct {
	ID       int64             `json:"id"`
	Workflow workflow.Workflow `json:"workflow"`
}
type testResponse struct {
	Data testData `json:"data"`
}

// func TestFetchBadResponse(t *testing.T) {
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintln(w, "")
// 	}))
// 	defer server.Close()

// 	_, err := workflow.FetchWorkflowFromAPI(server.URL, "bruh")
// 	assert.Error(t, err)
// }

// func TestFetchGoodResponse(t *testing.T) {
// 	response := testResponse{
// 		Data: testData{
// 			ID:       5,
// 			Workflow: workflow.Workflow{},
// 		},
// 	}
// 	secretKey := "secret"

// 	header := ""
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		header = r.Header.Get("X-Orchestrator-Key")
// 		data, _ := json.Marshal(response)
// 		w.Write(data)
// 	}))
// 	defer server.Close()

// 	workflow, err := workflow.FetchWorkflowFromAPI(server.URL, secretKey)
// 	assert.NoError(t, err)
// 	assert.Equal(t, response.Data.Workflow, *workflow)
// 	assert.Equal(t, header, secretKey)
// }

func manualTestStoreWorkflow(wf workflow.Workflow) {
	api := workflow.CreateAPIWorkflow()
	fmt.Println("storing workflow")
	// Store the workflow to the API
	err := api.StoreWorkflow("test", wf, []int64{}, 1)
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
	wf, err := api.FetchWorkflow(id)
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
	bob := "bob"
	err := api.UpdateWorkflow(id, &bob, &wf, nil, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Workflow updated successfully")
}

func maunalTestFetchActiveWorkflows() {
	api := workflow.CreateAPIWorkflow()
	fmt.Println("fetching active workflows")
	// Fetch the active workflows from the API
	activeWorkflows, err := api.FetchActiveWorkflows()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(activeWorkflows)
	if len(activeWorkflows) != 0 {
		for _, activeWorkflow := range activeWorkflows {
			var unmarshalledWorkflow workflow.Workflow
			err = json.Unmarshal(activeWorkflow.WorkflowDefinition, &unmarshalledWorkflow)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Println(unmarshalledWorkflow.GetWorkflowString())
		}
	}

	fmt.Println("Active workflows fetched successfully")
	fmt.Println(activeWorkflows)
}

func manualTestFetchActiveWorkflow(id int) {
	api := workflow.CreateAPIWorkflow()
	fmt.Println("fetching active workflow")
	// Fetch the active workflow from the API
	activeWorkflow, err := api.FetchActiveWorkflow(id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(activeWorkflow)
	var unmarshalledWorkflow workflow.Workflow
	err = json.Unmarshal(activeWorkflow.WorkflowDefinition, &unmarshalledWorkflow)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Active workflow fetched successfully")
	fmt.Println(activeWorkflow)
}
