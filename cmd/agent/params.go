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
	env.Parse(&config)

	if config.Address == "" {
		flag.StringVar(&config.Address, "a", defaultAddress, "port & address")
	}

	if config.PollInterval == 0 {
		flag.IntVar(&config.PollInterval, "p", defaultPollInterval, "polling interval")
	}

	if config.ReportInterval == 0 {
		flag.IntVar(&config.ReportInterval, "r", defaultReportInterval, "report interval")
	}

	flag.Parse()

	return config
}
