package workflow

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"orchestrator/env"

	// "orchestrator/env"
	"time"
)

const (
	// Dispute states
	StateDisputeCreated          = "dispute_created"
	StateComplaintRectification  = "complaint_rectification" // If complaint is not compliant
	StateDisputeFeeDue           = "dispute_fee_due"
	StateSuspended               = "suspended"
	StateDisputeCommenced        = "dispute_commenced" // Notification to be sent to the other party
	StateResponseDue             = "response_due"
	StateResponseCommunique      = "response_communique" // Notification to be sent to the other party
	StateReplyDue                = "reply_due"
	StateAppointAdjudicator      = "appoint_adjudicator"
	StateNoReplyDecision         = "no_reply_decision"
	StateDecisionDue             = "decision_due"
	StateDecisionCommunique      = "decision_communique" // Communicate DECISION to the Complainant and Registrant
	StateFinalDecisionCommunique = "final_decision_communique"
	StateDisputeArchived         = "dispute_archived"

	// Appeal states
	StateAppealSubmitted    = "appeal_submitted"
	StateAppealNoticeAndFee = "appeal_notice_and_fee"
	StateAppealCommunique   = "appeal_communique"
	StateAppealReplyDue     = "appeal_reply_due"
	StateAppointAppealPanel = "appoint_appeal_panel"
	StateAppealDecisionDue  = "appeal_decision_due"
)

const (
	// Dispute triggers
	TriggerComplaintNotCompliant = "complaint_not_compliant"
	TriggerFeeNotPaid            = "fee_not_paid"
	TriggerComplaintCompliant    = "complaint_compliant"
	TriggerTimedOut              = "timed_out"
	TriggerResponseReceived      = "response_received"
	TriggerResponseUndelivered   = "response_undelivered"
	TriggerNoAppeal              = "no_appeal"

	// Appeal triggers
	TriggerAppealSubmitted  = "appeal_submitted"
	TriggerAppealOmmission  = "appeal_ommission"
	TriggerAppealFeeNotPaid = "appeal_fee_not_paid"
)

// ----------------------------Timers--------------------------------
type Timer struct {
	// The duration that the timer should run for
	Duration durationWrapper `json:"duration"`

	// The ID of the trigger to fire when the timer expires
	OnExpire string `json:"on_expire"`
}

// Because time.Duration is not marshallable, we need to introduce
// a wrapper so that we can implement that ourselves
type durationWrapper struct {
	time.Duration
}

func CreateTimer(duration time.Duration, onExpire string) Timer {
	return Timer{Duration: durationWrapper{duration}, OnExpire: onExpire}
}

func (t *Timer) GetDuration() time.Duration {
	return t.Duration.Duration
}

func (t *Timer) GetDeadline() time.Time {
	return time.Now().Add(t.Duration.Duration)
}

func (d durationWrapper) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *durationWrapper) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	dur, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	d.Duration = dur
	return nil
}

// ----------------------------States--------------------------------
type State struct {
	// Human-readable label of the state
	Label string `json:"label"`

	// Human-readable description of the state, describing what the state means and all
	// the triggers that go from this state
	Description string `json:"description"`

	// All the outgoing triggers of the state, keyed by their IDs
	Triggers map[string]Trigger `json:"triggers,omitempty"`

	// The optional timer associated with a state
	Timer *Timer `json:"timer,omitempty"`
}

// ----------------------------Trigger--------------------------------
type Trigger struct {
	// Human-readable label of the trigger
	Label string `json:"label"`

	// The ID of the next state to transition to
	Next string `json:"next_state"`
}

func NewTrigger(label, next string) Trigger {
	return Trigger{Label: label, Next: next}
}

// ----------------------------Workflow--------------------------------
type Workflow struct {
	// The human-readable label for the workflow
	Label string `json:"label"`

	// The ID of the initial state of the workflow
	Initial string `json:"initial"`

	// All the states in the workflow, keyd by their ID
	States map[string]State `json:"states"`
}

// Factory method
func CreateWorkflow(label, initialId string, initial State) Workflow {
	w := Workflow{
		Label:   label,
		Initial: initialId,
		States:  make(map[string]State),
	}
	w.States[initialId] = initial
	return w
}

func (w *Workflow) GetInitialState() State {
	return w.States[w.Initial]
}

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

func FetchWorkflowFromAPI(apiURL string) (*Workflow, error) {
	// Create a new GET request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	key, err := env.Get("ORCHESTRATOR_KEY")
	if err != nil {
		return nil, err
	}

	// Set the X-Orchestrator-Key header
	req.Header.Set("X-Orchestrator-Key", key)

	// Perform the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch workflow: " + resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Define a temporary structure to extract the data field and ID
	var responseData struct {
		Data struct {
			ID                 int      `json:"ID"`
			WorkflowDefinition Workflow `json:"WorkflowDefinition"`
		} `json:"data"`
	}

	// Unmarshal the JSON response to extract the "data" field
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData.Data.WorkflowDefinition, nil
}

func StoreWorkflowToAPI(apiURL string, workflow Workflow, categories []int64, Author *int64) error {
	store := StoreWorkflowRequest{
		WorkflowDefinition: workflow,
		Category:           categories,
		Author:             Author,
	}
	storeJSON, err := json.Marshal(store)
	if err != nil {
		return err
	}

	// Create a new POST request with the workflow JSON as the body
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(storeJSON))
	if err != nil {
		return err
	}

	// Set the appropriate content-type header
	key, err := env.Get("ORCHESTRATOR_KEY")
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Orchestrator-Key", key)

	// Perform the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("failed to store workflow: " + resp.Status)
	}

	return nil
}

func UpdateWorkflowToAPI(apiURL string, workflow *Workflow, categories *[]int64, author *int64) error {
	// Prepare the update request structure
	var update UpdateWorkflowRequest
	update.WorkflowDefinition = workflow

	// Add categories if provided
	if categories != nil {
		update.Category = categories
	}

	// Add author if provided
	if author != nil {
		update.Author = author
	}

	// Marshal the update request object to JSON
	updateJSON, err := json.Marshal(update)
	if err != nil {
		return err
	}

	// Create a new PUT request with the update JSON as the body
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(updateJSON))
	if err != nil {
		return err
	}

	// Set the appropriate headers
	key, err := env.Get("ORCHESTRATOR_KEY")
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Orchestrator-Key", key)

	// Perform the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errors.New("failed to update workflow: " + resp.Status)
	}

	return nil
}
