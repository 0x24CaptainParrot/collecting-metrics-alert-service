package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

const (
	serverDefaultAddress   = "localhost:8080"
	logDefaultLevel        = "info"
	storeIntervalDefault   = 300
	fileStoragePathDefault = "../../internal/metrics_data_logs/metrics.json"
	restoreDefault         = true
)

type serverConfig struct {
	runServerAddrFlag string `env:"ADDRESS"`
	logLevel          string `env:"LOG_LEVEL"`
	storeInterval     uint   `env:"STORE_INTERVAL"`
	fileStoragePath   string `env:"FILE_STORAGE_PATH"`
	restore           bool   `env:"RESTORE"`
	dbDsn             string `env:"DATABASE_DSN"`
}

var serverCfg serverConfig

func parseServerFlags() {
	flag.StringVar(&serverCfg.runServerAddrFlag, "a", serverDefaultAddress, "server listens on this port")
	flag.StringVar(&serverCfg.logLevel, "l", logDefaultLevel, "log level")
	flag.UintVar(&serverCfg.storeInterval, "i", storeIntervalDefault, "interval in seconds for saving metrics")
	flag.StringVar(&serverCfg.fileStoragePath, "f", fileStoragePathDefault, "path to the file for saving metrics")
	flag.BoolVar(&serverCfg.restore, "r", restoreDefault, "whether or not to download metrics at server startup")
	flag.StringVar(&serverCfg.dbDsn, "d", "", "data source name")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Printf("Unknown flags %v\n", flag.Args())
		log.Fatal("Error: unknown flags were given")
	}

	err := env.Parse(&serverCfg)
	if err != nil {
		log.Printf("error occured while parsing server env variable: %v", err)
	}

	envRunServer := os.Getenv("ADDRESS")
	if envRunServer != "" {
		serverCfg.runServerAddrFlag = envRunServer
		log.Printf("Server configuration was changed via env variable.")
		log.Printf("ADDRESS was changed via env variable. (%s)", envRunServer)
	}

	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel != "" {
		serverCfg.logLevel = envLogLevel
		log.Printf("LOG_LEVEL was changed via env variable. (%s)", envLogLevel)
	}

	envStoreInterval := os.Getenv("STORE_INTERVAL")
	if envStoreInterval != "" {
		// serverCfg.storeInterval = envStoreInterval
		log.Printf("STORE_INTERVAL was changed via env variable. (%s)", envStoreInterval)
	}

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if envFileStoragePath != "" {
		serverCfg.fileStoragePath = envFileStoragePath
		log.Printf("FILE_STORAGE_PATH was changed via env variable. (%s)", envFileStoragePath)
	}

	envRestore := os.Getenv("RESTORE")
	if envRestore != "" {
		// serverCfg.restore = envRestore
		log.Printf("RESTORE was changed via env variable. (%s)", envRestore)
	}

	envDatabaseDsn := os.Getenv("DATABASE_DSN")
	if envDatabaseDsn != "" {
		serverCfg.dbDsn = envDatabaseDsn
		log.Printf("DATABASE_DSN was changed via env variable. (%s)", envDatabaseDsn)
	}

	log.Printf("Server will run on %s", serverCfg.runServerAddrFlag)
}
