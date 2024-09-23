package scheduler

import (
	"orchestrator/utilities"
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
	logger       *utilities.Logger
}

func New(pollInterval time.Duration) *Scheduler {
	return &Scheduler{
		timers:       make(map[string]Timer),
		lock:         &sync.RWMutex{},
		pollInterval: pollInterval,
	}
}

func NewScheduler() *Scheduler {

    return &Scheduler{}

}

func NewWithLogger(pollInterval time.Duration, logger *utilities.Logger) *Scheduler {
	return &Scheduler{
		timers:       make(map[string]Timer),
		lock:         &sync.RWMutex{},
		pollInterval: pollInterval,
		logger:       logger,
	}
}

// Adds a new timer that will execute the callback after the deadline, overwriting
// any existing timer with the same name.
func (s *Scheduler) AddTimer(name string, deadline time.Time, event func()) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.timers[name] = Timer{
		Name:     name,
		Deadline: deadline,
		Event:    event,
	}
	s.logger.Info("Timer added: ", name, ". Expires at ", deadline)
}

// Removes the timer with the passed-in name, returning whether that timer was
// found and successfully removed.
func (s *Scheduler) RemoveTimer(name string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, found := s.timers[name]
	delete(s.timers, name)

	if found {
		s.logger.Debug("Timer removed: ", name)
	} else {
		s.logger.Debug("Timer not found: ", name)
	}

	return found
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
				s.logger.Info("Scheduler stopped")
				return
			}
		}
	}()
}

// Triggers the callbacks of all expired timers
func (s *Scheduler) checkTimers(currentTime time.Time) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	s.logger.Info("Checking all timers")
	for _, timer := range s.timers {
		if currentTime.After(timer.Deadline) {
			s.logger.Info("Timer expired:", timer.Name)
			go timer.Event()

			// There is no way to upgrade an existing lock, so this is the way (unfortunately)
			s.lock.RUnlock()
			s.lock.Lock()

			delete(s.timers, timer.Name)
			s.logger.Debug("Timer removed: ", timer.Name)

			s.lock.Unlock()
			s.lock.RLock()
		}
	}
}
