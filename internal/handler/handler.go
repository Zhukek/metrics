package handler

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/go-chi/chi/v5"
)

func update(res http.ResponseWriter, req *http.Request, storage models.MemStorage) {
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

func get(res http.ResponseWriter, req *http.Request, storage models.MemStorage) {
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

func getList(res http.ResponseWriter, req *http.Request, storage models.MemStorage) {
	metrics := storage.GetList()
	markup := `
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
	router.Post("/update/{metricType}/{metricName}/{metricValue}", func(w http.ResponseWriter, r *http.Request) {
		update(w, r, storage)
	})
	router.Get("/value/{metricType}/{metricName}", func(w http.ResponseWriter, r *http.Request) {
		get(w, r, storage)
	})
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		getList(w, r, storage)
	})

	return router
}
