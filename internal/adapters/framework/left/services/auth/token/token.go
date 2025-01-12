package token

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
	"sync"
	"time"
)

type TokenStore struct {
	mu     sync.Mutex
	tokens map[string]string // email -> token
	ttl    map[string]time.Time
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]string),
		ttl:    make(map[string]time.Time),
	}
}

func (s *TokenStore) Set(email, token string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[email] = token
	s.ttl[email] = time.Now().Add(duration)
}

func (s *TokenStore) Get(email string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	token, exists := s.tokens[email]
	if !exists || time.Now().After(s.ttl[email]) {
		return "", false
	}
	return token, true
}

func (s *TokenStore) Delete(email string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, email)
	delete(s.ttl, email)
}

func GenerateToken() (string, error) {
	bytes := make([]byte, 3) // 3 bytes = 24 bits -> can encode to ~5 characters
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Use base32 encoding and trim padding
	token := base32.StdEncoding.EncodeToString(bytes)
	return strings.ToUpper(token[:5]), nil // Return the first 5 characters
}
