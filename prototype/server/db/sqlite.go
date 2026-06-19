package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
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
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	AgentType   string     `json:"agent_type"`
	Priority    int        `json:"priority"`
	Status      Status     `json:"status"`
	DependsOn   []string   `json:"depends_on"`
	Resources   []string   `json:"resources"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	AssignedTo  string     `json:"assigned_to,omitempty"`
}

type Checkpoint struct {
	ID               string     `json:"id"`
	TaskID           string     `json:"task_id"`
	AgentID          string     `json:"agent_id"`
	Output           string     `json:"output"`
	Diffs            []string   `json:"diffs"`
	Questions        []string   `json:"questions"`
	DecisionAction   string     `json:"decision_action,omitempty"`
	DecisionFeedback string     `json:"decision_feedback,omitempty"`
	DecisionBy       string     `json:"decision_by,omitempty"`
	DecisionAt       *time.Time `json:"decision_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
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
	Name        string   `json:"name"`
	Environment string   `json:"environment"`
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
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	Assignee    string     `json:"assignee"`
	Labels      []string   `json:"labels"`
	SprintID    string     `json:"sprint_id"`
	SLADueAt    *time.Time `json:"sla_due_at,omitempty"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
	Mu  sync.RWMutex
	db  *sql.DB
	dir string
}

// helpers
func marshalStrings(s []string) string {
	if s == nil {
		return "[]"
	}
	b, _ := json.Marshal(s)
	return string(b)
}

func unmarshalStrings(s string) []string {
	var out []string
	json.Unmarshal([]byte(s), &out)
	return out
}

func timeStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func timePtr(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, _ := time.Parse(time.RFC3339, s)
	return &t
}

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func NewStore(dir string) (*Store, error) {
	dbPath := filepath.Join(dir, "mitran.db")
	d, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	if err := runMigrations(d); err != nil {
		d.Close()
		return nil, err
	}
	return &Store{db: d, dir: dir}, nil
}

func (s *Store) path() string {
	return filepath.Join(s.dir, "mitran.db")
}

func (s *Store) Save() error {
	return nil // no-op, writes are immediate
}

// --- Projects ---

func (s *Store) AddProject(p Project) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	_, err := s.db.Exec(`INSERT OR REPLACE INTO projects (id,name,description,languages,team_size,created_at) VALUES (?,?,?,?,?,?)`,
		p.ID, p.Name, p.Description, marshalStrings(p.Languages), p.TeamSize, timeStr(p.CreatedAt))
	return err
}

func (s *Store) GetProject(id string) *Project {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	row := s.db.QueryRow(`SELECT id,name,description,languages,team_size,created_at FROM projects WHERE id=?`, id)
	var p Project
	var langs, ca string
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &langs, &p.TeamSize, &ca); err != nil {
		return nil
	}
	p.Languages = unmarshalStrings(langs)
	p.CreatedAt = parseTime(ca)
	return &p
}

// --- Tasks ---

func (s *Store) AddTask(t Task) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	return s.insertTask(t)
}

func (s *Store) AddTasks(tasks []Task) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	for _, t := range tasks {
		if _, err := tx.Exec(`INSERT OR REPLACE INTO tasks (id,project_id,title,description,agent_type,priority,status,depends_on,resources,created_by,created_at,started_at,completed_at,assigned_to) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			t.ID, t.ProjectID, t.Title, t.Description, t.AgentType, t.Priority, string(t.Status),
			marshalStrings(t.DependsOn), marshalStrings(t.Resources), t.CreatedBy,
			timeStr(t.CreatedAt), ptrTimeStr(t.StartedAt), ptrTimeStr(t.CompletedAt), t.AssignedTo); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) insertTask(t Task) error {
	_, err := s.db.Exec(`INSERT OR REPLACE INTO tasks (id,project_id,title,description,agent_type,priority,status,depends_on,resources,created_by,created_at,started_at,completed_at,assigned_to) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		t.ID, t.ProjectID, t.Title, t.Description, t.AgentType, t.Priority, string(t.Status),
		marshalStrings(t.DependsOn), marshalStrings(t.Resources), t.CreatedBy,
		timeStr(t.CreatedAt), ptrTimeStr(t.StartedAt), ptrTimeStr(t.CompletedAt), t.AssignedTo)
	return err
}

func ptrTimeStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func (s *Store) GetTask(id string) *Task {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.scanTask(s.db.QueryRow(`SELECT id,project_id,title,description,agent_type,priority,status,depends_on,resources,created_by,created_at,started_at,completed_at,assigned_to FROM tasks WHERE id=?`, id))
}

func (s *Store) scanTask(row *sql.Row) *Task {
	var t Task
	var status, deps, res, ca, sa, coa string
	if err := row.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.AgentType, &t.Priority, &status, &deps, &res, &t.CreatedBy, &ca, &sa, &coa, &t.AssignedTo); err != nil {
		return nil
	}
	t.Status = Status(status)
	t.DependsOn = unmarshalStrings(deps)
	t.Resources = unmarshalStrings(res)
	t.CreatedAt = parseTime(ca)
	t.StartedAt = timePtr(sa)
	t.CompletedAt = timePtr(coa)
	return &t
}

func (s *Store) UpdateTask(id string, fn func(*Task)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	row := s.db.QueryRow(`SELECT id,project_id,title,description,agent_type,priority,status,depends_on,resources,created_by,created_at,started_at,completed_at,assigned_to FROM tasks WHERE id=?`, id)
	var t Task
	var status, deps, res, ca, sa, coa string
	if err := row.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.AgentType, &t.Priority, &status, &deps, &res, &t.CreatedBy, &ca, &sa, &coa, &t.AssignedTo); err != nil {
		return fmt.Errorf("task %s not found", id)
	}
	t.Status = Status(status)
	t.DependsOn = unmarshalStrings(deps)
	t.Resources = unmarshalStrings(res)
	t.CreatedAt = parseTime(ca)
	t.StartedAt = timePtr(sa)
	t.CompletedAt = timePtr(coa)
	fn(&t)
	return s.insertTask(t)
}

func (s *Store) ListTasks(status Status, projectID string) []Task {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	q := `SELECT id,project_id,title,description,agent_type,priority,status,depends_on,resources,created_by,created_at,started_at,completed_at,assigned_to FROM tasks WHERE 1=1`
	var args []interface{}
	if status != "" {
		q += ` AND status=?`
		args = append(args, string(status))
	}
	if projectID != "" {
		q += ` AND project_id=?`
		args = append(args, projectID)
	}
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Task
	for rows.Next() {
		var t Task
		var st, deps, res, ca, sa, coa string
		rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.AgentType, &t.Priority, &st, &deps, &res, &t.CreatedBy, &ca, &sa, &coa, &t.AssignedTo)
		t.Status = Status(st)
		t.DependsOn = unmarshalStrings(deps)
		t.Resources = unmarshalStrings(res)
		t.CreatedAt = parseTime(ca)
		t.StartedAt = timePtr(sa)
		t.CompletedAt = timePtr(coa)
		result = append(result, t)
	}
	return result
}

// --- Checkpoints ---

func (s *Store) AddCheckpoint(c Checkpoint) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	_, err := s.db.Exec(`INSERT OR REPLACE INTO checkpoints (id,task_id,agent_id,output,diffs,questions,decision_action,decision_feedback,decision_by,decision_at,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		c.ID, c.TaskID, c.AgentID, c.Output, marshalStrings(c.Diffs), marshalStrings(c.Questions),
		c.DecisionAction, c.DecisionFeedback, c.DecisionBy, ptrTimeStr(c.DecisionAt), timeStr(c.CreatedAt))
	return err
}

func (s *Store) GetCheckpoint(id string) *Checkpoint {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	row := s.db.QueryRow(`SELECT id,task_id,agent_id,output,diffs,questions,decision_action,decision_feedback,decision_by,decision_at,created_at FROM checkpoints WHERE id=?`, id)
	var c Checkpoint
	var diffs, qs, da, ca string
	if err := row.Scan(&c.ID, &c.TaskID, &c.AgentID, &c.Output, &diffs, &qs, &c.DecisionAction, &c.DecisionFeedback, &c.DecisionBy, &da, &ca); err != nil {
		return nil
	}
	c.Diffs = unmarshalStrings(diffs)
	c.Questions = unmarshalStrings(qs)
	c.DecisionAt = timePtr(da)
	c.CreatedAt = parseTime(ca)
	return &c
}

func (s *Store) UpdateCheckpoint(id string, fn func(*Checkpoint)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	c := s.getCheckpointLocked(id)
	if c == nil {
		return fmt.Errorf("checkpoint %s not found", id)
	}
	fn(c)
	_, err := s.db.Exec(`UPDATE checkpoints SET decision_action=?,decision_feedback=?,decision_by=?,decision_at=? WHERE id=?`,
		c.DecisionAction, c.DecisionFeedback, c.DecisionBy, ptrTimeStr(c.DecisionAt), id)
	return err
}

func (s *Store) getCheckpointLocked(id string) *Checkpoint {
	row := s.db.QueryRow(`SELECT id,task_id,agent_id,output,diffs,questions,decision_action,decision_feedback,decision_by,decision_at,created_at FROM checkpoints WHERE id=?`, id)
	var c Checkpoint
	var diffs, qs, da, ca string
	if err := row.Scan(&c.ID, &c.TaskID, &c.AgentID, &c.Output, &diffs, &qs, &c.DecisionAction, &c.DecisionFeedback, &c.DecisionBy, &da, &ca); err != nil {
		return nil
	}
	c.Diffs = unmarshalStrings(diffs)
	c.Questions = unmarshalStrings(qs)
	c.DecisionAt = timePtr(da)
	c.CreatedAt = parseTime(ca)
	return &c
}

func (s *Store) ListPendingCheckpoints() []Checkpoint {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	rows, err := s.db.Query(`SELECT id,task_id,agent_id,output,diffs,questions,decision_action,decision_feedback,decision_by,decision_at,created_at FROM checkpoints WHERE decision_action=''`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Checkpoint
	for rows.Next() {
		var c Checkpoint
		var diffs, qs, da, ca string
		rows.Scan(&c.ID, &c.TaskID, &c.AgentID, &c.Output, &diffs, &qs, &c.DecisionAction, &c.DecisionFeedback, &c.DecisionBy, &da, &ca)
		c.Diffs = unmarshalStrings(diffs)
		c.Questions = unmarshalStrings(qs)
		c.DecisionAt = timePtr(da)
		c.CreatedAt = parseTime(ca)
		result = append(result, c)
	}
	return result
}

// --- Agents ---

func (s *Store) AddAgent(a Agent) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	_, err := s.db.Exec(`INSERT OR REPLACE INTO agents (id,name,type,status,callback_url,last_heartbeat,created_at) VALUES (?,?,?,?,?,?,?)`,
		a.ID, a.Name, a.Type, a.Status, a.CallbackURL, timeStr(a.LastHeartbeat), timeStr(a.CreatedAt))
	return err
}

func (s *Store) GetAgent(id string) *Agent {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	row := s.db.QueryRow(`SELECT id,name,type,status,callback_url,last_heartbeat,created_at FROM agents WHERE id=?`, id)
	var a Agent
	var lh, ca string
	if err := row.Scan(&a.ID, &a.Name, &a.Type, &a.Status, &a.CallbackURL, &lh, &ca); err != nil {
		return nil
	}
	a.LastHeartbeat = parseTime(lh)
	a.CreatedAt = parseTime(ca)
	return &a
}

func (s *Store) GetAgentByType(agentType string) *Agent {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	row := s.db.QueryRow(`SELECT id,name,type,status,callback_url,last_heartbeat,created_at FROM agents WHERE type=? AND status='online' LIMIT 1`, agentType)
	var a Agent
	var lh, ca string
	if err := row.Scan(&a.ID, &a.Name, &a.Type, &a.Status, &a.CallbackURL, &lh, &ca); err != nil {
		return nil
	}
	a.LastHeartbeat = parseTime(lh)
	a.CreatedAt = parseTime(ca)
	return &a
}

func (s *Store) UpdateAgentHeartbeat(id string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	res, err := s.db.Exec(`UPDATE agents SET last_heartbeat=?, status='online' WHERE id=?`, timeStr(time.Now()), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("agent %s not found", id)
	}
	return nil
}

// --- Environments ---

func (s *Store) SetEnvironments(projectID string, envs []Environment) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.db.Exec(`DELETE FROM environments WHERE project_id=?`, projectID)
	for _, e := range envs {
		if _, err := s.db.Exec(`INSERT INTO environments (project_id,name,auto_deploy,approval_required,rollback_enabled,created_at) VALUES (?,?,?,?,?,?)`,
			projectID, e.Name, boolToInt(e.AutoDeploy), boolToInt(e.ApprovalRequired), boolToInt(e.RollbackEnabled), timeStr(e.CreatedAt)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) GetEnvironments(projectID string) []Environment {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	rows, err := s.db.Query(`SELECT name,auto_deploy,approval_required,rollback_enabled,created_at FROM environments WHERE project_id=?`, projectID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Environment
	for rows.Next() {
		var e Environment
		var ad, ar, re int
		var ca string
		rows.Scan(&e.Name, &ad, &ar, &re, &ca)
		e.AutoDeploy = ad == 1
		e.ApprovalRequired = ar == 1
		e.RollbackEnabled = re == 1
		e.CreatedAt = parseTime(ca)
		result = append(result, e)
	}
	return result
}

func (s *Store) UpdateEnvironment(projectID, envName string, fn func(*Environment)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	row := s.db.QueryRow(`SELECT name,auto_deploy,approval_required,rollback_enabled,created_at FROM environments WHERE project_id=? AND name=?`, projectID, envName)
	var e Environment
	var ad, ar, re int
	var ca string
	if err := row.Scan(&e.Name, &ad, &ar, &re, &ca); err != nil {
		return fmt.Errorf("environment %s not found", envName)
	}
	e.AutoDeploy = ad == 1
	e.ApprovalRequired = ar == 1
	e.RollbackEnabled = re == 1
	e.CreatedAt = parseTime(ca)
	fn(&e)
	_, err := s.db.Exec(`UPDATE environments SET auto_deploy=?,approval_required=?,rollback_enabled=? WHERE project_id=? AND name=?`,
		boolToInt(e.AutoDeploy), boolToInt(e.ApprovalRequired), boolToInt(e.RollbackEnabled), projectID, envName)
	return err
}

// --- Pipelines ---

func (s *Store) AddPipeline(p Pipeline) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	stages, _ := json.Marshal(p.Stages)
	_, err := s.db.Exec(`INSERT OR REPLACE INTO pipelines (id,project_id,name,stages,created_at) VALUES (?,?,?,?,?)`,
		p.ID, p.ProjectID, p.Name, string(stages), timeStr(p.CreatedAt))
	return err
}

func (s *Store) ListPipelines(projectID string) []Pipeline {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	rows, err := s.db.Query(`SELECT id,project_id,name,stages,created_at FROM pipelines WHERE project_id=?`, projectID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Pipeline
	for rows.Next() {
		var p Pipeline
		var stages, ca string
		rows.Scan(&p.ID, &p.ProjectID, &p.Name, &stages, &ca)
		json.Unmarshal([]byte(stages), &p.Stages)
		p.CreatedAt = parseTime(ca)
		result = append(result, p)
	}
	return result
}

func (s *Store) AddPipelineRun(run PipelineRun) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	_, err := s.db.Exec(`INSERT OR REPLACE INTO pipeline_runs (id,pipeline_id,project_id,branch,commit_hash,status,started_at,finished_at) VALUES (?,?,?,?,?,?,?,?)`,
		run.ID, run.PipelineID, run.ProjectID, run.Branch, run.Commit, run.Status, timeStr(run.StartedAt), ptrTimeStr(run.FinishedAt))
	return err
}

func (s *Store) ListPipelineRuns(pipelineID string) []PipelineRun {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	rows, err := s.db.Query(`SELECT id,pipeline_id,project_id,branch,commit_hash,status,started_at,finished_at FROM pipeline_runs WHERE pipeline_id=?`, pipelineID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []PipelineRun
	for rows.Next() {
		var r PipelineRun
		var sa, fa string
		rows.Scan(&r.ID, &r.PipelineID, &r.ProjectID, &r.Branch, &r.Commit, &r.Status, &sa, &fa)
		r.StartedAt = parseTime(sa)
		r.FinishedAt = timePtr(fa)
		result = append(result, r)
	}
	return result
}

// --- Deployments ---

func (s *Store) AddDeployment(d Deployment) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	_, err := s.db.Exec(`INSERT OR REPLACE INTO deployments (id,project_id,environment,version,commit_hash,status,triggered_at,finished_at) VALUES (?,?,?,?,?,?,?,?)`,
		d.ID, d.ProjectID, d.Environment, d.Version, d.Commit, d.Status, timeStr(d.TriggeredAt), ptrTimeStr(d.FinishedAt))
	return err
}

func (s *Store) ListDeployments(projectID string) []Deployment {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	rows, err := s.db.Query(`SELECT id,project_id,environment,version,commit_hash,status,triggered_at,finished_at FROM deployments WHERE project_id=?`, projectID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Deployment
	for rows.Next() {
		var d Deployment
		var ta, fa string
		rows.Scan(&d.ID, &d.ProjectID, &d.Environment, &d.Version, &d.Commit, &d.Status, &ta, &fa)
		d.TriggeredAt = parseTime(ta)
		d.FinishedAt = timePtr(fa)
		result = append(result, d)
	}
	return result
}

// --- Tickets ---

func (s *Store) AddTicket(t Ticket) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	_, err := s.db.Exec(`INSERT OR REPLACE INTO tickets (id,title,description,status,priority,assignee,labels,sprint_id,sla_due_at,created_by,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`,
		t.ID, t.Title, t.Description, t.Status, t.Priority, t.Assignee, marshalStrings(t.Labels), t.SprintID, ptrTimeStr(t.SLADueAt), t.CreatedBy, timeStr(t.CreatedAt), timeStr(t.UpdatedAt))
	return err
}

func (s *Store) GetTicket(id string) *Ticket {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	row := s.db.QueryRow(`SELECT id,title,description,status,priority,assignee,labels,sprint_id,sla_due_at,created_by,created_at,updated_at FROM tickets WHERE id=?`, id)
	var t Ticket
	var labels, sla, ca, ua string
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Assignee, &labels, &t.SprintID, &sla, &t.CreatedBy, &ca, &ua); err != nil {
		return nil
	}
	t.Labels = unmarshalStrings(labels)
	t.SLADueAt = timePtr(sla)
	t.CreatedAt = parseTime(ca)
	t.UpdatedAt = parseTime(ua)
	return &t
}

func (s *Store) UpdateTicket(id string, fn func(*Ticket)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	t := s.getTicketLocked(id)
	if t == nil {
		return fmt.Errorf("ticket %s not found", id)
	}
	fn(t)
	_, err := s.db.Exec(`UPDATE tickets SET title=?,description=?,status=?,priority=?,assignee=?,labels=?,sprint_id=?,sla_due_at=?,updated_at=? WHERE id=?`,
		t.Title, t.Description, t.Status, t.Priority, t.Assignee, marshalStrings(t.Labels), t.SprintID, ptrTimeStr(t.SLADueAt), timeStr(t.UpdatedAt), id)
	return err
}

func (s *Store) getTicketLocked(id string) *Ticket {
	row := s.db.QueryRow(`SELECT id,title,description,status,priority,assignee,labels,sprint_id,sla_due_at,created_by,created_at,updated_at FROM tickets WHERE id=?`, id)
	var t Ticket
	var labels, sla, ca, ua string
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Assignee, &labels, &t.SprintID, &sla, &t.CreatedBy, &ca, &ua); err != nil {
		return nil
	}
	t.Labels = unmarshalStrings(labels)
	t.SLADueAt = timePtr(sla)
	t.CreatedAt = parseTime(ca)
	t.UpdatedAt = parseTime(ua)
	return &t
}

