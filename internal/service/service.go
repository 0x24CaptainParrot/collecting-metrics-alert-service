package service

import (
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/repository"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

type MetricStorage interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int64) error
	GetMetric(name string, metricType storage.MetricType) (interface{}, error)
	GetMetrics() map[string]interface{}
	SaveMetricsToFile(filePath string) error
	LoadMetricsFromFile(filePath string) error
}

type StorageDB interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int64) error
	GetMetric(name string, metricType storage.MetricType) (interface{}, error)
	GetMetrics() map[string]interface{}
	SaveMetricsToFile(filePath string) error
	LoadMetricsFromFile(filePath string) error
}

type Service struct {
	MetricStorage
	StorageDB
}

func NewService(repos *repository.Repository, st *storage.MemStorage) *Service {
	return &Service{
		MetricStorage: NewStorageService(st),
		StorageDB:     NewStorageDBService(repos.StorageDB),
	}
}
