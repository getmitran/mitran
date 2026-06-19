package hr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"github.com/getmitran/mitran/server/apierr"
)

type OffboardTask struct {
	ID, Title, Category, AssignedTo, Status, CompletedAt string
}

type OffboardingFlow struct {
	ID, EmployeeID, LastDay, Reason, Status, InitiatedBy, InitiatedAt string
	Tasks                                                             []OffboardTask
}

type OffboardingStore struct {
	offboardings map[string]*OffboardingFlow
	mu           sync.RWMutex
	counter      int
}

func NewOffboardingStore() *OffboardingStore {
	return &OffboardingStore{offboardings: map[string]*OffboardingFlow{}}
}

var offboardDefaults = []struct{ Title, Category string }{
	{"Revoke system access", "it"}, {"Return laptop and equipment", "it"},
	{"Exit interview", "hr"}, {"Final settlement processing", "finance"},
	{"Knowledge transfer", "team"}, {"Remove from all systems", "it"},
}

func (s *OffboardingStore) Initiate(empID, lastDay, reason, initiatedBy string) *OffboardingFlow {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	f := &OffboardingFlow{ID: fmt.Sprintf("off-%d", s.counter), EmployeeID: empID, LastDay: lastDay, Reason: reason, Status: "initiated", InitiatedBy: initiatedBy, InitiatedAt: time.Now().Format(time.RFC3339)}
	for i, t := range offboardDefaults {
		f.Tasks = append(f.Tasks, OffboardTask{ID: fmt.Sprintf("task-%d", i+1), Title: t.Title, Category: t.Category, Status: "pending"})
	}
	s.offboardings[f.ID] = f
	return f
}

func (s *OffboardingStore) CompleteOffboardTask(flowID, taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f := s.offboardings[flowID]
	if f == nil {
		return fmt.Errorf("flow not found")
	}
	for i := range f.Tasks {
		if f.Tasks[i].ID == taskID {
			f.Tasks[i].Status = "done"
			f.Tasks[i].CompletedAt = time.Now().Format(time.RFC3339)
			done := true
			for _, t := range f.Tasks {
				if t.Status != "done" {
					done = false
				}
			}
			if done {
				f.Status = "completed"
			} else {
				f.Status = "in-progress"
			}
			return nil
		}
	}
	return fmt.Errorf("task not found")
}

func RegisterOffboardingRoutes(mux *http.ServeMux) {
	store := NewOffboardingStore()
	mux.HandleFunc("POST /api/v1/hr/offboarding", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ EmployeeID, LastDay, Reason, InitiatedBy string }
		json.NewDecoder(r.Body).Decode(&req)
		json.NewEncoder(w).Encode(store.Initiate(req.EmployeeID, req.LastDay, req.Reason, req.InitiatedBy))
	})
	mux.HandleFunc("GET /api/v1/hr/offboarding/{id}", func(w http.ResponseWriter, r *http.Request) {
		store.mu.RLock()
		f := store.offboardings[r.PathValue("id")]
		store.mu.RUnlock()
		if f == nil {
			apierr.WriteError(w, apierr.NotFound("not found"))
			return
		}
		json.NewEncoder(w).Encode(f)
	})
	mux.HandleFunc("PUT /api/v1/hr/offboarding/{id}/tasks/{taskId}", func(w http.ResponseWriter, r *http.Request) {
		if err := store.CompleteOffboardTask(r.PathValue("id"), r.PathValue("taskId")); err != nil {
			apierr.WriteError(w, apierr.NotFound(err.Error()))
			return
		}
		w.WriteHeader(200)
	})
}
