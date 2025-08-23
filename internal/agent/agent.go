package agent

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"strings"
	"time"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/go-resty/resty/v2"
)

type StatsData struct {
	stat        runtime.MemStats
	counter     int64
	randomValue float64
}

type APIError struct {
	Code      int
	Message   string
	Timestapm time.Time
}

func GetBaseURL(URL string) string {
	if strings.Contains(URL, "://") {
		return URL
	} else {
		return "http://" + URL
	}
}

func buildURLCounter(metric string, value int64) string {
	return fmt.Sprintf("/update/%s/%s/%d", models.Counter, metric, value)
}

func buildURLGauge(metric string, value float64) string {
	return fmt.Sprintf("/update/%s/%s/%f", models.Gauge, metric, value)
}

func postUpdate(client *resty.Client, url string) {
	var responseErr APIError

	_, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetError(&responseErr).
		Post(url)

	if err != nil {
		fmt.Printf("Error: %s\n", responseErr.Message)
		return
	}
}

func Polling(data *StatsData) {
	data.randomValue = rand.Float64()
	runtime.ReadMemStats(&data.stat)
	data.counter++
}

func PostUpdates(client *resty.Client, data *StatsData) {
	postUpdate(client, buildURLCounter("counter", data.counter))
	postUpdate(client, buildURLGauge("RandomValue", data.randomValue))
	postUpdate(client, buildURLGauge("Alloc", float64(data.stat.Alloc)))
	postUpdate(client, buildURLGauge("BuckHashSys", float64(data.stat.BuckHashSys)))
	postUpdate(client, buildURLGauge("Frees", float64(data.stat.Frees)))
	postUpdate(client, buildURLGauge("GCCPUFraction", float64(data.stat.GCCPUFraction)))
	postUpdate(client, buildURLGauge("GCSys", float64(data.stat.GCSys)))
	postUpdate(client, buildURLGauge("HeapAlloc", float64(data.stat.HeapAlloc)))
	postUpdate(client, buildURLGauge("HeapIdle", float64(data.stat.HeapIdle)))
	postUpdate(client, buildURLGauge("HeapInuse", float64(data.stat.HeapInuse)))
	postUpdate(client, buildURLGauge("HeapObjects", float64(data.stat.HeapObjects)))
	postUpdate(client, buildURLGauge("HeapReleased", float64(data.stat.HeapReleased)))
	postUpdate(client, buildURLGauge("HeapSys", float64(data.stat.HeapSys)))
	postUpdate(client, buildURLGauge("LastGC", float64(data.stat.LastGC)))
	postUpdate(client, buildURLGauge("Lookups", float64(data.stat.Lookups)))
	postUpdate(client, buildURLGauge("MCacheInuse", float64(data.stat.MCacheInuse)))
	postUpdate(client, buildURLGauge("MCacheSys", float64(data.stat.MCacheSys)))
	postUpdate(client, buildURLGauge("MSpanInuse", float64(data.stat.MSpanInuse)))
	postUpdate(client, buildURLGauge("MSpanSys", float64(data.stat.MSpanSys)))
	postUpdate(client, buildURLGauge("Mallocs", float64(data.stat.Mallocs)))
	postUpdate(client, buildURLGauge("NextGC", float64(data.stat.NextGC)))
	postUpdate(client, buildURLGauge("NumForcedGC", float64(data.stat.NumForcedGC)))
	postUpdate(client, buildURLGauge("NumGC", float64(data.stat.NumGC)))
	postUpdate(client, buildURLGauge("OtherSys", float64(data.stat.OtherSys)))
	postUpdate(client, buildURLGauge("PauseTotalNs", float64(data.stat.PauseTotalNs)))
	postUpdate(client, buildURLGauge("StackInuse", float64(data.stat.StackInuse)))
	postUpdate(client, buildURLGauge("StackSys", float64(data.stat.StackSys)))
	postUpdate(client, buildURLGauge("Sys", float64(data.stat.Sys)))
	postUpdate(client, buildURLGauge("TotalAlloc", float64(data.stat.TotalAlloc)))

	data.counter = 0
}
