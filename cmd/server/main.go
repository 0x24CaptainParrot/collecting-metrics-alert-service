package main

import (
	"flag"
	"fmt"
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

var runServerAddrFlag string

func parseServerFlags() {
	flag.StringVar(&runServerAddrFlag, "a", ":8080", "server listens on this port")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Printf("Unknown flags %v\n", flag.Args())
		log.Fatal("Error: unknown flags were given")
	}
}

func main() {
	storage := storage.NewMemStorage()
	router := NewRouter(storage)

	parseServerFlags()
	log.Printf("starting server on %s", runServerAddrFlag)
	if err := http.ListenAndServe(runServerAddrFlag, router); err != nil {
		log.Fatalf("Error occured starting server: %v", err)
	}
}
