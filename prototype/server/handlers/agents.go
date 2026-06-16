package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/getmitran/mitran/server/db"
)

type AgentHandler struct {
	Store *db.Store
}

func (h *AgentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/agents")
	switch {
	case r.Method == "GET" && path == "":
		h.list(w, r)
	case r.Method == "POST" && path == "/register":
		h.register(w, r)
	case r.Method == "POST" && path == "/task-complete":
		h.taskComplete(w, r)
	case r.Method == "POST" && strings.HasSuffix(path, "/heartbeat"):
		id := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/heartbeat")
		h.heartbeat(w, r, id)
	default:
		http.NotFound(w, r)
	}
}

func (h *AgentHandler) list(w http.ResponseWriter, r *http.Request) {
	h.Store.Mu.RLock()
	defer h.Store.Mu.RUnlock()
	agents := h.Store.Agents
	if agents == nil {
		agents = []db.Agent{}
	}
	writeJSON(w, http.StatusOK, agents)
}

func (h *AgentHandler) register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		CallbackURL string `json:"callback_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	a := db.Agent{
		ID:            genID(),
		Name:          req.Name,
		Type:          req.Type,
		Status:        "online",
		CallbackURL:   req.CallbackURL,
		LastHeartbeat: time.Now(),
		CreatedAt:     time.Now(),
	}
	if err := h.Store.AddAgent(a); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, a)
}

func (h *AgentHandler) heartbeat(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Store.UpdateAgentHeartbeat(id); err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AgentHandler) taskComplete(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TaskID    string   `json:"task_id"`
		AgentID   string   `json:"agent_id"`
		Output    string   `json:"output"`
		Diffs     []string `json:"diffs"`
		Questions []string `json:"questions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	cp := db.Checkpoint{
		ID:        genID(),
		TaskID:    req.TaskID,
		AgentID:   req.AgentID,
		Output:    req.Output,
		Diffs:     req.Diffs,
		Questions: req.Questions,
		CreatedAt: time.Now(),
	}
	if err := h.Store.AddCheckpoint(cp); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Move task to checkpoint status
	h.Store.UpdateTask(req.TaskID, func(t *db.Task) {
		t.Status = db.StatusCheckpoint
	})

	writeJSON(w, http.StatusCreated, cp)
}
