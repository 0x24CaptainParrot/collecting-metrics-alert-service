package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	storage storage.MetricStorage
}

func NewHandler(storage storage.MetricStorage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	if storage.MetricType(metricType) != storage.Gauge && storage.MetricType(metricType) != storage.Counter {
		http.Error(w, "unknown type was given", http.StatusBadRequest)
		log.Printf("Error parsing url parameter")
	}
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if metricType == "" || metricName == "" || metricValue == "" {
		http.Error(w, "missing metric ID or value", http.StatusBadRequest)
		log.Printf("Metric ID, type or value is missing in the request")
		return
	}

	switch storage.MetricType(metricType) {
	case storage.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "invalid gauge value", http.StatusBadRequest)
			log.Printf("Error parsing gauge value %s: %v", metricValue, err)
			return
		}
		h.storage.UpdateGauge(metricName, value)
	case storage.Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter value", http.StatusBadRequest)
			log.Printf("Error parsing counter value %s: %v", metricValue, err)
			return
		}
		h.storage.UpdateCounter(metricName, value)
	default:
		http.Error(w, "invalid metric type", http.StatusNotFound)
		log.Printf("Invalid metric type: %s", metricType)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Metric %s updated successfully", metricName)
}

func (h *Handler) GetMetricValueHandler(w http.ResponseWriter, r *http.Request) {
	metricType := storage.MetricType(chi.URLParam(r, "type"))
	metricName := chi.URLParam(r, "name")

	if metricType == "" || metricName == "" {
		http.Error(w, "missing metric type or name", http.StatusBadRequest)
		log.Printf("Metric type or name is missing in the request")
		return
	}

	metric, err := h.storage.GetMetric(metricName, metricType)
	if err != nil {
		http.Error(w, "metric not found", http.StatusNotFound)
		log.Printf("Metric not found: %s %s", metricType, metricName)
		return
	}

	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v", metric)
	log.Printf("Metric retrieved: %s %s = %v", metricType, metricName, metric)
}

func (h *Handler) GetAllMetricsStatic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	metrics := h.storage.GetMetrics()
	fmt.Fprintln(w, "<html><body><h1>Metrics:</h1><ul>")
	for name, val := range metrics {
		fmt.Fprintf(w, "<li>%s: %v</li>", name, val)
	}
	fmt.Fprintln(w, "</ul></body></html>")
	log.Println("All metrics retrieved in HTML format")
}
