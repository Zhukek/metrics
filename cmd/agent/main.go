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

func buildUrlCounter(url, metric string, value int64) string {
	return fmt.Sprintf("%s/update/%s/%s/%d", url, models.Counter, metric, value)
}

func buildUrlGauge(url, metric string, value float64) string {
	return fmt.Sprintf("%s/update/%s/%s/%f", url, models.Gauge, metric, value)
}

func postUpdate(url string) {
	res, err := client.Post(url, "text/plain", nil)
	if err != nil {
		panic(err)
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
			postUpdate(buildUrlCounter(url, "counter", counter))
			postUpdate(buildUrlGauge(url, "RandomValue", RandomValue))
			postUpdate(buildUrlGauge(url, "Alloc", float64(stat.Alloc)))
			postUpdate(buildUrlGauge(url, "BuckHashSys", float64(stat.BuckHashSys)))
			postUpdate(buildUrlGauge(url, "Frees", float64(stat.Frees)))
			postUpdate(buildUrlGauge(url, "GCCPUFraction", float64(stat.GCCPUFraction)))
			postUpdate(buildUrlGauge(url, "GCSys", float64(stat.GCSys)))
			postUpdate(buildUrlGauge(url, "HeapAlloc", float64(stat.HeapAlloc)))
			postUpdate(buildUrlGauge(url, "HeapIdle", float64(stat.HeapIdle)))
			postUpdate(buildUrlGauge(url, "HeapInuse", float64(stat.HeapInuse)))
			postUpdate(buildUrlGauge(url, "HeapObjects", float64(stat.HeapObjects)))
			postUpdate(buildUrlGauge(url, "HeapReleased", float64(stat.HeapReleased)))
			postUpdate(buildUrlGauge(url, "HeapSys", float64(stat.HeapSys)))
			postUpdate(buildUrlGauge(url, "LastGC", float64(stat.LastGC)))
			postUpdate(buildUrlGauge(url, "Lookups", float64(stat.Lookups)))
			postUpdate(buildUrlGauge(url, "MCacheInuse", float64(stat.MCacheInuse)))
			postUpdate(buildUrlGauge(url, "MCacheSys", float64(stat.MCacheSys)))
			postUpdate(buildUrlGauge(url, "MSpanInuse", float64(stat.MSpanInuse)))
			postUpdate(buildUrlGauge(url, "MSpanSys", float64(stat.MSpanSys)))
			postUpdate(buildUrlGauge(url, "Mallocs", float64(stat.Mallocs)))
			postUpdate(buildUrlGauge(url, "NextGC", float64(stat.NextGC)))
			postUpdate(buildUrlGauge(url, "NumForcedGC", float64(stat.NumForcedGC)))
			postUpdate(buildUrlGauge(url, "NumGC", float64(stat.NumGC)))
			postUpdate(buildUrlGauge(url, "OtherSys", float64(stat.OtherSys)))
			postUpdate(buildUrlGauge(url, "PauseTotalNs", float64(stat.PauseTotalNs)))
			postUpdate(buildUrlGauge(url, "StackInuse", float64(stat.StackInuse)))
			postUpdate(buildUrlGauge(url, "StackSys", float64(stat.StackSys)))
			postUpdate(buildUrlGauge(url, "Sys", float64(stat.Sys)))
			postUpdate(buildUrlGauge(url, "TotalAlloc", float64(stat.TotalAlloc)))

			counter = 0
			time.Sleep(reportInterval * time.Second)
		}
	}()
	select {}
}
