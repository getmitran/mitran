package hr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"github.com/getmitran/mitran/server/apierr"
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
	jc, ac     int
}

func NewRecruitmentStore() *RecruitmentStore { return &RecruitmentStore{} }

func (s *RecruitmentStore) CreateJob(j *JobPosting) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jc++
	j.ID = fmt.Sprintf("job-%d", s.jc)
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
	s.ac++
	a.ID = fmt.Sprintf("app-%d", s.ac)
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

func RegisterRecruitmentRoutes(mux *http.ServeMux, s *RecruitmentStore) {
	mux.HandleFunc("POST /api/v1/hr/jobs", func(w http.ResponseWriter, r *http.Request) {
		var j JobPosting
		json.NewDecoder(r.Body).Decode(&j)
		s.CreateJob(&j)
		json.NewEncoder(w).Encode(j)
	})
	mux.HandleFunc("GET /api/v1/hr/jobs", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(s.ListJobs(r.URL.Query().Get("status")))
	})
	mux.HandleFunc("POST /api/v1/hr/jobs/{id}/applicants", func(w http.ResponseWriter, r *http.Request) {
		var a Applicant
		json.NewDecoder(r.Body).Decode(&a)
		a.JobID = r.PathValue("id")
		s.AddApplicant(&a)
		json.NewEncoder(w).Encode(a)
	})
	mux.HandleFunc("PUT /api/v1/hr/applicants/{id}/stage", func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Stage string }
		json.NewDecoder(r.Body).Decode(&body)
		if !s.MoveStage(r.PathValue("id"), body.Stage) {
			apierr.WriteError(w, apierr.NotFound("not found"))
			return
		}
		w.Write([]byte(`{"ok":true}`))
	})
}
