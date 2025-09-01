package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/go-chi/chi/v5"
)

func updatev1(res http.ResponseWriter, req *http.Request, storage models.MemStorage) {
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
		storage.UpdateCounter(metricName, int64(value))
	case models.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.UpdateGauge(metricName, value)
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func updatev2(res http.ResponseWriter, req *http.Request, storage models.MemStorage) {
	res.Header().Set("Content-Type", "application/json")

	var metric models.MetricsBody

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&metric); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch metric.MType {
	case models.Counter:
		storage.UpdateCounter(metric.ID, metric.Delta)
	case models.Gauge:
		storage.UpdateGauge(metric.ID, metric.Value)
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func getv1(res http.ResponseWriter, req *http.Request, storage models.MemStorage) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	value, err := storage.GetMetric(metricType, metricName)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	io.WriteString(res, value)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func getv2(res http.ResponseWriter, req *http.Request, storage models.MemStorage) {
	res.Header().Set("Content-Type", "application/json")
	var metric models.MetricsBody

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&metric); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := storage.GetMetricv2(metric)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(res)
	if err := encoder.Encode(value); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func getList(res http.ResponseWriter, req *http.Request, storage models.MemStorage) {
	metrics := storage.GetList()
	const markup = `
	<html>
	<body>
		<ul>
		{{range .}}
			<li>{{.}}</li>
		{{end}}
		</ul>
	</body>
	</html>`

	page, err := template.New("Response").Parse(markup)

	if err != nil {
		fmt.Println("Error:", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = page.Execute(res, metrics)

	if err != nil {
		fmt.Println("Error:", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func NewRouter(storage models.MemStorage) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		updatev2(w, r, storage)
	})
	router.Post("/update/{metricType}/{metricName}/{metricValue}", func(w http.ResponseWriter, r *http.Request) {
		updatev1(w, r, storage)
	})
	router.Get("/value/{metricType}/{metricName}", func(w http.ResponseWriter, r *http.Request) {
		getv1(w, r, storage)
	})
	router.Post("/value", func(w http.ResponseWriter, r *http.Request) {
		getv2(w, r, storage)
	})
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		getList(w, r, storage)
	})

	return router
}
