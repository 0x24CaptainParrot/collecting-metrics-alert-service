package config

import (
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	serverDefaultAddress   = "localhost:8080"
	logDefaultLevel        = "info"
	storeIntervalDefault   = 300
	fileStoragePathDefault = "../../internal/metrics_data_logs/metrics.json"
	restoreDefault         = true
)

type ServerConfig struct {
	RunServerAddrFlag string `env:"ADDRESS"`
	LogLevel          string `env:"LOG_LEVEL"`
	StoreInterval     uint   `env:"STORE_INTERVAL"`
	FileStoragePath   string `env:"FILE_STORAGE_PATH"`
	Restore           bool   `env:"RESTORE"`
	DbDsn             string `env:"DATABASE_DSN"`
}

var ServerCfg ServerConfig

func ParseServerFlags() {
	flag.StringVar(&ServerCfg.RunServerAddrFlag, "a", serverDefaultAddress, "server listens on this port")
	flag.StringVar(&ServerCfg.LogLevel, "l", logDefaultLevel, "log level")
	flag.UintVar(&ServerCfg.StoreInterval, "i", storeIntervalDefault, "interval in seconds for saving metrics")
	flag.StringVar(&ServerCfg.FileStoragePath, "f", fileStoragePathDefault, "path to the file for saving metrics")
	flag.BoolVar(&ServerCfg.Restore, "r", restoreDefault, "whether or not to download metrics at server startup")
	flag.StringVar(&ServerCfg.DbDsn, "d", "", "data source name")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Printf("Unknown flags %v\n", flag.Args())
		log.Fatal("Error: unknown flags were given")
	}

	err := env.Parse(&ServerCfg)
	if err != nil {
		log.Printf("error occured while parsing server env variable: %v", err)
	}

	log.Printf("Server will run on %s", ServerCfg.RunServerAddrFlag)
}
