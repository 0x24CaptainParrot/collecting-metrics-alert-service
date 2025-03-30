package main

import (
	"fmt"
	"time"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/config"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/metrics"
)

func main() {
	config.ParseAgentFlags()

	agent := metrics.NewAgent(
		config.AgentCfg.EndpointAddr,
		time.Duration(config.AgentCfg.PollInterval)*time.Second,
		time.Duration(config.AgentCfg.ReportInterval)*time.Second)

	fmt.Println("Starting agent")
	agent.Start()
}
