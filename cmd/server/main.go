package main

import (
	"fmt"
	"net/http"

	"github.com/Zhukek/metrics/internal/handler"
)

func main() {
	parseFlag()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Running server on", flagAddress)
	return http.ListenAndServe(flagAddress, handler.NewRouter())
}
