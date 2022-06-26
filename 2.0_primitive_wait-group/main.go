package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("start")
	// Go follows concurrency mode called 'fork-join' model.
	// fork refers to the fact that at any time in the program, it can split off a 'child'
	// branch of execution to be running concurrently with its 'parent'.
	// The 'join' refers to the fact that at some point in the future these concurrent
	// branches of execution can join back together.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		fmt.Println("VERY CRITICAL IMPORTANT MESSAGE!")
		wg.Done()
	}() // fork point
	wg.Wait() // join point
	fmt.Println("end")
}
