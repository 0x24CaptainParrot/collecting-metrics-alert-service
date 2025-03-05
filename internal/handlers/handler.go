package handlers

import (
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/logger"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/middlewares"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}

func (h *Handler) InitHandlerRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(logger.LoggingHttpMiddleware(logger.Log))
	r.Use(middlewares.GzipMiddleware)

	r.Get("/ping", h.PingDatabase)
	r.Post("/update/", h.UpdateMetricJSONHandler)
	r.Post("/value/", h.GetMetricJSONHandler)
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetricHandler)
	r.Get("/value/{type}/{name}", h.GetMetricValueHandler)
	r.Get("/", h.GetAllMetricsStatic)

	return r
}
