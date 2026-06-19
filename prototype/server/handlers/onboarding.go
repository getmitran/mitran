package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/getmitran/mitran/server/apierr"
)

// OnboardingHandler manages project onboarding flows.
type OnboardingHandler struct{}

// NewOnboardingHandler creates a new OnboardingHandler.
func NewOnboardingHandler() *OnboardingHandler {
	return &OnboardingHandler{}
}

func (h *OnboardingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status": "ready",
			"steps":  []string{"init", "configure", "deploy", "verify"},
		})
	default:
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
	}
}
