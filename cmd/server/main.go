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
	"github.com/Zhukek/metrics/internal/repository"
	"github.com/Zhukek/metrics/internal/repository/inmemory"
	pg "github.com/Zhukek/metrics/internal/repository/pgrep"
	"github.com/Zhukek/metrics/internal/service/fileworker"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	slogger, err := logger.NewSlogger()
	if err != nil {
		return err
	}
	defer slogger.Sync()

	var (
		storage repository.Repository
		db      *sql.DB
	)

	if pgConnect != "" {
		pgrep, err := pg.NewPgRepository(pgConnect)
		if err != nil {
			return err
		}

		db = pgrep.Db
	} else {

		var storageInitData []byte
		if filePath != "" {

			fileWroker, err := fileworker.NewFileWorker(filePath, interval == 0)

			if err != nil {
				return err
			}
			defer fileWroker.Close()

			if restore {
				data, err := fileWroker.ReadData()
				if err != nil {
					return err
				}
				storageInitData = data
			}

			switch interval {
			case 0:
				if err := fileWroker.WriteData(storage.GetAllMetrics()); err != nil {
					slogger.ErrLog(err)
				}
			default:
				ticker := time.NewTicker(time.Duration(interval) * time.Second)
				defer ticker.Stop()

				go func() {
					for range ticker.C {
						if err := fileWroker.WriteData(storage.GetAllMetrics()); err != nil {
							slogger.ErrLog(err)
						}
					}
				}()
			}
		}

		storage, err = inmemory.NewStorage(storageInitData)

		if err != nil {
			return err
		}
	}

	fmt.Println("Running server on", address)
	return http.ListenAndServe(address, slogger.WithLogging(compress.GzipMiddleware(handler.NewRouter(storage, db))))
}
