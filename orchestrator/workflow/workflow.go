package workflow

import (
	"errors"
	"time"
)

const (
	// Dispute states
	stateDisputeCreated         = "dispute_created"
	stateComplaintRectification = "complaint_rectification" // If complaint is not compliant
	stateDisputeFeeDue          = "dispute_fee_due"
	stateDisputeCommenced       = "dispute_commenced" // Notification to be sent to the other party
	stateResponseDue            = "response_due"
	stateResponseCommunique     = "response_communique" // Notification to be sent to the other party
	stateReplyDue               = "reply_due"
	stateAppointAdjudicator     = "appoint_adjudicator"
	stateNoReplyDecision        = "no_reply_decision"
	stateDecisionDue            = "decision_due"
	stateDecisionCommunique     = "decision_communique" // Communicate DECISION to the Complainant and Registrant
	stateFinalDecisionDue       = "final_decision_due"
	stateDisputeArchived        = "dispute_archived"

	// Appeal states
	stateAppealSubmitted    = "appeal_submitted"
	stateAppealNoticeAndFee = "appeal_notice_and_fee"
	stateAppealCommunique   = "appeal_communique"
	stateAppealReplyDue     = "appeal_reply_due"
	stateAppointAppealPanel = "appoint_appeal_panel"
	stateAppealDecisionDue  = "appeal_decision_due"
)

const (
	// Dispute triggers
	triggerComplaintNotCompliant = "complaint_not_compliant"
	triggerFeeNotPaid            = "fee_not_paid"
	triggerComplaintCompliant    = "complaint_compliant"
	triggerTimedOut              = "timed_out"
	triggerResponseReceived      = "response_received"
	triggerResponseUndelivered   = "response_undelivered"

	// Appeal triggers
	triggerAppealSubmitted  = "appeal_submitted"
	triggerAppealOmmission  = "appeal_ommission"
	triggerAppealFeeNotPaid = "appeal_fee_not_paid"
)

// ----------------------------Timers--------------------------------
type Timer struct {
	name        string
	duration    time.Duration
	willTrigger string
}

func (t *Timer) Name() string {
	return t.name
}

func (t *Timer) Duration() time.Duration {
	return t.duration
}

func (t *Timer) WillTrigger() string {
	return t.willTrigger
}

// ----------------------------States--------------------------------
type State struct {
	name   string
	timers map[string]Timer // timer name -> Timer
}

func (s *State) Name() string {
	return s.name
}

func (s *State) GetTimer(name string) Timer {
	return s.timers[name]
}

// ----------------------------Transitions--------------------------------
type Transition struct {
	name    string
	from    string
	to      string
	trigger string
}

func (t *Transition) Name() string {
	return t.name
}

func (t *Transition) From() string {
	return t.from
}

func (t *Transition) To() string {
	return t.to
}

func (t *Transition) Trigger() string {
	return t.trigger
}

// ----------------------------Workflow--------------------------------
type Workflow struct {
	id          uint32
	name        string
	initial     State
	states      map[string]State      // state name -> State
	transitions map[string]Transition // transition name -> Transition
}

func (w *Workflow) GetID() uint32 {
	return w.id
}

func (w *Workflow) GetName() string {
	return w.name
}

func (w *Workflow) GetInitialState() State {
	return w.initial
}

func (w *Workflow) HasState(name string) bool {
	_, ok := w.states[name]
	return ok
}

func (w *Workflow) AddTransition(t Transition) {
	w.transitions[t.Name()] = t
}

func (w *Workflow) GetTransitions() []Transition {
	transitions := make([]Transition, 0, len(w.transitions))
	for _, t := range w.transitions {
		transitions = append(transitions, t)
	}
	return transitions
}

func (w *Workflow) GetTransitionByTrigger(triggerstr string) (Transition, error) {
	for _, t := range w.transitions {
		if t.trigger == triggerstr {
			return t, nil
		}
	}
	return Transition{}, errors.New("Transition not found")
}

func (w *Workflow) GetTransitionByFrom(fromstr string) (Transition, error) {
	for _, t := range w.transitions {
		if t.from == fromstr {
			return t, nil
		}
	}
	return Transition{}, errors.New("Transition not found")
}

func (w *Workflow) GetTransitionByTo(tostr string) (Transition, error) {
	for _, t := range w.transitions {
		if t.to == tostr {
			return t, nil
		}
	}
	return Transition{}, errors.New("Transition not found")
}
