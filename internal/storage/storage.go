package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/models"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
	mu       sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (ms *MemStorage) UpdateGauge(name string, value float64) error {
	if value < 0 {
		return errors.New("given gauge value cannot be negative")
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.gauges[name] = value
	return nil
}

func (ms *MemStorage) UpdateCounter(name string, value int64) error {
	if value < 0 {
		return errors.New("given counter value cannot be negative")
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.counters[name]; !exists {
		ms.counters[name] = 0
	}
	ms.counters[name] += value
	return nil
}

func (ms *MemStorage) GetMetric(name string, metricType MetricType) (interface{}, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	switch metricType {
	case Gauge:
		val, ok := ms.gauges[name]
		if !ok {
			return nil, fmt.Errorf("metric not found")
		}
		return val, nil
	case Counter:
		val, ok := ms.counters[name]
		if !ok {
			return nil, fmt.Errorf("metric not found")
		}
		return val, nil
	default:
		return nil, fmt.Errorf("invalid metric type")
	}
}

func (ms *MemStorage) GetMetrics() (map[string]interface{}, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	metrics := make(map[string]interface{})
	for k, v := range ms.gauges {
		metrics[k] = v
	}

	for k, v := range ms.counters {
		metrics[k] = v
	}

	return metrics, nil
}

// implementation of saving and loading metrics from storage
func (ms *MemStorage) SaveMetricsToFile(filePath string) error {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	metrics := make([]models.Metrics, 0)
	for name, value := range ms.gauges {
		v := value
		metrics = append(metrics, models.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &v,
		})
	}

	for name, value := range ms.counters {
		v := value
		metrics = append(metrics, models.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &v,
		})
	}

	data, err := json.MarshalIndent(metrics, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func (ms *MemStorage) LoadMetricsFromFile(filePath string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var metrics []models.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}

	for _, m := range metrics {
		if m.MType == "gauge" && m.Value != nil {
			ms.gauges[m.ID] = *m.Value
		} else if m.MType == "counter" && m.Delta != nil {
			ms.counters[m.ID] = *m.Delta
		}
	}
	return nil
}
