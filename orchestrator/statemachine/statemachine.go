package statemachine

import (
	"context"
	"orchestrator/utilities"
	"orchestrator/workflow"
	"time"

	"github.com/qmuntal/stateless"
)

// Statemachine consumes a workflow and orchestrates the transitions between states

type IStateMachine interface {
	// initialise the state machine with a workflow
	Init(workflow workflow.IWorkflow)
}

type StateMachine struct {
	stateMachine *stateless.StateMachine
}

func (s *StateMachine) Init(workflow workflow.IWorkflow) {
	logger := utilities.NewLogger().LogWithCaller()

	initState := workflow.GetInitialState() // this whole sequence is a bit weird, but idk how else to do it
	initStatePtr := &initState              // without changing the workflow interface
	s.stateMachine = stateless.NewStateMachine(initStatePtr.GetName())

	// for every state in the workflow, add it to the state machine
	for _, state := range workflow.GetStates() {

		// for every related transition from the state, configure the state with the transition
		toTransitions := workflow.GetTransitionsByFrom(state.GetName())
		for _, transition := range toTransitions {
			s.stateMachine.Configure(state.GetName()).
				Permit(transition.GetTrigger(), transition.GetTo())
		}
		
		// for every timer in the state, add it to the state machine
		for _, timer := range state.GetTimers() {
			s.stateMachine.Configure(state.GetName()).
			OnEntry(func(ctx context.Context, args ...any) error {
					time.AfterFunc(timer.GetDuration(), func() {
						logger.WithField("Workflow", workflow.GetID()).Info("Timer expired")
						
					})
				return nil
			})
		}
	}
}

