package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/getmitran/mitran/server/db"
	"github.com/getmitran/mitran/server/queue"
)

type TaskHandler struct {
	Store *db.Store
	Queue *queue.TaskQueue
}

type createTaskReq struct {
	ProjectID   string   `json:"project_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	AgentType   string   `json:"agent_type"`
	Priority    int      `json:"priority"`
	DependsOn   []string `json:"depends_on"`
	Resources   []string `json:"resources"`
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks")
	switch {
	case r.Method == "GET" && path == "":
		h.list(w, r)
	case r.Method == "POST" && path == "":
		h.create(w, r)
	case r.Method == "PUT" && strings.HasSuffix(path, "/priority"):
		id := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/priority")
		h.updatePriority(w, r, id)
	case r.Method == "POST" && strings.HasSuffix(path, "/assign"):
		id := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/assign")
		h.assign(w, r, id)
	default:
		http.NotFound(w, r)
	}
}

func (h *TaskHandler) list(w http.ResponseWriter, r *http.Request) {
	status := db.Status(r.URL.Query().Get("status"))
	projectID := r.URL.Query().Get("project_id")
	tasks := h.Store.ListTasks(status, projectID)
	if tasks == nil {
		tasks = []db.Task{}
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	if h.Queue != nil && h.Queue.IsFull() {
		writeErr(w, http.StatusServiceUnavailable, "server busy, try again later")
		return
	}
	var req createTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	t := db.Task{
		ID:          genID(),
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
		AgentType:   req.AgentType,
		Priority:    req.Priority,
		Status:      db.StatusQueued,
		DependsOn:   req.DependsOn,
		Resources:   req.Resources,
		CreatedBy:   "human",
		CreatedAt:   time.Now(),
	}
	if err := h.Store.AddTask(t); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if h.Queue != nil {
		if err := h.Queue.Enqueue(queue.Task{ID: t.ID, Payload: t}); err != nil {
			writeErr(w, http.StatusServiceUnavailable, "server busy, try again later")
			return
		}
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *TaskHandler) updatePriority(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Priority int `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.Store.UpdateTask(id, func(t *db.Task) {
		t.Priority = req.Priority
	})
	if err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *TaskHandler) assign(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.Store.UpdateTask(id, func(t *db.Task) {
		t.AssignedTo = req.AgentID
		t.Status = db.StatusRunning
		now := time.Now()
		t.StartedAt = &now
	})
	if err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "assigned"})
}
