package statemachine_test

// import (
// 	"testing"
// 	"time"

// 	"orchestrator/scheduler"
// 	"orchestrator/statemachine"
// 	"orchestrator/workflow"

// 	"github.com/stretchr/testify/assert"
// )

// func TestStateMachineInit(t *testing.T) {
// 	// Create a mock workflow
// 	initialState := workflow.CreateState("Initial", "Initial state description")
// 	state2 := workflow.CreateState("State2", "Second state description")
// 	trigger := workflow.NewTrigger("to_state2", "state2")
// 	initialState.AddTrigger(trigger)
// 	timer := workflow.CreateTimer(10*time.Second, "timer_expired")
// 	state2.SetTimer(timer)

// 	wf := workflow.CreateWorkflow("initial_state", initialState)
// 	wf.States["state2"] = state2

// 	// Create a mock scheduler
// 	sch := scheduler.NewScheduler(1 * time.Second)

// 	// Initialize the state machine
// 	sm := statemachine.NewStateMachine()
// 	sm.Init("test_wf", wf, sch)

// 	// Check if the state machine is initialized correctly
// 	assert.NotNil(t, sm, "State machine should be initialized")
// 	assert.Equal(t, wf, sm.(*statemachine.stateMachine).workflow, "Workflow should be set correctly")
// 	assert.Equal(t, sch, sm.(*statemachine.stateMachine).scheduler, "Scheduler should be set correctly")

// 	// Check if the initial state is set correctly
// 	assert.Equal(t, "initial_state", sm.(*statemachine.stateMachine).stateMachine.State(), "Initial state should be set correctly")

// 	// Check if the state transitions are configured correctly
// 	stateConfig := sm.(*statemachine.stateMachine).stateMachine.Configure("initial_state")
// 	assert.NotNil(t, stateConfig, "State configuration should be set for initial state")
// 	assert.Contains(t, stateConfig.Triggers(), "to_state2", "Initial state should have the correct trigger")

// 	// Check if the timer is set correctly
// 	timerName := "test_wf_state2"
// 	sch.AddTimer(timerName, time.Now().Add(timer.GetDuration()), func() {
// 		sm.(*statemachine.stateMachine).stateMachine.Fire(timer.OnExpire)
// 	})
// 	assert.Contains(t, sch.Timers(), timerName, "Timer should be set for state2")
// }
