package integrations

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type GitHubWebhookHandler struct {
	secret      string
	agentRouter func(agentName string, payload map[string]interface{}) error
}

func NewGitHubWebhookHandler(agentRouter func(string, map[string]interface{}) error) *GitHubWebhookHandler {
	return &GitHubWebhookHandler{
		secret:      os.Getenv("GITHUB_WEBHOOK_SECRET"),
		agentRouter: agentRouter,
	}
}

func (h *GitHubWebhookHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/webhooks/github", h.HandleWebhook)
}

func (h *GitHubWebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if h.secret != "" {
		sig := r.Header.Get("X-Hub-Signature-256")
		if !h.verifySignature(body, sig) {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	event := r.Header.Get("X-GitHub-Event")
	if err := h.routeEvent(event, payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted", "event": event})
}

func (h *GitHubWebhookHandler) verifySignature(body []byte, signature string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	sig, err := hex.DecodeString(strings.TrimPrefix(signature, "sha256="))
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(body)
	return hmac.Equal(sig, mac.Sum(nil))
}

func (h *GitHubWebhookHandler) routeEvent(event string, payload map[string]interface{}) error {
	switch event {
	case "pull_request":
		return h.agentRouter("review", payload)
	case "push":
		return h.agentRouter("cicd", payload)
	case "issues":
		return h.agentRouter("ticketing", payload)
	case "ping":
		return nil
	default:
		return fmt.Errorf("unhandled event: %s", event)
	}
}
