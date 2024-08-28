package statemachine

import (
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
	Init(workflow workflow.IWorkflow)
	Start()
}

type stateMachine struct {
	id 		     uint32
	name  	     string
	stateMachine *stateless.StateMachine
	workflow	 workflow.IWorkflow
	scheduler    *scheduler.Scheduler
}

func NewStateMachine() IStateMachine {
	return &stateMachine{}
}

func (s *stateMachine) Init(wf workflow.IWorkflow) {
	logger := utilities.NewLogger().LogWithCaller()
	logger.Info("Initialising state machine")
	initState := wf.GetInitialState() // this whole sequence is a bit weird, but idk how else to do it
	initStatePtr := &initState              // without changing the workflow interface
	s.id = wf.GetID()
	s.name = wf.GetName()
	s.stateMachine = stateless.NewStateMachine(initStatePtr.GetName())
	s.workflow = wf
	s.scheduler = scheduler.NewScheduler(1 * time.Second) // 1 second interval

	// for every state in the workflow, add it to the state machine
	for _, state := range wf.GetStates() {
		// for every related transition from the state, configure the state with the transition
		toTransitions := wf.GetTransitionsByFrom(state.GetName())
		for _, transition := range toTransitions {
			s.stateMachine.Configure(state.GetName()).
			Permit(transition.GetTrigger(), transition.GetTo())
		}
		
		// for every timer in the state, add it to the scheduler
		for _, timer := range state.GetTimers() {
			timerName := fmt.Sprintf("%s_%s", state.GetName(), timer.WillTrigger())
			s.scheduler.AddTimer(timerName, time.Now().Add(timer.GetDuration()))
		
			// Store the expected state when the timer is set
			expectedState := state.GetName()
		
			// Define the event function to trigger the state transition
			s.scheduler.AddEvent(timerName, func() {
				// Check if the current state is still the expected state
				is_expected_state, err := s.stateMachine.IsInState(expectedState)
				if err != nil {
					logger.Error("Error checking if state machine is in state", expectedState, err)
					return
				}
				if is_expected_state {
					logger.Info("Timer expired for state", state.GetName(), ", triggering transition:", timer.WillTrigger())
					transition := wf.GetTransition(timer.WillTrigger())
					s.stateMachine.Fire(transition.GetTrigger())
				} else {
					logger.Info("Timer expired but state has already transitioned from", expectedState)
				}
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