package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/handlers"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/service"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
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

	h := handlers.NewHandler(service)
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetricHandler)
	r.Get("/value/{type}/{name}", h.GetMetricValueHandler)
	r.Get("/", h.GetAllMetricsStatic)

	return r
}

func main() {
	storage := storage.NewMemStorage()
	services := service.NewService(storage)
	router := NewRouter(services)

	srv := &Server{}

	parseServerFlags()
	log.Printf("starting server on %s", serverCfg.runServerAddrFlag)
	go func() {
		if err := srv.Run(serverCfg.runServerAddrFlag, router); err != nil {
			log.Fatalf("Error occured starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("collecting metrics alert service shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error occured on server shutting down: %s", err.Error())
	}
}
