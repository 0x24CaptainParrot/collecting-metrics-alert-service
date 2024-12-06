package metrics

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type Agent struct {
	pollInterval   time.Duration
	client         *resty.Client
	reportInterval time.Duration
	serverAddress  string
	pollCount      int64
}

func NewAgent(serverAddress string, pollInterval, reportInterval time.Duration) *Agent {
	return &Agent{
		client:         resty.New(),
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		serverAddress:  serverAddress,
	}
}

func (a *Agent) GetPollCount() int64 {
	return a.pollCount
}

func (a *Agent) SetPollCount(value int64) {
	a.pollCount = value
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

func (a *Agent) Start() {
	tickerPoll := time.NewTicker(a.pollInterval)
	tickerReport := time.NewTicker(a.reportInterval)

	metrics := make(map[string]interface{})

	for {
		select {
		case <-tickerPoll.C:
			newMetrics := a.CollectRuntimeMetrics()
			for k, v := range newMetrics {
				metrics[k] = v
			}
			fmt.Println("Metrics have been collected.")
		case <-tickerReport.C:
			a.SendMetrics(metrics)
			fmt.Println("Metrics have been sent.")
		}
	}
}