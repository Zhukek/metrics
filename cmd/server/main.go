package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zhukek/metrics/internal/handler"
	models "github.com/Zhukek/metrics/internal/model"
)

func main() {
	if err := run(); err != nil {
		log.Fatal("Critical:", err)
	}
}

func run() error {
	storage := models.NewStorage()
	params := getParams()

	fmt.Println("Running server on", params)
	return http.ListenAndServe(params, handler.NewRouter(storage))
}
