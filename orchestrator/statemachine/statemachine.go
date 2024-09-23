package statemachine

import (
	"context"
	"fmt"
	"orchestrator/scheduler"
	"orchestrator/utilities"
	"orchestrator/workflow"
	"time"

	"github.com/qmuntal/stateless"
)

// Statemachine consumes a workflow and orchestrates the transitions between states

type IStateMachine interface {
	// initialise the state machine with a workflow
	Init(wf_id string, workflow workflow.Workflow, scheduler *scheduler.Scheduler)
}

type StateMachine struct {
	Label                 string
	StatelessStateMachine *stateless.StateMachine
	Workflow              workflow.Workflow
	Scheduler             *scheduler.Scheduler
}

func NewStateMachine() IStateMachine {
	return &StateMachine{}
}

func (s *StateMachine) Init(wf_id string, wf workflow.Workflow, sch *scheduler.Scheduler) {
	logger := utilities.NewLogger()
	logger.Info("Initialising state machine")

	// initState := wf.GetInitialState() // this whole sequence is a bit weird, but idk how else to do it
	// initStatePtr := &initState        // without changing the workflow interface
	s.StatelessStateMachine = stateless.NewStateMachine(wf.Initial)
	s.Workflow = wf
	s.Scheduler = sch // 1 second interval

	// For every state in the workflow, add it to the state machine
	for state_id, state := range wf.States {
		// For every related transition from the state, configure the state with the transition
		stateConfig := s.StatelessStateMachine.Configure(state_id)

		for trigger_id, trigger := range state.Triggers {
			stateConfig.Permit(trigger_id, trigger.Next)
		}

		// Configure timer states
		if timer := state.Timer; timer != nil {
			timerName := fmt.Sprintf("%s_%s", wf_id, state_id)

			// If the current state is the initial state, start the timer
			if state_id == wf.Initial {
				s.Scheduler.AddTimer(timerName, time.Now().Add(timer.GetDuration()), func() {
					logger.Info("Timer expired for state", state_id)
					s.StatelessStateMachine.Fire(timer.OnExpire)
				})
			} else {
				// When the state is entered, start the timer
				stateConfig.OnEntry(func(_ context.Context, args ...any) error {
					logger.Info("New state entered")
					s.Scheduler.AddTimer(timerName, time.Now().Add(timer.GetDuration()), func() {
						logger.Info("Timer expired for state", state_id)
						s.StatelessStateMachine.Fire(timer.OnExpire)
					})
					return nil
				})
			}

			// When the state is exited, remove the timer.
			// WARNING: this may cause some kind of race condition when the exit
			stateConfig.OnExit(func(_ context.Context, args ...any) error {
				s.Scheduler.RemoveTimer(timerName)
				return nil
			})
		}
	}
	s.StatelessStateMachine.Fire(wf.Initial)
}
