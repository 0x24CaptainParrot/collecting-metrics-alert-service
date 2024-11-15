package handlers

import (
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}
