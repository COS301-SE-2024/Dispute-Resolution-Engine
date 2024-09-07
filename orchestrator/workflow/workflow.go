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

// Product Interface
type IWorkflow interface {
	GetID() uint32
	GetName() string
	GetInitialState() state
	GetState(name string) state
	GetStates() []state
	AddState(s state)
	HasState(name string) bool
	GetTransition(name string) transition
	AddTransition(t transition)
	GetTransitions() []transition
	GetTransitionsByTrigger(triggerstr string) []transition
	GetTransitionsByFrom(fromstr string) []transition
	GetTransitionsByTo(tostr string) []transition
}

// ----------------------------Timers--------------------------------
type timer struct {
	name        string
	duration    time.Duration
	willTrigger string
}

func CreateTimer(name string, duration time.Duration, willTrigger string) timer {
	return timer{name: name, duration: duration, willTrigger: willTrigger}
}

func (t *timer) GetName() string {
	return t.name
}

func (t *timer) GetDuration() time.Duration {
	return t.duration
}

func (t *timer) WillTrigger() string {
	return t.willTrigger
}

func (t *timer) GetDeadline() time.Time {
	return time.Now().Add(t.duration)
}

func (t *timer) HasDeadlinePassed() bool {
	return time.Now().After(t.GetDeadline())
}

// ----------------------------States--------------------------------
type state struct {
	name   string
	timers map[string]timer // timer name -> Timer
}

func CreateState(name string) state {
	return state{name: name, timers: make(map[string]timer)}
}

func (s *state) AddTimer(t timer) {
	s.timers[t.GetName()] = t
}

func (s *state) GetName() string {
	return s.name
}

func (s *state) GetTimer(name string) (timer, bool) {
	t, ok := s.timers[name]
	return t, ok
}

func (s *state) GetTimers() []timer {
	timers := make([]timer, 0, len(s.timers))
	for _, t := range s.timers {
		timers = append(timers, t)
	}
	return timers
}

// ----------------------------Transitions--------------------------------
type transition struct {
	name    string
	from    string
	to      string
	trigger string
}

func CreateTransition(name string, from string, to string, trigger string) transition {
	return transition{name: name, from: from, to: to, trigger: trigger}
}

func (t *transition) GetName() string {
	return t.name
}

func (t *transition) GetFrom() string {
	return t.from // name of state
}

func (t *transition) GetTo() string {
	return t.to // name of state
}

func (t *transition) GetTrigger() string {
	return t.trigger
}

// ----------------------------Workflow--------------------------------
// Concrete product
type Workflow struct {
	id          uint32 // from table primary key, ideally
	name        string
	initial     state
	states      map[string]state      // state name -> State
	transitions map[string]transition // transition name -> Transition
}

// Factory method
func CreateWorkflow(id uint32, name string, initial state) IWorkflow {
	w := &Workflow{
		id:          id,
		name:        name,
		initial:     initial,
		states:      make(map[string]state),
		transitions: make(map[string]transition),
	}
	w.AddState(initial)
	return w
}



// json representation of the workflow
// json representation of the workflow
func WorkFlowToJSON(w *Workflow) (string, error) {
	// Convert states to JSON
	convertStates := make([]json.RawMessage, 0, len(w.states))
	for _, s := range w.states {
		stateJSON, err := stateToJSON(s)
		if err != nil {
			return "", err
		}
		convertStates = append(convertStates, stateJSON)
	}

	// Convert transitions to JSON
	convertTransitions := make([]json.RawMessage, 0, len(w.transitions))
	for _, t := range w.transitions {
		transitionJSON, err := transitionToJSON(t)
		if err != nil {
			return "", err
		}
		convertTransitions = append(convertTransitions, transitionJSON)
	}

	// Create the final JSON representation of the workflow
	jsonWorkflow := map[string]interface{}{
		"id":          w.id,
		"name":        w.name,
		"initial":     w.initial.GetName(),
		"states":      convertStates,
		"transitions": convertTransitions,
	}

	// Convert to JSON string
	jsonWorkflowJSON, err := json.Marshal(jsonWorkflow)
	if err != nil {
		return "", err
	}

	return string(jsonWorkflowJSON), nil
}


type TimerJSON struct {
	Name        string `json:"name"`
	Duration    string `json:"duration"`
	WillTrigger string `json:"willTrigger"`
}

type StateJSON struct {
	Name   string      `json:"name"`
	Timers []TimerJSON `json:"timers"`
}

type TransitionJSON struct {
	Name    string `json:"name"`
	From    string `json:"from"`
	To      string `json:"to"`
	Trigger string `json:"trigger"`
}

type WorkflowJSON struct {
	ID          uint32           `json:"id"`
	Name        string           `json:"name"`
	Initial     string           `json:"initial"`
	States      []StateJSON      `json:"states"`
	Transitions []TransitionJSON `json:"transitions"`
}


