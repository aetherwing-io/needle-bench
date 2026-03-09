package main

import (
	"encoding/json"
	"os"
)

// Config holds the application configuration.
type Config struct {
	Name     string   `json:"name"`
	Port     int      `json:"port"`
	Version  string   `json:"version"`
	Features Features `json:"features"`
	Metrics  *MetricsConfig `json:"metrics,omitempty"`
}

// Features toggles optional functionality.
type Features struct {
	EnableMetrics bool `json:"enable_metrics"`
	EnableCache   bool `json:"enable_cache"`
	EnableRateLimit bool `json:"enable_rate_limit"`
}

// MetricsConfig holds metrics endpoint configuration.
type MetricsConfig struct {
	Endpoint    string `json:"endpoint"`
	IntervalSec int    `json:"interval_sec"`
}

// LoadConfig reads and parses a JSON config file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.Version == "" {
		cfg.Version = "0.0.1"
	}

	return &cfg, nil
}
