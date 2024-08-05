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

func CountDownTimer(ctx context.Context, duration time.Duration) error {
    timer := time.NewTimer(duration)
    defer timer.Stop()

    select {
    case <-timer.C:
        return nil // timer expired
    case <-ctx.Done():
        return ctx.Err() // context cancelled or timed out
    }
}

func (s *StateMachine) Init(wf workflow.IWorkflow) {
	logger := utilities.NewLogger().LogWithCaller()
	logger.Info("Initialising state machine")
	initState := wf.GetInitialState() // this whole sequence is a bit weird, but idk how else to do it
	initStatePtr := &initState              // without changing the workflow interface
	s.stateMachine = stateless.NewStateMachine(initStatePtr.GetName())

	// for every state in the workflow, add it to the state machine
	for _, state := range wf.GetStates() {
		// for every related transition from the state, configure the state with the transition
		toTransitions := wf.GetTransitionsByFrom(state.GetName())
		for _, transition := range toTransitions {
			s.stateMachine.Configure(state.GetName()).
			Permit(transition.GetTrigger(), transition.GetTo())
		}
		
		// for every timer in the state, add it to the state machine
		for _, timer := range state.GetTimers() {
			s.stateMachine.Configure(state.GetName()).OnEntry(func(ctx context.Context, args ...any) error {
				timerCtx, cancel := context.WithTimeout(ctx, timer.GetDuration())
                defer cancel()

                err := CountDownTimer(timerCtx, timer.GetDuration()) // this is where the timer is actually started
                if err != nil {
                    if err == context.DeadlineExceeded {
                        logger.WithField("state", state.GetName()).Warn("Timer expired")
                        // Handle timer expiration, e.g., transition to a timeout state
                        s.stateMachine.Fire("timeoutEvent")
                    } else {
                        logger.WithField("state", state.GetName()).Warn("Context cancelled")
                    }
                }
                return nil
			})	
		}
	}
	s.stateMachine.Configure(workflow).
        OnEntry(func(ctx context.Context, args ...any) error {
            logger.Warn("Entered timeout state")
            // Handle timeout scenario
            return nil
        })
}

