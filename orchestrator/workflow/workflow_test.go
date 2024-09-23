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

func TestMarshalTimer(t *testing.T) {
	duration := 20 * time.Second
	onExpire := "trigger_event"
	tmr := workflow.CreateTimer(duration, onExpire)

	marshaled, err := json.Marshal(tmr)
	assert.NoError(t, err, "Marshaling the timer should not return an error")

	expectedJSON := `{"duration":"20s","on_expire":"trigger_event"}`
	assert.JSONEq(t, expectedJSON, string(marshaled), "The marshaled JSON should match the expected structure")
}

func TestUnmarshalTimer(t *testing.T) {
	jsonStr := `{"duration":"10s","on_expire":"trigger_event"}`
	var tmr workflow.Timer

	err := json.Unmarshal([]byte(jsonStr), &tmr)
	assert.NoError(t, err, "Unmarshaling the timer should not return an error")

	expectedDuration := 10 * time.Second
	assert.Equal(t, expectedDuration, tmr.GetDuration(), "The unmarshaled timer should have a duration of 10s")
	assert.Equal(t, "trigger_event", tmr.OnExpire, "The OnExpire trigger should match the provided value")
}

// Mock Trigger for testing purposes
type Trigger struct {
	Label string
}
func TestCreateState(t *testing.T) {
	label := "Initial"
	description := "This is the initial state."
	state := workflow.CreateState(label, description)

	assert.Equal(t, label, state.Label, "The state label should match the provided value")
	assert.Equal(t, description, state.Description, "The state description should match the provided value")
	assert.Empty(t, state.Triggers, "The state should have no triggers initially")
	assert.Nil(t, state.Timer, "The state should have no timer initially")
}

func TestAddTrigger(t *testing.T) {
	state := workflow.CreateState("State1", "A test state")
	trigger := workflow.NewTrigger("Trigger1", "NextState")

	state.AddTrigger(trigger)

	assert.Contains(t, state.Triggers, trigger.Label, "The state should contain the added trigger")
	assert.Equal(t, trigger, state.Triggers[trigger.Label], "The added trigger should match the expected trigger")
}

func TestSetTimer(t *testing.T) {
	state := workflow.CreateState("State1", "A test state")
	timer := workflow.CreateTimer(10*time.Second, "trigger_event")

	state.SetTimer(timer)

	assert.NotNil(t, state.Timer, "The state should have a timer set")
	assert.Equal(t, timer, *state.Timer, "The state's timer should match the expected timer")
}