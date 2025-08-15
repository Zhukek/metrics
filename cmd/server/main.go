package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zhukek/metrics/internal/handler"
	models "github.com/Zhukek/metrics/internal/model"
)

func main() {
	parseFlag()
	if err := run(); err != nil {
		log.Fatal("Critical:", err)
	}
}

func run() error {
	storage := models.NewStorage()

	fmt.Println("Running server on", flagAddress)
	return http.ListenAndServe(flagAddress, handler.NewRouter(storage))
}
