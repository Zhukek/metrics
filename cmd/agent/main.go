package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"runtime"
	"time"

	models "github.com/Zhukek/metrics/internal/model"
)

const pollInterval = 2
const reportInterval = 10

var client = &http.Client{}

func buildURLCounter(url, metric string, value int64) string {
	return fmt.Sprintf("%s/update/%s/%s/%d", url, models.Counter, metric, value)
}

func buildURLGauge(url, metric string, value float64) string {
	return fmt.Sprintf("%s/update/%s/%s/%f", url, models.Gauge, metric, value)
}

func postUpdate(url string) {
	res, err := client.Post(url, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	res.Body.Close()
}

func main() {
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
			postUpdate(buildURLCounter(url, "counter", counter))
			postUpdate(buildURLGauge(url, "RandomValue", RandomValue))
			postUpdate(buildURLGauge(url, "Alloc", float64(stat.Alloc)))
			postUpdate(buildURLGauge(url, "BuckHashSys", float64(stat.BuckHashSys)))
			postUpdate(buildURLGauge(url, "Frees", float64(stat.Frees)))
			postUpdate(buildURLGauge(url, "GCCPUFraction", float64(stat.GCCPUFraction)))
			postUpdate(buildURLGauge(url, "GCSys", float64(stat.GCSys)))
			postUpdate(buildURLGauge(url, "HeapAlloc", float64(stat.HeapAlloc)))
			postUpdate(buildURLGauge(url, "HeapIdle", float64(stat.HeapIdle)))
			postUpdate(buildURLGauge(url, "HeapInuse", float64(stat.HeapInuse)))
			postUpdate(buildURLGauge(url, "HeapObjects", float64(stat.HeapObjects)))
			postUpdate(buildURLGauge(url, "HeapReleased", float64(stat.HeapReleased)))
			postUpdate(buildURLGauge(url, "HeapSys", float64(stat.HeapSys)))
			postUpdate(buildURLGauge(url, "LastGC", float64(stat.LastGC)))
			postUpdate(buildURLGauge(url, "Lookups", float64(stat.Lookups)))
			postUpdate(buildURLGauge(url, "MCacheInuse", float64(stat.MCacheInuse)))
			postUpdate(buildURLGauge(url, "MCacheSys", float64(stat.MCacheSys)))
			postUpdate(buildURLGauge(url, "MSpanInuse", float64(stat.MSpanInuse)))
			postUpdate(buildURLGauge(url, "MSpanSys", float64(stat.MSpanSys)))
			postUpdate(buildURLGauge(url, "Mallocs", float64(stat.Mallocs)))
			postUpdate(buildURLGauge(url, "NextGC", float64(stat.NextGC)))
			postUpdate(buildURLGauge(url, "NumForcedGC", float64(stat.NumForcedGC)))
			postUpdate(buildURLGauge(url, "NumGC", float64(stat.NumGC)))
			postUpdate(buildURLGauge(url, "OtherSys", float64(stat.OtherSys)))
			postUpdate(buildURLGauge(url, "PauseTotalNs", float64(stat.PauseTotalNs)))
			postUpdate(buildURLGauge(url, "StackInuse", float64(stat.StackInuse)))
			postUpdate(buildURLGauge(url, "StackSys", float64(stat.StackSys)))
			postUpdate(buildURLGauge(url, "Sys", float64(stat.Sys)))
			postUpdate(buildURLGauge(url, "TotalAlloc", float64(stat.TotalAlloc)))

			counter = 0
			time.Sleep(reportInterval * time.Second)
		}
	}()
	select {}
}
