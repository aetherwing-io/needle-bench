package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"timing-attack/auth"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var store *auth.UserStore

func init() {
	store = auth.NewUserStore()
	// Add some test users
	store.AddUser("admin", "correct-horse-battery-staple")
	store.AddUser("user1", "password123")
	store.AddUser("service", "s3rv1c3-acc0unt-t0k3n")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if store.Authenticate(req.Username, req.Password) {
		json.NewEncoder(w).Encode(loginResponse{
			Success: true,
			Message: "Login successful",
		})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(loginResponse{
			Success: false,
			Message: "Invalid credentials",
		})
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

// Start launches the HTTP server.
func Start() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/health", handleHealth)

	fmt.Println("Auth server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
