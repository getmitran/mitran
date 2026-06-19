package approval

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/getmitran/mitran/server/apierr"
)

type ApprovalMode string

const (
	ModeInteractive ApprovalMode = "interactive"
	ModeReads       ApprovalMode = "reads"
	ModeYolo        ApprovalMode = "yolo"
)

type RequestStatus string

const (
	StatusPending  RequestStatus = "pending"
	StatusApproved RequestStatus = "approved"
	StatusRejected RequestStatus = "rejected"
)

type ToolRequest struct {
	ID        string        `json:"id"`
	Agent     string        `json:"agent"`
	Tool      string        `json:"tool"`
	Args      any           `json:"args"`
	ReadOnly  bool          `json:"read_only"`
	Status    RequestStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}

type ToolApprovalService struct {
	mu       sync.RWMutex
	modes    map[string]ApprovalMode
	requests map[string]*ToolRequest
}

func NewToolApprovalService() *ToolApprovalService {
	return &ToolApprovalService{
		modes:    make(map[string]ApprovalMode),
		requests: make(map[string]*ToolRequest),
	}
}

func (s *ToolApprovalService) SetAgentMode(agent string, mode ApprovalMode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.modes[agent] = mode
}

func (s *ToolApprovalService) getMode(agent string) ApprovalMode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if m, ok := s.modes[agent]; ok {
		return m
	}
	return ModeInteractive
}

func (s *ToolApprovalService) Submit(agent, tool string, args any, readOnly bool) *ToolRequest {
	mode := s.getMode(agent)

	status := StatusPending
	if mode == ModeYolo || (mode == ModeReads && readOnly) {
		status = StatusApproved
	}

	req := &ToolRequest{
		ID:        uuid.New().String(),
		Agent:     agent,
		Tool:      tool,
		Args:      args,
		ReadOnly:  readOnly,
		Status:    status,
		CreatedAt: time.Now(),
	}

	if status == StatusPending {
		s.mu.Lock()
		s.requests[req.ID] = req
		s.mu.Unlock()
	}

	return req
}

func (s *ToolApprovalService) Approve(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	req, ok := s.requests[id]
	if !ok {
		return fmt.Errorf("request %s not found", id)
	}
	req.Status = StatusApproved
	delete(s.requests, id)
	return nil
}

func (s *ToolApprovalService) Reject(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	req, ok := s.requests[id]
	if !ok {
		return fmt.Errorf("request %s not found", id)
	}
	req.Status = StatusRejected
	delete(s.requests, id)
	return nil
}

func (s *ToolApprovalService) Pending() []*ToolRequest {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*ToolRequest, 0, len(s.requests))
	for _, r := range s.requests {
		result = append(result, r)
	}
	return result
}

func (s *ToolApprovalService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/tools/pending", s.handlePending)
	mux.HandleFunc("POST /api/v1/tools/{id}/approve", s.handleApprove)
	mux.HandleFunc("POST /api/v1/tools/{id}/reject", s.handleReject)
	mux.HandleFunc("PUT /api/v1/tools/config/{agent}", s.handleConfig)
}

func (s *ToolApprovalService) handlePending(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.Pending())
}

func (s *ToolApprovalService) handleApprove(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Approve(id); err != nil {
		apierr.WriteError(w, apierr.NotFound(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *ToolApprovalService) handleReject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Reject(id); err != nil {
		apierr.WriteError(w, apierr.NotFound(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *ToolApprovalService) handleConfig(w http.ResponseWriter, r *http.Request) {
	agent := r.PathValue("agent")
	var body struct {
		Mode ApprovalMode `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierr.WriteError(w, apierr.BadRequest("invalid body"))
		return
	}
	switch body.Mode {
	case ModeInteractive, ModeReads, ModeYolo:
		s.SetAgentMode(agent, body.Mode)
		w.WriteHeader(http.StatusOK)
	default:
		apierr.WriteError(w, apierr.BadRequest("invalid mode: use interactive, reads, or yolo"))
	}
}
