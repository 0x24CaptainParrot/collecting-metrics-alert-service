package main

import "flag"

var agentConfig struct {
	endpointAddr   string
	reportInterval int
	pollInterval   int
}

func parseAgentFlags() {
	flag.StringVar(&agentConfig.endpointAddr, "a", "http://localhost:8080", "endpoint addr of the server")
	flag.IntVar(&agentConfig.reportInterval, "r", 10, "report interval")
	flag.IntVar(&agentConfig.pollInterval, "p", 2, "poll interval")

	flag.Parse()
}
