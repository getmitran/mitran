package notify

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Channel string

const (
	ChannelWebSocket Channel = "websocket"
	ChannelSlack     Channel = "slack"
	ChannelWebhook   Channel = "webhook"
)

type Notification struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // task.completed, agent.error, checkpoint.pending
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Channel   Channel   `json:"channel"`
	Sent      bool      `json:"sent"`
	CreatedAt time.Time `json:"created_at"`
}

type WebhookConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Events  []string          `json:"events"`
}

type Service struct {
	mu          sync.RWMutex
	history     []Notification
	webhooks    []WebhookConfig
	slackURL    string
	wsBroadcast func([]byte) // injected WebSocket broadcaster
}

func New(wsBroadcast func([]byte)) *Service {
	return &Service{wsBroadcast: wsBroadcast}
}

func (s *Service) SetSlackWebhook(url string) { s.slackURL = url }
func (s *Service) AddWebhook(cfg WebhookConfig) {
	s.mu.Lock()
	s.webhooks = append(s.webhooks, cfg)
	s.mu.Unlock()
}

func (s *Service) Send(notifType, title, body string) {
	n := Notification{
		ID:        time.Now().Format("20060102150405"),
		Type:      notifType,
		Title:     title,
		Body:      body,
		CreatedAt: time.Now(),
	}

	// WebSocket
	if s.wsBroadcast != nil {
		data, _ := json.Marshal(map[string]interface{}{"type": notifType, "data": n})
		s.wsBroadcast(data)
		n.Channel = ChannelWebSocket
		n.Sent = true
	}

	// Slack
	if s.slackURL != "" {
		go s.sendSlack(title, body)
	}

	// Webhooks
	s.mu.RLock()
	for _, wh := range s.webhooks {
		for _, evt := range wh.Events {
			if evt == notifType || evt == "*" {
				go s.sendWebhook(wh, n)
				break
			}
		}
	}
	s.mu.RUnlock()

	s.mu.Lock()
	s.history = append(s.history, n)
	if len(s.history) > 500 {
		s.history = s.history[len(s.history)-500:]
	}
	s.mu.Unlock()
}

func (s *Service) History(limit int) []Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if limit <= 0 || limit > len(s.history) {
		limit = len(s.history)
	}
	return s.history[len(s.history)-limit:]
}

func (s *Service) sendSlack(title, body string) {
	payload, _ := json.Marshal(map[string]string{"text": "*" + title + "*\n" + body})
	http.Post(s.slackURL, "application/json", bytes.NewReader(payload))
}

func (s *Service) sendWebhook(cfg WebhookConfig, n Notification) {
	data, _ := json.Marshal(n)
	req, err := http.NewRequest("POST", cfg.URL, bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	if _, err := client.Do(req); err != nil {
		log.Printf("[notify] webhook error: %v", err)
	}
}
