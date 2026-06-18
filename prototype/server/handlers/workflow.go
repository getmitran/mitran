package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"strings"
)

type WorkflowStep struct {
	Name      string   `json:"name"`
	AgentID   string   `json:"agent_id"`
	Action    string   `json:"action"`
	DependsOn []string `json:"depends_on"`
}

type Workflow struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	TenantID  string         `json:"tenant_id"`
	Steps     []WorkflowStep `json:"steps"`
	CreatedAt time.Time      `json:"created_at"`
}

type WorkflowHandler struct {
	mu        sync.RWMutex
	workflows map[string]*Workflow
}

func NewWorkflowHandler() *WorkflowHandler {
	return &WorkflowHandler{workflows: make(map[string]*Workflow)}
}

func (h *WorkflowHandler) Create(w http.ResponseWriter, r *http.Request) {
	var wf Workflow
	if err := json.NewDecoder(r.Body).Decode(&wf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	wf.ID = uuid.New().String()
	wf.CreatedAt = time.Now()
	h.mu.Lock()
	h.workflows[wf.ID] = &wf
	h.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wf)
}

func (h *WorkflowHandler) List(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	result := make([]*Workflow, 0, len(h.workflows))
	for _, wf := range h.workflows {
		result = append(result, wf)
	}
	h.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *WorkflowHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	h.mu.RLock()
	wf, ok := h.workflows[id]
	h.mu.RUnlock()
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wf)
}

func (h *WorkflowHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	h.mu.Lock()
	if _, ok := h.workflows[id]; !ok {
		h.mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	delete(h.workflows, id)
	h.mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}


func extractID(path string) string {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	return parts[len(parts)-1]
}
