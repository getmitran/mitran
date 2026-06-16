package lock

import "testing"

func TestAcquireRelease(t *testing.T) {
	m := NewManager()
	if !m.Acquire("/src/main.go", "task1") {
		t.Fatal("expected acquire success")
	}
	if m.Acquire("/src/main.go", "task2") {
		t.Fatal("expected acquire fail for different task")
	}
	// same task re-acquire should succeed
	if !m.Acquire("/src/main.go", "task1") {
		t.Fatal("expected re-acquire by same task")
	}
	if !m.Release("/src/main.go", "task1") {
		t.Fatal("expected release success")
	}
	if m.Release("/src/main.go", "task1") {
		t.Fatal("expected release fail after already released")
	}
}

func TestIsLockedAndHeldBy(t *testing.T) {
	m := NewManager()
	if m.IsLocked("/a") {
		t.Fatal("expected not locked")
	}
	m.Acquire("/a", "t1")
	if !m.IsLocked("/a") {
		t.Fatal("expected locked")
	}
	if m.HeldBy("/a") != "t1" {
		t.Fatalf("expected t1, got %s", m.HeldBy("/a"))
	}
}

func TestReleaseWrongTask(t *testing.T) {
	m := NewManager()
	m.Acquire("/b", "t1")
	if m.Release("/b", "t2") {
		t.Fatal("expected release fail for wrong task")
	}
}
