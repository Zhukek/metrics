package main

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func getParams() Config {
	const (
		defaultAddress        = "localhost:8080"
		defaultPollInterval   = 2
		defaultReportInterval = 10
	)

	config := Config{}

	// Сначала определяем флаги
	flag.StringVar(&config.Address, "a", defaultAddress, "port & address")
	flag.IntVar(&config.PollInterval, "p", defaultPollInterval, "polling interval")
	flag.IntVar(&config.ReportInterval, "r", defaultReportInterval, "report interval")

	flag.Parse()

	// Затем парсим переменные окружения (они имеют приоритет)
	env.Parse(&config)

	return config
}
