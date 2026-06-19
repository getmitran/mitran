package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/getmitran/mitran/server/apierr"
)

var slackBotToken = os.Getenv("SLACK_BOT_TOKEN")

type SlackNotifyRequest struct {
	Channel  string `json:"channel"`
	Text     string `json:"text"`
	ThreadTS string `json:"thread_ts,omitempty"`
}

type SlackWebhookEvent struct {
	Type      string          `json:"type"`
	Challenge string          `json:"challenge,omitempty"`
	Event     json.RawMessage `json:"event,omitempty"`
}

// POST /api/v1/integrations/slack/notify
func HandleSlackNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}

	var req SlackNotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, apierr.BadRequest("invalid request body"))
		return
	}

	if req.Channel == "" || req.Text == "" {
		apierr.WriteError(w, apierr.BadRequest("channel and text required"))
		return
	}

	payload, _ := json.Marshal(map[string]string{
		"channel":   req.Channel,
		"text":      req.Text,
		"thread_ts": req.ThreadTS,
	})

	slackReq, _ := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewReader(payload))
	slackReq.Header.Set("Authorization", "Bearer "+slackBotToken)
	slackReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(slackReq)
	if err != nil {
		apierr.WriteError(w, apierr.Internal("slack api error"))
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// POST /api/v1/integrations/slack/webhook
func HandleSlackWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}

	var event SlackWebhookEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		apierr.WriteError(w, apierr.BadRequest("invalid payload"))
		return
	}

	// URL verification challenge
	if event.Type == "url_verification" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"challenge": event.Challenge})
		return
	}

	// Future: process event callbacks (slash commands, interactive messages)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "received"})
}

// GET /api/v1/integrations/slack/status
func HandleSlackStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}

	connected := slackBotToken != ""
	status := "disconnected"
	if connected {
		// Verify token with auth.test
		req, _ := http.NewRequest("POST", "https://slack.com/api/auth.test", nil)
		req.Header.Set("Authorization", "Bearer "+slackBotToken)
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			var result map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&result)
			if ok, _ := result["ok"].(bool); ok {
				status = "connected"
			} else {
				status = "invalid_token"
			}
		} else {
			status = "unreachable"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    status,
		"token_set": connected,
	})
}
