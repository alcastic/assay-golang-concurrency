package main

import (
	"fmt"
	"sync"
	"time"
)

// for-select pattern allows to iterate and wait for a possible external interruption
func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	done := make(chan struct{})

	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Nanosecond)
		close(done)
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
			}
			fmt.Println("Hi") // do something interesting
		}
	}()

	wg.Wait()

}
