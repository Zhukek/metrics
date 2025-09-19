package main

import (
	"fmt"
	"time"

	"github.com/Zhukek/metrics/internal/agent"
	"github.com/go-resty/resty/v2"
)

func main() {
	var client = resty.New()
	config := getParams()

	statsData := agent.StatsData{}

	baseURL := agent.GetBaseURL(config.Address)

	client.SetBaseURL(baseURL)
	fmt.Printf("Sending requests to %s\n", baseURL)

	reportTicker := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)

	defer reportTicker.Stop()
	defer pollTicker.Stop()

	go func() {
		for range pollTicker.C {
			agent.Polling(&statsData)
		}
	}()

	go func() {
		for range reportTicker.C {
			agent.PostUpdates(client, &statsData)
			agent.PostBatch(client, &statsData, nil)
		}
	}()
	select {}
}
