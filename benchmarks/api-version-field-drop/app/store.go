package main

import (
	"fmt"
	"sync"
	"time"
)

// Store is an in-memory user store.
type Store struct {
	mu    sync.RWMutex
	users map[string]*InternalUser
}

var store = &Store{
	users: make(map[string]*InternalUser),
}

func init() {
	now := time.Now()
	seedUsers := []*InternalUser{
		{
			ID: "usr-001", Name: "Alice Johnson", Email: "alice@example.com",
			Role: "admin", AvatarURL: "https://cdn.example.com/avatars/alice.jpg",
			Bio: "Engineering lead with 10 years of experience", Location: "San Francisco, CA",
			Department: "Engineering", PhoneNumber: "+1-555-0101", IsActive: true,
			CreatedAt: now.Add(-365 * 24 * time.Hour), UpdatedAt: now.Add(-24 * time.Hour),
		},
		{
			ID: "usr-002", Name: "Bob Smith", Email: "bob@example.com",
			Role: "user", AvatarURL: "https://cdn.example.com/avatars/bob.png",
			Bio: "Full-stack developer specializing in Go and React", Location: "New York, NY",
			Department: "Engineering", PhoneNumber: "+1-555-0102", IsActive: true,
			CreatedAt: now.Add(-200 * 24 * time.Hour), UpdatedAt: now.Add(-48 * time.Hour),
		},
		{
			ID: "usr-003", Name: "Carol Williams", Email: "carol@example.com",
			Role: "user", AvatarURL: "https://cdn.example.com/avatars/carol.jpg",
			Bio: "Product manager focused on developer tools", Location: "Austin, TX",
			Department: "Product", PhoneNumber: "+1-555-0103", IsActive: true,
			CreatedAt: now.Add(-100 * 24 * time.Hour), UpdatedAt: now.Add(-1 * time.Hour),
		},
		{
			ID: "usr-004", Name: "Dave Brown", Email: "dave@example.com",
			Role: "user", AvatarURL: "", Bio: "", Location: "Remote",
			Department: "Sales", PhoneNumber: "+1-555-0104", IsActive: false,
			CreatedAt: now.Add(-500 * 24 * time.Hour), UpdatedAt: now.Add(-30 * 24 * time.Hour),
		},
	}

	for _, u := range seedUsers {
		store.users[u.ID] = u
	}
}

// GetUser returns a user by ID.
func (s *Store) GetUser(id string) (*InternalUser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	return user, nil
}

// ListUsers returns all users.
func (s *Store) ListUsers() []*InternalUser {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*InternalUser, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}
