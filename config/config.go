package config

import (
	"os"
	"strconv"
)

// LLMConfig represents LLM configuration
type LLMConfig struct {
	BaseURL string
	Model   string
}

// Config represents the application configuration
type Config struct {
	LLM      LLMConfig
	DataPath string
}

// GetConfig returns the application configuration
func GetConfig() *Config {
	return &Config{
		LLM: LLMConfig{
			BaseURL: getEnvString("OLLAMA_BASE_URL", "http://localhost:11434"),
			Model:   getEnvString("OLLAMA_MODEL", "deepseek-r1:14b"),
		},
		DataPath: getEnvString("DATA_PATH", "./data"),
	}
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
