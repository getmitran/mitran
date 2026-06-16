package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/getmitran/mitran/server/db"
)

type PipelineHandler struct {
	Store *db.Store
}

func (h *PipelineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// /api/v1/projects/{id}/pipelines[/{pipeline_id}[/trigger|/runs]]
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[1] != "pipelines" {
		http.NotFound(w, r)
		return
	}

	projectID := parts[0]

	switch {
	case len(parts) == 2 && r.Method == "POST":
		h.create(w, r, projectID)
	case len(parts) == 2 && r.Method == "GET":
		h.list(w, r, projectID)
	case len(parts) == 4 && parts[3] == "trigger" && r.Method == "POST":
		h.trigger(w, r, projectID, parts[2])
	case len(parts) == 4 && parts[3] == "runs" && r.Method == "GET":
		h.runs(w, r, projectID, parts[2])
	default:
		http.NotFound(w, r)
	}
}

func (h *PipelineHandler) create(w http.ResponseWriter, r *http.Request, projectID string) {
	p := h.Store.GetProject(projectID)
	if p == nil {
		writeErr(w, http.StatusNotFound, "project not found")
		return
	}
	var req struct {
		Name   string           `json:"name"`
		Stages []db.PipelineStage `json:"stages"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	pipeline := db.Pipeline{
		ID:        genID(),
		ProjectID: projectID,
		Name:      req.Name,
		Stages:    req.Stages,
		CreatedAt: time.Now(),
	}
	if err := h.Store.AddPipeline(pipeline); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, pipeline)
}

func (h *PipelineHandler) list(w http.ResponseWriter, r *http.Request, projectID string) {
	pipelines := h.Store.ListPipelines(projectID)
	writeJSON(w, http.StatusOK, pipelines)
}

func (h *PipelineHandler) trigger(w http.ResponseWriter, r *http.Request, projectID, pipelineID string) {
	var req struct {
		Branch string `json:"branch"`
		Commit string `json:"commit"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	run := db.PipelineRun{
		ID:         genID(),
		PipelineID: pipelineID,
		ProjectID:  projectID,
		Branch:     req.Branch,
		Commit:     req.Commit,
		Status:     "running",
		StartedAt:  time.Now(),
	}
	if err := h.Store.AddPipelineRun(run); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, run)
}

func (h *PipelineHandler) runs(w http.ResponseWriter, r *http.Request, projectID, pipelineID string) {
	runs := h.Store.ListPipelineRuns(pipelineID)
	writeJSON(w, http.StatusOK, runs)
}
