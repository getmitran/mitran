package notifications

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var webhookURL = os.Getenv("MITRAN_SLACK_WEBHOOK")

type SlackMessage struct {
	Text string `json:"text"`
}

// Notify sends a formatted message to the configured Slack webhook.
func Notify(event string, payload map[string]string) error {
	if webhookURL == "" {
		return nil
	}
	var text string
	switch event {
	case "task_completed":
		text = "✅ Task completed: " + payload["title"] + " (agent: " + payload["agent"] + ")"
	case "task_failed":
		text = "❌ Task failed: " + payload["title"] + " — " + payload["error"]
	case "approval_needed":
		text = "⏳ Approval needed: " + payload["title"] + " (agent: " + payload["agent"] + ")"
	default:
		text = event + ": " + payload["title"]
	}
	body, _ := json.Marshal(SlackMessage{Text: text})
	client := &http.Client{Timeout: 5 * time.Second}
	_, err := client.Post(webhookURL, "application/json", bytes.NewReader(body))
	return err
}
