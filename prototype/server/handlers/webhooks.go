package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Webhook struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Event  string `json:"event"`
	Secret string `json:"secret,omitempty"`
}

type WebhookHandler struct {
	mu       sync.RWMutex
	webhooks map[string]Webhook
}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{webhooks: make(map[string]Webhook)}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.register(w, r)
	case http.MethodGet:
		h.list(w, r)
	case http.MethodDelete:
		h.remove(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *WebhookHandler) register(w http.ResponseWriter, r *http.Request) {
	var wh Webhook
	if err := json.NewDecoder(r.Body).Decode(&wh); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	wh.ID = uuid.New().String()
	h.mu.Lock()
	h.webhooks[wh.ID] = wh
	h.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wh)
}

func (h *WebhookHandler) list(w http.ResponseWriter, _ *http.Request) {
	h.mu.RLock()
	result := make([]Webhook, 0, len(h.webhooks))
	for _, wh := range h.webhooks {
		result = append(result, wh)
	}
	h.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *WebhookHandler) remove(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	id := parts[len(parts)-1]
	h.mu.Lock()
	if _, ok := h.webhooks[id]; !ok {
		h.mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	delete(h.webhooks, id)
	h.mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func (h *WebhookHandler) Dispatch(event string, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}
	h.mu.RLock()
	targets := make([]Webhook, 0)
	for _, wh := range h.webhooks {
		if wh.Event == event {
			targets = append(targets, wh)
		}
	}
	h.mu.RUnlock()

	for _, wh := range targets {
		req, err := http.NewRequest(http.MethodPost, wh.URL, bytes.NewReader(body))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		if wh.Secret != "" {
			mac := hmac.New(sha256.New, []byte(wh.Secret))
			mac.Write(body)
			req.Header.Set("X-Signature", hex.EncodeToString(mac.Sum(nil)))
		}
		go http.DefaultClient.Do(req)
	}
}
