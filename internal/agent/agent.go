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

func postUpdate(client *resty.Client, metric models.MetricsBody) {
	var responseErr APIError

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetError(&responseErr).
		SetBody(metric).
		Post("/update/")

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
	postUpdate(client, models.MetricsBody{ID: "counter", MType: models.Counter, Delta: data.counter})
	postUpdate(client, models.MetricsBody{ID: "RandomValue", MType: models.Gauge, Value: data.randomValue})
	postUpdate(client, models.MetricsBody{ID: "Alloc", MType: models.Gauge, Value: float64(data.stat.Alloc)})
	postUpdate(client, models.MetricsBody{ID: "BuckHashSys", MType: models.Gauge, Value: float64(data.stat.BuckHashSys)})
	postUpdate(client, models.MetricsBody{ID: "Frees", MType: models.Gauge, Value: float64(data.stat.Frees)})
	postUpdate(client, models.MetricsBody{ID: "GCCPUFraction", MType: models.Gauge, Value: float64(data.stat.GCCPUFraction)})
	postUpdate(client, models.MetricsBody{ID: "GCSys", MType: models.Gauge, Value: float64(data.stat.GCSys)})
	postUpdate(client, models.MetricsBody{ID: "HeapAlloc", MType: models.Gauge, Value: float64(data.stat.HeapAlloc)})
	postUpdate(client, models.MetricsBody{ID: "HeapIdle", MType: models.Gauge, Value: float64(data.stat.HeapIdle)})
	postUpdate(client, models.MetricsBody{ID: "HeapInuse", MType: models.Gauge, Value: float64(data.stat.HeapInuse)})
	postUpdate(client, models.MetricsBody{ID: "HeapObjects", MType: models.Gauge, Value: float64(data.stat.HeapObjects)})
	postUpdate(client, models.MetricsBody{ID: "HeapReleased", MType: models.Gauge, Value: float64(data.stat.HeapReleased)})
	postUpdate(client, models.MetricsBody{ID: "HeapSys", MType: models.Gauge, Value: float64(data.stat.HeapSys)})
	postUpdate(client, models.MetricsBody{ID: "LastGC", MType: models.Gauge, Value: float64(data.stat.LastGC)})
	postUpdate(client, models.MetricsBody{ID: "Lookups", MType: models.Gauge, Value: float64(data.stat.Lookups)})
	postUpdate(client, models.MetricsBody{ID: "MCacheInuse", MType: models.Gauge, Value: float64(data.stat.MCacheInuse)})
	postUpdate(client, models.MetricsBody{ID: "MCacheSys", MType: models.Gauge, Value: float64(data.stat.MCacheSys)})
	postUpdate(client, models.MetricsBody{ID: "MSpanInuse", MType: models.Gauge, Value: float64(data.stat.MSpanInuse)})
	postUpdate(client, models.MetricsBody{ID: "MSpanSys", MType: models.Gauge, Value: float64(data.stat.MSpanSys)})
	postUpdate(client, models.MetricsBody{ID: "Mallocs", MType: models.Gauge, Value: float64(data.stat.Mallocs)})
	postUpdate(client, models.MetricsBody{ID: "NextGC", MType: models.Gauge, Value: float64(data.stat.NextGC)})
	postUpdate(client, models.MetricsBody{ID: "NumForcedGC", MType: models.Gauge, Value: float64(data.stat.NumForcedGC)})
	postUpdate(client, models.MetricsBody{ID: "NumGC", MType: models.Gauge, Value: float64(data.stat.NumGC)})
	postUpdate(client, models.MetricsBody{ID: "OtherSys", MType: models.Gauge, Value: float64(data.stat.OtherSys)})
	postUpdate(client, models.MetricsBody{ID: "PauseTotalNs", MType: models.Gauge, Value: float64(data.stat.PauseTotalNs)})
	postUpdate(client, models.MetricsBody{ID: "StackInuse", MType: models.Gauge, Value: float64(data.stat.StackInuse)})
	postUpdate(client, models.MetricsBody{ID: "StackSys", MType: models.Gauge, Value: float64(data.stat.StackSys)})
	postUpdate(client, models.MetricsBody{ID: "Sys", MType: models.Gauge, Value: float64(data.stat.Sys)})
	postUpdate(client, models.MetricsBody{ID: "TotalAlloc", MType: models.Gauge, Value: float64(data.stat.TotalAlloc)})

	data.counter = 0
}
