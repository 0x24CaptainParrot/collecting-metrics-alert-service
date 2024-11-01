package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

type agentConfig struct {
	endpointAddr   string `env:"ADDRESS"`
	reportInterval int    `env:"REPORT_INTERVAL"`
	pollInterval   int    `env:"POLL_INTERVAL"`
}

var agentCfg agentConfig

func parseAgentFlags() {
	flag.StringVar(&agentCfg.endpointAddr, "a", "localhost:8080", "endpoint addr of the server")
	flag.IntVar(&agentCfg.reportInterval, "r", 10, "report interval")
	flag.IntVar(&agentCfg.pollInterval, "p", 2, "poll interval")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Printf("Unknown flags: %v\n", flag.Args())
		log.Fatal("Error: unnknown flags were given")
	}

	parsedUrl, err := url.Parse(agentCfg.endpointAddr)
	if err != nil || parsedUrl.Scheme == "" {
		agentCfg.endpointAddr = "http://" + agentCfg.endpointAddr
	}

	err = env.Parse(&agentCfg)
	if err != nil {
		log.Printf("error occured while parsing env variables: %v", err)
	}

	if envRunAgent := os.Getenv("ADDRESS"); envRunAgent != "" {
		agentCfg.endpointAddr = envRunAgent
	}
	if envReportInt := os.Getenv("REPORT_INTERVAL"); envReportInt != "" {
		agentCfg.reportInterval, err = strconv.Atoi(envReportInt)
		if err != nil {
			log.Fatalf("error occured: %v", err)
		}
	}
	if envPollInt := os.Getenv("POLL_INTERVAL"); envPollInt != "" {
		agentCfg.pollInterval, err = strconv.Atoi(envPollInt)
		if err != nil {
			log.Fatalf("error occured: %v", err)
		}
	}

	if !startsWithHTTP(agentCfg.endpointAddr) {
		agentCfg.endpointAddr = "http://" + agentCfg.endpointAddr
	}

	log.Printf("Agent will connect to %s", agentCfg.endpointAddr)
	log.Printf("Agent's reportInt: %d", agentCfg.reportInterval)
	log.Printf("Agent's pollInt: %d", agentCfg.pollInterval)
}

func startsWithHTTP(addr string) bool {
	return strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://")
}
