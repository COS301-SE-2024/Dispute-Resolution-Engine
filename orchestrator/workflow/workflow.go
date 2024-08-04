package workflow

import (
	"errors"
	"time"
)

const (
	// Dispute states
	StateDisputeCreated         = "dispute_created"
	StateComplaintRectification = "complaint_rectification" // If complaint is not compliant
	StateDisputeFeeDue          = "dispute_fee_due"
	StateDisputeCommenced       = "dispute_commenced" // Notification to be sent to the other party
	StateResponseDue            = "response_due"
	StateResponseCommunique     = "response_communique" // Notification to be sent to the other party
	StateReplyDue               = "reply_due"
	StateAppointAdjudicator     = "appoint_adjudicator"
	StateNoReplyDecision        = "no_reply_decision"
	StateDecisionDue            = "decision_due"
	StateDecisionCommunique     = "decision_communique" // Communicate DECISION to the Complainant and Registrant
	StateFinalDecisionDue       = "final_decision_due"
	StateDisputeArchived        = "dispute_archived"

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
	AddState(s state)
	HasState(name string) bool
	GetTransition(name string) transition
	AddTransition(t transition)
	GetTransitions() []transition
	GetTransitionByTrigger(triggerstr string) (transition, error)
	GetTransitionByFrom(fromstr string) (transition, error)
	GetTransitionByTo(tostr string) (transition, error)
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

func (t *timer) Name() string {
	return t.name
}

func (t *timer) Duration() time.Duration {
	return t.duration
}

func (t *timer) WillTrigger() string {
	return t.willTrigger
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
	s.timers[t.Name()] = t
}

func (s *state) Name() string {
	return s.name
}

func (s *state) GetTimer(name string) timer {
	return s.timers[name]
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

func (t *transition) Name() string {
	return t.name
}

func (t *transition) From() string {
	return t.from
}

func (t *transition) To() string {
	return t.to
}

func (t *transition) Trigger() string {
	return t.trigger
}

// ----------------------------workflow--------------------------------
// Concrete product
type workflow struct {
	id          uint32
	name        string
	initial     state
	states      map[string]state      // state name -> State
	transitions map[string]transition // transition name -> Transition
}

// Factory method
func CreateWorkflow(id uint32, name string, initial state) IWorkflow {
	return &workflow{
		id:          id,
		name:        name,
		initial:     initial,
		states:      make(map[string]state),
		transitions: make(map[string]transition),
	}
}

func (w *workflow) GetID() uint32 {
	return w.id
}

func (w *workflow) GetName() string {
	return w.name
}

func (w *workflow) GetInitialState() state {
	return w.initial
}

func (w *workflow) GetState(name string) state {
	return w.states[name]
}

func (w *workflow) AddState(s state) {
	w.states[s.Name()] = s
}

func (w *workflow) HasState(name string) bool {
	_, ok := w.states[name]
	return ok
}

func (w *workflow) GetTransition(name string) transition {
	return w.transitions[name]
}

func (w *workflow) AddTransition(t transition) {
	w.transitions[t.Name()] = t
}

func (w *workflow) GetTransitions() []transition {
	transitions := make([]transition, 0, len(w.transitions))
	for _, t := range w.transitions {
		transitions = append(transitions, t)
	}
	return transitions
}

func (w *workflow) GetTransitionByTrigger(triggerstr string) (transition, error) {
	for _, t := range w.transitions {
		if t.trigger == triggerstr {
			return t, nil
		}
	}
	return transition{}, errors.New("transition not found")
}

func (w *workflow) GetTransitionByFrom(fromstr string) (transition, error) {
	for _, t := range w.transitions {
		if t.from == fromstr {
			return t, nil
		}
	}
	return transition{}, errors.New("transition not found")
}

func (w *workflow) GetTransitionByTo(tostr string) (transition, error) {
	for _, t := range w.transitions {
		if t.to == tostr {
			return t, nil
		}
	}
	return transition{}, errors.New("transition not found")
}
