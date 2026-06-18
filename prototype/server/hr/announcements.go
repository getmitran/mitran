package hr

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Announcement struct {
	ID, Title, Body, Author, Priority, TargetDept, CreatedAt, ExpiresAt string
	Pinned                                                              bool
}

type AnnouncementStore struct {
	mu    sync.Mutex
	items []Announcement
}

func (s *AnnouncementStore) Create(a Announcement) Announcement {
	s.mu.Lock()
	defer s.mu.Unlock()
	a.ID, a.CreatedAt = uuid.NewString(), time.Now().UTC().Format(time.RFC3339)
	s.items = append([]Announcement{a}, s.items...)
	return a
}

func (s *AnnouncementStore) List(dept string, active bool) (out []Announcement) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC().Format(time.RFC3339)
	for _, a := range s.items {
		if (dept == "" || a.TargetDept == dept) && (!active || a.ExpiresAt == "" || a.ExpiresAt >= now) {
			out = append(out, a)
		}
	}
	return
}

func (s *AnnouncementStore) Pin(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Pinned = true; return
		}
	}
}

func (s *AnnouncementStore) Archive(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].ExpiresAt = time.Now().UTC().Format(time.RFC3339); return
		}
	}
}

func RegisterAnnouncementRoutes(r *mux.Router, s *AnnouncementStore) {
	r.HandleFunc("/api/v1/hr/announcements", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			var a Announcement
			json.NewDecoder(req.Body).Decode(&a)
			json.NewEncoder(w).Encode(s.Create(a))
		} else {
			json.NewEncoder(w).Encode(s.List(req.URL.Query().Get("dept"), req.URL.Query().Get("active") != "false"))
		}
	})
	r.HandleFunc("/api/v1/hr/announcements/{id}/pin", func(w http.ResponseWriter, req *http.Request) {
		s.Pin(mux.Vars(req)["id"]); w.WriteHeader(204)
	}).Methods("PUT")
}
