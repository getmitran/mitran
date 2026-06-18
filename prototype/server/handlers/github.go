package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	event := r.Header.Get("X-GitHub-Event")
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
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
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	org := r.URL.Query().Get("org")
	if org == "" {
		http.Error(w, "org query param required", http.StatusBadRequest)
		return
	}

	resp, err := githubRequest("GET", fmt.Sprintf("/orgs/%s/repos?per_page=100", org), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
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
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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
