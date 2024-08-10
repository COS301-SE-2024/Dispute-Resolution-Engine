package workflow

import (
	"encoding/json"
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

func stateToJSON(s state) map[string]interface{} {
	timers := make([]map[string]interface{}, 0, len(s.timers))
	for _, t := range s.timers {
		timers = append(timers, map[string]interface{}{
			"name":        t.GetName(),
			"duration":    t.GetDuration().String(),
			"willTrigger": t.WillTrigger(),
		})
	}

	return map[string]interface{}{
		"name":   s.GetName(),
		"timers": timers,
	}
}

// Convert transition to JSON-compatible format
func transitionToJSON(t transition) map[string]interface{} {
	return map[string]interface{}{
		"name":    t.GetName(),
		"from":    t.GetFrom(),
		"to":      t.GetTo(),
		"trigger": t.GetTrigger(),
	}
}

// json representation of the workflow
func WorkFlowToJSON(w *Workflow) (string, error) {

	convertStates := make([]map[string]interface{}, 0, len(w.states))
	for _, s := range w.states {
		convertStates = append(convertStates, stateToJSON(s))
	}

	convertTransitions := make([]map[string]interface{}, 0, len(w.transitions))
	for _, t := range w.transitions {
		convertTransitions = append(convertTransitions, transitionToJSON(t))
	}

	jsonWorkflow := map[string]interface{}{
		"id":          w.id,
		"name":        w.name,
		"initial":     w.initial.GetName(),
		"states":      convertStates,
		"transitions": convertTransitions,
	}

	//convert to json string
	jsonWorkflowJSON, err := json.Marshal(jsonWorkflow)
	if err != nil {
		return "", err
	}

	return string(jsonWorkflowJSON), nil
}

// Convert JSON to workflow
func JSONToWorkFlow(jsonWorkflow string) (*Workflow, error) {
	// Define a temporary structure to unmarshal the JSON data
	var temp struct {
		ID          uint32                   `json:"id"`
		Name        string                   `json:"name"`
		Initial     string                   `json:"initial"`
		States      []map[string]interface{} `json:"states"`
		Transitions []map[string]interface{} `json:"transitions"`
	}

	// Unmarshal JSON into the temporary structure
	err := json.Unmarshal([]byte(jsonWorkflow), &temp)
	if err != nil {
		return nil, err
	}

	// Create initial state
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
		stateName := s["name"].(string)
		newState := CreateState(stateName)

		timers := s["timers"].([]interface{})
		for _, t := range timers {
			timerMap := t.(map[string]interface{})
			duration, err := time.ParseDuration(timerMap["duration"].(string))
			if err != nil {
				return nil, err
			}
			newTimer := CreateTimer(timerMap["name"].(string), duration, timerMap["willTrigger"].(string))
			newState.AddTimer(newTimer)
		}

		w.AddState(newState)
	}

	// Add transitions to the workflow
	for _, t := range temp.Transitions {
		newTransition := CreateTransition(
			t["name"].(string),
			t["from"].(string),
			t["to"].(string),
			t["trigger"].(string),
		)
		w.AddTransition(newTransition)
	}

	return w, nil
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
