package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Zhukek/metrics/internal/handler"
	"github.com/Zhukek/metrics/internal/logger"
	compress "github.com/Zhukek/metrics/internal/middlewares"
	models "github.com/Zhukek/metrics/internal/model"
	"github.com/Zhukek/metrics/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	fileWroker, err := service.NewFileWorker(filePath, interval == 0)

	if err != nil {
		return err
	}
	defer fileWroker.Close()

	db, err := sql.Open("pgx", pgConnect)
	if err != nil {
		return err
	}
	defer db.Close()

	var storageInitData []byte

	if restore {
		data, err := fileWroker.ReadData()
		if err != nil {
			return err
		}
		storageInitData = data
	}

	storage, err := models.NewStorage(storageInitData)

	if err != nil {
		return err
	}

	slogger, err := logger.NewSlogger()

	if err != nil {
		return err
	}

	defer slogger.Sync()

	switch interval {
	case 0:
		if err := writeData(storage, fileWroker); err != nil {
			slogger.ErrLog(err)
		}
	default:
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		go func() {
			for range ticker.C {
				if err := writeData(storage, fileWroker); err != nil {
					slogger.ErrLog(err)
				}
			}
		}()
	}

	fmt.Println("Running server on", address)
	return http.ListenAndServe(address, slogger.WithLogging(compress.GzipMiddleware(handler.NewRouter(storage, db))))
}

func writeData(storage *models.MemStorage, fileWroker *service.FileWorker) error {
	data, err := storage.GetJSONStorage()
	if err != nil {
		return err
	}
	err = fileWroker.WriteData(data)
	if err != nil {
		return err
	}
	return nil
}
