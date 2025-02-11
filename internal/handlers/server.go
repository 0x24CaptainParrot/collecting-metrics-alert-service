package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/logger"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/middlewares"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(addr string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func NewRouter(service *service.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(logger.LoggingHttpMiddleware(logger.Log))
	r.Use(middlewares.GzipMiddleware)

	h := NewHandler(service)
	r.Post("/update/", h.UpdateMetricJSONHandler)
	// r.Post("/update", h.UpdateMetricJSONHandler)
	r.Post("/value/", h.GetMetricJSONHandler)
	// r.Post("/value", h.GetMetricJSONHandler)
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetricHandler)
	r.Get("/value/{type}/{name}", h.GetMetricValueHandler)
	r.Get("/", h.GetAllMetricsStatic)

	return r
}
