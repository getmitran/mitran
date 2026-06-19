package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/getmitran/mitran/server/apierr"
)

type UsageRecord struct {
	Model        string    `json:"model"`
	TenantID     string    `json:"tenant_id"`
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
	Cost         float64   `json:"cost"`
	Timestamp    time.Time `json:"timestamp"`
}

type CostHandler struct {
	mu      sync.Mutex
	records []UsageRecord
}

func NewCostHandler() *CostHandler {
	return &CostHandler{}
}

func (h *CostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/v1/usage/summary":
		if r.Method != http.MethodGet {
			apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
			return
		}
		h.summary(w, r)
	case "/api/v1/usage":
		switch r.Method {
		case http.MethodPost:
			h.record(w, r)
		case http.MethodGet:
			h.list(w, r)
		default:
			apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		}
	default:
		apierr.WriteError(w, apierr.NotFound("not found"))
	}
}

func (h *CostHandler) record(w http.ResponseWriter, r *http.Request) {
	var rec UsageRecord
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		apierr.WriteError(w, apierr.BadRequest("invalid json"))
		return
	}
	if rec.Model == "" {
		apierr.WriteError(w, apierr.BadRequest("model required"))
		return
	}
	if rec.Timestamp.IsZero() {
		rec.Timestamp = time.Now().UTC()
	}
	h.mu.Lock()
	h.records = append(h.records, rec)
	h.mu.Unlock()
	writeJSON(w, http.StatusCreated, rec)
}

func (h *CostHandler) list(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	h.mu.Lock()
	defer h.mu.Unlock()
	var result []UsageRecord
	for _, rec := range h.records {
		if tenantID == "" || rec.TenantID == tenantID {
			result = append(result, rec)
		}
	}
	if result == nil {
		result = []UsageRecord{}
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *CostHandler) summary(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	type ModelSummary struct {
		Model        string  `json:"model"`
		TotalInput   int     `json:"total_input_tokens"`
		TotalOutput  int     `json:"total_output_tokens"`
		TotalCost    float64 `json:"total_cost"`
		RequestCount int     `json:"request_count"`
	}
	agg := map[string]*ModelSummary{}
	for _, rec := range h.records {
		s, ok := agg[rec.Model]
		if !ok {
			s = &ModelSummary{Model: rec.Model}
			agg[rec.Model] = s
		}
		s.TotalInput += rec.InputTokens
		s.TotalOutput += rec.OutputTokens
		s.TotalCost += rec.Cost
		s.RequestCount++
	}
	result := make([]ModelSummary, 0, len(agg))
	for _, s := range agg {
		result = append(result, *s)
	}
	writeJSON(w, http.StatusOK, result)
}
