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
	Timers       map[string]Timer
	Lock         *sync.RWMutex
	PollInterval time.Duration
	Logger       *utilities.Logger
}

func New(pollInterval time.Duration) *Scheduler {
	return &Scheduler{
		Timers:       make(map[string]Timer),
		Lock:         &sync.RWMutex{},
		PollInterval: pollInterval,
	}
}

func NewScheduler() *Scheduler {

    return &Scheduler{}

}

func NewWithLogger(pollInterval time.Duration, logger *utilities.Logger) *Scheduler {
	return &Scheduler{
		Timers:       make(map[string]Timer),
		Lock:         &sync.RWMutex{},
		PollInterval: pollInterval,
		Logger:       logger,
	}
}

// Adds a new timer that will execute the callback after the deadline, overwriting
// any existing timer with the same name.
func (s *Scheduler) AddTimer(name string, deadline time.Time, event func()) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.Timers[name] = Timer{
		Name:     name,
		Deadline: deadline,
		Event:    event,
	}
	s.Logger.Info("Timer added: ", name, ". Expires at ", deadline)
}

// Removes the timer with the passed-in name, returning whether that timer was
// found and successfully removed.
func (s *Scheduler) RemoveTimer(name string) bool {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	_, found := s.Timers[name]
	delete(s.Timers, name)

	if found {
		s.Logger.Debug("Timer removed: ", name)
	} else {
		s.Logger.Debug("Timer not found: ", name)
	}

	return found
}

// Spawns a new goroutine that polls the timer registry for expired timers
func (s *Scheduler) Start(stop chan struct{}) {
	ticker := time.NewTicker(s.PollInterval)
	go func() {
		for {
			select {
			case currentTime := <-ticker.C:
				s.checkTimers(currentTime)
			case <-stop:
				s.Logger.Info("Scheduler stopped")
				return
			}
		}
	}()
}

// Triggers the callbacks of all expired timers
func (s *Scheduler) checkTimers(currentTime time.Time) {
	s.Lock.RLock()
	defer s.Lock.RUnlock()

	s.Logger.Info("Checking all timers")
	for _, timer := range s.Timers {
		if currentTime.After(timer.Deadline) {
			s.Logger.Info("Timer expired:", timer.Name)
			go timer.Event()

			// There is no way to upgrade an existing lock, so this is the way (unfortunately)
			s.Lock.RUnlock()
			s.Lock.Lock()

			delete(s.Timers, timer.Name)
			s.Logger.Debug("Timer removed: ", timer.Name)

			s.Lock.Unlock()
			s.Lock.RLock()
		}
	}
}
