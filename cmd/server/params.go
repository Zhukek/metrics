package main

import (
	"flag"
	"os"

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
		defaultFilePass  = "data.json"
		defaultRestore   = false
		defaultPGConnect = ""
	)

	//host=127.0.0.1 port=5432 user=postgres password=postgres dbname=test sslmode=disable

	config := Config{}
	env.Parse(&config)

	if _, exist := os.LookupEnv("ADDRESS"); !exist {
		flag.StringVar(&config.Address, "a", defaultAddress, "port & address")
	}
	if _, exist := os.LookupEnv("STORE_INTERVAL"); !exist {
		flag.IntVar(&config.Interval, "i", defaultInterval, "store interval")
	}
	if _, exist := os.LookupEnv("FILE_STORAGE_PATH"); !exist {
		flag.StringVar(&config.FilePath, "f", defaultFilePass, "storage file pass")
	}
	if _, exist := os.LookupEnv("RESTORE"); !exist {
		flag.BoolVar(&config.Restore, "r", defaultRestore, "restore")
	}
	if _, exist := os.LookupEnv("DATABASE_DSN"); !exist {
		flag.StringVar(&config.PGConnect, "d", defaultPGConnect, "db connect")
	}

	flag.Parse()

	return config
}
