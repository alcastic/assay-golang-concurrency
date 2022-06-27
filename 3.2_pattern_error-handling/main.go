package main

import (
	"fmt"
	"net/http"
)

func main() {
	// The most important question when working with error handling on concurrent code:
	// Who must be responsible for handling the error?
	// The key is to couple the potential result with the potential error together (See 'Result' Struct)
	// in order to allow the goroutine that spawn (has more context about the running program) others take
	// decisions about what to do with errors

	type Result struct {
		response *http.Response
		err      error
	}

	resultStream := func(urls ...string) chan *Result {
		results := make(chan *Result, len(urls))
		defer close(results)
		go func() {
			for _, url := range urls {
				go func(url string) {
					resp, err := http.Get(url) // child goroutine do not take decisions about the error
					results <- &Result{
						response: resp,
						err:      err,
					}
				}(url)
			}
		}()
		return results
	}

	urls := []string{"https://go.dev/", "hi", "https://www.gnu.org/"}
	for result := range resultStream(urls...) {
		if result.err != nil {
			fmt.Println("Error: %v", result.err)
		} else {
			fmt.Println("Success!")
		}
	}
}
