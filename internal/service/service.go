package service

import (
	"context"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/repository"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

type MetricsGetter interface {
	GetMetric(ctx context.Context, name string, metricType storage.MetricType) (interface{}, error)
	GetMetrics(ctx context.Context) (map[string]interface{}, error)
}

type MetricsSetter interface {
	UpdateGauge(ctx context.Context, name string, value float64) error
	UpdateCounter(ctx context.Context, name string, value int64) error
}

type MetricsSaverLoader interface {
	SaveLoadMetrics(filePath string, operation string) error
}

type Storage interface {
	MetricsGetter
	MetricsSetter
	MetricsSaverLoader
}

type Service struct {
	Storage Storage
}

func NewService(repos Storage, st Storage) *Service {
	service := &Service{
		Storage: NewStorageService(st),
	}

	if repos != nil && repos.(*repository.Repository).StorageDB != nil {
		service.Storage = NewStorageDBService(repos)
	}
	return service
}
