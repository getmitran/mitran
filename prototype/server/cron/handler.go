package cron

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/getmitran/mitran/server/apierr"
)

// CreateJobRequest is the payload for creating a new cron job
type CreateJobRequest struct {
	Name          string `json:"name"`
	AgentType     string `json:"agent_type"`
	Message       string `json:"message"`
	ScheduleType  string `json:"schedule_type"`  // "interval" or "cron"
	ScheduleValue string `json:"schedule_value"` // seconds string or cron expr
}

// UpdateJobRequest is the payload for updating a cron job
type UpdateJobRequest struct {
	Name          *string `json:"name,omitempty"`
	AgentType     *string `json:"agent_type,omitempty"`
	Message       *string `json:"message,omitempty"`
	ScheduleType  *string `json:"schedule_type,omitempty"`
	ScheduleValue *string `json:"schedule_value,omitempty"`
}

// RegisterRoutes registers cron REST endpoints on the mux
func RegisterRoutes(mux *http.ServeMux, scheduler *Scheduler) {
	// GET/POST /api/v1/projects/{id}/crons
	mux.HandleFunc("/api/v1/projects/", func(w http.ResponseWriter, r *http.Request) {
		// Parse: /api/v1/projects/{projectId}/crons[/{cronId}[/{action}]]
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
		parts := strings.Split(path, "/")

		if len(parts) < 2 || parts[1] != "crons" {
			return // not a cron route, let other handlers deal with it
		}

		projectID := parts[0]
		cronID := ""
		action := ""
		if len(parts) > 2 {
			cronID = parts[2]
		}
		if len(parts) > 3 {
			action = parts[3]
		}

		// Route dispatch
		switch {
		// LIST: GET /api/v1/projects/{id}/crons
		case r.Method == http.MethodGet && cronID == "":
			handleListJobs(w, scheduler, projectID)

		// CREATE: POST /api/v1/projects/{id}/crons
		case r.Method == http.MethodPost && cronID == "":
			handleCreateJob(w, r, scheduler, projectID)

		// DELETE: DELETE /api/v1/projects/{id}/crons/{cronId}
		case r.Method == http.MethodDelete && cronID != "" && action == "":
			handleDeleteJob(w, scheduler, cronID)

		// GET single: GET /api/v1/projects/{id}/crons/{cronId}
		case r.Method == http.MethodGet && cronID != "" && action == "":
			handleGetJob(w, scheduler, cronID)

		// UPDATE: PUT /api/v1/projects/{id}/crons/{cronId}
		case r.Method == http.MethodPut && cronID != "" && action == "":
			handleUpdateJob(w, r, scheduler, cronID)

		// TRIGGER: POST /api/v1/projects/{id}/crons/{cronId}/trigger
		case r.Method == http.MethodPost && action == "trigger":
			handleTriggerJob(w, scheduler, cronID)

		// PAUSE: POST /api/v1/projects/{id}/crons/{cronId}/pause
		case r.Method == http.MethodPost && action == "pause":
			handlePauseJob(w, scheduler, cronID)

		// RESUME: POST /api/v1/projects/{id}/crons/{cronId}/resume
		case r.Method == http.MethodPost && action == "resume":
			handleResumeJob(w, scheduler, cronID)

		default:
			apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		}
	})
}

func handleListJobs(w http.ResponseWriter, s *Scheduler, projectID string) {
	jobs, err := s.ListJobs(projectID)
	if err != nil {
		apierr.WriteError(w, apierr.Internal(err.Error()))
		return
	}
	if jobs == nil {
		jobs = []CronJob{}
	}
	respondJSON(w, http.StatusOK, jobs)
}

func handleGetJob(w http.ResponseWriter, s *Scheduler, cronID string) {
	job, err := s.GetJob(cronID)
	if err != nil {
		apierr.WriteError(w, apierr.NotFound("cron job not found"))
		return
	}
	respondJSON(w, http.StatusOK, job)
}

