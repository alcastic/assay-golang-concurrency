package main

import "fmt"

func main() {
	// A pipeline is a series of abstraction that take data in,
	// perform operations on it, and pass the data back out;
	// each of these operations are called a stage of the pipeline.
	// Next example shows a pipeline with two stages: 'multiply' and 'add'.

	// repeatGenerator and takeGenerator are pipeline generator functions.
	// Pipeline generator functions are any function that converts a set of
	// discrete values into a stream of values on a channel.
	repeatGenerator := func(
		done <-chan interface{},
		values ...int,
	) chan int {
		valueStream := make(chan int)
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	takeGenerator := func(
		done <-chan interface{},
		originStream <-chan int,
		n int,
	) chan int {
		valueStream := make(chan int)
		go func() {
			defer close(valueStream)
			for i := 0; i < n; i++ {
				select {
				case <-done:
					return
				case valueStream <- <-originStream:
				}
			}
		}()
		return valueStream
	}

	done := make(chan interface{})
	defer close(done)

	multiply := func(
		done <-chan interface{},
		originStream <-chan int,
		factor int,
	) chan int {
		valueStream := make(chan int)
		go func() {
			defer close(valueStream)
			for v := range originStream {
				select {
				case <-done:
					return
				case valueStream <- v * factor:
				}
			}
		}()
		return valueStream
	}

	add := func(
		done <-chan interface{},
		originStream <-chan int,
		number int,
	) chan int {
		valueStream := make(chan int)
		go func() {
			defer close(valueStream)
			for v := range originStream {
				select {
				case <-done:
					return
				case valueStream <- v + number:
				}
			}
		}()
		return valueStream
	}

	valueStream := takeGenerator(done, repeatGenerator(done, 1, 2, 3), 12)
	pipeline := multiply(
		done,
		add(done, valueStream, 1),
		2,
	)

	for v := range pipeline {
		fmt.Println(v)
	}
}
