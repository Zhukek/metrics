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
		res := BuildURLCounter(metric, counterValue)
		assert.Equal(t, counterExpected, res)
	})

	t.Run("gauge builder", func(t *testing.T) {
		res := BuildURLGauge(metric, gaugeValue)
		assert.Equal(t, gaugeExpected, res)
	})
}
