package main

import "fmt"

func main() {
	// Race condition occurs when the multiple concurrent branches of execution
	// does not warranty the expected correctness of the program.
	fmt.Println("start")

	// Here a "race condition": next go routine may not be executed before main goroutine ends.
	go func() {
		fmt.Println("THIS MESSAGE MUST BE DELIVERED")
	}()
	fmt.Println("end")
}
