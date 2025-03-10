package service

import (
	"database/sql"

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
	Storage MetricStorage
	DB      *sql.DB
	// StorageDB
}

func NewService(st *storage.MemStorage, db *sql.DB) *Service {
	return &Service{
		Storage: NewStorageService(st),
		DB:      db,
		// StorageDB: NewStorageDBService(repo),
	}
}
