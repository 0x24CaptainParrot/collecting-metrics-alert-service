package main

import (
	"log"
	"net/http"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/handlers"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type MetricStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetMetrics() map[string]interface{}
}

func main() {
	serveMux := http.NewServeMux()
	memStorage := storage.NewMemStorage()

	serveMux.HandleFunc("/update/", handlers.UpdateMetricHandler(memStorage))
	serveMux.HandleFunc("/", handlers.NotFoundHandler)

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", serveMux); err != nil {
		log.Fatalf("Error occured starting server: %v", err)
	}
}
