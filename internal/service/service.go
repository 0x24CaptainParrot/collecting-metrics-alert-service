package service

import "github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"

type Service struct {
	Storage storage.MetricStorage
}

func NewService(st *storage.MemStorage) *Service {
	return &Service{
		Storage: NewStorageService(st),
	}
}
