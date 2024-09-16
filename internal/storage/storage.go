package storage

import "sync"

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

func (ms *MemStorage) UpdateGauge(name string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.gauges[name] = value
}

func (ms *MemStorage) UpdateCounter(name string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.counters[name] += value
}

func (ms *MemStorage) GetMetrics() map[string]interface{} {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	metrics := make(map[string]interface{})
	for k, v := range ms.gauges {
		metrics[k] = v
	}

	for k, v := range ms.counters {
		metrics[k] = v
	}

	return metrics
}
