package main

import (
	"net/http"

	"github.com/Zhukek/metrics/internal/handler"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	return http.ListenAndServe(":8080", handler.NewRouter())
}
