package session

import (
	"errors"
	"os"
	"sync"
	"time"
)

type Store interface {
	Get(token string) (string, error)
	Set(token, email string, ttl time.Duration) error
	Delete(token string) error
}

// MemoryStore is a dev-only in-memory session store.
type MemoryStore struct {
	stopOnce sync.Once
	mu      sync.RWMutex
	entries map[string]memEntry
	done    chan struct{}
}

type memEntry struct {
	email  string
	expiry time.Time
}

func (m *MemoryStore) Get(token string) (string, error) {
	m.mu.RLock()
	e, ok := m.entries[token]
	m.mu.RUnlock()
	if !ok || time.Now().After(e.expiry) {
		if ok {
			m.Delete(token)
		}
		return "", errors.New("session not found")
	}
	return e.email, nil
}

func (m *MemoryStore) Set(token, email string, ttl time.Duration) error {
	m.mu.Lock()
	m.entries[token] = memEntry{email: email, expiry: time.Now().Add(ttl)}
	m.mu.Unlock()
	return nil
}

func (m *MemoryStore) Delete(token string) error {
	m.mu.Lock()
	delete(m.entries, token)
	m.mu.Unlock()
	return nil
}

// StartReaper spawns a background goroutine that periodically removes expired sessions.
func (m *MemoryStore) StartReaper(interval time.Duration) {
	m.done = make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				m.mu.Lock()
				for k, e := range m.entries {
					if now.After(e.expiry) {
						delete(m.entries, k)
					}
				}
				m.mu.Unlock()
			case <-m.done:
				return
			}
		}
	}()
}

// Stop halts the background reaper.
func (m *MemoryStore) Stop() {
	m.stopOnce.Do(func() {
		if m.done != nil {
			close(m.done)
		}
	})
}

// RedisStore is a placeholder for production Redis-backed sessions.
type RedisStore struct {
	Addr string
}

var errRedisNotConfigured = errors.New("redis not configured")

func (r *RedisStore) Get(token string) (string, error)              { return "", errRedisNotConfigured }
func (r *RedisStore) Set(token, email string, ttl time.Duration) error { return errRedisNotConfigured }
func (r *RedisStore) Delete(token string) error                     { return errRedisNotConfigured }

// NewStore returns a Store based on MITRAN_SESSION_BACKEND env var.
func NewStore() Store {
	if os.Getenv("MITRAN_SESSION_BACKEND") == "redis" {
		return &RedisStore{Addr: os.Getenv("MITRAN_REDIS_ADDR")}
	}
	ms := &MemoryStore{entries: make(map[string]memEntry)}
	ms.StartReaper(5 * time.Minute)
	return ms
}
