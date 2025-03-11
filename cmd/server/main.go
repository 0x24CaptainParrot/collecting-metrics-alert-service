package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/handlers"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/logger"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/repository"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/service"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

func main() {
	parseServerFlags()

	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()
	if db == nil {
		log.Println("Database is disabled. Running without db integration.")
	}

	storage := storage.NewMemStorage()
	repos := repository.NewRepository(db)
	services := service.NewService(repos, storage)
	handler := handlers.NewHandler(services)

	srv := &handlers.Server{}

	filePath := serverCfg.fileStoragePath
	restoreData := serverCfg.restore
	interval := int(serverCfg.storeInterval)

	handlers.StoreInterval = interval
	handlers.FileStoragePath = filePath

	if restoreData {
		if err := storage.LoadMetricsFromFile(filePath); err != nil {
			log.Printf("error loading metrics from from file: %v", err)
		} else {
			log.Println("Metrics successfully loaded from file.")
		}
	}

	stopSaving := make(chan struct{})
	go StartAutoSave(storage, filePath, interval, stopSaving)

	if err := logger.InitializeLogger(serverCfg.logLevel); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	log.Printf("starting server on %s", serverCfg.runServerAddrFlag)
	go func() {
		if err := srv.Run(serverCfg.runServerAddrFlag, handler.InitHandlerRoutes()); err != nil {
			log.Fatalf("Error occured starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("collecting metrics alert service shutting down")

	close(stopSaving)
	if err := storage.SaveMetricsToFile(filePath); err != nil {
		log.Printf("Failed to save metrics on shutdown: %v", err)
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error occured on server shutting down: %s", err.Error())
	}
}

func StartAutoSave(storage *storage.MemStorage, filePath string, interval int, stopChan chan struct{}) {
	if interval == 0 {
		log.Println("Auto-save disabled (STORE_INTERVAL=0)")
		return
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := storage.SaveMetricsToFile(filePath); err != nil {
				log.Printf("Failed to save metrics to file: %v", err)
			} else {
				log.Println("Metrics successfully saved to a file.")
			}
		case <-stopChan:
			return
		}
	}
}
