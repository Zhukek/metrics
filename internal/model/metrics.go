package models

import (
	"errors"
	"fmt"
	"strings"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

var ErrWrongMetric = errors.New("wrong metric")

type MetricsBody struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Delta int64   `json:"delta,omitempty"`
	Value float64 `json:"value,omitempty"`
}

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

type MemStorage struct {
	metrics map[string]Metrics
}

func (m *MemStorage) UpdateCounter(key string, value int64) {
	reskey := key + "_" + Counter
	v, ok := m.metrics[reskey]

	if !ok {
		m.metrics[reskey] = Metrics{
			ID:    reskey,
			MType: Counter,
			Delta: &value,
		}
	} else {
		*v.Delta += value
	}
}

func (m *MemStorage) UpdateGauge(key string, value float64) {
	reskey := key + "_" + Gauge
	v, ok := m.metrics[reskey]

	if !ok {
		m.metrics[reskey] = Metrics{
			ID:    reskey,
			MType: Gauge,
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
	case Counter:
		res = fmt.Sprint(*v.Delta)
	case Gauge:
		res = fmt.Sprint(*v.Value)
	default:
		err = ErrWrongMetric
	}

	return
}

func (m *MemStorage) GetMetricv2(body MetricsBody) (metricBody MetricsBody, err error) {
	reskey := body.ID + "_" + body.MType
	v, ok := m.metrics[reskey]

	if !ok {
		err = ErrWrongMetric
		return
	}

	metricBody = MetricsBody{
		ID:    v.ID,
		MType: v.MType,
		Delta: *v.Delta,
		Value: *v.Value,
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

/* var storage *MemStorage */

func NewStorage() MemStorage {
	storage := MemStorage{
		metrics: make(map[string]Metrics),
	}

	return storage
}
