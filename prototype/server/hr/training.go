package hr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Course struct {
	ID, Title, Description, Category, Duration, Instructor, Format string
	Required                                                       bool
}

type Enrollment struct {
	ID, EmployeeID, CourseID, Status, EnrolledAt, CompletedAt string
	Score                                                     int
}

type TrainingStore struct {
	mu          sync.Mutex
	courses     []Course
	enrollments []Enrollment
}

var ts = &TrainingStore{}

func (s *TrainingStore) AddCourse(c Course) { s.mu.Lock(); defer s.mu.Unlock(); c.ID = fmt.Sprintf("CRS-%d", len(s.courses)+1); s.courses = append(s.courses, c) }
func (s *TrainingStore) ListCourses(cat string) []Course {
	s.mu.Lock(); defer s.mu.Unlock()
	if cat == "" { return s.courses }
	var r []Course; for _, c := range s.courses { if c.Category == cat { r = append(r, c) } }; return r
}
func (s *TrainingStore) Enroll(empID, courseID string) Enrollment {
	s.mu.Lock(); defer s.mu.Unlock()
	e := Enrollment{ID: fmt.Sprintf("ENR-%d", len(s.enrollments)+1), EmployeeID: empID, CourseID: courseID, Status: "enrolled", EnrolledAt: time.Now().Format(time.RFC3339)}
	s.enrollments = append(s.enrollments, e); return e
}
func (s *TrainingStore) Complete(id string, score int) {
	s.mu.Lock(); defer s.mu.Unlock()
	for i := range s.enrollments { if s.enrollments[i].ID == id { s.enrollments[i].Status = "completed"; s.enrollments[i].Score = score; s.enrollments[i].CompletedAt = time.Now().Format(time.RFC3339) } }
}
func (s *TrainingStore) GetTranscript(empID string) []Enrollment {
	s.mu.Lock(); defer s.mu.Unlock()
	var r []Enrollment; for _, e := range s.enrollments { if e.EmployeeID == empID { r = append(r, e) } }; return r
}
func (s *TrainingStore) GetCompletionRate(courseID string) float64 {
	s.mu.Lock(); defer s.mu.Unlock()
	var total, done int; for _, e := range s.enrollments { if e.CourseID == courseID { total++; if e.Status == "completed" { done++ } } }
	if total == 0 { return 0 }; return float64(done) / float64(total) * 100
}

func RegisterTrainingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/hr/courses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var c Course; json.NewDecoder(r.Body).Decode(&c); ts.AddCourse(c); json.NewEncoder(w).Encode(c)
		} else {
			json.NewEncoder(w).Encode(ts.ListCourses(r.URL.Query().Get("category")))
		}
	})
	mux.HandleFunc("/api/v1/hr/training/enroll", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ EmployeeID, CourseID string }; json.NewDecoder(r.Body).Decode(&req)
		json.NewEncoder(w).Encode(ts.Enroll(req.EmployeeID, req.CourseID))
	})
	mux.HandleFunc("/api/v1/hr/training/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/hr/training/"), "/")
		if len(parts) >= 2 && parts[1] == "transcript" { json.NewEncoder(w).Encode(ts.GetTranscript(parts[0])) }
	})
}
