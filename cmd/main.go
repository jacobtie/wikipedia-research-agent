package main

import (
	"context"
	"fmt"

	"github.com/jacobtie/wikipedia-research-agent/internal/orchestrator"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/config"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg := config.New()
	results := orchestrator.New(cfg, "Design a distributed system that processes incoming messages in real time by taking events off a stream, enriching them with an external API call, and then deposits in many data warehouses.").Run(context.Background())
	for result := range results {
		if result.Error != nil {
			return result.Error
		}
		fmt.Println(result.Msg)
	}
	return nil
}
