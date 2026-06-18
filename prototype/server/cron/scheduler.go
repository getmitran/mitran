package cron

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Schedule string

const (
	Every1Min  Schedule = "1m"
	Every5Min  Schedule = "5m"
	Every15Min Schedule = "15m"
	EveryHour  Schedule = "1h"
	EveryDay   Schedule = "24h"
)

type Job struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Agent     string     `json:"agent"`
	Schedule  Schedule   `json:"schedule"`
	Payload   string     `json:"payload"`
	Active    bool       `json:"active"`
	LastRun   *time.Time `json:"last_run,omitempty"`
	NextRun   time.Time  `json:"next_run"`
	RunCount  int        `json:"run_count"`
	CreatedAt time.Time  `json:"created_at"`
}

type Scheduler struct {
	mu       sync.RWMutex
	jobs     []Job
	dir      string
	stopCh   chan struct{}
	callback func(job Job)
}

func New(dir string, callback func(Job)) (*Scheduler, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	s := &Scheduler{dir: dir, stopCh: make(chan struct{}), callback: callback}
	s.load()
	return s, nil
}

func (s *Scheduler) Start() {
	go s.loop()
	log.Println("[cron] scheduler started")
}

func (s *Scheduler) Stop() { close(s.stopCh) }

func (s *Scheduler) Add(name, agent string, schedule Schedule, payload string) *Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	j := Job{
		ID: time.Now().Format("20060102150405"), Name: name, Agent: agent,
		Schedule: schedule, Payload: payload, Active: true,
		NextRun: time.Now().Add(parseDuration(schedule)), CreatedAt: time.Now(),
	}
	s.jobs = append(s.jobs, j)
	s.save()
	return &s.jobs[len(s.jobs)-1]
}

func (s *Scheduler) Remove(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, j := range s.jobs {
		if j.ID == id {
			s.jobs = append(s.jobs[:i], s.jobs[i+1:]...)
			break
		}
	}
	s.save()
}

func (s *Scheduler) List() []Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Job{}, s.jobs...)
}

func (s *Scheduler) Pause(id string)  { s.setActive(id, false) }
func (s *Scheduler) Resume(id string) { s.setActive(id, true) }

func (s *Scheduler) setActive(id string, active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.jobs {
		if s.jobs[i].ID == id {
			s.jobs[i].Active = active
			break
		}
	}
	s.save()
}

func (s *Scheduler) Trigger(id string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, j := range s.jobs {
		if j.ID == id {
			go s.callback(j)
			return
		}
	}
}

func (s *Scheduler) loop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case now := <-ticker.C:
			s.mu.Lock()
			for i := range s.jobs {
				if s.jobs[i].Active && now.After(s.jobs[i].NextRun) {
					j := s.jobs[i]
					nowT := now
					s.jobs[i].LastRun = &nowT
					s.jobs[i].RunCount++
					s.jobs[i].NextRun = now.Add(parseDuration(j.Schedule))
					go s.callback(j)
				}
			}
			s.save()
			s.mu.Unlock()
		}
	}
}

func (s *Scheduler) load() {
	data, err := os.ReadFile(filepath.Join(s.dir, "cron_jobs.json"))
	if err == nil {
		json.Unmarshal(data, &s.jobs)
	}
}

func (s *Scheduler) save() {
	data, _ := json.MarshalIndent(s.jobs, "", "  ")
	os.WriteFile(filepath.Join(s.dir, "cron_jobs.json"), data, 0644)
}

func parseDuration(s Schedule) time.Duration {
	d, _ := time.ParseDuration(string(s))
	if d == 0 {
		d = time.Hour
	}
	return d
}
