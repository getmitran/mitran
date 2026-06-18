package cron

import (
	"encoding/json"
	"net/http"
	"strings"
)

type CreateJobRequest struct {
	Name     string   `json:"name"`
	Agent    string   `json:"agent"`
	Schedule Schedule `json:"schedule"`
	Payload  string   `json:"payload"`
}

func RegisterRoutes(mux *http.ServeMux, scheduler *Scheduler) {
	mux.HandleFunc("/api/v1/cron", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listJobs(w, scheduler)
		case http.MethodPost:
			createJob(w, r, scheduler)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/cron/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/cron/")
		parts := strings.Split(path, "/")
		if len(parts) == 0 || parts[0] == "" {
			http.Error(w, "missing job id", http.StatusBadRequest)
			return
		}
		id := parts[0]
		action := ""
		if len(parts) > 1 {
			action = parts[1]
		}

		switch {
		case r.Method == http.MethodDelete && action == "":
			scheduler.Remove(id)
			writeJSON(w, map[string]string{"status": "removed"})
		case r.Method == http.MethodPost && action == "pause":
			scheduler.Pause(id)
			writeJSON(w, map[string]string{"status": "paused"})
		case r.Method == http.MethodPost && action == "resume":
			scheduler.Resume(id)
			writeJSON(w, map[string]string{"status": "resumed"})
		case r.Method == http.MethodPost && action == "trigger":
			scheduler.Trigger(id)
			writeJSON(w, map[string]string{"status": "triggered"})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	})
}

func listJobs(w http.ResponseWriter, s *Scheduler) {
	writeJSON(w, s.List())
}

func createJob(w http.ResponseWriter, r *http.Request, s *Scheduler) {
	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Schedule == "" {
		http.Error(w, "name and schedule required", http.StatusBadRequest)
		return
	}
	job := s.Add(req.Name, req.Agent, req.Schedule, req.Payload)
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, job)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
