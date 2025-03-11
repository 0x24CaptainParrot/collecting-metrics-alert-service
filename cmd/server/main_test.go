package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/handlers"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/service"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricHandler(t *testing.T) {
	type want struct {
		code int
	}
	type testCase struct {
		name        string
		url         string
		metricType  string
		expectedVal interface{}
		want        want
	}
	tests := []testCase{
		{
			name:        "valid gauge metric",
			url:         "/update/gauge/testGauge/100.5",
			metricType:  "gauge",
			expectedVal: 100.5,
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:        "valid counter metric",
			url:         "/update/counter/testCounter/10",
			metricType:  "counter",
			expectedVal: int64(10),
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "invalid metric type",
			url:  "/update/unknown/testMetric/100",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "invalid gauge value",
			url:  "/update/gauge/testGauge/invalid",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "missing metric ID",
			url:  "/update/gauge/",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}

	storage := storage.NewMemStorage()
	services := service.NewService(nil, storage)
	// router := handlers.NewRouter(services)
	handler := handlers.NewHandler(services)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.url, nil)
			w := httptest.NewRecorder()
			handler.InitHandlerRoutes().ServeHTTP(w, req)

			res := w.Result()
			assert.Equal(t, tc.want.code, res.StatusCode)

			// Проверка правильности сохранения метрики в хранилище
			if tc.want.code == http.StatusOK {
				metrics := storage.GetMetrics()

				switch tc.metricType {
				case "gauge":
					assert.Equal(t, tc.expectedVal, metrics["testGauge"], "Expected gauge value to be updated")
				case "counter":
					assert.Equal(t, tc.expectedVal, metrics["testCounter"], "Expected counter value to be updated")
				}
			}
		})
	}
}

// func TestMetricStorage(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	ms := mocks.NewMockMetricStorage(ctrl)
// 	service := service.NewService(ms, nil)

// 	type want struct {
// 		code int
// 	}

// 	type testCase struct {
// 		name       string
// 		url        string
// 		metricType string
// 		expected   interface{}
// 		want       want
// 	}

// 	testCases := []testCase{}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			req := httptest.NewRequest(http.MethodPost, tc.url, nil)
// 		})
// 	}
// }
