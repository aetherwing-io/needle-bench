package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

var activeGoroutines int64

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/compute", computeHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/stats", statsHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// computeHandler spawns a background goroutine to do heavy computation
// and streams results back to the client.
func computeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ComputeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Iterations <= 0 {
		req.Iterations = 100
	}

	resultCh := make(chan ComputeResult, 1)

	// BUG: goroutine spawned without context cancellation awareness.
	// If the client disconnects, this goroutine keeps running forever.
	go func() {
		atomic.AddInt64(&activeGoroutines, 1)
		defer atomic.AddInt64(&activeGoroutines, -1)

		result := performComputation(req)
		resultCh <- result
	}()

	// Wait for result — but if client disconnects, we just return.
	// The goroutine above keeps running and leaks.
	select {
	case result := <-resultCh:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	case <-time.After(25 * time.Second):
		http.Error(w, "computation timed out", http.StatusGatewayTimeout)
		// goroutine still running — leaked
	}
}

func performComputation(req ComputeRequest) ComputeResult {
	start := time.Now()
	sum := 0.0

	for i := 0; i < req.Iterations; i++ {
		// Simulate heavy work
		for j := 0; j < 10000; j++ {
			sum += float64(i*j) / float64(j+1)
		}
		// This is where we SHOULD check for cancellation
		// but there's no context being passed
		time.Sleep(time.Duration(req.DelayMs) * time.Millisecond)
	}

	return ComputeResult{
		Value:    sum,
		Duration: time.Since(start).String(),
		Iters:    req.Iterations,
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"goroutines":        runtime.NumGoroutine(),
		"active_computes":   atomic.LoadInt64(&activeGoroutines),
		"heap_alloc_mb":     float64(getMemStats().HeapAlloc) / 1024 / 1024,
	})
}

func getMemStats() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}

type ComputeRequest struct {
	Iterations int `json:"iterations"`
	DelayMs    int `json:"delay_ms"`
}

type ComputeResult struct {
	Value    float64 `json:"value"`
	Duration string  `json:"duration"`
	Iters    int     `json:"iterations"`
}

func init() {
	// Log goroutine count periodically for debugging
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			count := runtime.NumGoroutine()
			active := atomic.LoadInt64(&activeGoroutines)
			if count > 10 {
				fmt.Fprintf(os.Stderr, "WARNING: %d goroutines running (%d active computes)\n", count, active)
			}
		}
	}()
}
