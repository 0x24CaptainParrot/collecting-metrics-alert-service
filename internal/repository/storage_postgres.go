package repository

import (
	"database/sql"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

type StoragePostgres struct {
	db *sql.DB
}

func NewStoragePostgres(db *sql.DB) *StoragePostgres {
	return &StoragePostgres{db: db}
}

func (sp *StoragePostgres) UpdateGauge(name string, value float64) error {
	return nil
}

func (sp *StoragePostgres) UpdateCounter(name string, value int64) error {
	return nil
}

func (sp *StoragePostgres) GetMetric(name string, metricType storage.MetricType) (interface{}, error) {
	return nil, nil
}

func (sp *StoragePostgres) GetMetrics() map[string]interface{} {
	return nil
}

func (sp *StoragePostgres) SaveMetricsToFile(filePath string) error {
	return storage.NewMemStorage().SaveMetricsToFile(filePath)
}

func (sp *StoragePostgres) LoadMetricsFromFile(filePath string) error {
	return storage.NewMemStorage().LoadMetricsFromFile(filePath)
}
