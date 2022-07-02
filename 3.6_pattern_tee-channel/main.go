package main

import "fmt"

func main() {
	// tee-channels pattern is used to take values from one channel and send them off into two separate areas of the code base
	// WARNING: writes to output channels are tightly couple as well the readings
	// INFO: tee-channels pattern takes its name from 'tee' unix command
	gen := func(nums ...int) <-chan int {
		stream := make(chan int)
		go func() {
			defer close(stream)
			for _, v := range nums {
				stream <- v
			}
		}()
		return stream
	}

	tee := func(done <-chan struct{}, origin <-chan int) (_, _ <-chan int) {
		out1, out2 := make(chan int), make(chan int)
		go func() {
			defer close(out1)
			defer close(out2)
			for {
				select {
				case <-done:
					return
				case val, ok := <-origin:
					if !ok {
						return
					}
					var out1, out2 = out1, out2
					for i := 0; i < 2; i++ {
						select {
						case out1 <- val:
							out1 = nil
						case out2 <- val:
							out2 = nil
						}
					}
				}
			}
		}()
		return out1, out2
	}

	done := make(chan struct{})
	defer close(done)

	nums1, nums2 := tee(done, gen(0, 1, 2, 3, 4, 5, 6, 7, 8, 9))
	_, _ = nums1, nums2

	for v := range nums1 {
		fmt.Printf("nums1: %d - nums2: %d\n", v, <-nums2)
	}

}
