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
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handler.Update)

	return http.ListenAndServe(":8080", mux)
}
