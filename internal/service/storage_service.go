package service

import (
	"context"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

type StorageService struct {
	st Storage
}

func NewStorageService(st Storage) *StorageService {
	return &StorageService{st: st}
}

func (s *StorageService) GetMetric(ctx context.Context, name string, metricType storage.MetricType) (interface{}, error) {
	return s.st.GetMetric(ctx, name, metricType)
}

func (s *StorageService) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
	return s.st.GetMetrics(ctx)
}

func (s *StorageService) UpdateGauge(ctx context.Context, name string, value float64) error {
	return s.st.UpdateGauge(ctx, name, value)
}

func (s *StorageService) UpdateCounter(ctx context.Context, name string, value int64) error {
	return s.st.UpdateCounter(ctx, name, value)
}

func (s *StorageService) SaveLoadMetrics(filePath string, operation string) error {
	return s.st.SaveLoadMetrics(filePath, operation)
}
