package inmemory

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/Zhukek/metrics/internal/repository"
)

var ErrWrongMetric = errors.New("wrong metric")

type MemStorage struct {
	metrics map[string]models.Metrics
}

func (m *MemStorage) UpdateCounter(key string, value int64) {
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
}

func (m *MemStorage) UpdateGauge(key string, value float64) {
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
}

func (m *MemStorage) GetMetric(metricType, metricName string) (res string, err error) {
	reskey := metricName + "_" + metricType
	v, ok := m.metrics[reskey]

	if !ok {
		err = ErrWrongMetric
		return
	}

	switch v.MType {
	case models.Counter:
		res = fmt.Sprint(*v.Delta)
	case models.Gauge:
		res = fmt.Sprint(*v.Value)
	default:
		err = ErrWrongMetric
	}

	return
}

func (m *MemStorage) GetMetricv2(body models.Metrics) (metricBody models.Metrics, err error) {
	reskey := body.ID + "_" + body.MType
	v, ok := m.metrics[reskey]

	if !ok {
		err = ErrWrongMetric
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

func (m *MemStorage) GetList() []string {
	var keys []string
	for key := range m.metrics {
		keys = append(keys, strings.Split(key, "_")[0])
	}
	return keys
}

func (m *MemStorage) GetAllMetrics() map[string]models.Metrics {
	return m.metrics
}

func NewStorage(data []byte) (repository.Repository, error) {
	metrics := make(map[string]models.Metrics)

	if len(data) > 0 {
		if err := json.Unmarshal(data, &metrics); err != nil {
			return nil, err
		}
	}

	storage := MemStorage{
		metrics: metrics,
	}

	return &storage, nil
}
