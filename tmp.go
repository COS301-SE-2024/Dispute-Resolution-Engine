package main

import (
	"fmt"
	"time"
)

type Timer struct {
	Name     string
	Deadline time.Time
	Event    func()
}

func main() {

	// The scheduler
	stop := make(chan struct{})
	timers := map[string]Timer{}
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case currentTime := <-ticker.C:
				for _, timer := range timers {
					if currentTime.After(timer.Deadline) {
						timer.Event()
					}
				}
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
