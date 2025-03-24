package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

type StorageDB interface {
	UpdateMetricValue(ctx context.Context, name string, value interface{}) error
	GetMetric(ctx context.Context, name string, metricType storage.MetricType) (interface{}, error)
	GetMetrics(ctx context.Context) (map[string]interface{}, error)
	SaveLoadMetrics(filePath string, operation string) error
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
