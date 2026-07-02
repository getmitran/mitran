package memory

import (
	"strings"
)

// SearchResult represents a unified search hit across memory types.
type SearchResult struct {
	Type    string  `json:"type"` // fact, lesson, episode
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

// SearchAll performs a LIKE-based search across facts, lessons, and episodes.
func (s *MemoryStore) SearchAll(projectID, query string) []SearchResult {
	var results []SearchResult
	q := strings.ToLower(query)

	// Search facts
	facts, _ := s.ListFacts(projectID)
	for _, f := range facts {
		if strings.Contains(strings.ToLower(f.Key), q) || strings.Contains(strings.ToLower(f.Value), q) {
			results = append(results, SearchResult{Type: "fact", Content: f.Key + ": " + f.Value, Score: 1.0})
		}
	}

	// Search lessons
	lessons, _ := s.ListLessons(projectID, "")
	for _, l := range lessons {
		if strings.Contains(strings.ToLower(l.Rule), q) || strings.Contains(strings.ToLower(l.Negative), q) {
			results = append(results, SearchResult{Type: "lesson", Content: l.Rule, Score: 0.9})
		}
	}

	// Search episodes
	episodes, _ := s.SearchEpisodes(projectID, query, 20)
	for _, e := range episodes {
		results = append(results, SearchResult{Type: "episode", Content: e.Content, Score: 0.8})
	}

	return results
}
