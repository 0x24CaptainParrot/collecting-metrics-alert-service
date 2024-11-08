package storage

type MetricStorage interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int64) error
	GetMetric(name string, metricType MetricType) (interface{}, error)
	GetMetrics() map[string]interface{}
}
