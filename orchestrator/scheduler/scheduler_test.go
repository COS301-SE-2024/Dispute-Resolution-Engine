package scheduler_test

import (
	"orchestrator/scheduler"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddTimer(t *testing.T) {
	s := scheduler.New(time.Second)
	s.AddTimer("timer", time.Now().Add(time.Second), func() {})
}

func TestRemoveTimer(t *testing.T) {
	s := scheduler.New(time.Second)
	s.AddTimer("timer", time.Now().Add(time.Second), func() {})
	assert.True(t, s.RemoveTimer("timer"))
}

func TestRemoveMissingTimer(t *testing.T) {
	s := scheduler.New(time.Second)
	assert.False(t, s.RemoveTimer("timer"))
}

func TestStart(t *testing.T) {
	stop := make(chan struct{})
	result := make(chan time.Time)

	s := scheduler.New(time.Second)
	s.AddTimer("timer", time.Now().Add(time.Second), func() {
		result <- time.Now()
	})

	currentTime := time.Now().Add(time.Second)
	s.Start(stop)
	nextTime := <-result
	assert.True(t, nextTime.After(currentTime))
	stop <- struct{}{}
}
