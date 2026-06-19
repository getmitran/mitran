package session

import (
	"testing"
	"time"
)

func TestNewStoreReturnsNonNil(t *testing.T) {
	s := NewStore()
	if s == nil {
		t.Fatal("NewStore returned nil")
	}
}

func TestMemoryStoreSetGet(t *testing.T) {
	s := &MemoryStore{entries: make(map[string]memEntry)}
	if err := s.Set("tok1", "user@example.com", time.Hour); err != nil {
		t.Fatalf("Set: %v", err)
	}
	email, err := s.Get("tok1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if email != "user@example.com" {
		t.Fatalf("got %q, want user@example.com", email)
	}
}

func TestMemoryStoreGetExpired(t *testing.T) {
	s := &MemoryStore{entries: make(map[string]memEntry)}
	s.Set("tok2", "expired@test.com", -time.Second)
	_, err := s.Get("tok2")
	if err == nil {
		t.Fatal("expected error for expired session")
	}
}

func TestMemoryStoreDelete(t *testing.T) {
	s := &MemoryStore{entries: make(map[string]memEntry)}
	s.Set("tok3", "del@test.com", time.Hour)
	s.Delete("tok3")
	_, err := s.Get("tok3")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}
