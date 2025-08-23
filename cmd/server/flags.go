package main

import "flag"

var flagAddress string

func parseFlag() {
	flag.StringVar(&flagAddress, "a", "localhost:8080", "port & address")

	flag.Parse()
}
