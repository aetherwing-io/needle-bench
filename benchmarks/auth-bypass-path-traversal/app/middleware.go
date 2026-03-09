package main

import (
	"net/http"
	"strings"
)

// authMiddleware checks for a valid token on protected routes.
// Routes under /api/ require authentication.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// BUG: Checks the raw r.URL.Path for "/api/" prefix.
		// The router cleans paths (normalizes //, .., etc.) before matching,
		// but this middleware sees the raw path. An attacker can use
		// "//api/admin" or "/x/../api/admin" — the raw path won't start
		// with "/api/" so the auth check is skipped, but the router
		// cleans it to "/api/admin" and serves the protected handler.
		if strings.HasPrefix(r.URL.Path, "/api/") {
			token := r.Header.Get("Authorization")
			if token == "" {
				jsonResponse(w, 401, map[string]string{"error": "unauthorized"})
				return
			}

			// Validate token
			if !isValidToken(token) {
				jsonResponse(w, 403, map[string]string{"error": "invalid token"})
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func isValidToken(token string) bool {
	// Strip "Bearer " prefix if present
	token = strings.TrimPrefix(token, "Bearer ")
	return token == "valid-token-abc123"
}
