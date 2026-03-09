package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/status", makeStatusHandler(cfg))
	mux.HandleFunc("/data", makeDataHandler(cfg))

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("server starting on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
		os.Exit(1)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func makeStatusHandler(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp := map[string]interface{}{
			"server":  cfg.Name,
			"port":    cfg.Port,
			"version": cfg.Version,
		}

		if cfg.Features.EnableMetrics {
			resp["metrics"] = cfg.Metrics.Endpoint
			resp["metrics_interval"] = cfg.Metrics.IntervalSec
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func makeDataHandler(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data := []map[string]interface{}{
			{"id": 1, "value": "alpha"},
			{"id": 2, "value": "bravo"},
			{"id": 3, "value": "charlie"},
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"items": data,
			"count": len(data),
		})
	}
}
