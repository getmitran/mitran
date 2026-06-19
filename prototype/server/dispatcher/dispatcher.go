package dispatcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/getmitran/mitran/server/db"
	"github.com/getmitran/mitran/server/logger"
	"github.com/getmitran/mitran/server/streaming"
)

type execRequest struct {
	TaskID  string            `json:"task_id"`
	Agent   string            `json:"agent"`
	Task    string            `json:"task"`
	Context map[string]string `json:"context"`
}

func workerURL() string {
	if u := os.Getenv("MITRAN_WORKER_URL"); u != "" {
		return u
	}
	return "http://localhost:8888/execute"
}

// Start polls the store for queued tasks and dispatches them to the worker.
// It blocks until ctx is cancelled.
func Start(ctx context.Context, store *db.Store, sse *streaming.Hub) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	client := &http.Client{Timeout: 30 * time.Second}
	url := workerURL()

	logger.Info("dispatcher started", "worker_url", url)

	for {
		select {
		case <-ctx.Done():
			logger.Info("dispatcher shutting down")
			return
		case <-ticker.C:
			dispatch(ctx, store, client, url, sse)
		}
	}
}

func dispatch(ctx context.Context, store *db.Store, client *http.Client, url string, sse *streaming.Hub) {
	tasks := store.ListTasks(db.StatusQueued, "")
	for _, t := range tasks {
		select {
		case <-ctx.Done():
			return
		default:
		}

		now := time.Now()
		_ = store.UpdateTask(t.ID, func(task *db.Task) {
			task.Status = db.StatusRunning
			task.StartedAt = &now
		})
		broadcastStatus(sse, t.ID, string(db.StatusRunning))

		body, _ := json.Marshal(execRequest{
			TaskID:  t.ID,
			Agent:   t.AgentType,
			Task:    t.Description,
			Context: map[string]string{},
		})

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			markFailed(store, t.ID, err, sse)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			markFailed(store, t.ID, err, sse)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&result)
			resp.Body.Close()

			summary := ""
			if s, ok := result["summary"].(string); ok {
				summary = s
			}

			completed := time.Now()
			_ = store.UpdateTask(t.ID, func(task *db.Task) {
				task.Status = "completed"
				task.CompletedAt = &completed
				task.Resources = append(task.Resources, summary)
			})
			broadcastStatus(sse, t.ID, "completed")
		} else {
			resp.Body.Close()
			markFailed(store, t.ID, fmt.Errorf("worker returned status %d", resp.StatusCode), sse)
		}
	}
}

func markFailed(store *db.Store, id string, err error, sse *streaming.Hub) {
	logger.Error("task dispatch failed", "task_id", id, "error", err)
	_ = store.UpdateTask(id, func(task *db.Task) {
		task.Status = db.StatusFailed
		task.Resources = append(task.Resources, "error: "+err.Error())
	})
	broadcastStatus(sse, id, string(db.StatusFailed))
}

func broadcastStatus(sse *streaming.Hub, taskID string, status string) {
	data, _ := json.Marshal(map[string]string{"task_id": taskID, "status": status})
	sse.Broadcast("task_update", string(data))
}
