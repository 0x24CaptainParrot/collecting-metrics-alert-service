package main

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/config"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func TestAgentFunctions(t *testing.T) {
	type want struct {
		randomValueCheck bool
		pollCount        int64
		contentType      string
		statusCode       int
	}
	type testCase struct {
		name       string
		metrics    map[string]interface{}
		want       want
		serverURL  string
		pollBefore int64
	}

	tests := []testCase{
		{
			name: "collect and send metrics",
			metrics: map[string]interface{}{
				"Alloc":       12345.67,
				"PollCount":   int64(1),
				"RandomValue": rand.Float64(),
			},
			want: want{
				randomValueCheck: true,
				pollCount:        2, // Ожидается, что PollCount увеличится после одного poll
				contentType:      "text/plain",
				statusCode:       http.StatusOK,
			},
			serverURL:  "http://localhost:8080",
			pollBefore: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Моковый сервер для тестирования отправки метрик
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.want.contentType, r.Header.Get("Content-Type"), "Content-Type should be text/plain")
				w.WriteHeader(tc.want.statusCode)
			}))
			defer ts.Close()

			// Запускаем агента с адресом тестового сервера
			agent := metrics.NewAgent(ts.URL, 2*time.Second, 10*time.Second, config.AgentCfg.RateLimit)

			// Устанавливаем начальный pollCount и проверяем его значение перед сбором метрик
			agent.SetPollCount(tc.pollBefore)
			assert.Equal(t, tc.pollBefore, agent.GetPollCount(), "Initial PollCount should be set correctly")

			// Собираем метрики
			metrics := agent.CollectRuntimeMetrics()

			// Проверка PollCount и RandomValue
			if tc.want.randomValueCheck {
				assert.Greater(t, metrics["RandomValue"].(float64), 0.0, "RandomValue should be greater than 0")
			}
			assert.Equal(t, tc.want.pollCount, agent.GetPollCount(), "PollCount should increment after each poll")

			// Отправка метрик
			agent.SendMetrics(tc.metrics)
		})
	}
}