func (s *Store) DeleteTicket(id string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	res, err := s.db.Exec(`DELETE FROM tickets WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("ticket %s not found", id)
	}
	return nil
}

func (s *Store) ListTickets(status, assignee, sprintID, labels string) []Ticket {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	q := `SELECT id,title,description,status,priority,assignee,labels,sprint_id,sla_due_at,created_by,created_at,updated_at FROM tickets WHERE 1=1`
	var args []interface{}
	if status != "" {
		q += ` AND status=?`
		args = append(args, status)
	}
	if assignee != "" {
		q += ` AND assignee=?`
		args = append(args, assignee)
	}
	if sprintID != "" {
		q += ` AND sprint_id=?`
		args = append(args, sprintID)
	}
	if labels != "" {
		q += ` AND labels LIKE ?`
		args = append(args, "%"+labels+"%")
	}
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Ticket
	for rows.Next() {
		var t Ticket
		var lb, sla, ca, ua string
		rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Assignee, &lb, &t.SprintID, &sla, &t.CreatedBy, &ca, &ua)
		t.Labels = unmarshalStrings(lb)
		t.SLADueAt = timePtr(sla)
		t.CreatedAt = parseTime(ca)
		t.UpdatedAt = parseTime(ua)
		result = append(result, t)
	}
	return result
}

func (s *Store) AddTicketComment(c TicketComment) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	_, err := s.db.Exec(`INSERT INTO ticket_comments (id,ticket_id,author,body,created_at) VALUES (?,?,?,?,?)`,
		c.ID, c.TicketID, c.Author, c.Body, timeStr(c.CreatedAt))
	return err
}

func (s *Store) ListTicketComments(ticketID string) []TicketComment {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	rows, err := s.db.Query(`SELECT id,ticket_id,author,body,created_at FROM ticket_comments WHERE ticket_id=?`, ticketID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []TicketComment
	for rows.Next() {
		var c TicketComment
		var ca string
		rows.Scan(&c.ID, &c.TicketID, &c.Author, &c.Body, &ca)
		c.CreatedAt = parseTime(ca)
		result = append(result, c)
	}
	return result
}

// --- Sprints ---

func (s *Store) AddSprint(sp Sprint) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	_, err := s.db.Exec(`INSERT OR REPLACE INTO sprints (id,name,start_date,end_date,created_at) VALUES (?,?,?,?,?)`,
		sp.ID, sp.Name, sp.StartDate, sp.EndDate, timeStr(sp.CreatedAt))
	return err
}

func (s *Store) ListSprints() []Sprint {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	rows, err := s.db.Query(`SELECT id,name,start_date,end_date,created_at FROM sprints`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Sprint
	for rows.Next() {
		var sp Sprint
		var ca string
		rows.Scan(&sp.ID, &sp.Name, &sp.StartDate, &sp.EndDate, &ca)
		sp.CreatedAt = parseTime(ca)
		result = append(result, sp)
	}
	return result
}

func (s *Store) UpdateSprint(id string, fn func(*Sprint)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	row := s.db.QueryRow(`SELECT id,name,start_date,end_date,created_at FROM sprints WHERE id=?`, id)
	var sp Sprint
	var ca string
	if err := row.Scan(&sp.ID, &sp.Name, &sp.StartDate, &sp.EndDate, &ca); err != nil {
		return fmt.Errorf("sprint %s not found", id)
	}
	sp.CreatedAt = parseTime(ca)
	fn(&sp)
	_, err := s.db.Exec(`UPDATE sprints SET name=?,start_date=?,end_date=? WHERE id=?`, sp.Name, sp.StartDate, sp.EndDate, id)
	return err
}

// --- Wiki ---

func (s *Store) AddWikiPage(p WikiPage) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	_, err := s.db.Exec(`INSERT OR REPLACE INTO wiki_pages (id,title,content,parent_id,slug,created_by,tags,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?)`,
		p.ID, p.Title, p.Content, p.ParentID, p.Slug, p.CreatedBy, marshalStrings(p.Tags), timeStr(p.CreatedAt), timeStr(p.UpdatedAt))
	return err
}

func (s *Store) GetWikiPage(id string) *WikiPage {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	row := s.db.QueryRow(`SELECT id,title,content,parent_id,slug,created_by,tags,created_at,updated_at FROM wiki_pages WHERE id=?`, id)
	var p WikiPage
	var tags, ca, ua string
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.ParentID, &p.Slug, &p.CreatedBy, &tags, &ca, &ua); err != nil {
		return nil
	}
	p.Tags = unmarshalStrings(tags)
	p.CreatedAt = parseTime(ca)
	p.UpdatedAt = parseTime(ua)
	return &p
}

func (s *Store) UpdateWikiPage(id string, fn func(*WikiPage)) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	p := s.getWikiPageLocked(id)
	if p == nil {
		return fmt.Errorf("wiki page %s not found", id)
	}
	fn(p)
	_, err := s.db.Exec(`UPDATE wiki_pages SET title=?,content=?,parent_id=?,slug=?,tags=?,updated_at=? WHERE id=?`,
		p.Title, p.Content, p.ParentID, p.Slug, marshalStrings(p.Tags), timeStr(p.UpdatedAt), id)
	return err
}

func (s *Store) getWikiPageLocked(id string) *WikiPage {
	row := s.db.QueryRow(`SELECT id,title,content,parent_id,slug,created_by,tags,created_at,updated_at FROM wiki_pages WHERE id=?`, id)
	var p WikiPage
	var tags, ca, ua string
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.ParentID, &p.Slug, &p.CreatedBy, &tags, &ca, &ua); err != nil {
		return nil
	}
	p.Tags = unmarshalStrings(tags)
	p.CreatedAt = parseTime(ca)
	p.UpdatedAt = parseTime(ua)
	return &p
}

func (s *Store) DeleteWikiPage(id string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	res, err := s.db.Exec(`DELETE FROM wiki_pages WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("wiki page %s not found", id)
	}
	return nil
}

