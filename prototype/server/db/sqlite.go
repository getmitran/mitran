package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

type Ticket struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Assignee    string    `json:"assignee"`
	Labels      []string  `json:"labels"`
	SprintID    string    `json:"sprint_id"`
	SLADueAt    *time.Time `json:"sla_due_at,omitempty"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TicketComment struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticket_id"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type Sprint struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
}

type WikiPage struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ParentID  string    `json:"parent_id"`
	Slug      string    `json:"slug"`
	CreatedBy string    `json:"created_by"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Tickets             []Ticket              `json:"tickets"`
	TicketComments      []TicketComment       `json:"ticket_comments"`
	Sprints             []Sprint              `json:"sprints"`
	WikiPages           []WikiPage            `json:"wiki_pages"`
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

// --- Ticket methods ---

func (s *Store) AddTicket(t Ticket) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Tickets = append(s.Tickets, t)
	return s.Save()
}

func (s *Store) GetTicket(id string) *Ticket {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	for i := range s.Tickets {
		if s.Tickets[i].ID == id {
			return &s.Tickets[i]
		}
	}
	return nil
}

func (s *Store) UpdateTicket(id string, fn func(*Ticket)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.Tickets {
		if s.Tickets[i].ID == id {
			fn(&s.Tickets[i])
			return s.Save()
		}
	}
	return fmt.Errorf("ticket %s not found", id)
}

func (s *Store) DeleteTicket(id string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.Tickets {
		if s.Tickets[i].ID == id {
			s.Tickets = append(s.Tickets[:i], s.Tickets[i+1:]...)
			return s.Save()
		}
	}
	return fmt.Errorf("ticket %s not found", id)
}

func (s *Store) ListTickets(status, assignee, sprintID, labels string) []Ticket {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	var result []Ticket
	for _, t := range s.Tickets {
		if status != "" && t.Status != status {
			continue
		}
		if assignee != "" && t.Assignee != assignee {
			continue
		}
		if sprintID != "" && t.SprintID != sprintID {
			continue
		}
		if labels != "" {
			found := false
			for _, l := range t.Labels {
				if l == labels {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		result = append(result, t)
	}
	return result
}

func (s *Store) AddTicketComment(c TicketComment) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.TicketComments = append(s.TicketComments, c)
	return s.Save()
}

func (s *Store) ListTicketComments(ticketID string) []TicketComment {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	var result []TicketComment
	for _, c := range s.TicketComments {
		if c.TicketID == ticketID {
			result = append(result, c)
		}
	}
	return result
}

// --- Sprint methods ---

func (s *Store) AddSprint(sp Sprint) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Sprints = append(s.Sprints, sp)
	return s.Save()
}

func (s *Store) ListSprints() []Sprint {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.Sprints
}

func (s *Store) UpdateSprint(id string, fn func(*Sprint)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.Sprints {
		if s.Sprints[i].ID == id {
			fn(&s.Sprints[i])
			return s.Save()
		}
	}
	return fmt.Errorf("sprint %s not found", id)
}

// --- Wiki methods ---

func (s *Store) AddWikiPage(p WikiPage) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.WikiPages = append(s.WikiPages, p)
	return s.Save()
}

func (s *Store) GetWikiPage(id string) *WikiPage {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	for i := range s.WikiPages {
		if s.WikiPages[i].ID == id {
			return &s.WikiPages[i]
		}
	}
	return nil
}

func (s *Store) UpdateWikiPage(id string, fn func(*WikiPage)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.WikiPages {
		if s.WikiPages[i].ID == id {
			fn(&s.WikiPages[i])
			return s.Save()
		}
	}
	return fmt.Errorf("wiki page %s not found", id)
}

func (s *Store) DeleteWikiPage(id string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.WikiPages {
		if s.WikiPages[i].ID == id {
			s.WikiPages = append(s.WikiPages[:i], s.WikiPages[i+1:]...)
			return s.Save()
		}
	}
	return fmt.Errorf("wiki page %s not found", id)
}

func (s *Store) ListWikiPages() []WikiPage {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.WikiPages
}

func (s *Store) SearchWikiPages(query string) []WikiPage {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	var result []WikiPage
	for _, p := range s.WikiPages {
		if containsLower(p.Title, query) || containsLower(p.Content, query) {
			result = append(result, p)
		}
	}
	return result
}

func containsLower(s, sub string) bool {
	return len(s) >= len(sub) && strings.Contains(strings.ToLower(s), sub)
}
