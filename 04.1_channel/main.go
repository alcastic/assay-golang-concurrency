package main

import "fmt"

func main() {
	// Channels are a primitive for synchronization. They can be used to synchronize access to memory, but they are
	// best used to communicate information between goroutines.
	// working with channels can lead to troubles, ex:
	//   * read from a nil channel (block)
	//   * write to a nil channel (block)
	//   * write to a closed channel (panic)
	//   * close a nil channel (panic)
	//   * close a closed channel (panic)
	// A good practice to avoid these 'troubles' is declaring 'channels owners'.
	// A goroutine that owns a chanel should:
	//   * Instantiate the channel
	//   * Perform writes, or pass ownership to another goroutine
	//   * Close the channel
	//   * Encapsulate previous points and expose a reader channel
	owner := func() <-chan int {
		var intStream = make(chan int, 0) // channel can be buffered or unbuffered; and unidirectional or bidirectional
		go func() {
			defer close(intStream)
			for i := 10; i > 0; i-- {
				intStream <- i
			}
		}()
		return intStream
	}

	// nest for loop is the consumer of the channel. As a channel consumer it as mainly 2 worries:
	//  * Knowing when a channel is closed
	//  * Responsibly for handling blocking (for any reason)
	for value := range owner() {
		fmt.Println(value)
	}
}
