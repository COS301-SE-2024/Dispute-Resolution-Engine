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
	Start()
}

type stateMachine struct {
	label         string
	stateMachine *stateless.StateMachine
	workflow     workflow.Workflow
	scheduler    *scheduler.Scheduler
}

func NewStateMachine() IStateMachine {
	return &stateMachine{}
}

func (s *stateMachine) Init(wf_id string,wf workflow.Workflow, sch *scheduler.Scheduler) {
	logger := utilities.NewLogger().LogWithCaller()
	logger.Info("Initialising state machine")

	// initState := wf.GetInitialState() // this whole sequence is a bit weird, but idk how else to do it
	// initStatePtr := &initState        // without changing the workflow interface

	s.label = wf.Label
	s.stateMachine = stateless.NewStateMachine(wf.Initial)
	s.workflow = wf
	s.scheduler = sch // 1 second interval

	// For every state in the workflow, add it to the state machine
	for state_id, state := range wf.States {
		// For every related transition from the state, configure the state with the transition
		stateConfig := s.stateMachine.Configure(state_id)

		for trigger_id, trigger := range state.Triggers {
			stateConfig.Permit(trigger_id, trigger.Next)
		}
		
		// Configure timer states
		if timer := state.Timer; timer != nil {
			timerName := fmt.Sprintf("%s_%s",wf_id ,state_id)
			// Start the timer once the state is entered
			stateConfig.OnEntry(func(_ context.Context, args ...any) error {
				logger.Info("New state entered")
				s.scheduler.AddTimer(timerName, time.Now().Add(timer.GetDuration()), func() {
					logger.Info("Timer expired for state", state_id)
					s.stateMachine.Fire(timer.OnExpire)
				})
				return nil
			})
			// When the state is exited, remove the timer.
			// WARNING: this may cause some kind of race condition when the exit
			stateConfig.OnExit(func(_ context.Context, args ...any) error {
				s.scheduler.RemoveTimer(timerName)
				return nil
			})
		}
	}
}

func (s *stateMachine) Start() {
	logger := utilities.NewLogger().LogWithCaller()
	logger.Info("Starting state machine")
	initState := s.workflow.GetInitialState()
	initStatePtr := &initState
	transition := s.workflow.GetTransitionsByFrom(initStatePtr.GetName())[0]
	s.stateMachine.Fire(transition.GetTrigger())
}
