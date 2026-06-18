package hr

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type JobPosting struct {
	ID, Title, Department, Description, Requirements, Location, Type, Status, PostedBy, PostedAt string
}

type Applicant struct {
	ID, JobID, Name, Email, Phone, ResumeURL, Stage, AppliedAt, Notes string
}

type RecruitmentStore struct {
	mu         sync.Mutex
	jobs       []JobPosting
	applicants []Applicant
}

var store = &RecruitmentStore{}

func (s *RecruitmentStore) CreateJob(j *JobPosting) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j.ID = uuid.New().String()
	j.PostedAt = time.Now().Format(time.RFC3339)
	if j.Status == "" {
		j.Status = "open"
	}
	s.jobs = append(s.jobs, *j)
}

func (s *RecruitmentStore) ListJobs(status string) []JobPosting {
	s.mu.Lock()
	defer s.mu.Unlock()
	if status == "" {
		return s.jobs
	}
	var out []JobPosting
	for _, j := range s.jobs {
		if j.Status == status {
			out = append(out, j)
		}
	}
	return out
}

func (s *RecruitmentStore) AddApplicant(a *Applicant) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a.ID = uuid.New().String()
	a.AppliedAt = time.Now().Format(time.RFC3339)
	if a.Stage == "" {
		a.Stage = "applied"
	}
	s.applicants = append(s.applicants, *a)
}

func (s *RecruitmentStore) MoveStage(id, stage string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.applicants {
		if s.applicants[i].ID == id {
			s.applicants[i].Stage = stage
			return true
		}
	}
	return false
}

func (s *RecruitmentStore) ListApplicants(jobID, stage string) []Applicant {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Applicant
	for _, a := range s.applicants {
		if a.JobID == jobID && (stage == "" || a.Stage == stage) {
			out = append(out, a)
		}
	}
	return out
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/hr/jobs", func(w http.ResponseWriter, r *http.Request) {
		var j JobPosting
		json.NewDecoder(r.Body).Decode(&j)
		store.CreateJob(&j)
		json.NewEncoder(w).Encode(j)
	})
	mux.HandleFunc("GET /api/v1/hr/jobs", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(store.ListJobs(r.URL.Query().Get("status")))
	})
	mux.HandleFunc("POST /api/v1/hr/jobs/{id}/applicants", func(w http.ResponseWriter, r *http.Request) {
		var a Applicant
		json.NewDecoder(r.Body).Decode(&a)
		a.JobID = r.PathValue("id")
		store.AddApplicant(&a)
		json.NewEncoder(w).Encode(a)
	})
	mux.HandleFunc("PUT /api/v1/hr/applicants/{id}/stage", func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Stage string }
		json.NewDecoder(r.Body).Decode(&body)
		if !store.MoveStage(r.PathValue("id"), body.Stage) {
			http.Error(w, "not found", 404)
			return
		}
		w.Write([]byte(`{"ok":true}`))
	})
}
