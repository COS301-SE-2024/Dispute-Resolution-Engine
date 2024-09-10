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
	type Alias State
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&s),
	})
}

func (s *State) UnmarshalJSON(b []byte) error {
	type Alias State
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	return json.Unmarshal(b, &aux)
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
	type Alias Transition
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&t),
	})
}

func (t *Transition) UnmarshalJSON(b []byte) error {
	type Alias Transition
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	return json.Unmarshal(b, &aux)
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

	// ID of the Initial state
	Initial string

	// Map of state IDs to the corresponding state info
	States map[string]State

	// Map of transition IDs to the corresponding transition info
	Transitions map[string]Transition
}

func (w Workflow) MarshalJSON() ([]byte, error) {
	type Alias Workflow
	return json.Marshal(&struct {
		ID          uint32                 `json:"id"`
		Name        string                 `json:"name"`
		Initial     string                 `json:"initial"`
		States      map[string]State       `json:"states"`
		Transitions map[string]Transition  `json:"transitions"`
		*Alias
	}{
		ID:          w.ID,
		Name:        w.Name,
		Initial:     w.Initial,
		States:      w.States,
		Transitions: w.Transitions,
		Alias:       (*Alias)(&w),
	})
}

func (w *Workflow) UnmarshalJSON(b []byte) error {
	type Alias Workflow
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(w),
	}
	return json.Unmarshal(b, &aux)
}

func (w *Workflow) AddState(s State) {
	w.States[s.ID] = s
}

func (w *Workflow) AddTransition(t Transition) {
	w.Transitions[t.ID] = t
}

func (w *Workflow) GetTransitions() []Transition {
	transitions := make([]Transition, 0, len(w.Transitions))
	for _, t := range w.Transitions {
		transitions = append(transitions, t)
	}
	return transitions
}

func (w *Workflow) GetTransitionsByFrom(state string) []Transition {
	if _, found := w.States[state]; !found {
		return nil
	}

	var transitions []Transition
	for _, t := range w.Transitions {
		if t.From == state {
			transitions = append(transitions, t)
		}
	}
	return transitions
}

func (w *Workflow) GetTransitionsByTo(state string) []Transition {
	if _, found := w.States[state]; !found {
		return nil
	}

	var transitions []Transition
	for _, t := range w.Transitions {
		if t.To == state {
			transitions = append(transitions, t)
		}
	}
	return transitions
}
