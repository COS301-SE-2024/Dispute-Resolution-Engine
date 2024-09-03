package scheduler

import (
	"fmt"
	"time"
)

var timers = map[string]Timer{}

type Timer struct {
	Name     string
	Deadline time.Time
	Event    func()
}

func AddTimer(name string, deadline time.Time, event func()) {
	timers[name] = Timer{
		Name:     name,
		Deadline: deadline,
		Event:    event,
	}
}

func RemoveTimer(name string) {
	delete(timers, name)
}

func Start(stop chan struct{}) {

	// The scheduler
	// timers := map[string]Timer{}
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case currentTime := <-ticker.C:
				checkTimers(timers, currentTime)
			case <-stop:
				fmt.Println("Stop scheduler")
				return
			}
		}
	}()

	for {
		var value string
		fmt.Print("Action: ")
		fmt.Scan(&value)
		if value == "quit" {
			stop <- struct{}{}
		}
	}
}

func checkTimers(timers map[string]Timer, currentTime time.Time) {
	for _, timer := range timers {
		if currentTime.After(timer.Deadline) {
			timer.Event()
		}
	}
}
