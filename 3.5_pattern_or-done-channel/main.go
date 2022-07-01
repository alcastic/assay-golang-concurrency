package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// genExternalChannelThatNoAssertionCanBeMade this method is not part of the pattern.
func genUnpredictableChannel() chan int {
	stream := make(chan int)
	go func() {
		for {
			stream <- rand.Intn(101)
		}
	}()
	return stream
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// Or Done Channel patter is used when working with channels that no assertions can be made
	// about their behaviour when our code is cancel via a "done channel".
	// This patter wraps our read from the external channel with a "select" statement that
	// also selects from a done channel
	orDone := func(done <-chan interface{}, external chan int) <-chan int {
		stream := make(chan int)
		go func() {
			defer close(stream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-external:
					if !ok {
						return
					}
					select {
					case stream <- v:
					case <-done:
					}
				}
			}
		}()
		return stream
	}

	done := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(5 * time.Second)
		defer wg.Done()
		defer fmt.Println("Interrupt signal sent to orDone")
		defer close(done)

	}()

	for v := range orDone(done, genUnpredictableChannel()) {
		fmt.Println(v)
	}
	wg.Wait()
}
