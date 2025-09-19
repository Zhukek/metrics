package main

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	Address   string `env:"ADDRESS"`
	Interval  int    `env:"STORE_INTERVAL"`
	FilePath  string `env:"FILE_STORAGE_PATH"`
	Restore   bool   `env:"RESTORE"`
	PGConnect string `env:"DATABASE_DSN"`
}

func getParams() Config {

	const (
		defaultAddress   = "localhost:8080"
		defaultInterval  = 300
		defaultFilePass  = ""
		defaultRestore   = false
		defaultPGConnect = ""
	)

	//host=127.0.0.1 port=5432 user=postgres password=postgres dbname=test sslmode=disable

	config := Config{}

	// Сначала определяем флаги
	flag.StringVar(&config.Address, "a", defaultAddress, "port & address")
	flag.IntVar(&config.Interval, "i", defaultInterval, "store interval")
	flag.StringVar(&config.FilePath, "f", defaultFilePass, "storage file pass")
	flag.BoolVar(&config.Restore, "r", defaultRestore, "restore")
	flag.StringVar(&config.PGConnect, "d", defaultPGConnect, "db connect")

	flag.Parse()

	// Затем парсим переменные окружения (они имеют приоритет)
	env.Parse(&config)

	return config
}
