package streaming

import "sync"

// Hub manages connected SSE clients.
type Hub struct {
	mu      sync.RWMutex
	clients map[chan Event]struct{}
}

// Event is a single SSE event.
type Event struct {
	Type string
	Data string
}

// NewHub creates an SSE hub.
func NewHub() *Hub {
	return &Hub{clients: make(map[chan Event]struct{})}
}

// Subscribe adds a client channel.
func (h *Hub) Subscribe() chan Event {
	ch := make(chan Event, 16)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

// Unsubscribe removes a client channel.
func (h *Hub) Unsubscribe(ch chan Event) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
	close(ch)
}

// Broadcast sends an event to all connected clients.
func (h *Hub) Broadcast(event, data string) {
	e := Event{Type: event, Data: data}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.clients {
		select {
		case ch <- e:
		default: // drop if client is slow
		}
	}
}
