package controller

import (
	"os"
	"os/signal"
	"syscall"
	"orchestrator/scheduler"
	"orchestrator/utilities"
	"orchestrator/workflow"
	"orchestrator/statemachine"
)

type Controller struct {
	StateMachineRegistry map[uint32]statemachine.IStateMachine
	logger utilities.Logger
	Scheduler *scheduler.Scheduler
}

// NewController creates a new controller instance
func NewController() *Controller {
	return &Controller{
		StateMachineRegistry: make(map[uint32]statemachine.IStateMachine),
		logger: *utilities.NewLogger().LogWithCaller(),
		Scheduler: scheduler.NewWithLogger(1, utilities.NewLogger().LogWithCaller()),
	}
}

// RegisterStateMachine registers a state machine with the controller AND starts it.
func (c *Controller) RegisterStateMachine(wf workflow.IWorkflow) {
	// Create and initilise a new state machine
	sm := statemachine.NewStateMachine()
	sm.Init(wf, c.Scheduler)
	c.StateMachineRegistry[wf.GetID()] = sm
	go c.StateMachineRegistry[wf.GetID()].Start()
}

// WaitForSignal blocks until an appropriate signal is received
func (c *Controller) WaitForSignal() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    sig := <-sigs
    c.logger.Info("Received signal: ", sig)
    c.Shutdown()
}

// Shutdown gracefully shuts down the controller
func (c *Controller) Shutdown() {
    c.logger.Info("Shutting down controller")
	// Functions for ceasing all statemachines and the scheduler
}