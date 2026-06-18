package events

import (
	"encoding/json"
	"sync"
	"time"
)

type EventType string

const (
	TaskCreated    EventType = "task.created"
	TaskUpdated    EventType = "task.updated"
	TaskCompleted  EventType = "task.completed"
	AgentStarted   EventType = "agent.started"
	AgentCompleted EventType = "agent.completed"
	AgentError     EventType = "agent.error"
	Checkpoint     EventType = "task.checkpoint"
	ChatMessage    EventType = "chat.message"
)

type Event struct {
	Type      EventType   `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	Source    string      `json:"source"`
}

type Handler func(Event)

type Bus struct {
	mu       sync.RWMutex
	handlers map[EventType][]Handler
	history  []Event
	maxHist  int
}

func New() *Bus {
	return &Bus{handlers: make(map[EventType][]Handler), maxHist: 1000}
}

func (b *Bus) Subscribe(eventType EventType, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *Bus) SubscribeAll(handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers["*"] = append(b.handlers["*"], handler)
}

func (b *Bus) Publish(eventType EventType, data interface{}, source string) {
	event := Event{Type: eventType, Data: data, Timestamp: time.Now(), Source: source}
	b.mu.Lock()
	b.history = append(b.history, event)
	if len(b.history) > b.maxHist {
		b.history = b.history[len(b.history)-b.maxHist:]
	}
	handlers := append([]Handler{}, b.handlers[eventType]...)
	handlers = append(handlers, b.handlers["*"]...)
	b.mu.Unlock()
	for _, h := range handlers {
		go h(event)
	}
}

func (b *Bus) History(limit int) []Event {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if limit <= 0 || limit > len(b.history) {
		limit = len(b.history)
	}
	return b.history[len(b.history)-limit:]
}

func (b *Bus) JSON() []byte {
	data, _ := json.Marshal(b.History(50))
	return data
}
