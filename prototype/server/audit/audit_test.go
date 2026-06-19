package audit

import (
	"os"
	"testing"
	"time"
)

func tempLog(t *testing.T) *AuditLog {
	t.Helper()
	f, err := os.CreateTemp("", "audit-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return &AuditLog{path: f.Name()}
}

func TestLogCreatesFile(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/new.jsonl"
	a := &AuditLog{path: path}
	err := a.Log(AuditEntry{Action: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestAppendWritesEvent(t *testing.T) {
	a := tempLog(t)
	err := a.Log(AuditEntry{Action: "create", Resource: "/api/tasks"})
	if err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(a.path)
	if info.Size() == 0 {
		t.Fatal("file is empty after append")
	}
}

func TestReadAllReturnsAppended(t *testing.T) {
	a := tempLog(t)
	a.Log(AuditEntry{Action: "a1"})
	a.Log(AuditEntry{Action: "a2"})
	a.Log(AuditEntry{Action: "a3"})

	a.mu.Lock()
	entries := a.readAllLocked()
	a.mu.Unlock()

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[1].Action != "a2" {
		t.Fatalf("expected action 'a2', got '%s'", entries[1].Action)
	}
}

func TestRoundTripFields(t *testing.T) {
	a := tempLog(t)
	now := time.Now().UTC().Truncate(time.Second)
	a.Log(AuditEntry{
		Timestamp: now,
		UserID:    "user-42",
		Action:    "DELETE",
		Resource:  "/api/v1/tasks/7",
		Details:   "removed task",
		IP:        "10.0.0.1",
	})

	a.mu.Lock()
	entries := a.readAllLocked()
	a.mu.Unlock()

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.UserID != "user-42" {
		t.Errorf("UserID: got %q", e.UserID)
	}
	if e.Action != "DELETE" {
		t.Errorf("Action: got %q", e.Action)
	}
	if e.Resource != "/api/v1/tasks/7" {
		t.Errorf("Resource: got %q", e.Resource)
	}
	if e.Details != "removed task" {
		t.Errorf("Details: got %q", e.Details)
	}
	if e.IP != "10.0.0.1" {
		t.Errorf("IP: got %q", e.IP)
	}
	if !e.Timestamp.Equal(now) {
		t.Errorf("Timestamp: got %v, want %v", e.Timestamp, now)
	}
}
