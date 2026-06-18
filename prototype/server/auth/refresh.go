package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type RefreshEntry struct {
	Email   string
	Expiry  time.Time
	Rotated bool
}

type RefreshStore struct {
	mu     sync.RWMutex
	tokens map[string]*RefreshEntry
}

func NewRefreshStore() *RefreshStore {
	return &RefreshStore{tokens: make(map[string]*RefreshEntry)}
}

func (s *RefreshStore) IssueRefreshToken(email string) string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	token := hex.EncodeToString(b)
	s.mu.Lock()
	s.tokens[token] = &RefreshEntry{
		Email:  email,
		Expiry: time.Now().Add(7 * 24 * time.Hour),
	}
	s.mu.Unlock()
	return token
}

func (s *RefreshStore) ValidateRefresh(token string) (string, error) {
	s.mu.RLock()
	entry, ok := s.tokens[token]
	s.mu.RUnlock()
	if !ok {
		return "", errors.New("invalid refresh token")
	}
	if entry.Rotated {
		return "", errors.New("refresh token already used")
	}
	if time.Now().After(entry.Expiry) {
		return "", errors.New("refresh token expired")
	}
	return entry.Email, nil
}

func (s *RefreshStore) RotateRefresh(oldToken string) (newAccess, newRefresh string, err error) {
	email, err := s.ValidateRefresh(oldToken)
	if err != nil {
		return "", "", err
	}
	s.mu.Lock()
	s.tokens[oldToken].Rotated = true
	s.mu.Unlock()

	newRefresh = s.IssueRefreshToken(email)

	// Generate access token (64-byte hex)
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	newAccess = hex.EncodeToString(b)

	return newAccess, newRefresh, nil
}
