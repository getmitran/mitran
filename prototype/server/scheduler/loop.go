package scheduler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/getmitran/mitran/server/db"
)

type Scheduler struct {
	Store        *db.Store
	WorkspaceDir string // base dir for project workspaces
}

func (s *Scheduler) Start() {
	go func() {
		for {
			time.Sleep(2 * time.Second)
			s.tick()
		}
	}()
}

func (s *Scheduler) tick() {
	s.Store.Mu.Lock()
	defer s.Store.Mu.Unlock()

	// Find approved task IDs for dep checking
	approved := make(map[string]bool)
	for _, t := range s.Store.Tasks {
		if t.Status == db.StatusApproved {
			approved[t.ID] = true
		}
	}

	// Find running task resources for conflict checking
	runningResources := make(map[string]bool)
	for _, t := range s.Store.Tasks {
		if t.Status == db.StatusRunning {
			for _, r := range t.Resources {
				runningResources[r] = true
			}
		}
	}

	// Find ready tasks (queued, all deps approved, no resource conflict)
	var ready []int
	for i, t := range s.Store.Tasks {
		if t.Status != db.StatusQueued {
			continue
		}
		// Check deps
		depsOK := true
		for _, dep := range t.DependsOn {
			if !approved[dep] {
				depsOK = false
				break
			}
		}
		if !depsOK {
			continue
		}
		// Check resources
		conflict := false
		for _, r := range t.Resources {
			if runningResources[r] {
				conflict = true
				break
			}
		}
		if conflict {
			continue
		}
		ready = append(ready, i)
	}

	if len(ready) == 0 {
		return
	}

	// Sort by priority (highest first)
	sort.Slice(ready, func(a, b int) bool {
		return s.Store.Tasks[ready[a]].Priority > s.Store.Tasks[ready[b]].Priority
	})

	// Pick highest priority task
	idx := ready[0]
	task := &s.Store.Tasks[idx]

	// Find agent for this task type
	var agent *db.Agent
	for i := range s.Store.Agents {
		if s.Store.Agents[i].Type == task.AgentType && s.Store.Agents[i].Status == "online" {
			agent = &s.Store.Agents[i]
			break
		}
	}

	if agent == nil {
		return // No agent available
	}

	// Mark as running
	now := time.Now()
	task.Status = db.StatusRunning
	task.AssignedTo = agent.ID
	task.StartedAt = &now
	s.Store.Save()

	// Notify agent via callback
	go s.notifyAgent(agent.CallbackURL, task)
}

func (s *Scheduler) notifyAgent(url string, task *db.Task) {
	if url == "" {
		return
	}

	// Find project for context
	var project *db.Project
	for i := range s.Store.Projects {
		if s.Store.Projects[i].ID == task.ProjectID {
			project = &s.Store.Projects[i]
			break
		}
	}

	workspacePath := ""
	if s.WorkspaceDir != "" && task.ProjectID != "" {
		workspacePath = s.WorkspaceDir + "/" + task.ProjectID
	}

	context := map[string]any{}
	if project != nil {
		context["project_name"] = project.Name
		context["description"] = project.Description
		context["languages"] = project.Languages
		context["team_size"] = project.TeamSize
	}

	payload, _ := json.Marshal(map[string]any{
		"task_id":        task.ID,
		"task":           task.Title,
		"agent":          task.AgentType,
		"project_id":     task.ProjectID,
		"workspace_path": workspacePath,
		"context":        context,
	})
	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Printf("scheduler: failed to notify agent at %s: %v", url, err)
		return
	}
	resp.Body.Close()
}
