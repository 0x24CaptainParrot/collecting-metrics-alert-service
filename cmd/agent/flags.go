package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var agentConfig struct {
	endpointAddr   string `env:"ADDRESS"`
	reportInterval int    `env:"REPORT_INTERVAL"`
	pollInterval   int    `env:"POLL_INTERVAL"`
}

func parseAgentFlags() {
	flag.StringVar(&agentConfig.endpointAddr, "a", "localhost:8080", "endpoint addr of the server")
	flag.IntVar(&agentConfig.reportInterval, "r", 10, "report interval")
	flag.IntVar(&agentConfig.pollInterval, "p", 2, "poll interval")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Printf("Unknown flags: %v\n", flag.Args())
		log.Fatal("Error: unnknown flags were given")
	}

	parsedUrl, err := url.Parse(agentConfig.endpointAddr)
	if err != nil || parsedUrl.Scheme == "" {
		agentConfig.endpointAddr = "http://" + agentConfig.endpointAddr
	}

	if envRunAgent := os.Getenv("ADDRESS"); envRunAgent != "" {
		agentConfig.endpointAddr = envRunAgent
	}
	if envReportInt := os.Getenv("REPORT_INTERVAL"); envReportInt != "" {
		agentConfig.reportInterval, err = strconv.Atoi(envReportInt)
		if err != nil {
			log.Fatalf("error occured: %v", err)
		}
	}
	if envPollInt := os.Getenv("POLL_INTERVAL"); envPollInt != "" {
		agentConfig.pollInterval, err = strconv.Atoi(envPollInt)
		if err != nil {
			log.Fatalf("error occured: %v", err)
		}
	}

	if !startsWithHTTP(agentConfig.endpointAddr) {
		agentConfig.endpointAddr = "http://" + agentConfig.endpointAddr
	}

	log.Printf("Agent will connect to %s", agentConfig.endpointAddr)
	log.Printf("Agent's reportInt: %d", agentConfig.reportInterval)
	log.Printf("Agent's pollInt: %d", agentConfig.pollInterval)
}

func startsWithHTTP(addr string) bool {
	return strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://")
}
