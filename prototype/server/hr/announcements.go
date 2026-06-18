package hr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Announcement struct {
	ID, Title, Body, Author, Priority, TargetDept, CreatedAt, ExpiresAt string
	Pinned                                                              bool
}

type AnnouncementStore struct {
	mu      sync.Mutex
	items   []Announcement
	counter int
}

func (s *AnnouncementStore) Create(a Announcement) Announcement {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	a.ID = fmt.Sprintf("ann-%d", s.counter)
	a.CreatedAt = time.Now().UTC().Format(time.RFC3339)
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
			s.items[i].Pinned = true
			return
		}
	}
}

func RegisterAnnouncementRoutes(mux *http.ServeMux, s *AnnouncementStore) {
	mux.HandleFunc("POST /api/v1/hr/announcements", func(w http.ResponseWriter, r *http.Request) {
		var a Announcement
		json.NewDecoder(r.Body).Decode(&a)
		json.NewEncoder(w).Encode(s.Create(a))
	})
	mux.HandleFunc("GET /api/v1/hr/announcements", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(s.List(r.URL.Query().Get("dept"), r.URL.Query().Get("active") != "false"))
	})
	mux.HandleFunc("PUT /api/v1/hr/announcements/{id}/pin", func(w http.ResponseWriter, r *http.Request) {
		s.Pin(r.PathValue("id"))
		w.WriteHeader(204)
	})
}
