package wf

import (
	"encoding/json"
	"time"
)

// ----------------------------Timers--------------------------------
type Timer struct {
	// How long the timer will run  for
	Duration time.Duration

	// The transition that will be triggered once the timer expires
	Trigger string
}

// Creates a new time with the passed-in duration and trigger
func NewTimer(duration time.Duration, trigger string) Timer {
	return Timer{duration, trigger}
}

// Because time.Duration cannot be marshalled, the simplest solution
// is to implement the custom marshaller for a timer
func (d Timer) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"duration": d.Duration.String(),
		"trigger":  d.Trigger,
	})
}

// Because time.Duration cannot be marshalled, the simplest solution
// is to implement the custom unmarshaller for a timer
func (d *Timer) UnmarshalJSON(b []byte) error {
	var result struct {
		Duration string `json:"duration"`
		Trigger  string `json:"trigger"`
	}
	if err := json.Unmarshal(b, &result); err != nil {
		return nil
	}
	dur, err := time.ParseDuration(result.Duration)
	if err != nil {
		return err
	}

	d.Duration = dur
	d.Trigger = result.Trigger
	return nil
}

// ----------------------------States--------------------------------
type State struct {
	// The unique name of the state
	ID string `json:"id"`

	// An optional timer associated with the state
	Timer *Timer `json:"timer,omitempty"`
}

func (s State) MarshalJSON() ([]byte, error) {
	timers, err := s.Timer.MarshalJSON();
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]string{
		"id":    s.ID,
		"timer": string(timers),
	})
}

func (s State) UnmarshalJSON(b []byte) (*State, error) {
	var result struct {
		ID    string `json:"id"`
		Timer Timer  `json:"timer"`
	}
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return &State{ID: result.ID, Timer: &result.Timer}, nil
}

// Initialises a new state with an empty timer
func NewState(id string) State {
	return State{ID: id, Timer: nil}
}

// Initialises a new state with the passed-in timer
func NewTimerState(id string, timer Timer) State {
	return State{ID: id, Timer: &timer}
}

// ----------------------------Transitions--------------------------------
type Transition struct {
	// Unique identifier of the transition
	ID string `json:"id"`

	// Name of the state that has transition as an outgoing one
	From string `json:"from"`

	// Name of the state that has transition as an incoming one
	To string `json:"to"`

	// Name of the event that will trigger the transition
	Trigger string `json:"trigger"`
}

func (t Transition) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"id":      t.ID,
		"from":    t.From,
		"to":      t.To,
		"trigger": t.Trigger,	
})
}

func (t Transition) UnmarshalJSON(b []byte) (*Transition, error) {
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func NewTransition(id, from, to, trigger string) Transition {
	return Transition{
		ID:      id,
		From:    from,
		To:      to,
		Trigger: trigger,
	}
}

// ----------------------------Workflow--------------------------------

type Workflow struct {
	// Unique identifier of a workflow (corresponds to ID in table)
	ID   uint32
	Name string

	// ID of the initial state
	initial string

	// Map of state IDs to the corresponding state info
	states map[string]State

	// Map of transition IDs to the corresponding transition info
	transitions map[string]Transition
}

func (w *Workflow) AddState(s State) {
	w.states[s.ID] = s
}

func (w *Workflow) AddTransition(t Transition) {
	w.transitions[t.ID] = t
}

func (w *Workflow) GetTransitions() []Transition {
	transitions := make([]Transition, 0, len(w.transitions))
	for _, t := range w.transitions {
		transitions = append(transitions, t)
	}
	return transitions
}

func (w *Workflow) GetTransitionsByFrom(state string) []Transition {
	if _, found := w.states[state]; !found {
		return nil
	}

	var transitions []Transition
	for _, t := range w.transitions {
		if t.From == state {
			transitions = append(transitions, t)
		}
	}
	return transitions
}

func (w *Workflow) GetTransitionsByTo(state string) []Transition {
	if _, found := w.states[state]; !found {
		return nil
	}

	var transitions []Transition
	for _, t := range w.transitions {
		if t.To == state {
			transitions = append(transitions, t)
		}
	}
	return transitions
}
