package hr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type LeaveRequest struct {
	ID, EmployeeID, Type, StartDate, EndDate, Reason, Status, ApproverID, ApprovedAt string
}

type LeaveStore struct {
	mu       sync.Mutex
	requests map[string]*LeaveRequest
	counter  int
}

func NewLeaveStore() *LeaveStore { return &LeaveStore{requests: make(map[string]*LeaveRequest)} }

func (s *LeaveStore) Submit(r *LeaveRequest) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	r.ID = fmt.Sprintf("LV-%d", s.counter)
	r.Status = "pending"
	s.requests[r.ID] = r
	return r.ID
}

func (s *LeaveStore) Approve(id, approverID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.requests[id]
	if !ok { return fmt.Errorf("not found") }
	r.Status, r.ApproverID, r.ApprovedAt = "approved", approverID, time.Now().Format(time.RFC3339)
	return nil
}

func (s *LeaveStore) Deny(id, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.requests[id]
	if !ok { return fmt.Errorf("not found") }
	r.Status, r.Reason = "denied", reason
	return nil
}

func (s *LeaveStore) GetBalance(employeeID string) map[string]int {
	bal := map[string]int{"vacation": 20, "sick": 12, "personal": 5, "parental": 90}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.requests {
		if r.EmployeeID == employeeID && r.Status == "approved" { bal[r.Type]-- }
	}
	return bal
}

func (s *LeaveStore) ListByEmployee(id string) (out []*LeaveRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.requests { if r.EmployeeID == id { out = append(out, r) } }
	return
}

func (s *LeaveStore) ListPending() (out []*LeaveRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.requests { if r.Status == "pending" { out = append(out, r) } }
	return
}

func (s *LeaveStore) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/hr/leave", func(w http.ResponseWriter, r *http.Request) {
		var req LeaveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}
		id := s.Submit(&req)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	})
	mux.HandleFunc("GET /api/v1/hr/leave", func(w http.ResponseWriter, r *http.Request) {
		eid := r.URL.Query().Get("employee_id")
		if eid != "" { json.NewEncoder(w).Encode(s.ListByEmployee(eid)); return }
		json.NewEncoder(w).Encode(s.ListPending())
	})
	mux.HandleFunc("PUT /api/v1/hr/leave/{id}/approve", func(w http.ResponseWriter, r *http.Request) {
		var b struct{ ApproverID string }
		json.NewDecoder(r.Body).Decode(&b)
		if err := s.Approve(r.PathValue("id"), b.ApproverID); err != nil { http.Error(w, err.Error(), 404) }
	})
	mux.HandleFunc("PUT /api/v1/hr/leave/{id}/deny", func(w http.ResponseWriter, r *http.Request) {
		var b struct{ Reason string }
		json.NewDecoder(r.Body).Decode(&b)
		if err := s.Deny(r.PathValue("id"), b.Reason); err != nil { http.Error(w, err.Error(), 404) }
	})
}
