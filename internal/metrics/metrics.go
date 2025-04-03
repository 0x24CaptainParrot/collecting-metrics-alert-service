package metrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/models"
)

func (a *Agent) SendMetrics(metrics map[string]interface{}) {
	for metricName, metricValue := range metrics {
		var metricType string
		var valueStr string

		switch v := metricValue.(type) {
		case float64:
			metricType = "gauge"
			valueStr = strconv.FormatFloat(v, 'f', -1, 64)
		case int64:
			metricType = "counter"
			valueStr = strconv.FormatInt(v, 10)
		default:
			continue
		}

		url := fmt.Sprintf("%s/update/%s/%s/%s", a.serverAddress, metricType, metricName, valueStr)
		resp, err := a.client.R().
			SetHeader("Content-Type", "text/plain").
			Post(url)

		if err != nil {
			fmt.Println("error while making the request", err)
			continue
		}

		if resp.IsSuccess() {
			fmt.Printf("Metric: %s has been sent successfully\n", metricName)
		} else {
			fmt.Printf("Failed to send metric %s: %s\n", metricName, resp.Status())
		}
	}
}

func (a *Agent) SendJSONMetrics(metrics map[string]interface{}) {
	for metricName, metricVal := range metrics {
		metric := models.Metrics{
			ID:    metricName,
			MType: "",
		}

		switch v := metricVal.(type) {
		case float64:
			metric.MType = "gauge"
			metric.Value = &v
		case int64:
			metric.MType = "counter"
			metric.Delta = &v
		default:
			continue
		}

		url := fmt.Sprintf("%s/update/", a.serverAddress)
		resp, err := a.client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(metric).
			Post(url)

		if err != nil {
			fmt.Println("error sending metric", err)
			continue
		}

		if resp.IsSuccess() {
			fmt.Printf("Metric: %s has been sent successfully in the JSON format.\n", metricName)
		} else {
			fmt.Printf("Failed to send metric %s: %s\n", metricName, resp.Status())
		}
	}
}

func (a *Agent) SendGzipJSONMetrics(metrics map[string]interface{}) {
	for metricName, metricVal := range metrics {
		metric := models.Metrics{
			ID:    metricName,
			MType: "",
		}

		switch v := metricVal.(type) {
		case float64:
			metric.MType = "gauge"
			metric.Value = &v
		case int64:
			metric.MType = "counter"
			metric.Delta = &v
		default:
			continue
		}

		var buff bytes.Buffer
		gz := gzip.NewWriter(&buff)
		if err := json.NewEncoder(gz).Encode(metric); err != nil {
			fmt.Println("Failed to encode metric:", err)
			continue
		}
		gz.Close()

		url := fmt.Sprintf("%s/update/", a.serverAddress)
		resp, err := a.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(buff.Bytes()).
			Post(url)

		if err != nil {
			fmt.Println("Error sending gzip metric:", err)
			continue
		}

		if resp.IsSuccess() {
			fmt.Printf("Gzip metric %s sent successfully.\n", metricName)
		} else {
			fmt.Printf("Failed to send gzip metric %s: %s\n", metricName, resp.Status())
		}
	}
}

func (a *Agent) SendBatchJSONMetrics(metrics map[string]interface{}) {
	var batchMetrics []models.Metrics

	for metricName, metricVal := range metrics {
		metric := models.Metrics{
			ID:    metricName,
			MType: "",
		}

		switch v := metricVal.(type) {
		case float64:
			metric.MType = "gauge"
			metric.Value = &v
		case int64:
			metric.MType = "counter"
			metric.Delta = &v
		default:
			continue
		}

		batchMetrics = append(batchMetrics, metric)
	}

	url := fmt.Sprintf("%s/updates/", a.serverAddress)
	resp, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(batchMetrics).
		Post(url)

	if err != nil {
		fmt.Println("error sending metric", err)
	}

	if resp.IsSuccess() {
		fmt.Printf("Batch of JSON metrics has been successfully sent.\n")
	} else {
		fmt.Printf("Failed to send batch of JSON metrics. Code: %s.", resp.Status())
	}
}
