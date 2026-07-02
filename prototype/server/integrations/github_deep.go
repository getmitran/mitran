package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/getmitran/mitran/server/apierr"
)

// GitHubDeepHandler provides deep GitHub integration beyond webhooks:
// creating PRs and triggering agent reviews on PR events.
type GitHubDeepHandler struct {
	token       string
	agentRouter func(agentName string, payload map[string]interface{}) error
}

func NewGitHubDeepHandler(agentRouter func(string, map[string]interface{}) error) *GitHubDeepHandler {
	return &GitHubDeepHandler{
		token:       os.Getenv("GITHUB_TOKEN"),
		agentRouter: agentRouter,
	}
}

func (h *GitHubDeepHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/integrations/github/pr", h.HandleCreatePR)
	mux.HandleFunc("/api/v1/integrations/github/review", h.HandleTriggerReview)
}

// CreatePRRequest is the request body for PR creation.
type CreatePRRequest struct {
	Repo   string `json:"repo"`
	Title  string `json:"title"`
	Branch string `json:"branch"`
	Base   string `json:"base,omitempty"`
	Body   string `json:"body"`
}

// HandleCreatePR creates a pull request via the GitHub API.
func (h *GitHubDeepHandler) HandleCreatePR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}

	if h.token == "" {
		apierr.WriteError(w, apierr.Internal("GITHUB_TOKEN not configured"))
		return
	}

	var req CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, apierr.BadRequest("invalid JSON: "+err.Error()))
		return
	}
	defer r.Body.Close()

	if req.Repo == "" || req.Title == "" || req.Branch == "" {
		apierr.WriteError(w, apierr.BadRequest("repo, title, and branch are required"))
		return
	}

	if req.Base == "" {
		req.Base = "main"
	}

	prPayload := map[string]string{
		"title": req.Title,
		"head":  req.Branch,
		"base":  req.Base,
		"body":  req.Body,
	}

	payloadBytes, err := json.Marshal(prPayload)
	if err != nil {
		apierr.WriteError(w, apierr.Internal("failed to marshal PR payload"))
		return
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", req.Repo)
	ghReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		apierr.WriteError(w, apierr.Internal("failed to create request"))
		return
	}

	ghReq.Header.Set("Authorization", "Bearer "+h.token)
	ghReq.Header.Set("Accept", "application/vnd.github+json")
	ghReq.Header.Set("Content-Type", "application/json")
	ghReq.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(ghReq)
	if err != nil {
		apierr.WriteError(w, apierr.Internal("GitHub API error: "+err.Error()))
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		apierr.WriteError(w, apierr.Internal(fmt.Sprintf("GitHub returned %d: %s", resp.StatusCode, string(respBody))))
		return
	}

	var prResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &prResponse); err != nil {
		w.WriteHeader(http.StatusCreated)
		w.Write(respBody)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "created",
		"number":  prResponse["number"],
		"url":     prResponse["html_url"],
		"api_url": prResponse["url"],
	})
}

// ReviewTriggerRequest is the request body for triggering an architect review.
type ReviewTriggerRequest struct {
	Repo     string `json:"repo"`
	PRNumber int    `json:"pr_number"`
	Action   string `json:"action,omitempty"`
}

// HandleTriggerReview triggers the architect agent to review a PR.
func (h *GitHubDeepHandler) HandleTriggerReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}

	var req ReviewTriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, apierr.BadRequest("invalid JSON: "+err.Error()))
		return
	}
	defer r.Body.Close()

	if req.Repo == "" || req.PRNumber == 0 {
		apierr.WriteError(w, apierr.BadRequest("repo and pr_number are required"))
		return
	}

	if req.Action == "" {
		req.Action = "opened"
	}

	payload := map[string]interface{}{
		"action": req.Action,
		"pull_request": map[string]interface{}{
			"number": req.PRNumber,
		},
		"repository": map[string]interface{}{
			"full_name": req.Repo,
		},
	}

	if err := h.agentRouter("architect", payload); err != nil {
		apierr.WriteError(w, apierr.Internal("failed to trigger review: "+err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "review_triggered",
		"agent":     "architect",
		"repo":      req.Repo,
		"pr_number": req.PRNumber,
	})
}
