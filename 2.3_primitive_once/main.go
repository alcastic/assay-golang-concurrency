package main

import (
	"fmt"
	"sync"
)

func main() {
	greet := func() {
		fmt.Println("hi!")
	}
	var once sync.Once
	for i := 0; i < 10; i++ {
		once.Do(greet) // greet function is going to be called just one time
	}
}
