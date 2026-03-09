package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

// UserStore holds user credentials.
type UserStore struct {
	users map[string]string // username -> hashed password
	mu    sync.RWMutex
}

// NewUserStore creates a new credential store.
func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]string),
	}
}

// AddUser adds a user with the given password.
func (s *UserStore) AddUser(username, password string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[username] = hashPassword(password)
}

// Authenticate checks if the given password matches the stored hash.
func (s *UserStore) Authenticate(username, password string) bool {
	s.mu.RLock()
	storedHash, exists := s.users[username]
	s.mu.RUnlock()

	if !exists {
		return false
	}

	inputHash := hashPassword(password)
	return compareHashes(storedHash, inputHash)
}

// GetStoredHash returns the stored hash for a user (for testing).
func (s *UserStore) GetStoredHash(username string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	hash, ok := s.users[username]
	return hash, ok
}

// hashPassword creates a SHA-256 hash of the password.
func hashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}
