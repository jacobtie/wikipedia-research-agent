package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/jacobtie/wikipedia-research-agent/internal/config"
	"github.com/jacobtie/wikipedia-research-agent/internal/orchestrator"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	inputFile := flag.String("file", "default.txt", "provide the name of your input file in the input directory")
	flag.Parse()
	cfg := config.New()
	if inputFile == nil || *inputFile == "" {
		return fmt.Errorf("failed to parse file flag")
	}
	input, err := readInput(*inputFile)
	if err != nil {
		return err
	}
	results, err := orchestrator.New(cfg, input).Run(context.Background())
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

func readInput(inputFile string) (string, error) {
	rawBytes, err := os.ReadFile(path.Join("input", inputFile))
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(rawBytes), nil
}
