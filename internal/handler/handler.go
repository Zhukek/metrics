package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/go-chi/chi/v5"
)

func update(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")

	switch metricType {
	case models.Counter:
		value, err := strconv.Atoi(metricValue)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		models.UpdateCounter(metricName, int64(value))
	case models.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		models.UpdateGauge(metricName, value)
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func get(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	value, err := models.GetMetric(metricType, metricName)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	io.WriteString(res, fmt.Sprintf("%s=%s", metricName, value))
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func getList(res http.ResponseWriter, req *http.Request) {
	metrics := models.GetList()
	var metricList string
	for _, k := range metrics {
		metricList += fmt.Sprintf("<li>%s</li>\n", k)
	}

	io.WriteString(res, fmt.Sprintf(`
	<html>
	<body>
		<ul>
		%s
		</ul>
	</body>
	</html>`, metricList))
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func NewRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/update/{metricType}/{metricName}/{metricValue}", update)
	router.Get("/value/{metricType}/{metricName}", get)
	router.Get("/", getList)

	return router
}
