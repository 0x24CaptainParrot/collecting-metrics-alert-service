package repository

import (
	"database/sql"
	"log"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

type StorageDB interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int64) error
	GetMetric(name string, metricType storage.MetricType) (interface{}, error)
	GetMetrics() map[string]interface{}
	SaveMetricsToFile(filePath string) error
	LoadMetricsFromFile(filePath string) error
}

type Repository struct {
	StorageDB
}

func NewRepository(db *sql.DB) *Repository {
	if db == nil {
		log.Println("[repository] Database is not set, using memory storage instead.")
		return &Repository{StorageDB: nil}
	}
	return &Repository{
		StorageDB: NewStoragePostgres(db),
	}
}