func JSONToWorkFlow(jsonWorkflow string) (*Workflow, error) {
	// Unmarshal the JSON into the WorkflowJSON struct
	var temp WorkflowJSON
	err := json.Unmarshal([]byte(jsonWorkflow), &temp)
	if err != nil {
		return nil, err
	}

	// Create the initial state
	initialState := CreateState(temp.Initial)

	// Create a new workflow
	w := &Workflow{
		id:          temp.ID,
		name:        temp.Name,
		initial:     initialState,
		states:      make(map[string]state),
		transitions: make(map[string]transition),
	}

	// Add states to the workflow
	for _, s := range temp.States {
		newState := CreateState(s.Name)

		for _, t := range s.Timers {
			duration, err := time.ParseDuration(t.Duration)
			if err != nil {
				return nil, err
			}
			newTimer := CreateTimer(t.Name, duration, t.WillTrigger)
			newState.AddTimer(newTimer)
		}

		w.AddState(newState)
	}

	// Add transitions to the workflow
	for _, t := range temp.Transitions {
		newTransition := CreateTransition(t.Name, t.From, t.To, t.Trigger)
		w.AddTransition(newTransition)
	}

	return w, nil
}

func stateToJSON(s state) ([]byte, error) {
	// Convert state timers to TimerJSON objects
	timers := make([]TimerJSON, 0, len(s.timers))
	for _, t := range s.timers {
		timers = append(timers, TimerJSON{
			Name:        t.GetName(),
			Duration:    t.GetDuration().String(),
			WillTrigger: t.WillTrigger(),
		})
	}

	// Create the StateJSON object
	stateJSON := StateJSON{
		Name:   s.GetName(),
		Timers: timers,
	}

	// Marshal the StateJSON object to JSON
	return json.Marshal(stateJSON)
}




// Convert transition to JSON-compatible format
func transitionToJSON(t transition) ([]byte, error) {
	// Create the TransitionJSON object
	transitionJSON := TransitionJSON{
		Name:    t.GetName(),
		From:    t.GetFrom(),
		To:      t.GetTo(),
		Trigger: t.GetTrigger(),
	}
	return json.Marshal(transitionJSON)
}



func JSONToMap(jsonStr string) (map[string]interface{}, error) {
	// Initialize an empty map to hold the JSON structure
	result := make(map[string]interface{})

	// Unmarshal the JSON string into the map
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (w *Workflow) GetID() uint32 {
	return w.id
}

func (w *Workflow) GetName() string {
	return w.name
}

func (w *Workflow) GetInitialState() state {
	return w.initial
}

func (w *Workflow) GetState(name string) state {
	return w.states[name]
}

func (w *Workflow) GetStates() []state {
	states := make([]state, 0, len(w.states))
	for _, s := range w.states {
		states = append(states, s)
	}
	return states
}

func (w *Workflow) AddState(s state) {
	w.states[s.GetName()] = s
}

func (w *Workflow) HasState(name string) bool {
	_, ok := w.states[name]
	return ok
}

func (w *Workflow) GetTransition(name string) transition {
	return w.transitions[name]
}

func (w *Workflow) AddTransition(t transition) {
	w.transitions[t.GetName()] = t
}

func (w *Workflow) GetTransitions() []transition {
	transitions := make([]transition, 0, len(w.transitions))
	for _, t := range w.transitions {
		transitions = append(transitions, t)
	}
	return transitions
}

func (w *Workflow) GetTransitionsByTrigger(triggerstr string) []transition {
	var transitions []transition
	for _, t := range w.transitions {
		if t.trigger == triggerstr {
			transitions = append(transitions, t)
		}
	}
	return transitions
}

func (w *Workflow) GetTransitionsByFrom(fromstr string) []transition {
	var transitions []transition
	for _, t := range w.transitions {
		if t.from == fromstr {
			transitions = append(transitions, t)
		}
	}
	return transitions
}

func (w *Workflow) GetTransitionsByTo(tostr string) []transition {
	var transitions []transition
	for _, t := range w.transitions {
		if t.to == tostr {
			transitions = append(transitions, t)
		}
	}
	return transitions
}

type StoreWorkflowRequest struct {
	WorkflowDefinition map[string]interface{} `json:"workflow_definition,omitempty"`
	Category           []int64                `json:"category,omitempty"`
	Author             *int64                 `json:"author,omitempty"`
}

type UpdateWorkflowRequest struct {
	WorkflowDefinition *map[string]interface{} `json:"workflow_definition,omitempty"`
	Category           *[]int64                `json:"category,omitempty"`
	Author             *int64                  `json:"author,omitempty"`
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
			ID                 int             `json:"ID"`
			WorkflowDefinition json.RawMessage `json:"WorkflowDefinition"`
		} `json:"data"`
	}

	// Unmarshal the JSON response to extract the "data" field
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	// Extract the ID field into a variable
	id := responseData.Data.ID

	// Convert the JSON response to a Workflow object using your existing function
	workflow, err := JSONToWorkFlow(string(responseData.Data.WorkflowDefinition))
	if err != nil {
		return nil, err
	}

	workflow.id = uint32(id)
	// Optionally, you can set the ID on the Workflow object
	// workflow.ID = id

	return workflow, nil
}

func StoreWorkflowToAPI(apiURL string, workflow IWorkflow, categories []int64, Author *int64) error {
	// Convert the workflow to JSON string
	workflowJSON, err := WorkFlowToJSON(workflow.(*Workflow))
	if err != nil {
		return err
	}

	workflowMap, err := JSONToMap(workflowJSON)
	if err != nil {
		return err
	}

	store := StoreWorkflowRequest{
		WorkflowDefinition: workflowMap,
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

	// Convert the workflow to a JSON map only if it's provided
	if workflow != nil {
		workflowJSON, err := WorkFlowToJSON(workflow)
		if err != nil {
			return err
		}

		workflowMap, err := JSONToMap(workflowJSON)
		if err != nil {
			return err
		}

		update.WorkflowDefinition = &workflowMap
	}

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
