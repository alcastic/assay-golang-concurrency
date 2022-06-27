package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("start process")
	// for-select pattern can be used to avoid goroutines leaks (no termination state reached).
	// A Convention: If a goroutine is responsible for creating a goroutine,
	// it is also responsible for ensuring it can stop the goroutine
	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case <-done:
					return
				case v := <-strings:
					fmt.Println(v) // do something interesting
				}
			}
		}()
		return terminated
	}

	var wg sync.WaitGroup
	wg.Add(1)

	done := make(chan interface{})
	var strings chan string

	completed := doWork(done, strings) // reading from nil channel will block goroutine execution
	go func() {
		time.Sleep(10 * time.Nanosecond)
		close(done)
	}()

	<-completed
	fmt.Println("end process")
}
