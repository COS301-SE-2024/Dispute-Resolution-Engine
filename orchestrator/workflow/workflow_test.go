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

