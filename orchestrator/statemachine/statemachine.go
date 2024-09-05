package statemachine

import (
	// "fmt"
	"orchestrator/scheduler"
	"orchestrator/utilities"
	"orchestrator/workflow"
	// "time"

	"github.com/qmuntal/stateless"
)

// Statemachine consumes a workflow and orchestrates the transitions between states

type IStateMachine interface {
	// initialise the state machine with a workflow
	Init(workflow workflow.IWorkflow, scheduler *scheduler.Scheduler)
	Start()
}

type stateMachine struct {
	id           uint32
	name         string
	stateMachine *stateless.StateMachine
	workflow     workflow.IWorkflow
	scheduler    *scheduler.Scheduler
}

func NewStateMachine() IStateMachine {
	return &stateMachine{}
}

func (s *stateMachine) Init(wf workflow.IWorkflow, sch *scheduler.Scheduler) {
	logger := utilities.NewLogger().LogWithCaller()
	logger.Info("Initialising state machine")

	initState := wf.GetInitialState() // this whole sequence is a bit weird, but idk how else to do it
	initStatePtr := &initState        // without changing the workflow interface

	s.id = wf.GetID()
	s.name = wf.GetName()
	s.stateMachine = stateless.NewStateMachine(initStatePtr.GetName())
	s.workflow = wf
	s.scheduler = sch // 1 second interval

	// For every state in the workflow, add it to the state machine
	for _, state := range wf.GetStates() {
		// For every related transition from the state, configure the state with the transition
		toTransitions := wf.GetTransitionsByFrom(state.GetName())
		for _, transition := range toTransitions {
			s.stateMachine.Configure(state.GetName()).
				Permit(transition.GetTrigger(), transition.GetTo())
		}

		/*
			// For every timer in the state, add it to the scheduler
			for _, timer := range state.GetTimers() {

				timerName := fmt.Sprintf("%s_%s", state.GetName(), timer.WillTrigger())

				// Add the timer to the scheduler
				s.scheduler.AddTimer(timerName, time.Now().Add(timer.GetDuration()), func() {
					logger.Info("Timer expired for state", state.GetName(), ", triggering transition:", timer.WillTrigger())
					transition := wf.GetTransition(timer.WillTrigger())
					s.stateMachine.Fire(transition.GetTrigger())
				})
			}
		*/
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
