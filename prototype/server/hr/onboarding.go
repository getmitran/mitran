package hr

import (
	"fmt"
	"sync"
	"time"
)

type OnboardingTask struct {
	ID, Title, Description, Category, AssignedTo string
	Required                                     bool
	DueInDays                                    int
}

type TaskStatus struct {
	OnboardingTask
	Status      string `json:"status"` // pending, done, skipped
	CompletedAt string `json:"completed_at,omitempty"`
}

type OnboardingFlow struct {
	ID         string       `json:"id"`
	EmployeeID string       `json:"employee_id"`
	Tasks      []TaskStatus `json:"tasks"`
	StartDate  string       `json:"start_date"`
}

type OnboardingStore struct {
	flows   map[string]*OnboardingFlow
	flowsMu sync.RWMutex
	counter int
}

func NewOnboardingStore() *OnboardingStore {
	return &OnboardingStore{flows: map[string]*OnboardingFlow{}}
}

var defaultTasks = []OnboardingTask{
	{"t1", "Setup laptop", "Configure dev machine with required tools", "it", "it-team", true, 1},
	{"t2", "Create email account", "Provision corporate email", "it", "it-team", true, 1},
	{"t3", "Setup code access", "Grant repo and CI/CD permissions", "it", "it-team", true, 2},
	{"t4", "Sign NDA", "Review and sign non-disclosure agreement", "compliance", "hr-team", true, 1},
	{"t5", "Sign code of conduct", "Acknowledge team policies", "compliance", "hr-team", true, 2},
	{"t6", "Complete security training", "Mandatory infosec awareness course", "compliance", "hr-team", true, 5},
	{"t7", "Meet your manager", "Intro call with direct manager", "team", "manager", true, 2},
	{"t8", "Meet the team", "Team lunch or virtual meet-and-greet", "team", "manager", false, 5},
	{"t9", "Setup benefits", "Enroll in health and retirement plans", "hr", "hr-team", true, 7},
	{"t10", "First week checkpoint", "30-min check-in with HR", "hr", "hr-team", false, 7},
}

func (s *OnboardingStore) CreateFlow(employeeID string) *OnboardingFlow {
	s.flowsMu.Lock()
	defer s.flowsMu.Unlock()
	s.counter++
	f := &OnboardingFlow{
		ID: fmt.Sprintf("onb-%d", s.counter), EmployeeID: employeeID,
		StartDate: time.Now().Format("2006-01-02"),
		Tasks:     make([]TaskStatus, len(defaultTasks)),
	}
	for i, t := range defaultTasks {
		f.Tasks[i] = TaskStatus{OnboardingTask: t, Status: "pending"}
	}
	s.flows[f.ID] = f
	return f
}

func (s *OnboardingStore) CompleteTask(flowID, taskID string) error {
	s.flowsMu.Lock()
	defer s.flowsMu.Unlock()
	f, ok := s.flows[flowID]
	if !ok {
		return fmt.Errorf("flow not found")
	}
	for i := range f.Tasks {
		if f.Tasks[i].ID == taskID {
			f.Tasks[i].Status = "done"
			f.Tasks[i].CompletedAt = time.Now().Format(time.RFC3339)
			return nil
		}
	}
	return fmt.Errorf("task not found")
}

func (s *OnboardingStore) GetProgress(flowID string) (done, total int, err error) {
	s.flowsMu.RLock()
	defer s.flowsMu.RUnlock()
	f, ok := s.flows[flowID]
	if !ok {
		return 0, 0, fmt.Errorf("flow not found")
	}
	for _, t := range f.Tasks {
		if t.Status == "done" {
			done++
		}
	}
	return done, len(f.Tasks), nil
}

func (s *OnboardingStore) GetFlow(flowID string) (*OnboardingFlow, bool) {
	s.flowsMu.RLock()
	defer s.flowsMu.RUnlock()
	f, ok := s.flows[flowID]
	return f, ok
}
