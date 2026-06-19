package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/getmitran/mitran/server/db"
)

type ProjectHandler struct {
	Store *db.Store
}

type createProjectReq struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Languages   []string `json:"languages"`
	TeamSize    int      `json:"team_size"`
}

func (h *ProjectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/projects")
	switch {
	case r.Method == "POST" && path == "":
		h.create(w, r)
	case r.Method == "GET" && path == "":
		h.list(w, r)
	case r.Method == "GET" && len(path) > 1:
		h.get(w, r, strings.TrimPrefix(path, "/"))
	default:
		http.NotFound(w, r)
	}
}

func (h *ProjectHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createProjectReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	p := db.Project{
		ID:          genID(),
		Name:        req.Name,
		Description: req.Description,
		Languages:   req.Languages,
		TeamSize:    req.TeamSize,
		CreatedAt:   time.Now(),
	}
	if err := h.Store.AddProject(p); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *ProjectHandler) list(w http.ResponseWriter, r *http.Request) {
	h.Store.Mu.RLock()
	defer h.Store.Mu.RUnlock()
	writeJSON(w, http.StatusOK, h.Store.ListProjects())
}

func (h *ProjectHandler) get(w http.ResponseWriter, r *http.Request, id string) {
	p := h.Store.GetProject(id)
	if p == nil {
		writeErr(w, http.StatusNotFound, "project not found")
		return
	}
	tasks := h.Store.ListTasks("", id)
	writeJSON(w, http.StatusOK, map[string]any{"project": p, "tasks": tasks})
}
