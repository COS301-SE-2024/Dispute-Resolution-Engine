package workflow

import (
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

// ----------------------------workflow--------------------------------
// Concrete product
type workflow struct {
	id          uint32 // from table primary key, ideally
	name        string
	initial     state
	states      map[string]state      // state name -> State
	transitions map[string]transition // transition name -> Transition
}

// Factory method
func CreateWorkflow(id uint32, name string, initial state) IWorkflow {
	w := &workflow{
		id:          id,
		name:        name,
		initial:     initial,
		states:      make(map[string]state),
		transitions: make(map[string]transition),
	}
	w.AddState(initial)
	return w
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

func (w *workflow) GetStates() []state {
	states := make([]state, 0, len(w.states))
	for _, s := range w.states {
		states = append(states, s)
	}
	return states
}

func (w *workflow) AddState(s state) {
	w.states[s.GetName()] = s
}

func (w *workflow) HasState(name string) bool {
	_, ok := w.states[name]
	return ok
}

func (w *workflow) GetTransition(name string) transition {
	return w.transitions[name]
}

func (w *workflow) AddTransition(t transition) {
	w.transitions[t.GetName()] = t
}

func (w *workflow) GetTransitions() []transition {
	transitions := make([]transition, 0, len(w.transitions))
	for _, t := range w.transitions {
		transitions = append(transitions, t)
	}
	return transitions
}

func (w *workflow) GetTransitionsByTrigger(triggerstr string) []transition {
	var transitions []transition
	for _, t := range w.transitions {
		if t.trigger == triggerstr {
			transitions = append(transitions, t)
		}
	}
	return transitions
}

func (w *workflow) GetTransitionsByFrom(fromstr string) []transition {
	var transitions []transition
	for _, t := range w.transitions {
		if t.from == fromstr {
			transitions = append(transitions, t)
		}
	}
	return transitions
}

func (w *workflow) GetTransitionsByTo(tostr string) []transition {
	var transitions []transition
	for _, t := range w.transitions {
		if t.to == tostr {
			transitions = append(transitions, t)
		}
	}
	return transitions
}
