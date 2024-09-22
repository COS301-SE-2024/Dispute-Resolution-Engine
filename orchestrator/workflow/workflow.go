package workflow

import (
	"encoding/json"
	"fmt"
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
type TimerInterface interface {
	GetDuration() time.Duration
	SetDuration(duration time.Duration)
	GetDeadline() time.Time
}

type DurationWrapperInterface interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
}

type Timer struct {
	// The duration that the timer should run for
	Duration DurationWrapper `json:"duration"`

	// The ID of the trigger to fire when the timer expires
	OnExpire string `json:"on_expire"`
}

// Because time.Duration is not marshallable, we need to introduce
// a wrapper so that we can implement that ourselves
type DurationWrapper struct {
	time.Duration
}

func CreateTimer(duration time.Duration, onExpire string) Timer {
	return Timer{Duration: DurationWrapper{duration}, OnExpire: onExpire}
}

func (t *Timer) GetDuration() time.Duration {
	return t.Duration.Duration
}

func (t *Timer) SetDuration(duration time.Duration) {
	t.Duration = DurationWrapper{duration}
}

func (t *Timer) GetDeadline() time.Time {
	return time.Now().Add(t.Duration.Duration)
}

func (d DurationWrapper) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *DurationWrapper) UnmarshalJSON(b []byte) error {
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
type StateInterface interface {
	AddTrigger(trigger Trigger)
	SetTimer(timer Timer)
}


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
type TriggerInterface interface {
}

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
type WorkflowInterface interface {
	GetInitialState() State
	GetWorkflowString() string
}


type Workflow struct {
	// The ID of the initial state of the workflow
	Initial string `json:"initial"`

	// All the states in the workflow, keyd by their ID
	States map[string]State `json:"states"`
}

// Factory method
func CreateWorkflow(initialId string, initial State) Workflow {
	w := Workflow{
		Initial: initialId,
		States:  make(map[string]State),
	}
	w.States[initialId] = initial
	return w
}

func (w *Workflow) GetInitialState() State {
	return w.States[w.Initial]
}


func (w *Workflow) GetWorkflowString() string {
	result := fmt.Sprintf("Initial State: %s\n", w.Initial)

	// Iterate through each state in the workflow
	for stateID, state := range w.States {
		result += fmt.Sprintf("\nState ID: %s\n", stateID)
		result += fmt.Sprintf("  Label: %s\n", state.Label)
		result += fmt.Sprintf("  Description: %s\n", state.Description)

		// Print triggers
		if len(state.Triggers) > 0 {
			result += "  Triggers:\n"
			for triggerID, trigger := range state.Triggers {
				result += fmt.Sprintf("    - ID: %s, Label: %s, Next State: %s\n", triggerID, trigger.Label, trigger.Next)
			}
		} else {
			result += "  No Triggers\n"
		}

		// Print timer if exists
		if state.Timer != nil {
			result += fmt.Sprintf("  Timer: Duration: %s, On Expire: %s\n", state.Timer.GetDuration().String(), state.Timer.OnExpire)
		} else {
			result += "  No Timer\n"
		}
	}

	return result
}
