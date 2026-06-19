package onboarding

import (
	"encoding/json"
	"net/http"
	"sync"
	"github.com/getmitran/mitran/server/apierr"
)

type OnboardingState struct {
	TenantID       string            `json:"tenant_id"`
	Step           int               `json:"step"`
	CompletedSteps []string          `json:"completed_steps"`
	Config         map[string]string `json:"config"`
}

var (
	store   = make(map[string]*OnboardingState)
	storeMu sync.RWMutex
)

type OnboardingHandler struct{}

func (h *OnboardingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/onboarding/start", h.Start)
	mux.HandleFunc("/api/v1/onboarding/step", h.Step)
	mux.HandleFunc("/api/v1/onboarding/status", h.Status)
}

func (h *OnboardingHandler) Start(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}
	var req struct {
		TenantID string `json:"tenant_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TenantID == "" {
		apierr.WriteError(w, apierr.BadRequest("tenant_id required"))
		return
	}
	state := &OnboardingState{
		TenantID:       req.TenantID,
		Step:           0,
		CompletedSteps: []string{},
		Config:         make(map[string]string),
	}
	storeMu.Lock()
	store[req.TenantID] = state
	storeMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (h *OnboardingHandler) Step(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}
	var req struct {
		TenantID  string            `json:"tenant_id"`
		StepName  string            `json:"step_name"`
		Config    map[string]string `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TenantID == "" {
		apierr.WriteError(w, apierr.BadRequest("tenant_id required"))
		return
	}
	storeMu.Lock()
	state, ok := store[req.TenantID]
	if !ok {
		storeMu.Unlock()
		apierr.WriteError(w, apierr.NotFound("onboarding not started"))
		return
	}
	state.Step++
	if req.StepName != "" {
		state.CompletedSteps = append(state.CompletedSteps, req.StepName)
	}
	for k, v := range req.Config {
		state.Config[k] = v
	}
	storeMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (h *OnboardingHandler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		apierr.WriteError(w, apierr.BadRequest("tenant_id required"))
		return
	}
	storeMu.RLock()
	state, ok := store[tenantID]
	storeMu.RUnlock()
	if !ok {
		apierr.WriteError(w, apierr.NotFound("onboarding not started"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}
