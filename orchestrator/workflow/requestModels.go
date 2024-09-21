package workflow

import "encoding/json"

//----request models----

type StoreWorkflowRequest struct {
	WorkflowDefinition Workflow `json:"workflow_definition,omitempty"`
	Category           []int64  `json:"category,omitempty"`
	Author             *int64   `json:"author,omitempty"`
}

type UpdateWorkflowRequest struct {
	WorkflowDefinition *Workflow `json:"workflow_definition,omitempty"`
	Category           *[]int64  `json:"category,omitempty"`
	Author             *int64    `json:"author,omitempty"`
}

//----response models----

type ActiveWorkflowsResponse struct {
	ID                 int64    `json:"id,omitempty"`
	WorkflowID         int64    `json:"workflow_id,omitempty"`
	WorkflowDefinition json.RawMessage `json:"workflow_definition,omitempty"`
	CurrentState       string   `json:"current_state,omitempty"`
	StateDeadline      string   `json:"state_deadline,omitempty"`
}
