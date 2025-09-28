package main

import (
	"context"
	"fmt"

	"github.com/jacobtie/wikipedia-research-agent/internal/config"
	"github.com/jacobtie/wikipedia-research-agent/internal/orchestrator"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg := config.New()
	results, err := orchestrator.New(cfg, "Write a research paper about Vlad the Impaler and how his policies impacted his kingdom").Run(context.Background())
	if err != nil {
		return err
	}
	for result := range results {
		if result.Error != nil {
			return result.Error
		}
		fmt.Println(result.Msg)
	}
	return nil
}
