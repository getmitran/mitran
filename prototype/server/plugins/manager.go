package plugins

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type AgentConfig struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Model        string   `json:"model"`
	SystemPrompt string   `json:"system_prompt"`
	Tools        []string `json:"tools"`
	MaxTokens    int      `json:"max_tokens"`
	Temperature  float64  `json:"temperature"`
	Enabled      bool     `json:"enabled"`
}

type Plugin struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	EntryPoint  string `json:"entry_point"`
	Enabled     bool   `json:"enabled"`
}

type Manager struct {
	mu      sync.RWMutex
	dir     string
	Agents  []AgentConfig `json:"agents"`
	Plugins []Plugin      `json:"plugins"`
}

func New(dir string) (*Manager, error) {
	os.MkdirAll(dir, 0755)
	m := &Manager{dir: dir}
	m.load()
	if len(m.Agents) == 0 {
		m.loadDefaults()
	}
	return m, nil
}

func (m *Manager) GetAgent(name string) *AgentConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for i := range m.Agents {
		if m.Agents[i].Name == name {
			return &m.Agents[i]
		}
	}
	return nil
}

func (m *Manager) UpdateAgent(cfg AgentConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.Agents {
		if m.Agents[i].Name == cfg.Name {
			m.Agents[i] = cfg
			m.save()
			return
		}
	}
	m.Agents = append(m.Agents, cfg)
	m.save()
}

func (m *Manager) ListAgents() []AgentConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Agents
}

func (m *Manager) ListPlugins() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Plugins
}

func (m *Manager) InstallPlugin(p Plugin) {
	m.mu.Lock()
	m.Plugins = append(m.Plugins, p)
	m.save()
	m.mu.Unlock()
}

func (m *Manager) RemovePlugin(id string) {
	m.mu.Lock()
	for i, p := range m.Plugins {
		if p.ID == id {
			m.Plugins = append(m.Plugins[:i], m.Plugins[i+1:]...)
			break
		}
	}
	m.save()
	m.mu.Unlock()
}

func (m *Manager) Reload() {
	m.mu.Lock()
	m.load()
	m.mu.Unlock()
	log.Println("[plugins] config reloaded")
}

func (m *Manager) load() {
	data, _ := os.ReadFile(filepath.Join(m.dir, "agents.json"))
	json.Unmarshal(data, &m.Agents)
	data2, _ := os.ReadFile(filepath.Join(m.dir, "plugins.json"))
	json.Unmarshal(data2, &m.Plugins)
}

func (m *Manager) save() {
	d1, _ := json.MarshalIndent(m.Agents, "", "  ")
	os.WriteFile(filepath.Join(m.dir, "agents.json"), d1, 0644)
	d2, _ := json.MarshalIndent(m.Plugins, "", "  ")
	os.WriteFile(filepath.Join(m.dir, "plugins.json"), d2, 0644)
}

func (m *Manager) loadDefaults() {
	defaults := []AgentConfig{
		{Name: "dev", Type: "dev", Model: "claude-sonnet-4-20250514", MaxTokens: 4096, Temperature: 0.7, Enabled: true},
		{Name: "docs", Type: "docs", Model: "claude-sonnet-4-20250514", MaxTokens: 4096, Temperature: 0.5, Enabled: true},
		{Name: "ops", Type: "ops", Model: "claude-sonnet-4-20250514", MaxTokens: 4096, Temperature: 0.3, Enabled: true},
		{Name: "review", Type: "review", Model: "claude-sonnet-4-20250514", MaxTokens: 4096, Temperature: 0.3, Enabled: true},
		{Name: "hr", Type: "hr", Model: "claude-sonnet-4-20250514", MaxTokens: 4096, Temperature: 0.5, Enabled: true},
		{Name: "cicd", Type: "cicd", Model: "claude-sonnet-4-20250514", MaxTokens: 4096, Temperature: 0.3, Enabled: true},
		{Name: "tickets", Type: "tickets", Model: "claude-sonnet-4-20250514", MaxTokens: 4096, Temperature: 0.5, Enabled: true},
		{Name: "wiki", Type: "wiki", Model: "claude-sonnet-4-20250514", MaxTokens: 4096, Temperature: 0.5, Enabled: true},
	}
	m.Agents = defaults
	m.save()
}

func RegisterRoutes(mux *http.ServeMux, mgr *Manager) {
	mux.HandleFunc("GET /api/v1/agents/config", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mgr.ListAgents())
	})
	mux.HandleFunc("PUT /api/v1/agents/config/{name}", func(w http.ResponseWriter, r *http.Request) {
		var cfg AgentConfig
		json.NewDecoder(r.Body).Decode(&cfg)
		cfg.Name = r.PathValue("name")
		mgr.UpdateAgent(cfg)
		w.WriteHeader(200)
	})
	mux.HandleFunc("GET /api/v1/plugins", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mgr.ListPlugins())
	})
	mux.HandleFunc("POST /api/v1/plugins", func(w http.ResponseWriter, r *http.Request) {
		var p Plugin
		json.NewDecoder(r.Body).Decode(&p)
		mgr.InstallPlugin(p)
		w.WriteHeader(201)
	})
	mux.HandleFunc("DELETE /api/v1/plugins/{id}", func(w http.ResponseWriter, r *http.Request) {
		mgr.RemovePlugin(r.PathValue("id"))
		w.WriteHeader(204)
	})
	mux.HandleFunc("POST /api/v1/agents/reload", func(w http.ResponseWriter, r *http.Request) {
		mgr.Reload()
		w.WriteHeader(200)
	})
}
