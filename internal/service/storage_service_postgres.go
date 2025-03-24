package service

import (
	"context"
	"database/sql"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/repository"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

type StorageDBService struct {
	repo repository.StorageDB
}

func NewStorageDBService(repo repository.StorageDB) *StorageDBService {
	return &StorageDBService{repo: repo}
}

func (sDBServ *StorageDBService) GetMetric(ctx context.Context, name string, metricType storage.MetricType) (interface{}, error) {
	return sDBServ.repo.GetMetric(ctx, name, metricType)
}

func (sDBServ *StorageDBService) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
	return sDBServ.repo.GetMetrics(ctx)
}

func (sDBServ *StorageDBService) DB() *sql.DB {
	return sDBServ.repo.(*repository.Repository).DB()
}

func (sDBServ *StorageDBService) Ping() error {
	return sDBServ.repo.(*repository.Repository).Ping()
}

func (sDBServ *StorageDBService) UpdateMetricValue(ctx context.Context, name string, value interface{}) error {
	return sDBServ.repo.UpdateMetricValue(ctx, name, value)
}

func (sDBServ *StorageDBService) SaveLoadMetrics(filePath string, operation string) error {
	return sDBServ.repo.SaveLoadMetrics(filePath, operation)
}
