package service

import (
	"context"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/repository"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

type Storage interface {
	UpdateMetricValue(ctx context.Context, name string, value interface{}) error
	GetMetric(ctx context.Context, name string, metricType storage.MetricType) (interface{}, error)
	GetMetrics(ctx context.Context) (map[string]interface{}, error)
	SaveLoadMetrics(filePath string, operation string) error
}

type Service struct {
	Storage Storage
}

func NewService(repos repository.StorageDB, st Storage) *Service {
	service := &Service{
		Storage: NewStorageService(st),
	}

	if repos != nil && repos.(*repository.Repository).StorageDB != nil {
		service.Storage = NewStorageDBService(repos)
	}

	return service
}
