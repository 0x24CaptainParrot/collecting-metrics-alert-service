package main

import (
	"fmt"
	"time"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/metrics"
)

func main() {
	parseAgentFlags()

	agent := metrics.NewAgent(
		agentCfg.endpointAddr,
		time.Duration(agentCfg.pollInterval)*time.Second,
		time.Duration(agentCfg.reportInterval)*time.Second)

	fmt.Println("Starting agent")
	agent.Start()
}
