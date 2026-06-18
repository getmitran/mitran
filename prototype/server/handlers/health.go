package handlers

import (
	"net/http"
	"os"
	"time"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("WORKER_URL")
	if url == "" {
		url = "http://localhost:8888"
	}
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url + "/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
