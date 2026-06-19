package session

import (
	crypto_rand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"` // user, assistant
	Content   string    `json:"content"`
	Agent     string    `json:"agent,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Agent     string    `json:"agent"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Active    bool      `json:"active"`
}

type Manager struct {
	mu       sync.RWMutex
	dir      string
	Sessions []Session `json:"sessions"`
}

func New(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	m := &Manager{dir: dir}
	m.load()
	return m, nil
}

func (m *Manager) Create(userID, agent string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	s := Session{
		ID: generateID(), UserID: userID, Agent: agent,
		Messages: []Message{}, CreatedAt: time.Now(), UpdatedAt: time.Now(), Active: true,
	}
	m.Sessions = append(m.Sessions, s)
	m.save()
	return &m.Sessions[len(m.Sessions)-1]
}

func (m *Manager) Get(id string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for i := range m.Sessions {
		if m.Sessions[i].ID == id {
			return &m.Sessions[i]
		}
	}
	return nil
}

func (m *Manager) AddMessage(sessionID string, msg Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.Sessions {
		if m.Sessions[i].ID == sessionID {
			m.Sessions[i].Messages = append(m.Sessions[i].Messages, msg)
			m.Sessions[i].UpdatedAt = time.Now()
			m.save()
			return
		}
	}
}

func (m *Manager) List(userID string) []Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []Session
	for _, s := range m.Sessions {
		if userID == "" || s.UserID == userID {
			result = append(result, s)
		}
	}
	return result
}

func (m *Manager) End(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.Sessions {
		if m.Sessions[i].ID == id {
			m.Sessions[i].Active = false
			m.save()
			return
		}
	}
}

func (m *Manager) load() {
	data, err := os.ReadFile(filepath.Join(m.dir, "sessions.json"))
	if err == nil {
		json.Unmarshal(data, &m.Sessions)
	}
}

func (m *Manager) save() {
	data, _ := json.MarshalIndent(m.Sessions, "", "  ")
	os.WriteFile(filepath.Join(m.dir, "sessions.json"), data, 0644)
}

func generateID() string {
	return time.Now().Format("20060102-150405") + "-" + randomHex(16)
}

func randomHex(n int) string {
	b := make([]byte, n/2)
	crypto_rand.Read(b)
	return hex.EncodeToString(b)
}
