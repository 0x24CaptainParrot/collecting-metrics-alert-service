package config

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/caarlos0/env/v6"
)

const (
	agentDefaultAddress   = "localhost:8080"
	agentDefaultReportInt = 10
	agentDefaultPollInt   = 2
)

type AgentConfig struct {
	EndpointAddr   string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

var AgentCfg AgentConfig

func ParseAgentFlags() {
	flag.StringVar(&AgentCfg.EndpointAddr, "a", agentDefaultAddress, "endpoint addr of the server")
	flag.IntVar(&AgentCfg.ReportInterval, "r", agentDefaultReportInt, "report interval")
	flag.IntVar(&AgentCfg.PollInterval, "p", agentDefaultPollInt, "poll interval")
	flag.StringVar(&AgentCfg.Key, "k", "", "secret key")
	flag.IntVar(&AgentCfg.RateLimit, "l", 5, "rate limit")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Printf("Unknown flags: %v\n", flag.Args())
		log.Fatal("Error: unnknown flags were given")
	}

	parsedUrl, err := url.Parse(AgentCfg.EndpointAddr)
	if err != nil || parsedUrl.Scheme == "" {
		AgentCfg.EndpointAddr = "http://" + AgentCfg.EndpointAddr
	}

	err = env.Parse(&AgentCfg)
	if err != nil {
		log.Printf("error occured while parsing agent env variables: %v", err)
	}

	if !startsWithHTTP(AgentCfg.EndpointAddr) {
		AgentCfg.EndpointAddr = "http://" + AgentCfg.EndpointAddr
	}

	log.Printf("Agent will connect to %s", AgentCfg.EndpointAddr)
	log.Printf("Agent's reportInt: %d", AgentCfg.ReportInterval)
	log.Printf("Agent's pollInt: %d", AgentCfg.PollInterval)
	log.Printf("Agent's rate limit: %d", AgentCfg.RateLimit)
}

func startsWithHTTP(addr string) bool {
	return strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://")
}
