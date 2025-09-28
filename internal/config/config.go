package config

import (
	"os"
	"strconv"
)

type Config struct {
	ModelType     string
	MaxIterations int
	// Ollama
	OllamaBaseEndpoint string
	OllamaModelID      string
	// Gemini
	GeminiAPIKey string
}

func New() *Config {
	return &Config{
		ModelType:          getEnvString("MODEL_TYPE", "ollama"),
		MaxIterations:      getEnvInt("MAX_ITERATIONS", 10),
		OllamaBaseEndpoint: getEnvString("OLLAMA_BASE_ENDPOINT", "http://localhost:11434"),
		OllamaModelID:      getEnvString("OLLAMA_MODEL_ID", "llama3.2"),
		GeminiAPIKey:       getEnvString("GEMINI_API_KEY", ""),
	}
}

func getEnvString(name, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		value = fallback
	}
	return value
}

func getEnvInt(name string, fallback int) int {
	rawValue := os.Getenv(name)
	if rawValue == "" {
		return fallback
	}
	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0
	}
	return value
}
