package main

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"strings"
	"time"

	"github.com/Zhukek/metrics/internal/agent"
	"github.com/go-resty/resty/v2"
)

func getBaseURL(URL string) string {
	if strings.Contains(flags.address, "://") {
		return URL
	} else {
		return "http://" + URL
	}
}

func main() {
	var client = resty.New()
	parseFlags()

	var stat runtime.MemStats
	var counter int64 = 0
	var RandomValue float64 = 0

	baseURL := getBaseURL(flags.address)

	client.SetBaseURL(baseURL)
	fmt.Printf("Sending requests to %s\n", baseURL)

	go func() {
		for {
			RandomValue = rand.Float64()
			runtime.ReadMemStats(&stat)
			counter++
			time.Sleep(time.Duration(flags.pollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			agent.PostUpdate(client, agent.BuildURLCounter("counter", counter))
			agent.PostUpdate(client, agent.BuildURLGauge("RandomValue", RandomValue))
			agent.PostUpdate(client, agent.BuildURLGauge("Alloc", float64(stat.Alloc)))
			agent.PostUpdate(client, agent.BuildURLGauge("BuckHashSys", float64(stat.BuckHashSys)))
			agent.PostUpdate(client, agent.BuildURLGauge("Frees", float64(stat.Frees)))
			agent.PostUpdate(client, agent.BuildURLGauge("GCCPUFraction", float64(stat.GCCPUFraction)))
			agent.PostUpdate(client, agent.BuildURLGauge("GCSys", float64(stat.GCSys)))
			agent.PostUpdate(client, agent.BuildURLGauge("HeapAlloc", float64(stat.HeapAlloc)))
			agent.PostUpdate(client, agent.BuildURLGauge("HeapIdle", float64(stat.HeapIdle)))
			agent.PostUpdate(client, agent.BuildURLGauge("HeapInuse", float64(stat.HeapInuse)))
			agent.PostUpdate(client, agent.BuildURLGauge("HeapObjects", float64(stat.HeapObjects)))
			agent.PostUpdate(client, agent.BuildURLGauge("HeapReleased", float64(stat.HeapReleased)))
			agent.PostUpdate(client, agent.BuildURLGauge("HeapSys", float64(stat.HeapSys)))
			agent.PostUpdate(client, agent.BuildURLGauge("LastGC", float64(stat.LastGC)))
			agent.PostUpdate(client, agent.BuildURLGauge("Lookups", float64(stat.Lookups)))
			agent.PostUpdate(client, agent.BuildURLGauge("MCacheInuse", float64(stat.MCacheInuse)))
			agent.PostUpdate(client, agent.BuildURLGauge("MCacheSys", float64(stat.MCacheSys)))
			agent.PostUpdate(client, agent.BuildURLGauge("MSpanInuse", float64(stat.MSpanInuse)))
			agent.PostUpdate(client, agent.BuildURLGauge("MSpanSys", float64(stat.MSpanSys)))
			agent.PostUpdate(client, agent.BuildURLGauge("Mallocs", float64(stat.Mallocs)))
			agent.PostUpdate(client, agent.BuildURLGauge("NextGC", float64(stat.NextGC)))
			agent.PostUpdate(client, agent.BuildURLGauge("NumForcedGC", float64(stat.NumForcedGC)))
			agent.PostUpdate(client, agent.BuildURLGauge("NumGC", float64(stat.NumGC)))
			agent.PostUpdate(client, agent.BuildURLGauge("OtherSys", float64(stat.OtherSys)))
			agent.PostUpdate(client, agent.BuildURLGauge("PauseTotalNs", float64(stat.PauseTotalNs)))
			agent.PostUpdate(client, agent.BuildURLGauge("StackInuse", float64(stat.StackInuse)))
			agent.PostUpdate(client, agent.BuildURLGauge("StackSys", float64(stat.StackSys)))
			agent.PostUpdate(client, agent.BuildURLGauge("Sys", float64(stat.Sys)))
			agent.PostUpdate(client, agent.BuildURLGauge("TotalAlloc", float64(stat.TotalAlloc)))

			counter = 0
			time.Sleep(time.Duration(flags.reportInterval) * time.Second)
		}
	}()
	select {}
}
