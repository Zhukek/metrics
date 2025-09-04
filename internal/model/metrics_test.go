package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var initData []byte
var testStorage, err = NewStorage(initData)

func TestNewStorage(t *testing.T) {
	t.Run("create storage", func(t *testing.T) {
		require.NoError(t, err)
	})
}

func TestUpdateCounter(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		key   string
		want  string
	}{
		{
			name:  "first counter add",
			value: 64,
			key:   "counter1",
			want:  "64",
		},
		{
			name:  "first counter update",
			value: 15,
			key:   "counter1",
			want:  "79",
		},
		{
			name:  "second counter add",
			value: 20,
			key:   "counter2",
			want:  "20",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testStorage.UpdateCounter(test.key, test.value)
			res, err := testStorage.GetMetric(Counter, test.key)
			require.NoError(t, err)
			require.Equal(t, test.want, res)
		})
	}
}

func TestGauge(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		key   string
		want  string
	}{
		{
			name:  "first gauge add",
			value: 6.124342,
			key:   "gauge1",
			want:  "6.124342",
		},
		{
			name:  "first gauge update",
			value: 15.214234,
			key:   "gauge1",
			want:  "15.214234",
		},
		{
			name:  "second gauge add",
			value: 4.34534,
			key:   "gauge2",
			want:  "4.34534",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testStorage.UpdateGauge(test.key, test.value)
			res, err := testStorage.GetMetric(Gauge, test.key)
			require.NoError(t, err)
			require.Equal(t, test.want, res)
		})
	}
}
