package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type EpisodicEntry struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

type Correction struct {
	ID       string `json:"id"`
	Rule     string `json:"rule"`
	Negative string `json:"negative,omitempty"`
	Category string `json:"category"`
}

type MemoryStore struct {
	mu          sync.RWMutex
	dir         string
	Facts       map[string]string `json:"facts"`
	Episodes    []EpisodicEntry   `json:"episodes"`
	Corrections []Correction      `json:"corrections"`
}

func NewStore(dataDir string) (*MemoryStore, error) {
	dir := filepath.Join(dataDir, "memory")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	s := &MemoryStore{
		dir:         dir,
		Facts:       make(map[string]string),
		Episodes:    []EpisodicEntry{},
		Corrections: []Correction{},
	}
	s.load()
	return s, nil
}

func (s *MemoryStore) SetFact(key, value string) {
	s.mu.Lock()
	s.Facts[key] = value
	s.mu.Unlock()
	s.saveFacts()
}

func (s *MemoryStore) GetFact(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.Facts[key]
	return v, ok
}

func (s *MemoryStore) SearchFacts(query string) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q := strings.ToLower(query)
	result := make(map[string]string)
	for k, v := range s.Facts {
		if strings.Contains(strings.ToLower(k), q) || strings.Contains(strings.ToLower(v), q) {
			result[k] = v
		}
	}
	return result
}

func (s *MemoryStore) AddEpisode(text string) EpisodicEntry {
	e := EpisodicEntry{ID: genID(), Text: text, Timestamp: time.Now()}
	s.mu.Lock()
	s.Episodes = append(s.Episodes, e)
	s.mu.Unlock()
	s.saveEpisodes()
	return e
}

func (s *MemoryStore) SearchEpisodes(query string) []EpisodicEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q := strings.ToLower(query)
	var results []EpisodicEntry
	for _, e := range s.Episodes {
		if strings.Contains(strings.ToLower(e.Text), q) {
			results = append(results, e)
		}
	}
	return results
}

func (s *MemoryStore) AddCorrection(rule, negative, category string) Correction {
	c := Correction{ID: genID(), Rule: rule, Negative: negative, Category: category}
	s.mu.Lock()
	s.Corrections = append(s.Corrections, c)
	s.mu.Unlock()
	s.saveCorrections()
	return c
}

func (s *MemoryStore) ListCorrections() []Correction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Corrections
}

func (s *MemoryStore) RemoveCorrection(id string) bool {
	s.mu.Lock()
	for i, c := range s.Corrections {
		if c.ID == id {
			s.Corrections = append(s.Corrections[:i], s.Corrections[i+1:]...)
			s.mu.Unlock()
			s.saveCorrections()
			return true
		}
	}
	s.mu.Unlock()
	return false
}

func (s *MemoryStore) load() {
	s.loadFile("semantic.json", &s.Facts)
	s.loadFile("episodes.json", &s.Episodes)
	s.loadFile("corrections.json", &s.Corrections)
}

func (s *MemoryStore) loadFile(name string, target interface{}) {
	data, err := os.ReadFile(filepath.Join(s.dir, name))
	if err != nil {
		return
	}
	json.Unmarshal(data, target)
}

func (s *MemoryStore) saveFacts() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.writeFile("semantic.json", s.Facts)
}

func (s *MemoryStore) saveEpisodes() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.writeFile("episodes.json", s.Episodes)
}

func (s *MemoryStore) saveCorrections() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.writeFile("corrections.json", s.Corrections)
}

func (s *MemoryStore) writeFile(name string, data interface{}) {
	b, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(filepath.Join(s.dir, name), b, 0644)
}

func genID() string {
	return time.Now().Format("20060102150405") + fmt.Sprintf("%04d", time.Now().Nanosecond()%10000)
}
