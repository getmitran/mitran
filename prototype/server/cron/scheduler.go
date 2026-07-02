package cron

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/getmitran/mitran/server/db"
)

// CronJob represents a scheduled recurring job
type CronJob struct {
	ID            string     `json:"id"`
	ProjectID     string     `json:"project_id"`
	Name          string     `json:"name"`
	AgentType     string     `json:"agent_type"`
	Message       string     `json:"message"`
	ScheduleType  string     `json:"schedule_type"`  // "interval" or "cron"
	ScheduleValue string     `json:"schedule_value"` // seconds for interval, cron expr for cron
	NextRun       time.Time  `json:"next_run"`
	LastRun       *time.Time `json:"last_run,omitempty"`
	Status        string     `json:"status"` // "active" or "paused"
	CreatedAt     time.Time  `json:"created_at"`
}

// Scheduler manages cron jobs with SQLite persistence and task dispatch
type Scheduler struct {
	mu    sync.RWMutex
	store *db.Store
	done  chan struct{}
}

// NewScheduler creates a scheduler backed by the given store
func NewScheduler(store *db.Store) *Scheduler {
	return &Scheduler{
		store: store,
		done:  make(chan struct{}),
	}
}

// Start begins the scheduler loop, checking every 30s for due jobs
func (s *Scheduler) Start(ctx context.Context) {
	s.ensureTable()
	go s.loop(ctx)
	log.Println("[cron] scheduler started (30s tick)")
}

// Stop halts the scheduler loop
func (s *Scheduler) Stop() {
	close(s.done)
	log.Println("[cron] scheduler stopped")
}

