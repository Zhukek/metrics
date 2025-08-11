package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCounter(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		key   string
		want  int64
	}{
		{
			name:  "first counter add",
			value: 64,
			key:   "counter1",
			want:  64,
		},
		{
			name:  "first counter update",
			value: 15,
			key:   "counter1",
			want:  79,
		},
		{
			name:  "second counter add",
			value: 20,
			key:   "counter2",
			want:  20,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			UpdateCounter(test.key, test.value)
			assert.Equal(t, test.want, *storage.metrics[test.key].Delta)
		})
	}
}

func TestGauge(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		key   string
		want  float64
	}{
		{
			name:  "first gauge add",
			value: 6.124342,
			key:   "gauge1",
			want:  6.124342,
		},
		{
			name:  "first gauge update",
			value: 15.214234,
			key:   "gauge1",
			want:  15.214234,
		},
		{
			name:  "second gauge add",
			value: 4.34534,
			key:   "gauge2",
			want:  4.34534,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			UpdateGauge(test.key, test.value)

			assert.Equal(t, test.want, *storage.metrics[test.key].Value)
		})
	}
}
