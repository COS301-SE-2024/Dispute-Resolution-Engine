package scheduler

import (
	"fmt"
	"sync"
	"time"
)

type Timer struct {
	Name     string
	Deadline time.Time
	Event    func()
}

type Scheduler struct {
	timers       map[string]Timer
	lock         *sync.RWMutex
	pollInterval time.Duration
}

func New(pollInterval time.Duration) Scheduler {
	return Scheduler{
		timers:       make(map[string]Timer),
		lock:         &sync.RWMutex{},
		pollInterval: pollInterval,
	}
}

// Adds a new timer that will execute the callback after the deadline
func (s *Scheduler) AddTimer(name string, deadline time.Time, event func()) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.timers[name] = Timer{
		Name:     name,
		Deadline: deadline,
		Event:    event,
	}
}

// Removes the timer with the passed-in name
func (s *Scheduler) RemoveTimer(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.timers, name)
}

// Spawns a new goroutine that polls the timer registry for expired timers
func (s *Scheduler) Start(stop chan struct{}) {
	ticker := time.NewTicker(s.pollInterval)
	go func() {
		for {
			select {
			case currentTime := <-ticker.C:
				s.checkTimers(currentTime)
			case <-stop:
				fmt.Println("Stop scheduler")
				return
			}
		}
	}()
}

// Triggers the callbacks of all expired timers
func (s *Scheduler) checkTimers(currentTime time.Time) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, timer := range s.timers {
		if currentTime.After(timer.Deadline) {
			timer.Event()

			// There is no way to upgrade an existing lock, so this is the way (unfortunately)
			s.lock.RUnlock()
			s.lock.Lock()

			delete(s.timers, timer.Name)

			s.lock.Unlock()
			s.lock.RLock()
		}
	}
}
