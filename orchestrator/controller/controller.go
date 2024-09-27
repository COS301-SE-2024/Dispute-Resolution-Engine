package controller

import (
	"orchestrator/scheduler"
	"orchestrator/statemachine"
	"orchestrator/utilities"
	"orchestrator/workflow"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Controller struct {
	StateMachineRegistry map[string]statemachine.IStateMachine
	logger               utilities.Logger
	Scheduler            *scheduler.Scheduler
}

// NewController creates a new controller instance
func NewController() *Controller {
	return &Controller{
		StateMachineRegistry: make(map[string]statemachine.IStateMachine),
		logger:               *utilities.NewLogger(),
		Scheduler:            scheduler.NewWithLogger(time.Second, utilities.NewLogger()),
	}
}

// Start starts the controller (and scheduler)
func (c *Controller) Start() {
	c.logger.Info("Starting controller...")
	c.logger.Info("Starting scheduler...")
	c.Scheduler.Start(make(chan struct{})) // idk why we need this struct chan...
}

// RegisterStateMachine registers a state machine with the controller AND starts it.
func (c *Controller) RegisterStateMachine(wfID string, wf workflow.Workflow) {
	// Create and initilise a new state machine
	sm := statemachine.NewStateMachine()
	// Remove any active timers for the state machine
	for state_id := range wf.States {
		timerName := wfID + "_" + state_id
		c.Scheduler.RemoveTimer(timerName)
	}
	sm.Init(wfID, wf, c.Scheduler)
	c.StateMachineRegistry[wfID] = sm
}

// Transition to the next state in the specified state machine
func (c *Controller) FireTrigger(wfID string, trigger string) (string, time.Time) {
	/*
	This function drives the transition from state to state given a trigger casued by an event
	*/
	// Check if the state machine exists
	_, ok := c.StateMachineRegistry[wfID]
	if !ok {
		c.logger.Error("No state machine found for workflow ", wfID)
		return "", time.Time{}
	}

	// Fire the state's trigger
	err := c.StateMachineRegistry[wfID].TriggerTransition(trigger)
	if err != nil {
		c.logger.Error("Error firing trigger: ", err)
		return "", time.Time{}
	}

	// NOW we return the new state and deadline

	// Get the current state of the state machine
	current_state, err := c.StateMachineRegistry[wfID].GetCurrentState()
	if err != nil {
		c.logger.Error("Error getting current state: ", err)
		return "", time.Time{}
	}

	// Get the deadline of the current state
	stateDeadline, err := c.StateMachineRegistry[wfID].GetStateDeadline()
	// If timer is empty but no error, then teh state just doesn't have a timer
	if err != nil {
		c.logger.Error("Error getting state deadline: ", err)
		return "", time.Time{}
	}

	return current_state, stateDeadline
}

// WaitForSignal blocks until an appropriate signal is received
func (c *Controller) WaitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	c.logger.Info("Received signal: ", sig)
	c.shutdown()
}

// shutdown gracefully shuts down the controller
func (c *Controller) shutdown() {
	c.logger.Info("Shutting down controller")
	// Functions for ceasing all statemachines and the scheduler
}
