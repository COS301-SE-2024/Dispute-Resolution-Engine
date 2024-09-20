package workflow_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"orchestrator/workflow"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	ID       int64             `json:"id"`
	Workflow workflow.Workflow `json:"workflow"`
}
type testResponse struct {
	Data testData `json:"data"`
}

func TestFetchBadResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	}))
	defer server.Close()

	_, err := workflow.FetchWorkflowFromAPI(server.URL, "bruh")
	assert.Error(t, err)
}

func TestFetchGoodResponse(t *testing.T) {
	response := testResponse{
		Data: testData{
			ID:       5,
			Workflow: workflow.Workflow{},
		},
	}
	secretKey := "secret"

	header := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header = r.Header.Get("X-Orchestrator-Key")
		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer server.Close()

	workflow, err := workflow.FetchWorkflowFromAPI(server.URL, secretKey)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Workflow, *workflow)
	assert.Equal(t, header, secretKey)
}