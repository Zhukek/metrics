package main

import (
	"flag"
	"os"
)

func getParams() string {
	address, ok := os.LookupEnv("ADDRESS")

	if !ok {
		flag.StringVar(&address, "a", "localhost:8080", "port & address")
		flag.Parse()
	}

	return address
}
