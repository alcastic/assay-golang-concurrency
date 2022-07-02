package main

import "fmt"

func orDone(done <-chan struct{}, stream <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-stream:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-done:
				}
			}
		}
	}()
	return out
}

func main() {
	// 'bridge channel pattern is used for consume values from a sequence of channels;
	// the sequence of channels suggest an ordered write (this fact make it different from fan-in)
	genSequence := func() <-chan <-chan int {
		sequence := make(chan (<-chan int))
		go func() {
			defer close(sequence)
			for i := 0; i < 10; i++ {
				stream := make(chan int, 1)
				stream <- i
				close(stream)
				sequence <- stream
			}
		}()
		return sequence
	}

	bridge := func(done <-chan struct{}, channels <-chan <-chan int) <-chan int {
		stream := make(chan int)
		go func() {
			defer close(stream)
			for {
				select {
				case <-done:
					return
				case ts, ok := <-channels:
					if !ok {
						return
					}
					for v := range orDone(done, ts) {
						select {
						case stream <- v:
						case <-done:
						}
					}
				}
			}
		}()
		return stream
	}
	done := make(chan struct{})
	defer close(done)
	for v := range bridge(done, genSequence()) {
		fmt.Println(v)
	}
}
