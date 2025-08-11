package main

import (
	"net/http"

	"github.com/Zhukek/metrics/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	router := chi.NewRouter()
	router.Post("/update/{metricType}/{metricName}/{metricValue}", handler.Update)
	router.Get("/value/{metricType}/{metricName}", handler.Get)
	router.Get("/", handler.GetList)

	return http.ListenAndServe(":8080", router)
}
