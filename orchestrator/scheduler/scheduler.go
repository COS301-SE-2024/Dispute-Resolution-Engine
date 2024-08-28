package scheduler

import (
    "time"
    "sync"
)

type Timer struct {
    Deadline time.Time
    Name     string
	Event   func()
}

type IScheduler interface {
    AddTimer(name string, deadline time.Time)
    AddEvent(name string, event func())
    CancelTimer(name string)
    Start()
    Stop()
}

type Scheduler struct {
    Timers   map[string]Timer
    Interval time.Duration
    mu       sync.Mutex
    stopChan chan struct{}
}

func NewScheduler(interval time.Duration) *Scheduler {
    return &Scheduler{
        Timers:   make(map[string]Timer),
        Interval: interval,
        stopChan: make(chan struct{}),
    }
}

func (s *Scheduler) AddTimer(name string, deadline time.Time) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Timers[name] = Timer{Deadline: deadline, Name: name, Event: nil}
}

func (s *Scheduler) AddEvent(name string, event func()) {
    s.mu.Lock()
    defer s.mu.Unlock()
    if timer, exists := s.Timers[name]; exists {
        timer.Event = event
        s.Timers[name] = timer
    }
}

func (s *Scheduler) CancelTimer(name string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.Timers, name)
}

func (s *Scheduler) Start() {
    ticker := time.NewTicker(s.Interval)
    go func() {
        for {
            select {
            case <-ticker.C:
                s.checkTimers()
            case <-s.stopChan:
                ticker.Stop()
                return
            }
        }
    }()
}

func (s *Scheduler) Stop() {
    close(s.stopChan)
}

func (s *Scheduler) checkTimers() {
    s.mu.Lock()
    defer s.mu.Unlock()
    now := time.Now()
    for name, timer := range s.Timers {
        if now.After(timer.Deadline) {
            if timer.Event != nil {
                timer.Event()
            }
            delete(s.Timers, name)
        }
    }
}