func handleCreateJob(w http.ResponseWriter, r *http.Request, s *Scheduler, projectID string) {
	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, apierr.BadRequest("invalid json"))
		return
	}
	if req.Name == "" {
		apierr.WriteError(w, apierr.BadRequest("name is required"))
		return
	}
	if req.ScheduleType == "" || req.ScheduleValue == "" {
		apierr.WriteError(w, apierr.BadRequest("schedule_type and schedule_value are required"))
		return
	}
	if req.ScheduleType != "interval" && req.ScheduleType != "cron" {
		apierr.WriteError(w, apierr.BadRequest("schedule_type must be 'interval' or 'cron'"))
		return
	}

	job := CronJob{
		ProjectID:     projectID,
		Name:          req.Name,
		AgentType:     req.AgentType,
		Message:       req.Message,
		ScheduleType:  req.ScheduleType,
		ScheduleValue: req.ScheduleValue,
	}
	if job.AgentType == "" {
		job.AgentType = "general"
	}

	created, err := s.AddJob(job)
	if err != nil {
		apierr.WriteError(w, apierr.BadRequest(err.Error()))
		return
	}
	respondJSON(w, http.StatusCreated, created)
}

func handleUpdateJob(w http.ResponseWriter, r *http.Request, s *Scheduler, cronID string) {
	var req UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, apierr.BadRequest("invalid json"))
		return
	}

	job, err := s.GetJob(cronID)
	if err != nil {
		apierr.WriteError(w, apierr.NotFound("cron job not found"))
		return
	}

	if req.Name != nil {
		job.Name = *req.Name
	}
	if req.AgentType != nil {
		job.AgentType = *req.AgentType
	}
	if req.Message != nil {
		job.Message = *req.Message
	}
	if req.ScheduleType != nil && req.ScheduleValue != nil {
		job.ScheduleType = *req.ScheduleType
		job.ScheduleValue = *req.ScheduleValue
	}

	// Recompute next run if schedule changed
	if req.ScheduleType != nil || req.ScheduleValue != nil {
		baseTime := job.CreatedAt
		if job.LastRun != nil {
			baseTime = *job.LastRun
		}
		nextRun, err := computeNextRun(job.ScheduleType, job.ScheduleValue, baseTime)
		if err != nil {
			apierr.WriteError(w, apierr.BadRequest(err.Error()))
			return
		}
		job.NextRun = nextRun
	}

	s.mu.Lock()
	_, err = s.store.DB().Exec(
		`UPDATE cron_jobs SET name=?, agent_type=?, message=?, schedule_type=?, schedule_value=?, next_run=? WHERE id=?`,
		job.Name, job.AgentType, job.Message, job.ScheduleType, job.ScheduleValue,
		job.NextRun.UTC().Format("2006-01-02T15:04:05Z07:00"), cronID,
	)
	s.mu.Unlock()

	if err != nil {
		apierr.WriteError(w, apierr.Internal(err.Error()))
		return
	}
	respondJSON(w, http.StatusOK, job)
}

func handleDeleteJob(w http.ResponseWriter, s *Scheduler, cronID string) {
	if err := s.RemoveJob(cronID); err != nil {
		apierr.WriteError(w, apierr.Internal(err.Error()))
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "removed", "id": cronID})
}

func handleTriggerJob(w http.ResponseWriter, s *Scheduler, cronID string) {
	if err := s.TriggerJob(cronID); err != nil {
		apierr.WriteError(w, apierr.NotFound(err.Error()))
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "triggered", "id": cronID})
}

func handlePauseJob(w http.ResponseWriter, s *Scheduler, cronID string) {
	if err := s.PauseJob(cronID); err != nil {
		apierr.WriteError(w, apierr.Internal(err.Error()))
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "paused", "id": cronID})
}

func handleResumeJob(w http.ResponseWriter, s *Scheduler, cronID string) {
	if err := s.ResumeJob(cronID); err != nil {
		apierr.WriteError(w, apierr.Internal(err.Error()))
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "resumed", "id": cronID})
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
