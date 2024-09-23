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

func TestStart(t *testing.T) {
	logger := utilities.NewLogger()
	s := scheduler.NewWithLogger(time.Second, logger)

	stop := make(chan struct{})
	result := make(chan bool)

	timerName := "timer1"
	deadline := time.Now().Add(500 * time.Millisecond)
	event := func() {
		result <- true
	}

	s.AddTimer(timerName, deadline, event)
	go s.Start(stop)

	select {
	case <-result:
		// Timer event was triggered
		s.Lock.RLock()
		_, exists := s.Timers[timerName]
		s.Lock.RUnlock()
		assert.False(t, exists)
	case <-time.After(2 * time.Second):
		// Timer event was not triggered within the expected time
		t.Fatal("Timer event was not triggered")
	}

	stop <- struct{}{}
}
func TestStartWithMultipleTimers(t *testing.T) {
	logger := utilities.NewLogger()
	s := scheduler.NewWithLogger(time.Second, logger)

	stop := make(chan struct{})
	result1 := make(chan bool)
	result2 := make(chan bool)

	timerName1 := "timer1"
	deadline1 := time.Now().Add(500 * time.Millisecond)
	event1 := func() {
		result1 <- true
	}

	timerName2 := "timer2"
	deadline2 := time.Now().Add(1 * time.Second)
	event2 := func() {
		result2 <- true
	}

	s.AddTimer(timerName1, deadline1, event1)
	s.AddTimer(timerName2, deadline2, event2)
	go s.Start(stop)

	select {
	case <-result1:
		// Timer1 event was triggered
		s.Lock.RLock()
		_, exists := s.Timers[timerName1]
		s.Lock.RUnlock()
		assert.False(t, exists)
	case <-time.After(2 * time.Second):
		// Timer1 event was not triggered within the expected time
		t.Fatal("Timer1 event was not triggered")
	}

	select {
	case <-result2:
		// Timer2 event was triggered
		s.Lock.RLock()
		_, exists := s.Timers[timerName2]
		s.Lock.RUnlock()
		assert.False(t, exists)
	case <-time.After(2 * time.Second):
		// Timer2 event was not triggered within the expected time
		t.Fatal("Timer2 event was not triggered")
	}

	stop <- struct{}{}
}
func TestCheckTimers(t *testing.T) {
	logger := utilities.NewLogger()
	s := scheduler.NewWithLogger(time.Second, logger)

	result := make(chan bool)

	timerName := "timer1"
	deadline := time.Now().Add(-time.Second) // Set deadline in the past to trigger immediately
	event := func() {
		result <- true
	}

	s.AddTimer(timerName, deadline, event)
	s.CheckTimers(time.Now())

	select {
	case <-result:
		// Timer event was triggered
		s.Lock.RLock()
		_, exists := s.Timers[timerName]
		s.Lock.RUnlock()
		assert.False(t, exists)
	case <-time.After(2 * time.Second):
		// Timer event was not triggered within the expected time
		t.Fatal("Timer event was not triggered")
	}
}
