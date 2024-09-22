package workflow_test

import (
	"encoding/json"
	"testing"
	"time"

	"orchestrator/workflow"

	"github.com/stretchr/testify/assert"
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

func TestMarshalDurationWrapper(t *testing.T) {
	duration := 15 * time.Second
	durWrapper := workflow.DurationWrapper{Duration: duration}

	marshaled, err := json.Marshal(durWrapper)
	assert.NoError(t, err, "Marshaling should not return an error")

	expected := `"15s"`
	assert.JSONEq(t, expected, string(marshaled), "The marshaled JSON should be the string '15s'")
}

func TestUnmarshalDurationWrapper(t *testing.T) {
	jsonStr := `"30s"`
	var durWrapper workflow.DurationWrapper

	err := json.Unmarshal([]byte(jsonStr), &durWrapper)
	assert.NoError(t, err, "Unmarshaling should not return an error")

	expectedDuration := 30 * time.Second
	assert.Equal(t, expectedDuration, durWrapper.Duration, "The unmarshaled duration should be 30s")
}



