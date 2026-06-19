package hr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"github.com/getmitran/mitran/server/apierr"
)

type ReviewCycle struct {
	ID, Name, StartDate, EndDate, Status string
}
type Review struct {
	ID, CycleID, EmployeeID, ManagerID string
	SelfRating, ManagerRating           int
	SelfComments, ManagerComments       string
	Goals, Status                       string
}
type ReviewStore struct {
	mu      sync.Mutex
	cycles  []ReviewCycle
	reviews []Review
}

func NewReviewStore() *ReviewStore { return &ReviewStore{} }

func (s *ReviewStore) CreateCycle(name, start, end string) ReviewCycle {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := ReviewCycle{ID: fmt.Sprintf("cycle-%d", time.Now().UnixNano()), Name: name, StartDate: start, EndDate: end, Status: "active"}
	s.cycles = append(s.cycles, c)
	return c
}
func (s *ReviewStore) ListCycles() []ReviewCycle { s.mu.Lock(); defer s.mu.Unlock(); return s.cycles }
func (s *ReviewStore) SubmitSelfReview(r Review) Review {
	s.mu.Lock()
	defer s.mu.Unlock()
	r.ID = fmt.Sprintf("rev-%d", time.Now().UnixNano())
	r.Status = "draft"
	s.reviews = append(s.reviews, r)
	return r
}
func (s *ReviewStore) SubmitManagerReview(id string, rating int, comments string) *Review {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.reviews {
		if s.reviews[i].ID == id {
			s.reviews[i].ManagerRating = rating
			s.reviews[i].ManagerComments = comments
			s.reviews[i].Status = "completed"
			return &s.reviews[i]
		}
	}
	return nil
}
func (s *ReviewStore) GetByEmployee(empID, cycleID string) []Review {
	s.mu.Lock()
	defer s.mu.Unlock()
	var res []Review
	for _, r := range s.reviews {
		if r.EmployeeID == empID && (cycleID == "" || r.CycleID == cycleID) {
			res = append(res, r)
		}
	}
	return res
}
func RegisterReviewRoutes(mux *http.ServeMux, store *ReviewStore) {
	mux.HandleFunc("POST /api/v1/hr/reviews/cycles", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Name, StartDate, EndDate string }
		json.NewDecoder(r.Body).Decode(&req)
		json.NewEncoder(w).Encode(store.CreateCycle(req.Name, req.StartDate, req.EndDate))
	})
	mux.HandleFunc("GET /api/v1/hr/reviews/cycles", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(store.ListCycles())
	})
	mux.HandleFunc("GET /api/v1/hr/reviews", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(store.GetByEmployee(r.URL.Query().Get("employee_id"), r.URL.Query().Get("cycle_id")))
	})
	mux.HandleFunc("POST /api/v1/hr/reviews", func(w http.ResponseWriter, r *http.Request) {
		var rev Review
		if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
			apierr.WriteError(w, apierr.BadRequest("invalid request body"))
			return
		}
		json.NewEncoder(w).Encode(store.SubmitSelfReview(rev))
	})
	mux.HandleFunc("PUT /api/v1/hr/reviews/{id}", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ ManagerRating int; ManagerComments string }
		json.NewDecoder(r.Body).Decode(&req)
		if rev := store.SubmitManagerReview(r.PathValue("id"), req.ManagerRating, req.ManagerComments); rev != nil {
			json.NewEncoder(w).Encode(rev)
		} else {
			apierr.WriteError(w, apierr.NotFound("not found"))
		}
	})
}
