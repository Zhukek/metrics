package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zhukek/metrics/internal/handler"
	"github.com/Zhukek/metrics/internal/logger"
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
	slogger, err := logger.NewSlogger()

	if err != nil {
		return err
	}

	defer slogger.Sync()

	fmt.Println("Running server on", params)
	return http.ListenAndServe(params, slogger.WithLogging(handler.NewRouter(&storage)))
}
