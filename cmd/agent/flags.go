package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
)

var agentConfig struct {
	endpointAddr   string
	reportInterval int
	pollInterval   int
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

	if !startsWithHTTP(agentConfig.endpointAddr) {
		agentConfig.endpointAddr = "http://" + agentConfig.endpointAddr
	}

	log.Printf("Agent will connect to %s", agentConfig.endpointAddr)
}

func startsWithHTTP(addr string) bool {
	return strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://")
}