// AddJob persists a new cron job and returns it with computed NextRun
func (s *Scheduler) AddJob(job CronJob) (*CronJob, error) {
	if job.ID == "" {
		job.ID = fmt.Sprintf("cron_%d", time.Now().UnixNano())
	}
	if job.Status == "" {
		job.Status = "active"
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now()
	}
	nextRun, err := computeNextRun(job.ScheduleType, job.ScheduleValue, time.Now())
	if err != nil {
		return nil, fmt.Errorf("invalid schedule: %w", err)
	}
	job.NextRun = nextRun

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err = s.store.DB().Exec(
		`INSERT INTO cron_jobs (id, project_id, name, agent_type, message, schedule_type, schedule_value, next_run, last_run, status, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		job.ID, job.ProjectID, job.Name, job.AgentType, job.Message,
		job.ScheduleType, job.ScheduleValue,
		job.NextRun.UTC().Format(time.RFC3339),
		nil, job.Status, job.CreatedAt.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// RemoveJob deletes a cron job by ID
func (s *Scheduler) RemoveJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.store.DB().Exec(`DELETE FROM cron_jobs WHERE id = ?`, id)
	return err
}

// PauseJob sets a job's status to paused
func (s *Scheduler) PauseJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.store.DB().Exec(`UPDATE cron_jobs SET status = 'paused' WHERE id = ?`, id)
	return err
}

// ResumeJob sets a job's status to active
func (s *Scheduler) ResumeJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.store.DB().Exec(`UPDATE cron_jobs SET status = 'active' WHERE id = ?`, id)
	return err
}

// TriggerJob immediately dispatches a job as a task
func (s *Scheduler) TriggerJob(id string) error {
	s.mu.RLock()
	job, err := s.getJobByID(id)
	s.mu.RUnlock()
	if err != nil {
		return err
	}
	return s.dispatchJob(*job)
}

// ListJobs returns all jobs for a project (or all if projectID is empty)
func (s *Scheduler) ListJobs(projectID string) ([]CronJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rows *sql.Rows
	var err error
	if projectID == "" {
		rows, err = s.store.DB().Query(`SELECT id, project_id, name, agent_type, message, schedule_type, schedule_value, next_run, last_run, status, created_at FROM cron_jobs ORDER BY created_at DESC`)
	} else {
		rows, err = s.store.DB().Query(`SELECT id, project_id, name, agent_type, message, schedule_type, schedule_value, next_run, last_run, status, created_at FROM cron_jobs WHERE project_id = ? ORDER BY created_at DESC`, projectID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []CronJob
	for rows.Next() {
		var j CronJob
		var nextRun, createdAt string
		var lastRunNull sql.NullString
		if err := rows.Scan(&j.ID, &j.ProjectID, &j.Name, &j.AgentType, &j.Message,
			&j.ScheduleType, &j.ScheduleValue, &nextRun, &lastRunNull, &j.Status, &createdAt); err != nil {
			return nil, err
		}
		j.NextRun, _ = time.Parse(time.RFC3339, nextRun)
		j.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		if lastRunNull.Valid {
			t, _ := time.Parse(time.RFC3339, lastRunNull.String)
			j.LastRun = &t
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

// GetJob returns a single job by ID
func (s *Scheduler) GetJob(id string) (*CronJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getJobByID(id)
}

func (s *Scheduler) getJobByID(id string) (*CronJob, error) {
	row := s.store.DB().QueryRow(`SELECT id, project_id, name, agent_type, message, schedule_type, schedule_value, next_run, last_run, status, created_at FROM cron_jobs WHERE id = ?`, id)
	var j CronJob
	var nextRun, createdAt string
	var lastRunNull sql.NullString
	if err := row.Scan(&j.ID, &j.ProjectID, &j.Name, &j.AgentType, &j.Message,
		&j.ScheduleType, &j.ScheduleValue, &nextRun, &lastRunNull, &j.Status, &createdAt); err != nil {
		return nil, fmt.Errorf("job not found: %s", id)
	}
	j.NextRun, _ = time.Parse(time.RFC3339, nextRun)
	j.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	if lastRunNull.Valid {
		t, _ := time.Parse(time.RFC3339, lastRunNull.String)
		j.LastRun = &t
	}
	return &j, nil
}

func (s *Scheduler) loop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			return
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			s.tick(now)
		}
	}
}

func (s *Scheduler) tick(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.store.DB().Query(
		`SELECT id, project_id, name, agent_type, message, schedule_type, schedule_value, next_run, last_run, status, created_at
		 FROM cron_jobs WHERE status = 'active' AND next_run <= ?`,
		now.UTC().Format(time.RFC3339),
	)
	if err != nil {
		log.Printf("[cron] tick query error: %v", err)
		return
	}
	defer rows.Close()

	var dueJobs []CronJob
	for rows.Next() {
		var j CronJob
		var nextRun, createdAt string
		var lastRunNull sql.NullString
		if err := rows.Scan(&j.ID, &j.ProjectID, &j.Name, &j.AgentType, &j.Message,
			&j.ScheduleType, &j.ScheduleValue, &nextRun, &lastRunNull, &j.Status, &createdAt); err != nil {
			continue
		}
		j.NextRun, _ = time.Parse(time.RFC3339, nextRun)
		j.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		dueJobs = append(dueJobs, j)
	}

	for _, j := range dueJobs {
		go s.dispatchJob(j)

		nextRun, _ := computeNextRun(j.ScheduleType, j.ScheduleValue, now)
		s.store.DB().Exec(
			`UPDATE cron_jobs SET last_run = ?, next_run = ? WHERE id = ?`,
			now.UTC().Format(time.RFC3339),
			nextRun.UTC().Format(time.RFC3339),
			j.ID,
		)
	}
}

func (s *Scheduler) dispatchJob(job CronJob) error {
	task := db.Task{
		ID:          fmt.Sprintf("cron_task_%d", time.Now().UnixNano()),
		ProjectID:   job.ProjectID,
		Title:       fmt.Sprintf("[cron] %s", job.Name),
		Description: job.Message,
		AgentType:   job.AgentType,
		Priority:    5,
		Status:      db.StatusQueued,
		CreatedBy:   "scheduler",
		CreatedAt:   time.Now(),
	}
	return s.store.AddTask(task)
}

func (s *Scheduler) ensureTable() {
	s.store.DB().Exec(`CREATE TABLE IF NOT EXISTS cron_jobs (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		name TEXT NOT NULL,
		agent_type TEXT NOT NULL DEFAULT 'general',
		message TEXT NOT NULL DEFAULT '',
		schedule_type TEXT NOT NULL,
		schedule_value TEXT NOT NULL,
		next_run TEXT NOT NULL,
		last_run TEXT,
		status TEXT NOT NULL DEFAULT 'active',
		created_at TEXT NOT NULL
	)`)
}

// computeNextRun calculates the next execution time based on schedule
func computeNextRun(scheduleType, scheduleValue string, from time.Time) (time.Time, error) {
	switch scheduleType {
	case "interval":
		secs, err := strconv.Atoi(scheduleValue)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid interval: %s", scheduleValue)
		}
		if secs < 60 {
			secs = 60
		}
		return from.Add(time.Duration(secs) * time.Second), nil
	case "cron":
		return nextCronTime(scheduleValue, from)
	default:
		return time.Time{}, fmt.Errorf("unknown schedule type: %s", scheduleType)
	}
}

// nextCronTime parses a 5-field cron expression and computes the next run time
// Fields: minute hour day-of-month month day-of-week
func nextCronTime(expr string, from time.Time) (time.Time, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return time.Time{}, fmt.Errorf("cron expression must have 5 fields, got %d", len(fields))
	}

	minutes, err := parseCronField(fields[0], 0, 59)
	if err != nil {
		return time.Time{}, fmt.Errorf("minute field: %w", err)
	}
	hours, err := parseCronField(fields[1], 0, 23)
	if err != nil {
		return time.Time{}, fmt.Errorf("hour field: %w", err)
	}
	doms, err := parseCronField(fields[2], 1, 31)
	if err != nil {
		return time.Time{}, fmt.Errorf("day-of-month field: %w", err)
	}
	months, err := parseCronField(fields[3], 1, 12)
	if err != nil {
		return time.Time{}, fmt.Errorf("month field: %w", err)
	}
	dows, err := parseCronField(fields[4], 0, 6)
	if err != nil {
		return time.Time{}, fmt.Errorf("day-of-week field: %w", err)
	}

	// Brute-force search for next matching time (max 366 days ahead)
	candidate := from.Add(time.Minute).Truncate(time.Minute)
	limit := from.Add(366 * 24 * time.Hour)

	for candidate.Before(limit) {
		if contains(months, int(candidate.Month())) &&
			contains(doms, candidate.Day()) &&
			contains(dows, int(candidate.Weekday())) &&
			contains(hours, candidate.Hour()) &&
			contains(minutes, candidate.Minute()) {
			return candidate, nil
		}
		candidate = candidate.Add(time.Minute)
	}
	return time.Time{}, fmt.Errorf("no matching time found within 366 days")
}

// parseCronField parses a single cron field into a list of valid values
func parseCronField(field string, min, max int) ([]int, error) {
	if field == "*" {
		return makeRange(min, max), nil
	}

	var result []int
	parts := strings.Split(field, ",")
	for _, part := range parts {
		// Handle step values: */5 or 1-30/5
		if strings.Contains(part, "/") {
			stepParts := strings.SplitN(part, "/", 2)
			step, err := strconv.Atoi(stepParts[1])
			if err != nil || step <= 0 {
				return nil, fmt.Errorf("invalid step: %s", part)
			}
			rangeStart, rangeEnd := min, max
			if stepParts[0] != "*" {
				if strings.Contains(stepParts[0], "-") {
					rp := strings.SplitN(stepParts[0], "-", 2)
					rangeStart, _ = strconv.Atoi(rp[0])
					rangeEnd, _ = strconv.Atoi(rp[1])
				} else {
					rangeStart, _ = strconv.Atoi(stepParts[0])
				}
			}
			for i := rangeStart; i <= rangeEnd; i += step {
				result = append(result, i)
			}
		} else if strings.Contains(part, "-") {
			// Handle ranges: 1-5, MON-FRI
			rp := strings.SplitN(part, "-", 2)
			start, err1 := parseCronValue(rp[0], min, max)
			end, err2 := parseCronValue(rp[1], min, max)
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else {
			// Single value
			v, err := parseCronValue(part, min, max)
			if err != nil {
				return nil, fmt.Errorf("invalid value: %s", part)
			}
			result = append(result, v)
		}
	}
	return result, nil
}

var dowNames = map[string]int{
	"SUN": 0, "MON": 1, "TUE": 2, "WED": 3, "THU": 4, "FRI": 5, "SAT": 6,
}

func parseCronValue(s string, min, max int) (int, error) {
	upper := strings.ToUpper(s)
	if v, ok := dowNames[upper]; ok {
		return v, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if v < min || v > max {
		return 0, fmt.Errorf("%d out of range [%d-%d]", v, min, max)
	}
	return v, nil
}

func makeRange(min, max int) []int {
	r := make([]int, 0, max-min+1)
	for i := min; i <= max; i++ {
		r = append(r, i)
	}
	return r
}

func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
