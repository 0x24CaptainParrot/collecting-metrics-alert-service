package main

import (
	"fmt"
	"time"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/config"
	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/metrics"
)

func main() {
	agentCfg := config.ParseAgentFlags()

	agent := metrics.NewAgent(
		agentCfg.EndpointAddr,
		time.Duration(agentCfg.PollInterval)*time.Second,
		time.Duration(agentCfg.ReportInterval)*time.Second,
		agentCfg.RateLimit,
		agentCfg.Key)

	fmt.Println("Starting agent")
	agent.Start()
}
