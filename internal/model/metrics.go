package models

import "errors"

type MType string

func (m MType) String() string {
	return string(m)
}

const (
	Counter MType = "counter"
	Gauge   MType = "gauge"
)

var ErrWrongMetric = errors.New("wrong metric")

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metrics struct {
	ID    string   `json:"id"`
	MType MType    `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}
