package service

import "github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"

type StorageService struct {
	st *storage.MemStorage
}

func NewStorageService(st *storage.MemStorage) *StorageService {
	return &StorageService{st: st}
}

func (s *StorageService) UpdateGauge(name string, value float64) error {
	return s.st.UpdateGauge(name, value)
}

func (s *StorageService) UpdateCounter(name string, value int64) error {
	return s.st.UpdateCounter(name, value)
}

func (s *StorageService) GetMetric(name string, metricType storage.MetricType) (interface{}, error) {
	return s.st.GetMetric(name, metricType)
}

func (s *StorageService) GetMetrics() (map[string]interface{}, error) {
	return s.st.GetMetrics()
}

func (s *StorageService) SaveMetricsToFile(filePath string) error {
	return s.st.SaveMetricsToFile(filePath)
}

func (s *StorageService) LoadMetricsFromFile(filePath string) error {
	return s.st.LoadMetricsFromFile(filePath)
}
