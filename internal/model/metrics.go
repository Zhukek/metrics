package models

import "fmt"

const (
	Counter = "counter"
	Gauge   = "gauge"
)

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

func (m *MemStorage) UpdateGauge(key string, value float64) {
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

	fmt.Println(m.metrics[key])
}

var Storage *MemStorage

func init() {
	Storage = &MemStorage{
		metrics: make(map[string]Metrics),
	}
}
