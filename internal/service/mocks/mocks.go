// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/service (interfaces: MetricStorage)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	storage "github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
	gomock "github.com/golang/mock/gomock"
)

// MockMetricStorage is a mock of MetricStorage interface.
type MockMetricStorage struct {
	ctrl     *gomock.Controller
	recorder *MockMetricStorageMockRecorder
}

// MockMetricStorageMockRecorder is the mock recorder for MockMetricStorage.
type MockMetricStorageMockRecorder struct {
	mock *MockMetricStorage
}

// NewMockMetricStorage creates a new mock instance.
func NewMockMetricStorage(ctrl *gomock.Controller) *MockMetricStorage {
	mock := &MockMetricStorage{ctrl: ctrl}
	mock.recorder = &MockMetricStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricStorage) EXPECT() *MockMetricStorageMockRecorder {
	return m.recorder
}

// GetMetric mocks base method.
func (m *MockMetricStorage) GetMetric(arg0 string, arg1 storage.MetricType) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetric", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetric indicates an expected call of GetMetric.
func (mr *MockMetricStorageMockRecorder) GetMetric(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetric", reflect.TypeOf((*MockMetricStorage)(nil).GetMetric), arg0, arg1)
}

// GetMetrics mocks base method.
func (m *MockMetricStorage) GetMetrics() map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetrics")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// GetMetrics indicates an expected call of GetMetrics.
func (mr *MockMetricStorageMockRecorder) GetMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetrics", reflect.TypeOf((*MockMetricStorage)(nil).GetMetrics))
}

// LoadMetricsFromFile mocks base method.
func (m *MockMetricStorage) LoadMetricsFromFile(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadMetricsFromFile", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadMetricsFromFile indicates an expected call of LoadMetricsFromFile.
func (mr *MockMetricStorageMockRecorder) LoadMetricsFromFile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadMetricsFromFile", reflect.TypeOf((*MockMetricStorage)(nil).LoadMetricsFromFile), arg0)
}

// SaveMetricsToFile mocks base method.
func (m *MockMetricStorage) SaveMetricsToFile(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMetricsToFile", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveMetricsToFile indicates an expected call of SaveMetricsToFile.
func (mr *MockMetricStorageMockRecorder) SaveMetricsToFile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMetricsToFile", reflect.TypeOf((*MockMetricStorage)(nil).SaveMetricsToFile), arg0)
}

// UpdateCounter mocks base method.
func (m *MockMetricStorage) UpdateCounter(arg0 string, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCounter", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCounter indicates an expected call of UpdateCounter.
func (mr *MockMetricStorageMockRecorder) UpdateCounter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCounter", reflect.TypeOf((*MockMetricStorage)(nil).UpdateCounter), arg0, arg1)
}

// UpdateGauge mocks base method.
func (m *MockMetricStorage) UpdateGauge(arg0 string, arg1 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGauge", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGauge indicates an expected call of UpdateGauge.
func (mr *MockMetricStorageMockRecorder) UpdateGauge(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGauge", reflect.TypeOf((*MockMetricStorage)(nil).UpdateGauge), arg0, arg1)
}
