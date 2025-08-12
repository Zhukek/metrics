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

func (m *MemStorage) updateCounter(key string, value int64) {
	v, ok := m.metrics[key]

	if !ok {
		m.metrics[key] = Metrics{
			ID:    key,
			MType: Counter,
			Delta: &value,
		}
	} else {
		*v.Delta += value
	}
}

func (m *MemStorage) updateGauge(key string, value float64) {
	v, ok := m.metrics[key]

	if !ok {
		m.metrics[key] = Metrics{
			ID:    key,
			MType: Gauge,
			Value: &value,
		}
	} else {
		*v.Value = value
	}
}

func (m *MemStorage) getMetric(key string) (res string, err error) {
	v, ok := m.metrics[key]

	if !ok {
		err = ErrWrongMetric
		return
	}

	switch v.MType {
	case Counter:
		res = fmt.Sprintln(*v.Delta)
	case Gauge:
		res = fmt.Sprintln(*v.Value)
	default:
		err = ErrWrongMetric
	}

	return
}

func (m *MemStorage) getMetricsList() []string {
	var keys []string
	for key := range m.metrics {
		keys = append(keys, strings.Split(key, "_")[0])
	}
	return keys
}

var storage *MemStorage

func init() {
	storage = &MemStorage{
		metrics: make(map[string]Metrics),
	}
}

func UpdateCounter(key string, value int64) {
	reskey := key + "_" + Counter
	storage.updateCounter(reskey, value)
}

func UpdateGauge(key string, value float64) {
	reskey := key + "_" + Gauge
	storage.updateGauge(reskey, value)
}

func GetMetric(metricType, metricName string) (res string, err error) {
	reskey := metricName + "_" + metricType
	return storage.getMetric(reskey)
}

func GetList() []string {
	return storage.getMetricsList()
}
