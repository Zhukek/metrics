package repository

import models "github.com/Zhukek/metrics/internal/model"

type Repository interface {
	GetList() []string
	GetMetric(metricType string, metricName string) (res string, err error)
	GetMetricv2(body models.Metrics) (metricBody models.Metrics, err error)
	UpdateCounter(key string, value int64)
	UpdateGauge(key string, value float64)
	GetAllMetrics() map[string]models.Metrics
}
