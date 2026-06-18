package hr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type OffboardTask struct {
	ID, Title, Category, AssignedTo, Status, CompletedAt string
}

type OffboardingFlow struct {
	ID, EmployeeID, LastDay, Reason, Status, InitiatedBy, InitiatedAt string
	Tasks                                                             []OffboardTask
}

var (
	offboardings = map[string]*OffboardingFlow{}
	obMu         sync.RWMutex
	obCounter    int
)

var offboardDefaults = []struct{ Title, Category string }{
	{"Revoke system access", "it"}, {"Return laptop and equipment", "it"},
	{"Exit interview", "hr"}, {"Final settlement processing", "finance"},
	{"Knowledge transfer", "team"}, {"Remove from all systems", "it"},
}

func Initiate(empID, lastDay, reason, initiatedBy string) *OffboardingFlow {
	obMu.Lock()
	defer obMu.Unlock()
	obCounter++
	f := &OffboardingFlow{ID: fmt.Sprintf("off-%d", obCounter), EmployeeID: empID, LastDay: lastDay, Reason: reason, Status: "initiated", InitiatedBy: initiatedBy, InitiatedAt: time.Now().Format(time.RFC3339)}
	for i, t := range offboardDefaults {
		f.Tasks = append(f.Tasks, OffboardTask{ID: fmt.Sprintf("task-%d", i+1), Title: t.Title, Category: t.Category, Status: "pending"})
	}
	offboardings[f.ID] = f
	return f
}

func CompleteOffboardTask(flowID, taskID string) error {
	obMu.Lock()
	defer obMu.Unlock()
	f := offboardings[flowID]
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
	mux.HandleFunc("POST /api/v1/hr/offboarding", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ EmployeeID, LastDay, Reason, InitiatedBy string }
		json.NewDecoder(r.Body).Decode(&req)
		json.NewEncoder(w).Encode(Initiate(req.EmployeeID, req.LastDay, req.Reason, req.InitiatedBy))
	})
	mux.HandleFunc("GET /api/v1/hr/offboarding/{id}", func(w http.ResponseWriter, r *http.Request) {
		obMu.RLock()
		f := offboardings[r.PathValue("id")]
		obMu.RUnlock()
		if f == nil {
			http.Error(w, "not found", 404)
			return
		}
		json.NewEncoder(w).Encode(f)
	})
	mux.HandleFunc("PUT /api/v1/hr/offboarding/{id}/tasks/{taskId}", func(w http.ResponseWriter, r *http.Request) {
		if err := CompleteOffboardTask(r.PathValue("id"), r.PathValue("taskId")); err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		w.WriteHeader(200)
	})
}
