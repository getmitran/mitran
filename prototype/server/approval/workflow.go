package approval

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Status string

const (
	Pending  Status = "pending"
	Approved Status = "approved"
	Rejected Status = "rejected"
)

type Checkpoint struct {
	ID           string     `json:"id"`
	TaskID       string     `json:"task_id"`
	Agent        string     `json:"agent"`
	Description  string     `json:"description"`
	FilesChanged []string   `json:"files_changed"`
	Status       Status     `json:"status"`
	ReviewedBy   string     `json:"reviewed_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
}

type Store struct {
	mu   sync.RWMutex
	dir  string
	data []Checkpoint
}

func NewStore(dir string) (*Store, error) {
	os.MkdirAll(dir, 0755)
	s := &Store{dir: dir}
	s.load()
	return s, nil
}

func (s *Store) Create(taskID, agent, desc string, files []string) *Checkpoint {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := Checkpoint{
		ID: time.Now().Format("20060102150405"), TaskID: taskID, Agent: agent,
		Description: desc, FilesChanged: files, Status: Pending, CreatedAt: time.Now(),
	}
	s.data = append(s.data, cp)
	s.save()
	return &s.data[len(s.data)-1]
}

func (s *Store) Approve(id, reviewer string) { s.setStatus(id, Approved, reviewer) }
func (s *Store) Reject(id, reviewer string)  { s.setStatus(id, Rejected, reviewer) }

func (s *Store) setStatus(id string, status Status, reviewer string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.data {
		if s.data[i].ID == id {
			s.data[i].Status = status
			s.data[i].ReviewedBy = reviewer
			now := time.Now()
			s.data[i].ReviewedAt = &now
			break
		}
	}
	s.save()
}

func (s *Store) List() []Checkpoint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Checkpoint{}, s.data...)
}

func (s *Store) Pending() []Checkpoint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var r []Checkpoint
	for _, cp := range s.data {
		if cp.Status == Pending {
			r = append(r, cp)
		}
	}
	return r
}

func (s *Store) load() {
	data, _ := os.ReadFile(filepath.Join(s.dir, "checkpoints.json"))
	json.Unmarshal(data, &s.data)
}

func (s *Store) save() {
	data, _ := json.MarshalIndent(s.data, "", "  ")
	os.WriteFile(filepath.Join(s.dir, "checkpoints.json"), data, 0644)
}

func RegisterRoutes(mux *http.ServeMux, store *Store) {
	mux.HandleFunc("GET /api/v1/checkpoints", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(store.List())
	})
	mux.HandleFunc("GET /api/v1/checkpoints/pending", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(store.Pending())
	})
	mux.HandleFunc("POST /api/v1/checkpoints/{id}/approve", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		reviewer := r.Header.Get("X-User-ID")
		store.Approve(id, reviewer)
		w.WriteHeader(200)
	})
	mux.HandleFunc("POST /api/v1/checkpoints/{id}/reject", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		reviewer := r.Header.Get("X-User-ID")
		store.Reject(id, reviewer)
		w.WriteHeader(200)
	})
}
