package agent

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"runtime"
	"strings"
	"time"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/Zhukek/metrics/internal/service/gzip"
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

var responseErr APIError

func GetBaseURL(URL string) string {
	if strings.Contains(URL, "://") {
		return URL
	} else {
		return "http://" + URL
	}
}

func postUpdate(client *resty.Client, metric models.MetricsBody) {
	data, err := json.Marshal(metric)

	if err != nil {
		fmt.Print("Error: Marshal json")
		return
	}

	data, err = gzip.GzipCompress(data)

	if err != nil {
		fmt.Print("Error: Gzip Compress")
		return
	}

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetError(&responseErr).
		SetBody(data).
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
	postUpdate(client, models.MetricsBody{ID: "PollCount", MType: models.Counter, Delta: data.counter})
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

func PostBatch(client *resty.Client, data *StatsData) {
	metrics := []models.MetricsBody{
		models.MetricsBody{ID: "PollCount", MType: models.Counter, Delta: data.counter},
		models.MetricsBody{ID: "counter", MType: models.Counter, Delta: data.counter},
		models.MetricsBody{ID: "RandomValue", MType: models.Gauge, Value: data.randomValue},
		models.MetricsBody{ID: "Alloc", MType: models.Gauge, Value: float64(data.stat.Alloc)},
		models.MetricsBody{ID: "BuckHashSys", MType: models.Gauge, Value: float64(data.stat.BuckHashSys)},
		models.MetricsBody{ID: "Frees", MType: models.Gauge, Value: float64(data.stat.Frees)},
		models.MetricsBody{ID: "GCCPUFraction", MType: models.Gauge, Value: float64(data.stat.GCCPUFraction)},
		models.MetricsBody{ID: "GCSys", MType: models.Gauge, Value: float64(data.stat.GCSys)},
		models.MetricsBody{ID: "HeapAlloc", MType: models.Gauge, Value: float64(data.stat.HeapAlloc)},
		models.MetricsBody{ID: "HeapIdle", MType: models.Gauge, Value: float64(data.stat.HeapIdle)},
		models.MetricsBody{ID: "HeapInuse", MType: models.Gauge, Value: float64(data.stat.HeapInuse)},
		models.MetricsBody{ID: "HeapObjects", MType: models.Gauge, Value: float64(data.stat.HeapObjects)},
		models.MetricsBody{ID: "HeapReleased", MType: models.Gauge, Value: float64(data.stat.HeapReleased)},
		models.MetricsBody{ID: "HeapSys", MType: models.Gauge, Value: float64(data.stat.HeapSys)},
		models.MetricsBody{ID: "LastGC", MType: models.Gauge, Value: float64(data.stat.LastGC)},
		models.MetricsBody{ID: "Lookups", MType: models.Gauge, Value: float64(data.stat.Lookups)},
		models.MetricsBody{ID: "MCacheInuse", MType: models.Gauge, Value: float64(data.stat.MCacheInuse)},
		models.MetricsBody{ID: "MCacheSys", MType: models.Gauge, Value: float64(data.stat.MCacheSys)},
		models.MetricsBody{ID: "MSpanInuse", MType: models.Gauge, Value: float64(data.stat.MSpanInuse)},
		models.MetricsBody{ID: "MSpanSys", MType: models.Gauge, Value: float64(data.stat.MSpanSys)},
		models.MetricsBody{ID: "Mallocs", MType: models.Gauge, Value: float64(data.stat.Mallocs)},
		models.MetricsBody{ID: "NextGC", MType: models.Gauge, Value: float64(data.stat.NextGC)},
		models.MetricsBody{ID: "NumForcedGC", MType: models.Gauge, Value: float64(data.stat.NumForcedGC)},
		models.MetricsBody{ID: "NumGC", MType: models.Gauge, Value: float64(data.stat.NumGC)},
		models.MetricsBody{ID: "OtherSys", MType: models.Gauge, Value: float64(data.stat.OtherSys)},
		models.MetricsBody{ID: "PauseTotalNs", MType: models.Gauge, Value: float64(data.stat.PauseTotalNs)},
		models.MetricsBody{ID: "StackInuse", MType: models.Gauge, Value: float64(data.stat.StackInuse)},
		models.MetricsBody{ID: "StackSys", MType: models.Gauge, Value: float64(data.stat.StackSys)},
		models.MetricsBody{ID: "Sys", MType: models.Gauge, Value: float64(data.stat.Sys)},
		models.MetricsBody{ID: "TotalAlloc", MType: models.Gauge, Value: float64(data.stat.TotalAlloc)},
	}

	body, err := json.Marshal(metrics)

	if err != nil {
		fmt.Print("Error: Marshal json")
		return
	}

	body, err = gzip.GzipCompress(body)
	if err != nil {
		fmt.Print("Error: Gzip Compress")
		return
	}

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetError(&responseErr).
		SetBody(body).
		Post("/updates/")

	if err != nil {
		fmt.Printf("Error: %s\n", responseErr.Message)
		return
	}
}
