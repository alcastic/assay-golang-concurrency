package main

import (
	"fmt"
	"sync"
)

type Performer interface {
	perform()
}
type Resource struct{}

func (r *Resource) perform() {}

func main() {
	var countCallsToNew int

	// Pool is a concurrent safe implementaton of the object pool pattern.
	// It's commonly used to constrain the creations of 'things' and encourage
	// the reuse of available instances in the pool
	var pool = sync.Pool{
		New: func() interface{} {
			countCallsToNew++
			fmt.Printf("Calls to New: %d\n", countCallsToNew)
			return &Resource{} // return some 'hard' resource to create
		},
	}

	var wg sync.WaitGroup
	nroGoroutines := 100
	wg.Add(nroGoroutines)
	for i := nroGoroutines; i > 0; i-- {
		go func() {
			defer wg.Done()
			resource := pool.Get()
			resource.(Performer).perform()
			pool.Put(resource)
		}()
	}
	wg.Wait()

	fmt.Printf("Final Calls to New: %d\n", countCallsToNew)
}
