package metrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/models"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/utils"
	"github.com/go-resty/resty/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Agent struct {
	pollInterval   time.Duration
	client         *resty.Client
	reportInterval time.Duration
	serverAddress  string
	pollCount      int64
	key            string
	rateLimit      int
	metricQueue    chan map[string]interface{}
}

func NewAgent(serverAddress string, pollInterval, reportInterval time.Duration, rateLimit int, key string) *Agent {
	return &Agent{
		client:         resty.New(),
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		serverAddress:  serverAddress,
		rateLimit:      rateLimit,
		key:            key,
		metricQueue:    make(chan map[string]interface{}, 100),
	}
}

func (a *Agent) GetPollCount() int64 {
	return a.pollCount
}

func (a *Agent) SetPollCount(value int64) {
	a.pollCount = value
}

func (a *Agent) CollectGopsutilMetrics() map[string]interface{} {
	result := make(map[string]interface{})

	if vmStat, err := mem.VirtualMemory(); err == nil {
		result["TotalMemory"] = float64(vmStat.Total)
		result["FreeMemory"] = float64(vmStat.Free)
	}

	if cpuPercents, err := cpu.Percent(0, true); err == nil {
		for i, perc := range cpuPercents {
			result[fmt.Sprintf("CPUutilization%d", i+1)] = perc
		}
	}

	return result
}

func (a *Agent) CollectRuntimeMetrics() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := map[string]interface{}{
		"Alloc":         float64(memStats.Alloc),
		"BuckHashSys":   float64(memStats.BuckHashSys),
		"Frees":         float64(memStats.Frees),
		"GCCPUFraction": float64(memStats.GCCPUFraction),
		"GCSys":         float64(memStats.GCSys),
		"HeapAlloc":     float64(memStats.HeapAlloc),
		"HeapIdle":      float64(memStats.HeapIdle),
		"HeapInuse":     float64(memStats.HeapInuse),
		"HeapObjects":   float64(memStats.HeapObjects),
		"HeapReleased":  float64(memStats.HeapReleased),
		"HeapSys":       float64(memStats.HeapSys),
		"LastGC":        float64(memStats.LastGC),
		"Lookups":       float64(memStats.Lookups),
		"MCacheInuse":   float64(memStats.MCacheInuse),
		"MCacheSys":     float64(memStats.MCacheSys),
		"MSpanInuse":    float64(memStats.MSpanInuse),
		"MSpanSys":      float64(memStats.MSpanSys),
		"Mallocs":       float64(memStats.Mallocs),
		"NextGC":        float64(memStats.NextGC),
		"NumForcedGC":   float64(memStats.NumForcedGC),
		"NumGC":         float64(memStats.NumGC),
		"OtherSys":      float64(memStats.OtherSys),
		"PauseTotalNs":  float64(memStats.PauseTotalNs),
		"StackInuse":    float64(memStats.StackInuse),
		"StackSys":      float64(memStats.StackSys),
		"Sys":           float64(memStats.Sys),
		"TotalAlloc":    float64(memStats.TotalAlloc),
		"PollCount":     a.pollCount,
		"RandomValue":   rand.Float64(),
	}

	a.pollCount++
	return metrics
}

func (a *Agent) SendMetricsRetry(metrics map[string]interface{}) {
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

		headers := map[string]string{}
		if a.key != "" {
			hash, err := utils.ComputeSHA256("", a.key)
			if err != nil {
				log.Fatal("error computing the hash")
			}
			headers["HashSHA256"] = hash
		}

		url := fmt.Sprintf("%s/update/%s/%s/%s", a.serverAddress, metricType, metricName, valueStr)
		DoRequestWithRetry(func() error {
			resp, err := a.client.R().
				SetHeader("Content-Type", "text/plain").
				SetHeaders(headers).
				Post(url)

			if err != nil {
				return err
			}

			if !resp.IsSuccess() {
				if resp.StatusCode() >= 500 && resp.StatusCode() < 600 {
					return fmt.Errorf("server returned %d", resp.StatusCode())
				}
				return fmt.Errorf("server returned %d (non-retriable)", resp.StatusCode())
			}

			fmt.Printf("Metric: %s has been sent successfully\n", metricName)
			return nil
		})
	}
}

func (a *Agent) SendJSONMetricsRetry(metrics map[string]interface{}) {
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

		headers := map[string]string{}
		if a.key != "" {
			data, err := json.Marshal(metric)
			if err != nil {
				log.Fatal("error encoding metric")
			}
			hash, err := utils.ComputeSHA256(string(data), a.key)
			if err != nil {
				log.Fatal("error computing the hash")
			}
			headers["HashSHA256"] = hash
		}

		url := fmt.Sprintf("%s/update/", a.serverAddress)
		DoRequestWithRetry(func() error {
			resp, err := a.client.R().
				SetHeader("Content-Type", "application/json").
				SetHeaders(headers).
				SetBody(metric).
				Post(url)

			if err != nil {
				return err
			}

			if !resp.IsSuccess() {
				if resp.StatusCode() >= 500 && resp.StatusCode() < 600 {
					return fmt.Errorf("server returned %d", resp.StatusCode())
				}
				return fmt.Errorf("server returned %d (non-retriable)", resp.StatusCode())
			}

			fmt.Printf("Metric: %s has been sent successfully in the JSON format.\n", metricName)
			return nil
		})
	}
}

