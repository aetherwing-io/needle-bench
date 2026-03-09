package main

import (
	"os"
	"strconv"
)

// Config holds server configuration loaded from environment variables.
type Config struct {
	Port           string
	MaxConcurrent  int
	ComputeTimeout int // seconds
}

// LoadConfig reads configuration from environment with defaults.
func LoadConfig() Config {
	cfg := Config{
		Port:           "8080",
		MaxConcurrent:  100,
		ComputeTimeout: 25,
	}

	if p := os.Getenv("PORT"); p != "" {
		cfg.Port = p
	}

	if mc := os.Getenv("MAX_CONCURRENT"); mc != "" {
		if v, err := strconv.Atoi(mc); err == nil && v > 0 {
			cfg.MaxConcurrent = v
		}
	}

	if ct := os.Getenv("COMPUTE_TIMEOUT"); ct != "" {
		if v, err := strconv.Atoi(ct); err == nil && v > 0 {
			cfg.ComputeTimeout = v
		}
	}

	return cfg
}
