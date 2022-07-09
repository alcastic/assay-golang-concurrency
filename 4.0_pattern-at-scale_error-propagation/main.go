package main

import (
	"fmt"
	"github.com/alcastic/assay-golang-concurrency/4.0_pattern-at-scale_error-propagation/service"
	"log"
	"net/http"
)

func handleError(err error) string {
	// here we ensure proper user message depending on the type/level error
	msg := "Unexpected Error"
	if _, ok := err.(service.ServiceError); ok {
		msg = err.Error()
	}
	return msg
}

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	err := service.Echo("Alex")
	if err != nil {
		fmt.Fprint(w, handleError(err))
		return
	}

	fmt.Fprint(w, "ok")
	return
}

func main() {
	if err := http.ListenAndServe(":8080", http.HandlerFunc(homeHandler)); err != nil {
		log.Fatal(err)
	}
}
