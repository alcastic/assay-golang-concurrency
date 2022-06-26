package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	value int
	mu    sync.Mutex
}

func main() {

	var wg sync.WaitGroup
	// Mutex is primitive for mutual exclusion memory access synchronization
	var counter Counter

	inc := func() {
		counter.mu.Lock()
		defer counter.mu.Unlock()
		counter.value++
		fmt.Printf("Increment to: %d\n", counter.value)
	}

	dec := func() {
		counter.mu.Lock()
		defer counter.mu.Unlock()
		counter.value--
		fmt.Printf("Decrement to: %d\n", counter.value)
	}

	nInc, nDec := 10, 10

	// INCREMENT
	wg.Add(nInc)
	for i := 0; i < nInc; i++ {
		go func() {
			defer wg.Done()
			inc()
		}()
	}
	// DECREMENT
	wg.Add(nDec)
	for i := 0; i < nDec; i++ {
		go func() {
			defer wg.Done()
			dec()
		}()
	}

	wg.Wait()
}
