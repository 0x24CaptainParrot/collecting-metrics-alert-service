package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/models"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
	"github.com/go-chi/chi/v5"
)

var StoreInterval int
var FileStoragePath string

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
		if err := h.services.Storage.UpdateGauge(metricName, value); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Failed to update gauge: %s: %v", metricName, err)
			return
		}
	case storage.Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter value", http.StatusBadRequest)
			log.Printf("Error parsing counter value %s: %v", metricValue, err)
			return
		}
		if err := h.services.Storage.UpdateCounter(metricName, value); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Failed to update counter: %s: %v", metricName, err)
			return
		}
	default:
		http.Error(w, "invalid metric type", http.StatusNotFound)
		log.Printf("Invalid metric type: %s", metricType)
		return
	}

	if StoreInterval == 0 || StoreInterval > 0 {
		if err := h.services.Storage.SaveMetricsToFile(FileStoragePath); err != nil {
			log.Printf("Failed to save metrics to file: %v", err)
		}
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

	metric, err := h.services.Storage.GetMetric(metricName, metricType)
	if err != nil {
		http.Error(w, "metric not found", http.StatusNotFound)
		log.Printf("Metric not found: %s %s", metricType, metricName)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v", metric)
	log.Printf("Metric retrieved: %s %s = %v", metricType, metricName, metric)
}

func (h *Handler) GetAllMetricsStatic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	metrics := h.services.Storage.GetMetrics()
	fmt.Fprintln(w, "<html><body><h1>Metrics:</h1><ul>")
	for name, val := range metrics {
		fmt.Fprintf(w, "<li>%s: %v</li>", name, val)
	}
	fmt.Fprintln(w, "</ul></body></html>")
	log.Println("All metrics retrieved in HTML format")
}

// json metric handlers/methods
func (h *Handler) UpdateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		http.Error(w, "empty request body", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid media type was given", http.StatusUnsupportedMediaType)
		return
	}

	var metric models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "invalid data was given", http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			http.Error(w, "missing value for gauge type", http.StatusBadRequest)
			return
		}
		if err := h.services.Storage.UpdateGauge(metric.ID, *metric.Value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "counter":
		if metric.Delta == nil {
			http.Error(w, "missing value for counter type", http.StatusBadRequest)
			return
		}
		if err := h.services.Storage.UpdateCounter(metric.ID, *metric.Delta); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid metric type", http.StatusBadRequest)
		return
	}

	if StoreInterval == 0 || StoreInterval > 0 {
		if err := h.services.Storage.SaveMetricsToFile(FileStoragePath); err != nil {
			log.Printf("Failed to save metrics to file: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metric)
}

func (h *Handler) GetMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		http.Error(w, "empty request body", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid media type was given", http.StatusUnsupportedMediaType)
		return
	}

	var metric models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resultMetric models.Metrics
	resultMetric.ID = metric.ID
	resultMetric.MType = metric.MType

	switch metric.MType {
	case "gauge":
		val, err := h.services.Storage.GetMetric(metric.ID, storage.Gauge)
		if err != nil {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		value := val.(float64)
		resultMetric.Value = &value
	case "counter":
		val, err := h.services.Storage.GetMetric(metric.ID, storage.Counter)
		if err != nil {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		delta := val.(int64)
		resultMetric.Delta = &delta
	default:
		http.Error(w, "invalid metric type was given", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultMetric)
}

func (h *Handler) PingDatabase(w http.ResponseWriter, r *http.Request) {
	if err := h.services.DB.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Database ping error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
