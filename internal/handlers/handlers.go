package handlers

import (
	"bytes"
	"crypto/hmac"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/config"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/models"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/service"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/utils"
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	if !verifyRequestSignature(r, body) {
		http.Error(w, "invalid hash", http.StatusBadRequest)
		log.Fatal("invalid hash. not equal")
		return
	}

	ctx := r.Context()
	switch storage.MetricType(metricType) {
	case storage.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "invalid gauge value", http.StatusBadRequest)
			log.Printf("Error parsing gauge value %s: %v", metricValue, err)
			return
		}
		if err := h.services.Storage.UpdateGauge(ctx, metricName, value); err != nil {
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
		if err := h.services.Storage.UpdateCounter(ctx, metricName, value); err != nil {
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
		if err := h.services.Storage.SaveLoadMetrics(FileStoragePath, "save"); err != nil {
			log.Printf("Failed to save metrics to file: %v", err)
		}
	}

	if config.ServerCfg.Key != "" {
		hash, err := utils.ComputeSHA256("", config.ServerCfg.Key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("HashSHA256", hash)
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

	ctx := r.Context()
	metric, err := h.services.Storage.GetMetric(ctx, metricName, metricType)
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

	ctx := r.Context()
	metrics, err := h.services.Storage.GetMetrics(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	if !verifyRequestSignature(r, body) {
		http.Error(w, "invalid hash", http.StatusBadRequest)
		log.Fatal("invalid hash. not equal")
		return
	}

	ctx := r.Context()
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
		if err := h.services.Storage.UpdateGauge(ctx, metric.ID, *metric.Value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "counter":
		if metric.Delta == nil {
			http.Error(w, "missing value for counter type", http.StatusBadRequest)
			return
		}
		if err := h.services.Storage.UpdateCounter(ctx, metric.ID, *metric.Delta); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid metric type", http.StatusBadRequest)
		return
	}

	if StoreInterval == 0 || StoreInterval > 0 {
		if err := h.services.Storage.SaveLoadMetrics(FileStoragePath, "save"); err != nil {
			log.Printf("Failed to save metrics to file: %v", err)
		}
	}

	if config.ServerCfg.Key != "" {
		data, err := json.Marshal(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hash, err := utils.ComputeSHA256(data, config.ServerCfg.Key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("HashSHA256", hash)
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

	ctx := r.Context()
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
		val, err := h.services.Storage.GetMetric(ctx, metric.ID, storage.Gauge)
		if err != nil {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		value := val.(float64)
		resultMetric.Value = &value
	case "counter":
		val, err := h.services.Storage.GetMetric(ctx, metric.ID, storage.Counter)
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

func (h *Handler) UpdateBatchMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		http.Error(w, "empty request body", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "accepts json content-type only", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	if !verifyRequestSignature(r, body) {
		http.Error(w, "invalid hash", http.StatusBadRequest)
		log.Fatal("invalid hash. not equal")
		return
	}

	ctx := r.Context()
	var metrics []models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, metric := range metrics {
		switch metric.MType {
		case "gauge":
			if metric.Value == nil {
				http.Error(w, "missing value for gauge type", http.StatusBadRequest)
				return
			}
			if err := h.services.Storage.UpdateGauge(ctx, metric.ID, *metric.Value); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		case "counter":
			if metric.Delta == nil {
				http.Error(w, "missing value for counter type", http.StatusBadRequest)
				return
			}
			if err := h.services.Storage.UpdateCounter(ctx, metric.ID, *metric.Delta); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "invalid metric type", http.StatusBadRequest)
			return
		}
	}

	if config.ServerCfg.Key != "" {
		data, err := json.Marshal(metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hash, err := utils.ComputeSHA256(data, config.ServerCfg.Key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("HashSHA256", hash)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

func (h *Handler) PingDatabase(w http.ResponseWriter, r *http.Request) {
	if h.services.Storage == nil {
		http.Error(w, "database is not configured", http.StatusServiceUnavailable)
		log.Println("Ping request received, but databse is not configured.")
		return
	}

	dbStorage, ok := h.services.Storage.(*service.StorageDBService)
	if !ok || dbStorage.DB() == nil {
		http.Error(w, "database connection is disabled", http.StatusServiceUnavailable)
		log.Println("Ping request received, but database is not connected.")
		return
	}

	if err := dbStorage.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Database ping error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func verifyRequestSignature(r *http.Request, body []byte) bool {
	expected := r.Header.Get("HashSHA256")
	if config.ServerCfg.Key == "" || expected == "" {
		return true
	}
	actual, _ := utils.ComputeSHA256(string(body), config.ServerCfg.Key)
	return hmac.Equal([]byte(actual), []byte(expected))
}