func (a *Agent) SendGzipJSONMetricsRetry(metrics map[string]interface{}) {
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

		headers := map[string]string{}
		data, err := json.Marshal(metric)
		if err != nil {
			fmt.Println("Failed to encode metric:", err)
		}
		if a.key != "" {
			hash, err := utils.ComputeSHA256(string(data), a.key)
			if err != nil {
				log.Fatal("error computing the hash")
			}
			headers["HashSHA256"] = hash
		}

		var buff bytes.Buffer
		gz := gzip.NewWriter(&buff)
		if _, err := gz.Write(data); err != nil {
			fmt.Println("Failed to compress metric:", err)
			continue
		}
		gz.Close()

		url := fmt.Sprintf("%s/update/", a.serverAddress)
		DoRequestWithRetry(func() error {
			resp, err := a.client.R().
				SetHeader("Content-Type", "application/json").
				SetHeader("Content-Encoding", "gzip").
				SetHeaders(headers).
				SetBody(buff.Bytes()).
				Post(url)

			if err != nil {
				return err
			}

			if !resp.IsSuccess() {
				if resp.StatusCode() >= 500 && resp.StatusCode() < 600 {
					return fmt.Errorf("server returned %d", resp.StatusCode())
				}
				return fmt.Errorf("server returned %d (non-retriable)", resp.StatusCode())
			}

			fmt.Printf("Gzip metric %s sent successfully.\n", metricName)
			return nil
		})
	}
}

func (a *Agent) SendBatchJSONMetricsRetry(metrics map[string]interface{}) {
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

	headers := map[string]string{}
	if a.key != "" {
		data, err := json.Marshal(batchMetrics)
		if err != nil {
			log.Fatal("error encoding batch of metrics")
		}
		hash, err := utils.ComputeSHA256(string(data), a.key)
		if err != nil {
			log.Fatal("error computing the hash")
		}
		headers["HashSHA256"] = hash
	}

	url := fmt.Sprintf("%s/updates/", a.serverAddress)
	DoRequestWithRetry(func() error {
		resp, err := a.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeaders(headers).
			SetBody(batchMetrics).
			Post(url)

		if err != nil {
			return err
		}

		if !resp.IsSuccess() {
			if resp.StatusCode() >= 500 && resp.StatusCode() < 600 {
				return fmt.Errorf("server returned %d", resp.StatusCode())
			}
			return fmt.Errorf("server returned %d (non-retriable)", resp.StatusCode())
		}

		fmt.Printf("Batch of JSON metrics has been successfully sent.\n")
		return nil
	})
}

func (a *Agent) worker() {
	for metrics := range a.metricQueue {
		a.SendMetricsRetry(metrics)
		a.SendJSONMetricsRetry(metrics)
		a.SendGzipJSONMetricsRetry(metrics)
		a.SendBatchJSONMetricsRetry(metrics)
	}
}

func (a *Agent) Start() {
	var mu sync.Mutex
	metricsMap := make(map[string]interface{})

	for i := 0; i < a.rateLimit; i++ {
		go a.worker()
	}

	go func() {
		tickerPoll := time.NewTicker(a.pollInterval)
		defer tickerPoll.Stop()

		for range tickerPoll.C {
			mu.Lock()
			metrics := a.CollectRuntimeMetrics()
			for k, v := range metrics {
				metricsMap[k] = v
			}
			mu.Unlock()
			fmt.Println("Metrics have been collected.")
		}
	}()

	go func() {
		tickerPoll := time.NewTicker(a.pollInterval)
		defer tickerPoll.Stop()

		for range tickerPoll.C {
			extraMetrics := a.CollectGopsutilMetrics()
			mu.Lock()
			for k, v := range extraMetrics {
				metricsMap[k] = v
			}
			mu.Unlock()
			fmt.Println("Extra metrics have been collected")
		}
	}()

	tickerReport := time.NewTicker(a.reportInterval)
	defer tickerReport.Stop()
	for range tickerReport.C {
		mu.Lock()
		snapshot := make(map[string]interface{}, len(metricsMap))
		for k, v := range metricsMap {
			snapshot[k] = v
		}
		mu.Unlock()
		a.metricQueue <- snapshot
	}
}

func DoRequestWithRetry(fn func() error) error {
	var backoffs = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	var lastErr error
	for i := 0; i < len(backoffs)+1; i++ {
		if err := fn(); err != nil {
			if IsRetriableNetworkErr(err) && i < len(backoffs) {
				lastErr = err
				time.Sleep(backoffs[i])
				continue
			}
			return err
		}
		return nil
	}

	return lastErr
}

func IsRetriableNetworkErr(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr)
}
