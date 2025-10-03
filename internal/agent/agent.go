package agent

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/Zhukek/metrics/internal/service/gzip"
	"github.com/Zhukek/metrics/internal/service/hash"
	"github.com/go-resty/resty/v2"
)

type StatsData struct {
	stat        runtime.MemStats
	TotalMemory float64
	FreeMemory  float64
	CPU         []float64
	counter     int64
	randomValue float64
}

func (s *StatsData) SetCounter(value int64) {
	s.counter = value
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
		SetError(&responseErr).
		SetBody(data).
		Post("/update/")

	if err != nil {
		if *iter < 3 {
			intervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
			await := intervals[*iter]
			*iter += 1

			time.AfterFunc(await, func() { postUpdate(client, metric, iter) })
		} else {
			fmt.Print("no response")
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
		TotalMemory   = float64(data.TotalMemory)
		FreeMemory    = float64(data.FreeMemory)
	)

	var metrics []models.Metrics

	metrics = append(metrics, models.Metrics{ID: "PollCount", MType: models.Counter, Delta: &counter})
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
	metrics = append(metrics, models.Metrics{ID: "TotalMemory", MType: models.Gauge, Value: &TotalMemory})
	metrics = append(metrics, models.Metrics{ID: "FreeMemory", MType: models.Gauge, Value: &FreeMemory})
	for i, v := range data.CPU {
		name := fmt.Sprintf("CPUutilization%d", i)
		metrics = append(metrics, models.Metrics{ID: name, MType: models.Gauge, Value: &v})
	}

	return metrics
}

func PollMemoryData(data *StatsData) {
	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Print("reading virtual mem error")
	}
	data.FreeMemory = float64(v.Free)
	data.TotalMemory = float64(v.Total)

	CPU, err := cpu.Percent(time.Second, true)
	if err != nil {
		fmt.Print("reading cpu percent error")
	}
	data.CPU = CPU
}

func Polling(data *StatsData) {
	data.randomValue = rand.Float64()
	runtime.ReadMemStats(&data.stat)
	data.counter++
}

func worker(wg *sync.WaitGroup, jobs <-chan models.Metrics, results chan<- models.Metrics, client *resty.Client) {
	defer wg.Done()
	for j := range jobs {
		postUpdate(client, j, nil)
		results <- j
	}
}

func PostUpdates(client *resty.Client, data *StatsData, rateLimit int) {

	metrics := getDataSlice(data)
	jobs := make(chan models.Metrics, len(metrics))
	results := make(chan models.Metrics, len(metrics))
	var wg sync.WaitGroup

	for w := 1; w <= rateLimit; w++ {
		wg.Add(1)
		go worker(&wg, jobs, results, client)
	}

	for _, v := range metrics {
		jobs <- v
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	close(jobs)

	for a := 1; a <= len(metrics); a++ {
		<-results
	}
}

func PostBatch(client *resty.Client, data *StatsData, iter *int, hasher *hash.Hasher) {
	metrics := getDataSlice(data)

	if len(metrics) == 0 {
		return
	}

	if iter == nil {
		i := 0
		iter = &i
	}

	body, err := json.Marshal(metrics)

	if err != nil {
		fmt.Print("marshal json")
		return
	}

	zipedBody, err := gzip.GzipCompress(body)
	if err != nil {
		fmt.Print("gzip Compress")
		return
	}

	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetError(&responseErr).
		SetBody(zipedBody)

	if hasher != nil {
		hashValue, err := hasher.Sign(body)
		if err != nil {
			fmt.Print("hasher sign")
			return
		}
		request.SetHeader("HashSHA256", hashValue)
	}

	_, err = request.Post("/updates/")

	if err != nil {
		if *iter < 3 {
			intervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
			await := intervals[*iter]
			*iter += 1

			time.AfterFunc(await, func() { PostBatch(client, data, iter, hasher) })
		} else {
			fmt.Print("no response")
			return
		}
	}
}
