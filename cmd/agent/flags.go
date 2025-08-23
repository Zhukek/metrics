package main

import "flag"

type Flags struct {
	address        string
	reportInterval int
	pollInterval   int
}

func parseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&flags.address, "a", "localhost:8080", "port & address")
	flag.IntVar(&flags.pollInterval, "p", 2, "polling interval")
	flag.IntVar(&flags.reportInterval, "r", 10, "report interval")

	flag.Parse()

	return flags
}
