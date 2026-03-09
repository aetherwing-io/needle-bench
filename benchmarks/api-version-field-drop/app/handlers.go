package main

import (
	"net/http"
	"strings"
)

// v1 handlers

func v1UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users := store.ListUsers()
	var result []UserV1
	for _, u := range users {
		result = append(result, u.ToV1())
	}

	jsonResp(w, 200, map[string]interface{}{
		"users":   result,
		"count":   len(result),
		"version": "v1",
	})
}

func v1UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	userID := parts[3]

	user, err := store.GetUser(userID)
	if err != nil {
		jsonResp(w, 404, map[string]string{"error": err.Error()})
		return
	}

	v1User := user.ToV1()
	jsonResp(w, 200, v1User)
}

// v2 handlers

func v2UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users := store.ListUsers()
	var result []UserV2
	for _, u := range users {
		result = append(result, u.ToV2())
	}

	jsonResp(w, 200, map[string]interface{}{
		"users":   result,
		"count":   len(result),
		"version": "v2",
	})
}

func v2UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	userID := parts[3]

	user, err := store.GetUser(userID)
	if err != nil {
		jsonResp(w, 404, map[string]string{"error": err.Error()})
		return
	}

	v2User := user.ToV2()
	jsonResp(w, 200, v2User)
}
