package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
)

func main() {
	repeatFnGen := func(done <-chan interface{}, fn func() int) <-chan int {
		stream := make(chan int)
		go func() {
			defer close(stream)
			for {
				select {
				case <-done:
					return
				case stream <- fn():
				}
			}
		}()
		return stream
	}

	takeGen := func(done <-chan interface{}, stream <-chan int, n int) <-chan int {
		takeStream := make(chan int)
		go func() {
			defer close(takeStream)
			for i := 0; i < n; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-stream:
				}
			}
		}()
		return takeStream
	}

	random := func() int {
		return rand.Intn(1001)
	}

	isEven := func(n int) bool {
		return n%2 == 0
	}
	filter := func(done <-chan interface{}, stream <-chan int, fn func(n int) bool) <-chan int {
		filterStream := make(chan int)
		go func() {
			defer close(filterStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-stream:
					if !ok {
						return
					} else if fn(v) {
						filterStream <- v
					}
				}
			}
		}()
		return filterStream
	}

	done := make(chan interface{})
	defer close(done)

	fanOut := func(done chan interface{}, numbersDataIn <-chan int) []<-chan int {
		numCpu := runtime.NumCPU()
		filters := make([]<-chan int, numCpu)
		for i := 0; i < numCpu; i++ {
			filters[i] = filter(done, numbersDataIn, isEven)
		}
		return filters
	}

	fanIn := func(
		done <-chan interface{},
		channels ...<-chan int,
	) <-chan int {
		var wg sync.WaitGroup
		multiplexedChannel := make(chan int)
		multiplex := func(channel <-chan int) {
			defer wg.Done()
			for v := range channel {
				select {
				case <-done:
					return
				case multiplexedChannel <- v:
				}
			}
		}
		for _, c := range channels {
			wg.Add(1)
			go multiplex(c)
		}
		go func() {
			wg.Wait()
			close(multiplexedChannel)
		}()
		return multiplexedChannel
	}

	numbersDataIn := takeGen(done, repeatFnGen(done, random), 1000)
	pipeline := fanIn(done, fanOut(done, numbersDataIn)...)
	for v := range pipeline {
		fmt.Println(v)
	}
}
