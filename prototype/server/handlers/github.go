package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/getmitran/mitran/server/apierr"
)

var githubAPIURL = "https://api.github.com"

func init() {
	if url := os.Getenv("GITHUB_API_URL"); url != "" {
		githubAPIURL = url
	}
}

func githubRequest(method, path string, body io.Reader) (*http.Response, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN not set")
	}
	req, err := http.NewRequest(method, githubAPIURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

// POST /api/v1/integrations/github/webhook
func HandleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		apierr.WriteError(w, apierr.BadRequest("failed to read body"))
		return
	}
	defer r.Body.Close()

	event := r.Header.Get("X-GitHub-Event")
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		apierr.WriteError(w, apierr.BadRequest("invalid JSON"))
		return
	}

	// Log and acknowledge - downstream processing happens async
	fmt.Printf("[github-webhook] event=%s action=%v\n", event, payload["action"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"received": true,
		"event":    event,
	})
}

// GET /api/v1/integrations/github/repos
func HandleGitHubListRepos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}
	org := r.URL.Query().Get("org")
	if org == "" {
		apierr.WriteError(w, apierr.BadRequest("org query param required"))
		return
	}

	resp, err := githubRequest("GET", fmt.Sprintf("/orgs/%s/repos?per_page=100", org), nil)
	if err != nil {
		apierr.WriteError(w, apierr.Internal("github api error"))
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// GET /api/v1/integrations/github/status
func HandleGitHubStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"connected": false,
			"error":     "GITHUB_TOKEN not configured",
		})
		return
	}

	resp, err := githubRequest("GET", "/user", nil)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"connected": false,
			"error":     err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"connected":      resp.StatusCode == 200,
		"authenticated":  resp.StatusCode == 200,
		"user":           user["login"],
		"api_url":        githubAPIURL,
		"rate_remaining": resp.Header.Get("X-RateLimit-Remaining"),
	})
}
