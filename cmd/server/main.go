package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/config"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/handlers"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/logger"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/repository"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/service"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
)

func main() {
	config.ParseServerFlags()

	db, err := repository.NewPostgresDB(os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()

	if db != nil {
		if err := repository.RunMigrations(db, getMigrationsPath()); err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}
	}

	if db == nil {
		log.Println("Database is disabled. Running without db integration.")
	}

	storage := storage.NewMemStorage()
	repos := repository.NewRepository(db)
	services := service.NewService(repos, storage)
	handler := handlers.NewHandler(services)

	srv := &handlers.Server{}

	filePath := config.ServerCfg.FileStoragePath
	restoreData := config.ServerCfg.Restore
	interval := int(config.ServerCfg.StoreInterval)

	handlers.StoreInterval = interval
	handlers.FileStoragePath = filePath

	if restoreData {
		if err := storage.SaveLoadMetrics(filePath, "load"); err != nil {
			log.Printf("error loading metrics from from file: %v", err)
		} else {
			log.Println("Metrics successfully loaded from file.")
		}
	}

	stopSaving := make(chan struct{})
	go StartAutoSave(storage, filePath, interval, stopSaving)

	if err := logger.InitializeLogger(config.ServerCfg.LogLevel); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	log.Printf("starting server on %s", config.ServerCfg.RunServerAddrFlag)
	go func() {
		if err := srv.Run(config.ServerCfg.RunServerAddrFlag, handler.InitHandlerRoutes()); err != nil {
			log.Fatalf("Error occured starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("collecting metrics alert service shutting down")

	close(stopSaving)
	if err := storage.SaveLoadMetrics(filePath, "save"); err != nil {
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
			if err := storage.SaveLoadMetrics(filePath, "save"); err != nil {
				log.Printf("Failed to save metrics to file: %v", err)
			} else {
				log.Println("Metrics successfully saved to a file.")
			}
		case <-stopChan:
			return
		}
	}
}

func getMigrationsPath() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working dir: %v", err)
	}

	rootMarker := "collecting-metrics-alert-service"
	idx := strings.Index(wd, rootMarker)
	if idx == -1 {
		log.Fatalf("project root marker '%s' not found in: %s", rootMarker, wd)
	}

	rootPath := wd[:idx+len(rootMarker)]
	migrationsPath := filepath.Join(rootPath, "internal", "schema")

	return migrationsPath
}
