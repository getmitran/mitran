package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type SearchResult struct {
	Type    string  `json:"type"` // fact, episode, correction
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

type SearchableMemory struct {
	mu  sync.RWMutex
	dir string
}

func NewSearchableMemory(dir string) (*SearchableMemory, error) {
	os.MkdirAll(dir, 0755)
	sm := &SearchableMemory{dir: dir}
	return sm, nil
}

func (sm *SearchableMemory) Search(query string, store *MemoryStore) []SearchResult {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	var results []SearchResult
	q := strings.ToLower(query)

	// Search facts
	for k, v := range store.Facts {
		if strings.Contains(strings.ToLower(k), q) || strings.Contains(strings.ToLower(v), q) {
			results = append(results, SearchResult{Type: "fact", Content: k + ": " + v, Score: 1.0})
		}
	}
	// Search episodes
	for _, e := range store.Episodes {
		if strings.Contains(strings.ToLower(e.Text), q) {
			results = append(results, SearchResult{Type: "episode", Content: e.Text, Score: 0.8})
		}
	}
	// Search corrections
	for _, c := range store.Corrections {
		if strings.Contains(strings.ToLower(c.Rule), q) || strings.Contains(strings.ToLower(c.Negative), q) {
			results = append(results, SearchResult{Type: "correction", Content: c.Rule, Score: 0.9})
		}
	}
	return results
}

func (sm *SearchableMemory) AddCorrection(store *MemoryStore, rule, negative, category string) Correction {
	return store.AddCorrection(rule, negative, category)
}

func (sm *SearchableMemory) ListCorrections(store *MemoryStore) []Correction {
	return store.ListCorrections()
}

func (sm *SearchableMemory) load() {
	// No-op: corrections now live in MemoryStore
}

func (sm *SearchableMemory) save() {
	// No-op: corrections now live in MemoryStore
}

// LoadSearchIndex loads any supplemental search index data
func (sm *SearchableMemory) LoadSearchIndex() error {
	indexPath := filepath.Join(sm.dir, "search_index.json")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return nil
	}
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}
	_ = data // reserved for future index expansion
	return nil
}

// SaveSearchIndex persists supplemental search index data
func (sm *SearchableMemory) SaveSearchIndex() error {
	indexPath := filepath.Join(sm.dir, "search_index.json")
	data, _ := json.MarshalIndent(map[string]interface{}{}, "", "  ")
	return os.WriteFile(indexPath, data, 0644)
}
