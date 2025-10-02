package main

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	Key            string `env:"KEY"`
	Rate           int    `env:"RATE_LIMIT"`
}

func getParams() Config {
	const (
		defaultAddress        = "localhost:8080"
		defaultPollInterval   = 2
		defaultReportInterval = 10
		defaultKey            = ""
		defaultRate           = 1
	)

	config := Config{}

	flag.StringVar(&config.Address, "a", defaultAddress, "port & address")
	flag.IntVar(&config.PollInterval, "p", defaultPollInterval, "polling interval")
	flag.IntVar(&config.ReportInterval, "r", defaultReportInterval, "report interval")
	flag.StringVar(&config.Key, "k", defaultKey, "hash key")
	flag.IntVar(&config.Rate, "l", defaultRate, "rate limit")

	flag.Parse()

	env.Parse(&config)

	return config
}
