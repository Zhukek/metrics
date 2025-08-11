package main

import (
	"math/rand/v2"
	"runtime"
	"time"

	"github.com/Zhukek/metrics/internal/agent"
	"github.com/go-resty/resty/v2"
)

func main() {
	var client = resty.New()

	const pollInterval = 2
	const reportInterval = 10
	var stat runtime.MemStats
	var counter int64 = 0
	var RandomValue float64 = 0
	const url = "http://localhost:8080"

	go func() {
		for {
			RandomValue = rand.Float64()
			runtime.ReadMemStats(&stat)
			counter++
			time.Sleep(pollInterval * time.Second)
		}
	}()

	go func() {
		for {
			agent.PostUpdate(client, agent.BuildURLCounter(url, "counter", counter))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "RandomValue", RandomValue))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "Alloc", float64(stat.Alloc)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "BuckHashSys", float64(stat.BuckHashSys)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "Frees", float64(stat.Frees)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "GCCPUFraction", float64(stat.GCCPUFraction)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "GCSys", float64(stat.GCSys)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "HeapAlloc", float64(stat.HeapAlloc)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "HeapIdle", float64(stat.HeapIdle)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "HeapInuse", float64(stat.HeapInuse)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "HeapObjects", float64(stat.HeapObjects)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "HeapReleased", float64(stat.HeapReleased)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "HeapSys", float64(stat.HeapSys)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "LastGC", float64(stat.LastGC)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "Lookups", float64(stat.Lookups)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "MCacheInuse", float64(stat.MCacheInuse)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "MCacheSys", float64(stat.MCacheSys)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "MSpanInuse", float64(stat.MSpanInuse)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "MSpanSys", float64(stat.MSpanSys)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "Mallocs", float64(stat.Mallocs)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "NextGC", float64(stat.NextGC)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "NumForcedGC", float64(stat.NumForcedGC)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "NumGC", float64(stat.NumGC)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "OtherSys", float64(stat.OtherSys)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "PauseTotalNs", float64(stat.PauseTotalNs)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "StackInuse", float64(stat.StackInuse)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "StackSys", float64(stat.StackSys)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "Sys", float64(stat.Sys)))
			agent.PostUpdate(client, agent.BuildURLGauge(url, "TotalAlloc", float64(stat.TotalAlloc)))

			counter = 0
			time.Sleep(reportInterval * time.Second)
		}
	}()
	select {}
}
