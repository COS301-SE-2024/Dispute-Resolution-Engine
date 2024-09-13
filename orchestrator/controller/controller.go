package controller

import (
	"os"
	"os/signal"
	"syscall"
	"orchestrator/scheduler"
	"orchestrator/utilities"
	"orchestrator/workflow"
	"orchestrator/statemachine"
	"time"
)

type Controller struct {
	StateMachineRegistry map[string]statemachine.IStateMachine
	logger utilities.Logger
	Scheduler *scheduler.Scheduler
}

// NewController creates a new controller instance
func NewController() *Controller {
	return &Controller{
		StateMachineRegistry: make(map[string]statemachine.IStateMachine),
		logger: *utilities.NewLogger().LogWithCaller(),
		Scheduler: scheduler.NewWithLogger(time.Second, utilities.NewLogger()),
	}
}

func (c *Controller) Start() {
	c.logger.Info("Starting controller...")
	c.logger.Info("Starting scheduler...")
	c.Scheduler.Start(make(chan struct{}))// idk why we need this struct chan...
}

// RegisterStateMachine registers a state machine with the controller AND starts it.
func (c *Controller) RegisterStateMachine(wfID string, wf workflow.Workflow) {
	// Create and initilise a new state machine
	sm := statemachine.NewStateMachine()
	sm.Init(wfID, wf, c.Scheduler)
	c.StateMachineRegistry[wfID] = sm
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