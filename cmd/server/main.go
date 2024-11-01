package main

import (
	"log"
	"net/http"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/handlers"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(storage storage.MetricStorage) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	h := handlers.NewHandler(storage)
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetricHandler)
	r.Get("/value/{type}/{name}", h.GetMetricValueHandler)
	r.Get("/", h.GetAllMetricsStatic)

	return r
}

func main() {
	storage := storage.NewMemStorage()
	router := NewRouter(storage)

	parseServerFlags()
	log.Printf("starting server on %s", serverCfg.runServerAddrFlag)
	if err := http.ListenAndServe(serverCfg.runServerAddrFlag, router); err != nil {
		log.Fatalf("Error occured starting server: %v", err)
	}
}
