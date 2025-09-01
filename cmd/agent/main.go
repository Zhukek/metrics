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

	go func() {
		for {
			agent.Polling(&statsData)
			time.Sleep(time.Duration(config.PollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(config.ReportInterval) * time.Second)
			agent.PostUpdates(client, &statsData)
		}
	}()
	select {}
}
