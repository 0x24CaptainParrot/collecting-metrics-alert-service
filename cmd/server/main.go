package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/config"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/handlers"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/logger"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/repository"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/service"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/storage"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/utils"
)

func main() {
	srvCfg := config.ParseServerFlags()

	db, err := repository.NewPostgresDB(os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()

	if db != nil {
		if err := repository.RunMigrations(db, utils.GetMigrationsPath()); err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}
	}

	if db == nil {
		log.Println("Database is disabled. Running without db integration.")
	}

	storage := storage.NewMemStorage()
	repos := repository.NewRepository(db)
	services := service.NewService(repos, storage)
	handler := handlers.NewHandler(services, srvCfg)

	srv := &handlers.Server{}

	if srvCfg.Restore {
		if err := storage.SaveLoadMetrics(srvCfg.FileStoragePath, "load"); err != nil {
			log.Printf("error loading metrics from from file: %v", err)
		} else {
			log.Println("Metrics successfully loaded from file.")
		}
	}

	stopSaving := make(chan struct{})
	go utils.StartAutoSave(storage, srvCfg.FileStoragePath, int(srvCfg.StoreInterval), stopSaving)

	if err := logger.InitializeLogger(srvCfg.LogLevel); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	log.Printf("starting server on %s", srvCfg.RunServerAddrFlag)
	go func() {
		if err := srv.Run(srvCfg.RunServerAddrFlag, handler.InitHandlerRoutes()); err != nil {
			log.Fatalf("Error occured starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("collecting metrics alert service shutting down")

	close(stopSaving)
	if err := storage.SaveLoadMetrics(srvCfg.FileStoragePath, "save"); err != nil {
		log.Printf("Failed to save metrics on shutdown: %v", err)
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error occured on server shutting down: %s", err.Error())
	}
}