func (s *Store) ListWikiPages() []WikiPage {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	rows, err := s.db.Query(`SELECT id,title,content,parent_id,slug,created_by,tags,created_at,updated_at FROM wiki_pages`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []WikiPage
	for rows.Next() {
		var p WikiPage
		var tags, ca, ua string
		rows.Scan(&p.ID, &p.Title, &p.Content, &p.ParentID, &p.Slug, &p.CreatedBy, &tags, &ca, &ua)
		p.Tags = unmarshalStrings(tags)
		p.CreatedAt = parseTime(ca)
		p.UpdatedAt = parseTime(ua)
		result = append(result, p)
	}
	return result
}

func (s *Store) SearchWikiPages(query string) []WikiPage {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	rows, err := s.db.Query(`SELECT id,title,content,parent_id,slug,created_by,tags,created_at,updated_at FROM wiki_pages WHERE title LIKE ? OR content LIKE ?`,
		"%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []WikiPage
	for rows.Next() {
		var p WikiPage
		var tags, ca, ua string
		rows.Scan(&p.ID, &p.Title, &p.Content, &p.ParentID, &p.Slug, &p.CreatedBy, &tags, &ca, &ua)
		p.Tags = unmarshalStrings(tags)
		p.CreatedAt = parseTime(ca)
		p.UpdatedAt = parseTime(ua)
		result = append(result, p)
	}
	return result
}

func containsLower(s, sub string) bool {
	return len(s) >= len(sub) && strings.Contains(strings.ToLower(s), sub)
}

func (s *Store) ListAgents() []Agent {
	rows, err := s.db.Query("SELECT id, name, type, status, callback_url, last_heartbeat, created_at FROM agents")
	if err != nil { return nil }
	defer rows.Close()
	var agents []Agent
	for rows.Next() {
		var a Agent
		rows.Scan(&a.ID, &a.Name, &a.Type, &a.Status, &a.CallbackURL, &a.LastHeartbeat, &a.CreatedAt)
		agents = append(agents, a)
	}
	return agents
}

func (s *Store) ListProjects() []Project {
	rows, err := s.db.Query("SELECT id, name, description, languages, team_size, created_at FROM projects")
	if err != nil { return nil }
	defer rows.Close()
	var projects []Project
	for rows.Next() {
		var p Project
		var langs string
		rows.Scan(&p.ID, &p.Name, &p.Description, &langs, &p.TeamSize, &p.CreatedAt)
		projects = append(projects, p)
	}
	return projects
}
