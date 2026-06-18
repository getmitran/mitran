package memory

import (
	"os"
	"testing"
)

func TestStore_AddAndGetFact(t *testing.T) {
	dir, _ := os.MkdirTemp("", "mem-test")
	defer os.RemoveAll(dir)
	s, err := New(dir)
	if err != nil { t.Fatal(err) }
	s.AddFact("color", "blue")
	if len(s.Facts) != 1 { t.Fatalf("expected 1 fact, got %d", len(s.Facts)) }
	if s.Facts[0].Value != "blue" { t.Errorf("expected blue, got %s", s.Facts[0].Value) }
}

func TestStore_AddEpisode(t *testing.T) {
	dir, _ := os.MkdirTemp("", "mem-test")
	defer os.RemoveAll(dir)
	s, _ := New(dir)
	s.AddEpisode("user said hello", "greeting")
	if len(s.Episodes) != 1 { t.Fatalf("expected 1 episode, got %d", len(s.Episodes)) }
}

func TestSearchableMemory_Search(t *testing.T) {
	dir, _ := os.MkdirTemp("", "search-test")
	defer os.RemoveAll(dir)
	store, _ := New(dir)
	store.AddFact("language", "golang")
	store.AddEpisode("built the server in golang", "dev")
	sm, _ := NewSearchableMemory(dir)
	results := sm.Search("golang", store)
	if len(results) < 2 { t.Errorf("expected at least 2 results, got %d", len(results)) }
}

func TestSearchableMemory_Corrections(t *testing.T) {
	dir, _ := os.MkdirTemp("", "corr-test")
	defer os.RemoveAll(dir)
	sm, _ := NewSearchableMemory(dir)
	sm.AddCorrection("old way", "new way", "better", "agent")
	if len(sm.Corrections) != 1 { t.Fatal("expected 1 correction") }
}
