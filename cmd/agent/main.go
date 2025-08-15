package main

import (
	"fmt"
	"time"

	"github.com/Zhukek/metrics/internal/agent"
	"github.com/go-resty/resty/v2"
)

func main() {
	var client = resty.New()
	flags := parseFlags()

	statsData := agent.StatsData{}

	baseURL := agent.GetBaseURL(flags.address)

	client.SetBaseURL(baseURL)
	fmt.Printf("Sending requests to %s\n", baseURL)

	go func() {
		for {
			agent.Polling(&statsData)
			time.Sleep(time.Duration(flags.pollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			agent.PostUpdates(client, &statsData)
			time.Sleep(time.Duration(flags.reportInterval) * time.Second)
		}
	}()
	select {}
}
