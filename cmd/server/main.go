package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zhukek/metrics/internal/handler"
	"github.com/Zhukek/metrics/internal/logger"
	compress "github.com/Zhukek/metrics/internal/middlewares"
	"github.com/Zhukek/metrics/internal/repository"
	"github.com/Zhukek/metrics/internal/repository/inmemory"
	pg "github.com/Zhukek/metrics/internal/repository/pgrep"
)

func main() {
	if err := run(); err != nil {
		log.Fatal("Critical:", err)
	}
}

func run() error {
	params := getParams()
	var (
		address   = params.Address
		filePath  = params.FilePath
		interval  = params.Interval
		restore   = params.Restore
		pgConnect = params.PGConnect
	)

	slogger, err := logger.NewSlogger()
	if err != nil {
		return err
	}
	defer slogger.Sync()

	var storage repository.Repository

	if pgConnect != "" {
		pgrep, err := pg.NewPgRepository(pgConnect)
		if err != nil {
			return err
		}
		defer pgrep.Close()

		storage = pgrep
	} else {

		storage, err = inmemory.NewStorage(filePath, interval, restore)

		if err != nil {
			return err
		}
	}

	fmt.Println("Running server on", address)
	return http.ListenAndServe(address, slogger.WithLogging(compress.GzipMiddleware(handler.NewRouter(storage))))
}
