package agent

import (
	"fmt"
	"time"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/go-resty/resty/v2"
)

type APIError struct {
	Code      int
	Message   string
	Timestapm time.Time
}

func BuildURLCounter(url, metric string, value int64) string {
	return fmt.Sprintf("%s/update/%s/%s/%d", url, models.Counter, metric, value)
}

func BuildURLGauge(url, metric string, value float64) string {
	return fmt.Sprintf("%s/update/%s/%s/%f", url, models.Gauge, metric, value)
}

func PostUpdate(client *resty.Client, url string) {
	var responseErr APIError

	_, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetError(&responseErr).
		Post(url)

	if err != nil {
		fmt.Println(responseErr)
		return
	}
}
