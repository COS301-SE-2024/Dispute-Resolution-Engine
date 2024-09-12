package workflow

import (
	"encoding/json"

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

func CreateState(label, description string) State {
	return State{
		Label:       label,
		Description: description,
		Triggers:    make(map[string]Trigger),
	}
}

func (s *State) AddTrigger(trigger Trigger) {
	s.Triggers[trigger.Label] = trigger
}

func (s *State) SetTimer(timer Timer) {
	s.Timer = &timer
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

