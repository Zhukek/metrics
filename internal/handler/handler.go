package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/Zhukek/metrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

func updatev1(res http.ResponseWriter, req *http.Request, storage repository.Repository) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")

	switch metricType {
	case models.Counter.String():
		value, err := strconv.Atoi(metricValue)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = storage.UpdateCounter(metricName, int64(value)); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	case models.Gauge.String():
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = storage.UpdateGauge(metricName, value); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func updatev2(res http.ResponseWriter, req *http.Request, storage repository.Repository) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	var metric models.Metrics

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&metric); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch metric.MType {
	case models.Counter:
		if metric.Delta == nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := storage.UpdateCounter(metric.ID, *metric.Delta); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	case models.Gauge:
		if metric.Value == nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := storage.UpdateGauge(metric.ID, *metric.Value); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func updates(res http.ResponseWriter, req *http.Request, storage repository.Repository) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	var metrics []models.Metrics

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&metrics); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// Не обрабатываем пустые батчи
	if len(metrics) == 0 {
		res.WriteHeader(http.StatusOK)
		return
	}

	if err := storage.Updates(metrics); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func getv1(res http.ResponseWriter, req *http.Request, storage repository.Repository) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	value, err := storage.GetMetric(models.MType(metricType), metricName)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	io.WriteString(res, value)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func getv2(res http.ResponseWriter, req *http.Request, storage repository.Repository) {
	res.Header().Set("Content-Type", "application/json")
	var metric models.Metrics

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
		return
	}
}

func getList(res http.ResponseWriter, req *http.Request, storage repository.Repository) {
	metrics, err := storage.GetList()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

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
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = page.Execute(res, metrics)

	if err != nil {
		fmt.Println("Error:", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ping(res http.ResponseWriter, req *http.Request, storage repository.Repository) {
	ctx, cancel := context.WithTimeout(context.Background(), 11*time.Second)
	defer cancel()
	if err := storage.Ping(ctx); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func NewRouter(storage repository.Repository) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/update/", func(w http.ResponseWriter, r *http.Request) {
		updatev2(w, r, storage)
	})
	router.Post("/updates/", func(w http.ResponseWriter, r *http.Request) {
		updates(w, r, storage)
	})
	router.Post("/update/{metricType}/{metricName}/{metricValue}", func(w http.ResponseWriter, r *http.Request) {
		updatev1(w, r, storage)
	})
	router.Get("/value/{metricType}/{metricName}", func(w http.ResponseWriter, r *http.Request) {
		getv1(w, r, storage)
	})
	router.Post("/value/", func(w http.ResponseWriter, r *http.Request) {
		getv2(w, r, storage)
	})
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		ping(w, r, storage)
	})
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		getList(w, r, storage)
	})

	return router
}
