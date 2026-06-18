package artifacts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

const MaxVersions = 50

type Artifact struct {
	Slug      string   `json:"slug"`
	Name      string   `json:"name"`
	Kind      string   `json:"kind"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	Version   int      `json:"version"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type ArtifactVersion struct {
	Version   int    `json:"version"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type Store struct {
	mu       sync.RWMutex
	dir      string
	artifacts map[string]*Artifact
	versions  map[string][]ArtifactVersion
}

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = slugRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 80 {
		s = s[:80]
	}
	return s
}

func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	s := &Store{dir: dir, artifacts: make(map[string]*Artifact), versions: make(map[string][]ArtifactVersion)}
	s.load()
	return s, nil
}

func (s *Store) path() string { return filepath.Join(s.dir, "artifacts.json") }

func (s *Store) load() {
	data, err := os.ReadFile(s.path())
	if err != nil {
		return
	}
	var state struct {
		Artifacts map[string]*Artifact          `json:"artifacts"`
		Versions  map[string][]ArtifactVersion  `json:"versions"`
	}
	if json.Unmarshal(data, &state) == nil {
		if state.Artifacts != nil {
			s.artifacts = state.Artifacts
		}
		if state.Versions != nil {
			s.versions = state.Versions
		}
	}
}

func (s *Store) persist() {
	state := struct {
		Artifacts map[string]*Artifact          `json:"artifacts"`
		Versions  map[string][]ArtifactVersion  `json:"versions"`
	}{s.artifacts, s.versions}
	data, _ := json.MarshalIndent(state, "", "  ")
	os.WriteFile(s.path(), data, 0644)
}

func (s *Store) Save(name, kind, content string, tags []string) (*Artifact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slug := slugify(name)
	if _, exists := s.artifacts[slug]; exists {
		return nil, fmt.Errorf("artifact %q already exists", slug)
	}
	now := time.Now().UTC().Format(time.RFC3339)
	a := &Artifact{Slug: slug, Name: name, Kind: kind, Content: content, Tags: tags, Version: 1, CreatedAt: now, UpdatedAt: now}
	s.artifacts[slug] = a
	s.versions[slug] = []ArtifactVersion{{Version: 1, Content: content, CreatedAt: now}}
	s.persist()
	return a, nil
}

func (s *Store) Update(slug, content string) (*Artifact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.artifacts[slug]
	if !ok {
		return nil, fmt.Errorf("artifact %q not found", slug)
	}
	a.Version++
	a.Content = content
	a.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	vers := s.versions[slug]
	vers = append(vers, ArtifactVersion{Version: a.Version, Content: content, CreatedAt: a.UpdatedAt})
	if len(vers) > MaxVersions {
		vers = vers[len(vers)-MaxVersions:]
	}
	s.versions[slug] = vers
	s.persist()
	return a, nil
}

func (s *Store) Get(slug string) (*Artifact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.artifacts[slug]
	if !ok {
		return nil, fmt.Errorf("artifact %q not found", slug)
	}
	return a, nil
}

func (s *Store) GetVersion(slug string, ver int) (*ArtifactVersion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	vers, ok := s.versions[slug]
	if !ok {
		return nil, fmt.Errorf("artifact %q not found", slug)
	}
	for _, v := range vers {
		if v.Version == ver {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("version %d not found for %q", ver, slug)
}

func (s *Store) List(kind, tag, query string) []*Artifact {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*Artifact
	for _, a := range s.artifacts {
		if kind != "" && a.Kind != kind {
			continue
		}
		if tag != "" {
			found := false
			for _, t := range a.Tags {
				if t == tag {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if query != "" && !strings.Contains(strings.ToLower(a.Name), strings.ToLower(query)) {
			continue
		}
		result = append(result, a)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].UpdatedAt > result[j].UpdatedAt })
	return result
}

func (s *Store) Delete(slug string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.artifacts[slug]; !ok {
		return fmt.Errorf("artifact %q not found", slug)
	}
	delete(s.artifacts, slug)
	delete(s.versions, slug)
	s.persist()
	return nil
}

func (s *Store) Versions(slug string) ([]int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	vers, ok := s.versions[slug]
	if !ok {
		return nil, fmt.Errorf("artifact %q not found", slug)
	}
	var result []int
	for _, v := range vers {
		result = append(result, v.Version)
	}
	return result, nil
}

// REST handlers

func (s *Store) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/artifacts", s.handleArtifacts)
	mux.HandleFunc("/api/v1/artifacts/", s.handleArtifactBySlug)
}

func (s *Store) handleArtifacts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		kind := r.URL.Query().Get("kind")
		tag := r.URL.Query().Get("tag")
		query := r.URL.Query().Get("q")
		json.NewEncoder(w).Encode(s.List(kind, tag, query))
	case http.MethodPost:
		var req struct {
			Name    string   `json:"name"`
			Kind    string   `json:"kind"`
			Content string   `json:"content"`
			Tags    []string `json:"tags"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		a, err := s.Save(req.Name, req.Kind, req.Content, req.Tags)
		if err != nil {
			http.Error(w, err.Error(), 409)
			return
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(a)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func (s *Store) handleArtifactBySlug(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/artifacts/")
	parts := strings.SplitN(path, "/", 2)
	slug := parts[0]

	if len(parts) == 2 && parts[1] == "versions" {
		vers, err := s.Versions(slug)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		json.NewEncoder(w).Encode(vers)
		return
	}

	switch r.Method {
	case http.MethodGet:
		a, err := s.Get(slug)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		json.NewEncoder(w).Encode(a)
	case http.MethodPut:
		var req struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		a, err := s.Update(slug, req.Content)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		json.NewEncoder(w).Encode(a)
	case http.MethodDelete:
		if err := s.Delete(slug); err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		w.WriteHeader(204)
	default:
		http.Error(w, "method not allowed", 405)
	}
}
