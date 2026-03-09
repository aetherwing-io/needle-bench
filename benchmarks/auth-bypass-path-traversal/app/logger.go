package main

import (
	"log"
	"net/http"
	"time"
)

// requestLogger logs each request with method, path, status, and duration.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &responseCapture{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(sw, r)

		log.Printf("[%s] %s %s -> %d (%s)",
			r.RemoteAddr, r.Method, r.URL.Path,
			sw.statusCode, time.Since(start))
	})
}

type responseCapture struct {
	http.ResponseWriter
	statusCode int
}

func (rc *responseCapture) WriteHeader(code int) {
	rc.statusCode = code
	rc.ResponseWriter.WriteHeader(code)
}
