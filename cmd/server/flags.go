package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

type serverConfig struct {
	runServerAddrFlag string `env:"ADDRESS"`
}

var serverCfg serverConfig

func parseServerFlags() {
	flag.StringVar(&serverCfg.runServerAddrFlag, "a", "localhost:8080", "server listens on this port")
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
	}

	log.Printf("Server will run on %s", serverCfg.runServerAddrFlag)
}
