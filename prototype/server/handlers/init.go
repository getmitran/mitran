package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/getmitran/mitran/server/db"
)

type InitHandler struct {
	Store *db.Store
}

type initReq struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Languages   []string `json:"languages"`
	TeamSize    int      `json:"team_size"`
}

type taskPlan struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	AgentType string   `json:"agent_type"`
	DependsOn []string `json:"depends_on"`
	Priority  int      `json:"priority"`
}

func (h *InitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	var req initReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create project
	project := db.Project{
		ID:          genID(),
		Name:        req.Name,
		Description: req.Description,
		Languages:   req.Languages,
		TeamSize:    req.TeamSize,
		CreatedAt:   time.Now(),
	}
	if err := h.Store.AddProject(project); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Generate task plan with DAG dependencies
	plan := generatePlan(project)

	// Insert tasks
	var tasks []db.Task
	for _, p := range plan {
		tasks = append(tasks, db.Task{
			ID:          p.ID,
			ProjectID:   project.ID,
			Title:       p.Title,
			AgentType:   p.AgentType,
			Priority:    p.Priority,
			Status:      db.StatusQueued,
			DependsOn:   p.DependsOn,
			Resources:   []string{},
			CreatedBy:   "mitran-init",
			CreatedAt:   time.Now(),
		})
	}

	if err := h.Store.AddTasks(tasks); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"project": project,
		"plan":    plan,
		"tasks":   len(tasks),
	})
}

func generatePlan(p db.Project) []taskPlan {
	langs := "the project"
	if len(p.Languages) > 0 {
		langs = p.Languages[0]
	}

	// IDs for dependency linking
	ids := make([]string, 12)
	for i := range ids {
		ids[i] = genID()
	}

	return []taskPlan{
		{ID: ids[0], Title: fmt.Sprintf("Initialize %s repository structure", langs), AgentType: "dev", Priority: 100, DependsOn: nil},
		{ID: ids[1], Title: "Set up CI/CD pipeline", AgentType: "cicd", Priority: 90, DependsOn: []string{ids[0]}},
		{ID: ids[2], Title: "Create project wiki with architecture overview", AgentType: "wiki", Priority: 80, DependsOn: []string{ids[0]}},
		{ID: ids[3], Title: "Generate API documentation templates", AgentType: "docs", Priority: 70, DependsOn: []string{ids[0]}},
		{ID: ids[4], Title: "Set up monitoring and alerting dashboards", AgentType: "ops", Priority: 75, DependsOn: []string{ids[1]}},
		{ID: ids[5], Title: "Create ticket templates and workflow", AgentType: "tickets", Priority: 60, DependsOn: nil},
		{ID: ids[6], Title: "Configure code review policies", AgentType: "review", Priority: 65, DependsOn: []string{ids[1]}},
		{ID: ids[7], Title: fmt.Sprintf("Scaffold core %s service", langs), AgentType: "dev", Priority: 95, DependsOn: []string{ids[0], ids[1]}},
		{ID: ids[8], Title: "Set up team onboarding documentation", AgentType: "hr", Priority: 40, DependsOn: []string{ids[2]}},
		{ID: ids[9], Title: "Implement health check endpoints", AgentType: "dev", Priority: 85, DependsOn: []string{ids[7]}},
		{ID: ids[10], Title: "Create runbook for incident response", AgentType: "ops", Priority: 55, DependsOn: []string{ids[4], ids[7]}},
		{ID: ids[11], Title: "Generate test coverage baseline", AgentType: "review", Priority: 50, DependsOn: []string{ids[7], ids[9]}},
	}
}
