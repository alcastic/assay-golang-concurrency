package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rnd := func() int {
		return rand.Intn(101)
	}

	genFun := func(ctx context.Context, fn func() int) <-chan int {
		stream := make(chan int)
		go func() {
			defer close(stream)
			for {
				select {
				case <-ctx.Done():
					return
				case stream <- fn():
				}

			}
		}()
		return stream
	}

	// context package provides feature for cancellations (deadlines, timeouts) and carry values through stack-frames
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for v := range genFun(ctx, rnd) {
		fmt.Println(v)
	}

}
