package handlers

import (
	"net/http"
	"os"
	"time"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func checkHealth(url string) bool {
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Get(url)
	return err == nil && resp.StatusCode == http.StatusOK
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	workerURL := os.Getenv("WORKER_URL")
	if workerURL == "" {
		workerURL = "http://localhost:8888"
	}
	chromaURL := os.Getenv("CHROMADB_URL")
	if chromaURL == "" {
		chromaURL = "http://localhost:8000"
	}
	checks := map[string]bool{
		"worker":   checkHealth(workerURL + "/health"),
		"chromadb": checkHealth(chromaURL + "/api/v1/heartbeat"),
	}
	status, code := "ready", http.StatusOK
	for _, ok := range checks {
		if !ok {
			status, code = "not_ready", http.StatusServiceUnavailable
			break
		}
	}
	writeJSON(w, code, map[string]interface{}{"status": status, "checks": checks})
}
