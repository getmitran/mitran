package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Status string

const (
	StatusQueued     Status = "queued"
	StatusRunning    Status = "running"
	StatusCheckpoint Status = "checkpoint"
	StatusApproved   Status = "approved"
	StatusRejected   Status = "rejected"
	StatusFailed     Status = "failed"
)

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Languages   []string  `json:"languages"`
	TeamSize    int       `json:"team_size"`
	CreatedAt   time.Time `json:"created_at"`
}

type Task struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AgentType   string    `json:"agent_type"`
	Priority    int       `json:"priority"`
	Status      Status    `json:"status"`
	DependsOn   []string  `json:"depends_on"`
	Resources   []string  `json:"resources"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	AssignedTo  string    `json:"assigned_to,omitempty"`
}

type Checkpoint struct {
	ID               string    `json:"id"`
	TaskID           string    `json:"task_id"`
	AgentID          string    `json:"agent_id"`
	Output           string    `json:"output"`
	Diffs            []string  `json:"diffs"`
	Questions        []string  `json:"questions"`
	DecisionAction   string    `json:"decision_action,omitempty"`
	DecisionFeedback string    `json:"decision_feedback,omitempty"`
	DecisionBy       string    `json:"decision_by,omitempty"`
	DecisionAt       *time.Time `json:"decision_at,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type Agent struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	CallbackURL   string    `json:"callback_url"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	CreatedAt     time.Time `json:"created_at"`
}

type Environment struct {
	Name             string    `json:"name"`
	AutoDeploy       bool      `json:"auto_deploy"`
	ApprovalRequired bool      `json:"approval_required"`
	RollbackEnabled  bool      `json:"rollback_enabled"`
	CreatedAt        time.Time `json:"created_at"`
}

type PipelineStage struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Steps       []string `json:"steps"`
}

type Pipeline struct {
	ID        string          `json:"id"`
	ProjectID string          `json:"project_id"`
	Name      string          `json:"name"`
	Stages    []PipelineStage `json:"stages"`
	CreatedAt time.Time       `json:"created_at"`
}

type PipelineRun struct {
	ID         string     `json:"id"`
	PipelineID string     `json:"pipeline_id"`
	ProjectID  string     `json:"project_id"`
	Branch     string     `json:"branch"`
	Commit     string     `json:"commit"`
	Status     string     `json:"status"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type Deployment struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Environment string     `json:"environment"`
	Version     string     `json:"version"`
	Commit      string     `json:"commit"`
	Status      string     `json:"status"`
	TriggeredAt time.Time  `json:"triggered_at"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
}

type ProjectEnvironments struct {
	ProjectID    string        `json:"project_id"`
	Environments []Environment `json:"environments"`
}

type Store struct {
	Mu                  sync.RWMutex // exported via accessor below
	dir                 string
	Projects            []Project             `json:"projects"`
	Tasks               []Task                `json:"tasks"`
	Checkpoints         []Checkpoint          `json:"checkpoints"`
	Agents              []Agent               `json:"agents"`
	Config              map[string]string     `json:"config"`
	ProjectEnvironments []ProjectEnvironments `json:"project_environments"`
	Pipelines           []Pipeline            `json:"pipelines"`
	PipelineRuns        []PipelineRun         `json:"pipeline_runs"`
	Deployments         []Deployment          `json:"deployments"`
}

func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	s := &Store{
		dir:    dir,
		Config: make(map[string]string),
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

func (s *Store) path() string {
	return filepath.Join(s.dir, "state.json")
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}

func (s *Store) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(), data, 0644)
}

func (s *Store) AddProject(p Project) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Projects = append(s.Projects, p)
	return s.Save()
}

func (s *Store) GetProject(id string) *Project {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	for i := range s.Projects {
		if s.Projects[i].ID == id {
			return &s.Projects[i]
		}
	}
	return nil
}

func (s *Store) AddTask(t Task) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Tasks = append(s.Tasks, t)
	return s.Save()
}

func (s *Store) AddTasks(tasks []Task) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Tasks = append(s.Tasks, tasks...)
	return s.Save()
}

func (s *Store) GetTask(id string) *Task {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.Tasks {
		if s.Tasks[i].ID == id {
			return &s.Tasks[i]
		}
	}
	return nil
}

func (s *Store) UpdateTask(id string, fn func(*Task)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.Tasks {
		if s.Tasks[i].ID == id {
			fn(&s.Tasks[i])
			return s.Save()
		}
	}
	return fmt.Errorf("task %s not found", id)
}

