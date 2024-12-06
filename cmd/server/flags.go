package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

const (
	serverDefaultAddress = "localhost:8080"
	logDefaultLevel      = "info"
)

type serverConfig struct {
	runServerAddrFlag string `env:"ADDRESS"`
	logLevel          string `env:"LOG_LEVEL"`
}

var serverCfg serverConfig

func parseServerFlags() {
	flag.StringVar(&serverCfg.runServerAddrFlag, "a", serverDefaultAddress, "server listens on this port")
	flag.StringVar(&serverCfg.logLevel, "l", logDefaultLevel, "log level")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Printf("Unknown flags %v\n", flag.Args())
		log.Fatal("Error: unknown flags were given")
	}

	err := env.Parse(&serverCfg)
	if err != nil {
		log.Printf("error accured while parsing server env variable: %v", err)
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

	log.Printf("Server will run on %s", serverCfg.runServerAddrFlag)
}
