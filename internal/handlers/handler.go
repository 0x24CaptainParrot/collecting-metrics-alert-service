package handlers

import (
	"net/http"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/config"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/logger"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/middlewares"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	services *service.Service
	srvCfg   *config.ServerConfig
}

func NewHandler(service *service.Service, srvCfg *config.ServerConfig) *Handler {
	return &Handler{
		services: service,
		srvCfg:   srvCfg,
	}
}

func (h *Handler) InitHandlerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(logger.LoggingHttpMiddleware(logger.Log))
	r.Use(middlewares.GzipMiddleware)

	r.Get("/ping", h.PingDatabase)
	r.Post("/updates/", h.UpdateBatchMetricsJSONHandler)
	r.Post("/update/", h.UpdateMetricJSONHandler)
	r.Post("/value/", h.GetMetricJSONHandler)
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetricHandler)
	r.Get("/value/{type}/{name}", h.GetMetricValueHandler)
	r.Get("/", h.GetAllMetricsStatic)

	return r
}