func (s *Store) ListTasks(status Status, projectID string) []Task {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	var result []Task
	for _, t := range s.Tasks {
		if status != "" && t.Status != status {
			continue
		}
		if projectID != "" && t.ProjectID != projectID {
			continue
		}
		result = append(result, t)
	}
	return result
}

func (s *Store) AddCheckpoint(c Checkpoint) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Checkpoints = append(s.Checkpoints, c)
	return s.Save()
}

func (s *Store) GetCheckpoint(id string) *Checkpoint {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.Checkpoints {
		if s.Checkpoints[i].ID == id {
			return &s.Checkpoints[i]
		}
	}
	return nil
}

func (s *Store) UpdateCheckpoint(id string, fn func(*Checkpoint)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.Checkpoints {
		if s.Checkpoints[i].ID == id {
			fn(&s.Checkpoints[i])
			return s.Save()
		}
	}
	return fmt.Errorf("checkpoint %s not found", id)
}

func (s *Store) ListPendingCheckpoints() []Checkpoint {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	var result []Checkpoint
	for _, c := range s.Checkpoints {
		if c.DecisionAction == "" {
			result = append(result, c)
		}
	}
	return result
}

func (s *Store) AddAgent(a Agent) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.Agents {
		if s.Agents[i].ID == a.ID {
			s.Agents[i] = a
			return s.Save()
		}
	}
	s.Agents = append(s.Agents, a)
	return s.Save()
}

func (s *Store) GetAgent(id string) *Agent {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	for i := range s.Agents {
		if s.Agents[i].ID == id {
			return &s.Agents[i]
		}
	}
	return nil
}

func (s *Store) GetAgentByType(agentType string) *Agent {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	for i := range s.Agents {
		if s.Agents[i].Type == agentType && s.Agents[i].Status == "online" {
			return &s.Agents[i]
		}
	}
	return nil
}

func (s *Store) UpdateAgentHeartbeat(id string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.Agents {
		if s.Agents[i].ID == id {
			s.Agents[i].LastHeartbeat = time.Now()
			s.Agents[i].Status = "online"
			return s.Save()
		}
	}
	return fmt.Errorf("agent %s not found", id)
}

// --- Environment methods ---

func (s *Store) SetEnvironments(projectID string, envs []Environment) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.ProjectEnvironments {
		if s.ProjectEnvironments[i].ProjectID == projectID {
			s.ProjectEnvironments[i].Environments = envs
			return s.Save()
		}
	}
	s.ProjectEnvironments = append(s.ProjectEnvironments, ProjectEnvironments{ProjectID: projectID, Environments: envs})
	return s.Save()
}

func (s *Store) GetEnvironments(projectID string) []Environment {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	for _, pe := range s.ProjectEnvironments {
		if pe.ProjectID == projectID {
			return pe.Environments
		}
	}
	return nil
}

func (s *Store) UpdateEnvironment(projectID, envName string, fn func(*Environment)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.ProjectEnvironments {
		if s.ProjectEnvironments[i].ProjectID == projectID {
			for j := range s.ProjectEnvironments[i].Environments {
				if s.ProjectEnvironments[i].Environments[j].Name == envName {
					fn(&s.ProjectEnvironments[i].Environments[j])
					return s.Save()
				}
			}
			return fmt.Errorf("environment %s not found", envName)
		}
	}
	return fmt.Errorf("no environments configured for project %s", projectID)
}

// --- Pipeline methods ---

func (s *Store) AddPipeline(p Pipeline) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Pipelines = append(s.Pipelines, p)
	return s.Save()
}

func (s *Store) ListPipelines(projectID string) []Pipeline {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	var result []Pipeline
	for _, p := range s.Pipelines {
		if p.ProjectID == projectID {
			result = append(result, p)
		}
	}
	return result
}

func (s *Store) AddPipelineRun(run PipelineRun) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.PipelineRuns = append(s.PipelineRuns, run)
	return s.Save()
}

func (s *Store) ListPipelineRuns(pipelineID string) []PipelineRun {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	var result []PipelineRun
	for _, r := range s.PipelineRuns {
		if r.PipelineID == pipelineID {
			result = append(result, r)
		}
	}
	return result
}

// --- Deployment methods ---

func (s *Store) AddDeployment(d Deployment) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Deployments = append(s.Deployments, d)
	return s.Save()
}

func (s *Store) ListDeployments(projectID string) []Deployment {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	var result []Deployment
	for _, d := range s.Deployments {
		if d.ProjectID == projectID {
			result = append(result, d)
		}
	}
	return result
}
