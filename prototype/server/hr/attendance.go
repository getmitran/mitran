package hr

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ClockEntry struct {
	ID, EmployeeID, Date, ClockIn, ClockOut, Status string
	HoursWorked                                     float64
}

type AttendanceStore struct {
	mu      sync.Mutex
	entries map[string][]ClockEntry
}

func NewAttendanceStore() *AttendanceStore {
	return &AttendanceStore{entries: make(map[string][]ClockEntry)}
}

func (s *AttendanceStore) ClockIn(empID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.entries[empID] = append(s.entries[empID], ClockEntry{
		ID: empID + "-" + now.Format("20060102-150405"), EmployeeID: empID,
		Date: now.Format("2006-01-02"), ClockIn: now.Format(time.RFC3339), Status: "present",
	})
}

func (s *AttendanceStore) ClockOut(empID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := len(s.entries[empID]) - 1; i >= 0; i-- {
		e := &s.entries[empID][i]
		if e.ClockOut == "" {
			e.ClockOut = time.Now().Format(time.RFC3339)
			if ci, err := time.Parse(time.RFC3339, e.ClockIn); err == nil {
				e.HoursWorked = time.Since(ci).Hours()
			}
			return
		}
	}
}

func (s *AttendanceStore) GetTimesheet(empID, month string) []ClockEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	var res []ClockEntry
	for _, e := range s.entries[empID] {
		if strings.HasPrefix(e.Date, month) {
			res = append(res, e)
		}
	}
	return res
}

func (s *AttendanceStore) GetMonthlyHours(empID, month string) float64 {
	var total float64
	for _, e := range s.GetTimesheet(empID, month) {
		total += e.HoursWorked
	}
	return total
}

func (s *AttendanceStore) MarkAbsent(empID, date string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[empID] = append(s.entries[empID], ClockEntry{
		ID: empID + "-abs-" + date, EmployeeID: empID, Date: date, Status: "absent",
	})
}

func (s *AttendanceStore) Handler(w http.ResponseWriter, r *http.Request) {
	emp := r.URL.Query().Get("employee_id")
	switch r.Method + " " + r.URL.Path {
	case "POST /api/v1/hr/attendance/clock-in":
		json.NewDecoder(r.Body).Decode(&struct{ EmployeeID *string }{&emp})
		s.ClockIn(emp)
	case "POST /api/v1/hr/attendance/clock-out":
		json.NewDecoder(r.Body).Decode(&struct{ EmployeeID *string }{&emp})
		s.ClockOut(emp)
	default:
		json.NewEncoder(w).Encode(s.GetTimesheet(emp, r.URL.Query().Get("month")))
		return
	}
	w.WriteHeader(http.StatusOK)
}
