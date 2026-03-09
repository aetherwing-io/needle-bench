package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	router := NewRouter()

	// Public routes
	router.Handle("/health", http.HandlerFunc(healthHandler))
	router.Handle("/login", http.HandlerFunc(loginHandler))

	// Protected API routes — should require auth
	router.Handle("/api/users", http.HandlerFunc(usersHandler))
	router.Handle("/api/admin", http.HandlerFunc(adminHandler))
	router.Handle("/api/data", http.HandlerFunc(dataHandler))
	router.Handle("/api/settings", http.HandlerFunc(settingsHandler))

	// Wrap with auth middleware
	handler := authMiddleware(router)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Router is a custom HTTP router that cleans paths before matching.
type Router struct {
	routes map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]http.Handler)}
}

func (rt *Router) Handle(pattern string, handler http.Handler) {
	rt.routes[pattern] = handler
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean the path before matching routes.
	// This normalizes /../api/admin to /api/admin.
	cleanPath := path.Clean(r.URL.Path)
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}

	handler, ok := rt.routes[cleanPath]
	if !ok {
		jsonResponse(w, 404, map[string]string{"error": "not found"})
		return
	}

	handler.ServeHTTP(w, r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, map[string]string{"status": "ok"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		jsonResponse(w, 400, map[string]string{"error": "invalid body"})
		return
	}

	if creds.Username == "admin" && creds.Password == "secret123" {
		jsonResponse(w, 200, map[string]string{
			"token": "valid-token-abc123",
			"user":  "admin",
		})
		return
	}

	jsonResponse(w, 401, map[string]string{"error": "invalid credentials"})
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, map[string]interface{}{
		"users": []map[string]string{
			{"id": "1", "name": "Alice", "email": "alice@example.com", "role": "admin"},
			{"id": "2", "name": "Bob", "email": "bob@example.com", "role": "user"},
			{"id": "3", "name": "Carol", "email": "carol@example.com", "role": "user"},
		},
	})
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, map[string]interface{}{
		"admin_panel": true,
		"secrets":     "database-password-xyz",
		"api_keys":    []string{"sk-prod-abc123", "sk-prod-def456"},
	})
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, map[string]interface{}{
		"records": []map[string]interface{}{
			{"id": 1, "data": "sensitive-record-1"},
			{"id": 2, "data": "sensitive-record-2"},
		},
	})
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, map[string]interface{}{
		"settings": map[string]string{
			"db_host": "10.0.1.5",
			"db_name": "production",
		},
	})
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func init() {
	fmt.Fprintln(os.Stderr, "API server initialized")
}
