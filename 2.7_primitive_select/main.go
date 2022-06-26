package main

import (
	"fmt"
	"time"
)

func main() {
	msgStream1 := make(chan string)
	msgStream2 := make(chan string)

	go func() {
		time.Sleep(time.Second * 1)
		msgStream1 <- "Hi!"
	}()
	go func() {
		time.Sleep(time.Second * 1)
		msgStream2 <- "Hey!"
	}()

	end := false
	for !end {
		// Select is used to compose channels.
		// It is "the glue to bind channels together"
		select {
		// Also is possible to add a "default" case for select that will
		// be executed when all other cases get blocked
		case msg := <-msgStream1:
			fmt.Println(msg)
		case msg := <-msgStream2:
			fmt.Println(msg)
		case <-time.After(5 * time.Second):
			fmt.Println("time out")
			end = true
		}
	}
}
