package main

import (
	"fmt"
	"sync"
	"time"
)

type FireButton struct {
	// sync.Cond allows goroutines pause theirs execution waiting for an external evert occurs.
	// Such as event is an arbitrary signal between two or more goroutines and carries no information other
	// than the fact that it has occurred
	trigger   *sync.Cond
	actionsWg *sync.WaitGroup
}

func callFirefighters() {
	fmt.Println("Firefighter on their way")
}

func activateAudibleEmergencyAlarm() {
	fmt.Println("Alarm ringing loud loud")
}

func unlockAllDoors() {
	fmt.Println("Alarm ringing loud loud")
}

func activeEmergencyLight() {
	fmt.Println("Emergency light on")
}

func main() {
	fb := &FireButton{
		trigger:   sync.NewCond(&sync.Mutex{}),
		actionsWg: &sync.WaitGroup{},
	}

	subscribe := func(fb *FireButton, fn func()) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			fb.actionsWg.Add(1)
			wg.Done()
			defer fb.actionsWg.Done()
			fb.trigger.L.Lock()
			fb.trigger.Wait()
			fb.trigger.L.Unlock()
			fn()
		}()
		wg.Wait()
	}

	subscribe(fb, callFirefighters)
	subscribe(fb, activateAudibleEmergencyAlarm)
	subscribe(fb, unlockAllDoors)
	subscribe(fb, activeEmergencyLight)

	// Signal method notify the first waiting goroutine for the Cond in a FIFO order.
	fb.trigger.Signal() // It will make "callFirefighters" goroutine stop waiting.

	time.Sleep(time.Second) // add delay just to showcase difference between Signal and Broadcast

	// Broadcast method notify all waiting goroutines for the Cond
	fb.trigger.Broadcast()
	fb.actionsWg.Wait()
}
