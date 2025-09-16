package inmemory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/Zhukek/metrics/internal/service/fileworker"
)

type MemStorage struct {
	metrics map[string]models.Metrics
}

func (m *MemStorage) UpdateCounter(key string, value int64) error {
	reskey := key + "_" + models.Counter
	v, ok := m.metrics[reskey]

	if !ok {
		m.metrics[reskey] = models.Metrics{
			ID:    reskey,
			MType: models.Counter,
			Delta: &value,
		}
	} else {
		*v.Delta += value
	}
	return nil
}

func (m *MemStorage) UpdateGauge(key string, value float64) error {
	reskey := key + "_" + models.Gauge
	v, ok := m.metrics[reskey]

	if !ok {
		m.metrics[reskey] = models.Metrics{
			ID:    reskey,
			MType: models.Gauge,
			Value: &value,
		}
	} else {
		*v.Value = value
	}
	return nil
}

func (m *MemStorage) GetMetric(metricType, metricName string) (res string, err error) {
	reskey := metricName + "_" + metricType
	v, ok := m.metrics[reskey]

	if !ok {
		err = models.ErrWrongMetric
		return
	}

	switch v.MType {
	case models.Counter:
		res = fmt.Sprint(*v.Delta)
	case models.Gauge:
		res = fmt.Sprint(*v.Value)
	default:
		err = models.ErrWrongMetric
	}

	return
}

func (m *MemStorage) GetMetricv2(body models.Metrics) (metricBody models.Metrics, err error) {
	reskey := body.ID + "_" + body.MType
	v, ok := m.metrics[reskey]

	if !ok {
		err = models.ErrWrongMetric
		return
	}

	metricBody = models.Metrics{
		ID:    body.ID,
		MType: body.MType,
		Delta: v.Delta,
		Value: v.Value,
	}

	return
}

func (m *MemStorage) GetList() ([]string, error) {
	var keys []string
	for key := range m.metrics {
		keys = append(keys, strings.Split(key, "_")[0])
	}
	return keys, nil
}

func (m *MemStorage) GetAllMetrics() map[string]models.Metrics {
	return m.metrics
}

func (m *MemStorage) Ping(context.Context) error {
	return fmt.Errorf("in memory")
}

func NewStorage(filePath string, interval int, restore bool) (*MemStorage, error) {
	metrics := make(map[string]models.Metrics)
	var (
		data       []byte
		fileWorker *fileworker.FileWorker
		err        error
	)

	if filePath != "" {
		fileWorker, err = fileworker.NewFileWorker(filePath, interval == 0)

		if err != nil {
			return nil, err
		}
		defer fileWorker.Close()

		if restore {
			data, err = fileWorker.ReadData()
			if err != nil {
				return nil, err
			}
		}
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &metrics); err != nil {
			return nil, err
		}
	}

	storage := MemStorage{
		metrics: metrics,
	}

	if fileWorker != nil {
		switch interval {
		case 0:
			if err := fileWorker.WriteData(storage.GetAllMetrics()); err != nil {
				fmt.Println(err)
				// #TODO сделать нормальное логирование
			}
		default:
			ticker := time.NewTicker(time.Duration(interval) * time.Second)
			defer ticker.Stop()

			go func() {
				for range ticker.C {
					if err := fileWorker.WriteData(storage.GetAllMetrics()); err != nil {
						fmt.Println(err)
						// #TODO сделать нормальное логирование
					}
				}
			}()
		}
	}

	return &storage, nil
}
