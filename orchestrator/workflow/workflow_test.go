package workflow_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"orchestrator/workflow"
)

func TestCreateTimer(t *testing.T) {
	duration := 10 * time.Second
	onExpire := "trigger_event"
	
	tmr := workflow.CreateTimer(duration, onExpire)

	assert.Equal(t, duration, tmr.GetDuration(), "The duration should be 10s")
	assert.Equal(t, onExpire, tmr.OnExpire, "The OnExpire trigger should match the provided value")
}

func TestSetDuration(t *testing.T) {
	tmr := workflow.CreateTimer(5*time.Second, "trigger_event")
	
	newDuration := 20 * time.Second
	tmr.SetDuration(newDuration)

	assert.Equal(t, newDuration, tmr.GetDuration(), "The duration should be updated to 20s")
}

func TestGetDeadline(t *testing.T) {
	duration := 5 * time.Second
	tmr := workflow.CreateTimer(duration, "trigger_event")
	
	deadline := tmr.GetDeadline()

	// The deadline should be roughly now + 5 seconds, we add some tolerance
	assert.WithinDuration(t, time.Now().Add(duration), deadline, time.Millisecond*100, "The deadline should be current time + 5s")
}

