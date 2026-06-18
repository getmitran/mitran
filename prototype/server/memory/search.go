package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type SearchResult struct {
	Type    string  `json:"type"` // fact, episode, correction
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

type Correction struct {
	ID        string    `json:"id"`
	Original  string    `json:"original"`
	Corrected string    `json:"corrected"`
	Reason    string    `json:"reason"`
	AppliedTo string   `json:"applied_to"`
	CreatedAt time.Time `json:"created_at"`
}

type SearchableMemory struct {
	mu          sync.RWMutex
	dir         string
	Corrections []Correction `json:"corrections"`
}

func NewSearchableMemory(dir string) (*SearchableMemory, error) {
	os.MkdirAll(dir, 0755)
	sm := &SearchableMemory{dir: dir}
	sm.load()
	return sm, nil
}

func (sm *SearchableMemory) Search(query string, store *Store) []SearchResult {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	var results []SearchResult
	q := strings.ToLower(query)

	// Search facts
	for _, f := range store.Facts {
		if strings.Contains(strings.ToLower(f.Key), q) || strings.Contains(strings.ToLower(f.Value), q) {
			results = append(results, SearchResult{Type: "fact", Content: f.Key + ": " + f.Value, Score: 1.0})
		}
	}
	// Search episodes
	for _, e := range store.Episodes {
		if strings.Contains(strings.ToLower(e.Content), q) {
			results = append(results, SearchResult{Type: "episode", Content: e.Content, Score: 0.8})
		}
	}
	// Search corrections
	for _, c := range sm.Corrections {
		if strings.Contains(strings.ToLower(c.Original), q) || strings.Contains(strings.ToLower(c.Corrected), q) {
			results = append(results, SearchResult{Type: "correction", Content: c.Original + " -> " + c.Corrected, Score: 0.9})
		}
	}
	return results
}

func (sm *SearchableMemory) AddCorrection(original, corrected, reason, appliedTo string) *Correction {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	c := Correction{
		ID: time.Now().Format("20060102150405"), Original: original,
		Corrected: corrected, Reason: reason, AppliedTo: appliedTo, CreatedAt: time.Now(),
	}
	sm.Corrections = append(sm.Corrections, c)
	sm.save()
	return &sm.Corrections[len(sm.Corrections)-1]
}

func (sm *SearchableMemory) ListCorrections() []Correction {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return append([]Correction{}, sm.Corrections...)
}

func (sm *SearchableMemory) load() {
	data, _ := os.ReadFile(filepath.Join(sm.dir, "corrections.json"))
	json.Unmarshal(data, &sm.Corrections)
}

func (sm *SearchableMemory) save() {
	data, _ := json.MarshalIndent(sm.Corrections, "", "  ")
	os.WriteFile(filepath.Join(sm.dir, "corrections.json"), data, 0644)
}
