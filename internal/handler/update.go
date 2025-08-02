package handler

import (
	"net/http"
	"strconv"
	"strings"

	models "github.com/Zhukek/metrics/internal/model"
)

func Update(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(req.URL.Path, "/")

	if len(parts) < 5 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	metricType := parts[2]
	metricName := parts[3]
	metricValue := parts[4]

	switch metricType {
	case models.Counter:
		value, err := strconv.Atoi(metricValue)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		models.Storage.UpdateCounter(metricName, int64(value))
	case models.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		models.Storage.UpdateGauge(metricName, value)
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}
