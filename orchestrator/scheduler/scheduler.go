package scheduler

import (
    "time"
    "sync"
)

type Timer struct {
    Deadline time.Time
    Name     string
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
    s.Timers[name] = Timer{Deadline: deadline, Name: name}
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
            // Fire the event (this is just a placeholder, implement actual event firing)
            println("Event fired for timer:", name)
            // Remove the expired timer
            delete(s.Timers, name)
        }
    }
}