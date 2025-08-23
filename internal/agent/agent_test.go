package agent

import (
	"fmt"
	"testing"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestBuilders(t *testing.T) {
	metric := "testMetric"
	var counterValue int64 = 45
	var gaugeValue = 24.2344432
	counterExpected := fmt.Sprintf("/update/%s/%s/%d", models.Counter, metric, counterValue)
	gaugeExpected := fmt.Sprintf("/update/%s/%s/%f", models.Gauge, metric, gaugeValue)

	t.Run("counter builder", func(t *testing.T) {
		res := buildURLCounter(metric, counterValue)
		assert.Equal(t, counterExpected, res)
	})

	t.Run("gauge builder", func(t *testing.T) {
		res := buildURLGauge(metric, gaugeValue)
		assert.Equal(t, gaugeExpected, res)
	})
}

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name string
		URL  string
		want string
	}{
		{
			name: "http protocol",
			URL:  "http://localhost:3000",
			want: "http://localhost:3000",
		},
		{
			name: "https protocol",
			URL:  "https://localhost:3000",
			want: "https://localhost:3000",
		},
		{
			name: "without protocol",
			URL:  "localhost:3000",
			want: "http://localhost:3000",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, GetBaseURL(test.URL))
		})
	}
}
