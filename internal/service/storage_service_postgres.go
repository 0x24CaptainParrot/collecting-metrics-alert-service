package service

import (
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

func (sDBServ *StorageDBService) UpdateGauge(name string, value float64) error {
	return sDBServ.repo.UpdateGauge(name, value)
}

func (sDBServ *StorageDBService) UpdateCounter(name string, value int64) error {
	return sDBServ.repo.UpdateCounter(name, value)
}

func (sDBServ *StorageDBService) GetMetric(name string, metricType storage.MetricType) (interface{}, error) {
	return sDBServ.repo.GetMetric(name, metricType)
}

func (sDBServ *StorageDBService) GetMetrics() (map[string]interface{}, error) {
	return sDBServ.repo.GetMetrics()
}

func (sDBServ *StorageDBService) DB() *sql.DB {
	return sDBServ.repo.(*repository.StoragePostgres).DB()
}

func (sDBServ *StorageDBService) Ping() error {
	return sDBServ.repo.(*repository.StoragePostgres).Ping()
}

func (sDBServ *StorageDBService) SaveMetricsToFile(filePath string) error {
	return sDBServ.repo.SaveMetricsToFile(filePath)
}

func (sDBServ *StorageDBService) LoadMetricsFromFile(filePath string) error {
	return sDBServ.repo.LoadMetricsFromFile(filePath)
}
