package repository

import (
	"context"

	models "github.com/Zhukek/metrics/internal/model"
)

type Repository interface {
	GetList() ([]string, error)
	GetMetric(metricType models.MType, metricName string) (res string, err error)
	GetMetricv2(body models.Metrics) (metricBody models.Metrics, err error)
	UpdateCounter(key string, value int64) error
	UpdateGauge(key string, value float64) error
	Updates([]models.MetricsBody) error
	Ping(ctx context.Context) error
	Close()
}
