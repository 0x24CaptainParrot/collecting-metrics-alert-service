package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

type UpdateMetrics interface {
	UpdateGauge(ctx context.Context, name string, value float64) error
	UpdateCounter(ctx context.Context, name string, value int64) error
}

type ReadMetrics interface {
	GetMetric(ctx context.Context, name string, metricType storage.MetricType) (interface{}, error)
	GetMetrics(ctx context.Context) (map[string]interface{}, error)
}

type BackUpMetrics interface {
	SaveLoadMetrics(filePath string, operation string) error
}

type StorageDB interface {
	UpdateMetrics
	ReadMetrics
	BackUpMetrics
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

func (r *Repository) DB() *sql.DB {
	return r.StorageDB.(*StoragePostgres).DB()
}

func (r *Repository) Ping() error {
	return r.StorageDB.(*StoragePostgres).Ping()
}
