package agent

import (
	"fmt"
	"testing"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestBuilders(t *testing.T) {
	const url = "http://localhost:8080"
	metric := "testMetric"
	var counterValue int64 = 45
	var gaugeValue = 24.2344432
	counterExpected := fmt.Sprintf("%s/update/%s/%s/%d", url, models.Counter, metric, counterValue)
	gaugeExpected := fmt.Sprintf("%s/update/%s/%s/%f", url, models.Gauge, metric, gaugeValue)

	t.Run("counter builder", func(t *testing.T) {
		res := BuildURLCounter(url, metric, counterValue)
		assert.Equal(t, counterExpected, res)
	})

	t.Run("gauge builder", func(t *testing.T) {
		res := BuildURLGauge(url, metric, gaugeValue)
		assert.Equal(t, gaugeExpected, res)
	})
}
