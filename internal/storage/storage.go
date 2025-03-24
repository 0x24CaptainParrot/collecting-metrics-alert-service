package storage

import (
	"context"
	"encoding/json"
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

func (ms *MemStorage) UpdateMetricValue(ctx context.Context, name string, value interface{}) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	switch v := value.(type) {
	case float64:
		ms.gauges[name] = v
		return nil
	case int64:
		if _, exists := ms.counters[name]; !exists {
			ms.counters[name] = 0
		}
		ms.counters[name] += v
		return nil
	default:
		return fmt.Errorf("invalid metric type was given: %v", v)
	}
}

func (ms *MemStorage) GetMetric(ctx context.Context, name string, metricType MetricType) (interface{}, error) {
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

func (ms *MemStorage) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
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

func (ms *MemStorage) SaveLoadMetrics(filePath string, operation string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	switch operation {
	case "save":
		metrics := make([]models.Metrics, 0)

		for name, value := range ms.gauges {
			metrics = append(metrics, models.Metrics{
				ID:    name,
				MType: "gauge",
				Value: &value,
			})
		}
		for name, value := range ms.counters {
			metrics = append(metrics, models.Metrics{
				ID:    name,
				MType: "counter",
				Delta: &value,
			})
		}

		data, err := json.MarshalIndent(metrics, "", "   ")
		if err != nil {
			return err
		}
		return os.WriteFile(filePath, data, 0644)

	case "load":
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
	default:
		return fmt.Errorf("invalid operation or incorrect filepath was given")
	}
}
