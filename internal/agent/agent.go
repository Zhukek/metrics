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

func GetBaseURL(URL string) string {
	if strings.Contains(URL, "://") {
		return URL
	} else {
		return "http://" + URL
	}
}

func postUpdate(client *resty.Client, metric models.Metrics, iter *int) {
	if iter == nil {
		i := 0
		iter = &i
	}

	data, err := json.Marshal(metric)

	if err != nil {
		fmt.Print("marshal json")
		return
	}

	data, err = gzip.GzipCompress(data)

	if err != nil {
		fmt.Print("gzip Compress")
		return
	}

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(data).
		Post("/update/")

	if err != nil {
		if *iter < 3 {
			await := (*iter * 2) + 1
			*iter += 1
			time.AfterFunc(time.Duration(await)*time.Second, func() { postUpdate(client, metric, iter) })
		} else {
			fmt.Println("no response")
			return
		}
	}
}

func getDataSlice(data *StatsData) []models.Metrics {
	var (
		counter       = data.counter
		RandomValue   = data.randomValue
		Alloc         = float64(data.stat.Alloc)
		BuckHashSys   = float64(data.stat.BuckHashSys)
		Frees         = float64(data.stat.Frees)
		GCCPUFraction = float64(data.stat.GCCPUFraction)
		GCSys         = float64(data.stat.GCSys)
		HeapAlloc     = float64(data.stat.HeapAlloc)
		HeapIdle      = float64(data.stat.HeapIdle)
		HeapInuse     = float64(data.stat.HeapInuse)
		HeapObjects   = float64(data.stat.HeapObjects)
		HeapReleased  = float64(data.stat.HeapReleased)
		HeapSys       = float64(data.stat.HeapSys)
		LastGC        = float64(data.stat.LastGC)
		Lookups       = float64(data.stat.Lookups)
		MCacheInuse   = float64(data.stat.MCacheInuse)
		MCacheSys     = float64(data.stat.MCacheSys)
		MSpanInuse    = float64(data.stat.MSpanInuse)
		MSpanSys      = float64(data.stat.MSpanSys)
		Mallocs       = float64(data.stat.Mallocs)
		NextGC        = float64(data.stat.NextGC)
		NumForcedGC   = float64(data.stat.NumForcedGC)
		NumGC         = float64(data.stat.NumGC)
		OtherSys      = float64(data.stat.OtherSys)
		PauseTotalNs  = float64(data.stat.PauseTotalNs)
		StackInuse    = float64(data.stat.StackInuse)
		StackSys      = float64(data.stat.StackSys)
		Sys           = float64(data.stat.Sys)
		TotalAlloc    = float64(data.stat.TotalAlloc)
	)

	var metrics []models.Metrics

	metrics = append(metrics, models.Metrics{ID: "PollCount", MType: models.Counter, Delta: &counter})
	metrics = append(metrics, models.Metrics{ID: "counter", MType: models.Counter, Delta: &counter})
	metrics = append(metrics, models.Metrics{ID: "RandomValue", MType: models.Gauge, Value: &RandomValue})
	metrics = append(metrics, models.Metrics{ID: "Alloc", MType: models.Gauge, Value: &Alloc})
	metrics = append(metrics, models.Metrics{ID: "BuckHashSys", MType: models.Gauge, Value: &BuckHashSys})
	metrics = append(metrics, models.Metrics{ID: "Frees", MType: models.Gauge, Value: &Frees})
	metrics = append(metrics, models.Metrics{ID: "GCCPUFraction", MType: models.Gauge, Value: &GCCPUFraction})
	metrics = append(metrics, models.Metrics{ID: "GCSys", MType: models.Gauge, Value: &GCSys})
	metrics = append(metrics, models.Metrics{ID: "HeapAlloc", MType: models.Gauge, Value: &HeapAlloc})
	metrics = append(metrics, models.Metrics{ID: "HeapIdle", MType: models.Gauge, Value: &HeapIdle})
	metrics = append(metrics, models.Metrics{ID: "HeapInuse", MType: models.Gauge, Value: &HeapInuse})
	metrics = append(metrics, models.Metrics{ID: "HeapObjects", MType: models.Gauge, Value: &HeapObjects})
	metrics = append(metrics, models.Metrics{ID: "HeapReleased", MType: models.Gauge, Value: &HeapReleased})
	metrics = append(metrics, models.Metrics{ID: "HeapSys", MType: models.Gauge, Value: &HeapSys})
	metrics = append(metrics, models.Metrics{ID: "LastGC", MType: models.Gauge, Value: &LastGC})
	metrics = append(metrics, models.Metrics{ID: "Lookups", MType: models.Gauge, Value: &Lookups})
	metrics = append(metrics, models.Metrics{ID: "MCacheInuse", MType: models.Gauge, Value: &MCacheInuse})
	metrics = append(metrics, models.Metrics{ID: "MCacheSys", MType: models.Gauge, Value: &MCacheSys})
	metrics = append(metrics, models.Metrics{ID: "MSpanInuse", MType: models.Gauge, Value: &MSpanInuse})
	metrics = append(metrics, models.Metrics{ID: "MSpanSys", MType: models.Gauge, Value: &MSpanSys})
	metrics = append(metrics, models.Metrics{ID: "Mallocs", MType: models.Gauge, Value: &Mallocs})
	metrics = append(metrics, models.Metrics{ID: "NextGC", MType: models.Gauge, Value: &NextGC})
	metrics = append(metrics, models.Metrics{ID: "NumForcedGC", MType: models.Gauge, Value: &NumForcedGC})
	metrics = append(metrics, models.Metrics{ID: "NumGC", MType: models.Gauge, Value: &NumGC})
	metrics = append(metrics, models.Metrics{ID: "OtherSys", MType: models.Gauge, Value: &OtherSys})
	metrics = append(metrics, models.Metrics{ID: "PauseTotalNs", MType: models.Gauge, Value: &PauseTotalNs})
	metrics = append(metrics, models.Metrics{ID: "StackInuse", MType: models.Gauge, Value: &StackInuse})
	metrics = append(metrics, models.Metrics{ID: "StackSys", MType: models.Gauge, Value: &StackSys})
	metrics = append(metrics, models.Metrics{ID: "Sys", MType: models.Gauge, Value: &Sys})
	metrics = append(metrics, models.Metrics{ID: "TotalAlloc", MType: models.Gauge, Value: &TotalAlloc})

	return metrics
}

func Polling(data *StatsData) {
	data.randomValue = rand.Float64()
	runtime.ReadMemStats(&data.stat)
	data.counter++
}

func PostUpdates(client *resty.Client, data *StatsData) {

	metrics := getDataSlice(data)

	for _, v := range metrics {
		postUpdate(client, v, nil)
	}

	data.counter = 0
}

func PostBatch(client *resty.Client, data *StatsData, iter *int) {
	metrics := getDataSlice(data)

	if len(metrics) == 0 {
		fmt.Println("emty batch")
		return
	}

	if iter == nil {
		i := 0
		iter = &i
	}

	body, err := json.Marshal(metrics)

	if err != nil {
		fmt.Println("marshal json")
		return
	}

	body, err = gzip.GzipCompress(body)
	if err != nil {
		fmt.Println("gzip Compress")
		return
	}

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(body).
		Post("/updates/")

	if err != nil {
		if *iter < 3 {
			await := (*iter * 2) + 1
			*iter += 1
			time.AfterFunc(time.Duration(await)*time.Second, func() { PostBatch(client, data, iter) })
		} else {
			fmt.Println("no response")
			return
		}
	}
}
