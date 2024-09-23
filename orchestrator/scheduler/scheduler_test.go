package scheduler_test

import (
	"orchestrator/scheduler"
	"orchestrator/utilities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// func TestAddTimer(t *testing.T) {
// 	s := scheduler.New(time.Second)
// 	s.AddTimer("timer", time.Now().Add(time.Second), func() {})
// }

// func TestRemoveTimer(t *testing.T) {
// 	s := scheduler.New(time.Second)
// 	s.AddTimer("timer", time.Now().Add(time.Second), func() {})
// 	assert.True(t, s.RemoveTimer("timer"))
// }

// func TestRemoveMissingTimer(t *testing.T) {
// 	s := scheduler.New(time.Second)
// 	assert.False(t, s.RemoveTimer("timer"))
// }

// func TestStart(t *testing.T) {
// 	stop := make(chan struct{})
// 	result := make(chan time.Time)

// 	s := scheduler.New(time.Second)
// 	s.AddTimer("timer", time.Now().Add(time.Second), func() {
// 		result <- time.Now()
// 	})

// 	currentTime := time.Now().Add(time.Second)
// 	s.Start(stop)
// 	nextTime := <-result
// 	assert.True(t, nextTime.After(currentTime))
// 	stop <- struct{}{}
// }

func TestNewWithLogger(t *testing.T) {
	logger := utilities.NewLogger()
	s := scheduler.NewWithLogger(time.Second, logger)

	assert.NotNil(t, s)
	assert.Equal(t, time.Second, s.PollInterval)
	assert.Equal(t, logger, s.Logger)
	assert.NotNil(t, s.Timers)
	assert.NotNil(t, s.Lock)
}
func TestAddTimer(t *testing.T) {
	logger := utilities.NewLogger()
	s := scheduler.NewWithLogger(time.Second, logger)
	
	timerName := "timer1"
	deadline := time.Now().Add(time.Second)
	eventTriggered := false
	event := func() {
		eventTriggered = true
	}

	s.AddTimer(timerName, deadline, event)

	s.Lock.RLock()
	timer, exists := s.Timers[timerName]
	s.Lock.RUnlock()

	assert.True(t, exists)
	assert.Equal(t, timerName, timer.Name)
	assert.Equal(t, deadline, timer.Deadline)
	assert.NotNil(t, timer.Event)
	assert.False(t, eventTriggered)
}
func TestRemoveTimer(t *testing.T) {
	logger := utilities.NewLogger()
	s := scheduler.NewWithLogger(time.Second, logger)

	timerName := "timer1"
	deadline := time.Now().Add(time.Second)
	event := func() {}

	s.AddTimer(timerName, deadline, event)
	removed := s.RemoveTimer(timerName)

	assert.True(t, removed)

	s.Lock.RLock()
	_, exists := s.Timers[timerName]
	s.Lock.RUnlock()

	assert.False(t, exists)
}

func TestRemoveMissingTimer(t *testing.T) {
	logger := utilities.NewLogger()
	s := scheduler.NewWithLogger(time.Second, logger)

	timerName := "nonexistent_timer"
	removed := s.RemoveTimer(timerName)

	assert.False(t, removed)
}