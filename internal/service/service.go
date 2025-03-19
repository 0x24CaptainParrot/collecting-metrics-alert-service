package service

import (
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/repository"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

type MetricStorage interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int64) error
	GetMetric(name string, metricType storage.MetricType) (interface{}, error)
	GetMetrics() (map[string]interface{}, error)
	SaveMetricsToFile(filePath string) error
	LoadMetricsFromFile(filePath string) error
}

type StorageDB interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int64) error
	GetMetric(name string, metricType storage.MetricType) (interface{}, error)
	GetMetrics() (map[string]interface{}, error)
	SaveMetricsToFile(filePath string) error
	LoadMetricsFromFile(filePath string) error
}

type MetricStorageDB interface {
	MetricStorage
	StorageDB
}

type Service struct {
	MetricStorageDB MetricStorageDB
}

func NewService(repos *repository.Repository, st *storage.MemStorage) *Service {
	service := &Service{
		MetricStorageDB: NewStorageService(st),
	}
	if repos != nil && repos.StorageDB != nil {
		service.MetricStorageDB = NewStorageDBService(repos.StorageDB)
	}
	return service
}
