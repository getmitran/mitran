package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/getmitran/mitran/server/db"
)

type EnvironmentHandler struct {
	Store *db.Store
}

func (h *EnvironmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Routes: /api/v1/projects/{id}/environments[/{name}[/deploy]]
	// Also:   /api/v1/projects/{id}/deployments
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}

	projectID := parts[0]
	resource := parts[1]

	if resource == "deployments" && r.Method == "GET" {
		h.listDeployments(w, r, projectID)
		return
	}

	if resource != "environments" {
		http.NotFound(w, r)
		return
	}

	switch {
	case len(parts) == 2 && r.Method == "POST":
		h.setEnvironments(w, r, projectID)
	case len(parts) == 2 && r.Method == "GET":
		h.getEnvironments(w, r, projectID)
	case len(parts) == 3 && r.Method == "PUT":
		h.updateEnvironment(w, r, projectID, parts[2])
	case len(parts) == 4 && parts[3] == "deploy" && r.Method == "POST":
		h.triggerDeploy(w, r, projectID, parts[2])
	default:
		http.NotFound(w, r)
	}
}

func (h *EnvironmentHandler) setEnvironments(w http.ResponseWriter, r *http.Request, projectID string) {
	p := h.Store.GetProject(projectID)
	if p == nil {
		writeErr(w, http.StatusNotFound, "project not found")
		return
	}
	var envs []db.Environment
	if err := json.NewDecoder(r.Body).Decode(&envs); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	for i := range envs {
		if envs[i].CreatedAt.IsZero() {
			envs[i].CreatedAt = time.Now()
		}
	}
	if err := h.Store.SetEnvironments(projectID, envs); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, envs)
}

func (h *EnvironmentHandler) getEnvironments(w http.ResponseWriter, r *http.Request, projectID string) {
	envs := h.Store.GetEnvironments(projectID)
	writeJSON(w, http.StatusOK, envs)
}

func (h *EnvironmentHandler) updateEnvironment(w http.ResponseWriter, r *http.Request, projectID, envName string) {
	var update struct {
		AutoDeploy       *bool `json:"auto_deploy"`
		ApprovalRequired *bool `json:"approval_required"`
		RollbackEnabled  *bool `json:"rollback_enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.Store.UpdateEnvironment(projectID, envName, func(e *db.Environment) {
		if update.AutoDeploy != nil {
			e.AutoDeploy = *update.AutoDeploy
		}
		if update.ApprovalRequired != nil {
			e.ApprovalRequired = *update.ApprovalRequired
		}
		if update.RollbackEnabled != nil {
			e.RollbackEnabled = *update.RollbackEnabled
		}
	})
	if err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *EnvironmentHandler) triggerDeploy(w http.ResponseWriter, r *http.Request, projectID, envName string) {
	var req struct {
		Version string `json:"version"`
		Commit  string `json:"commit"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	deploy := db.Deployment{
		ID:          genID(),
		ProjectID:   projectID,
		Environment: envName,
		Version:     req.Version,
		Commit:      req.Commit,
		Status:      "pending",
		TriggeredAt: time.Now(),
	}
	if err := h.Store.AddDeployment(deploy); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, deploy)
}

func (h *EnvironmentHandler) listDeployments(w http.ResponseWriter, r *http.Request, projectID string) {
	deployments := h.Store.ListDeployments(projectID)
	writeJSON(w, http.StatusOK, deployments)
}
