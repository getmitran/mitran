package mcp

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Registry struct {
	mu      sync.RWMutex
	dir     string
	Servers []ServerConfig `json:"servers"`
}

func NewRegistry(dir string) (*Registry, error) {
	os.MkdirAll(dir, 0755)
	r := &Registry{dir: dir}
	r.load()
	return r, nil
}

func (r *Registry) Add(cfg ServerConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cfg.Enabled = true
	r.Servers = append(r.Servers, cfg)
	r.save()
}

func (r *Registry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, s := range r.Servers {
		if s.ID == id {
			r.Servers = append(r.Servers[:i], r.Servers[i+1:]...)
			break
		}
	}
	r.save()
}

func (r *Registry) Toggle(id string, enabled bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.Servers {
		if r.Servers[i].ID == id {
			r.Servers[i].Enabled = enabled
			break
		}
	}
	r.save()
}

func (r *Registry) List() []ServerConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]ServerConfig{}, r.Servers...)
}

func (r *Registry) Enabled() []ServerConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ServerConfig
	for _, s := range r.Servers {
		if s.Enabled {
			result = append(result, s)
		}
	}
	return result
}

func (r *Registry) load() {
	data, _ := os.ReadFile(filepath.Join(r.dir, "mcp_servers.json"))
	json.Unmarshal(data, &r.Servers)
}

func (r *Registry) save() {
	data, _ := json.MarshalIndent(r.Servers, "", "  ")
	os.WriteFile(filepath.Join(r.dir, "mcp_servers.json"), data, 0644)
}

func RegisterRoutes(mux *http.ServeMux, reg *Registry) {
	mux.HandleFunc("GET /api/v1/mcp/servers", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(reg.List())
	})
	mux.HandleFunc("POST /api/v1/mcp/servers", func(w http.ResponseWriter, r *http.Request) {
		var cfg ServerConfig
		json.NewDecoder(r.Body).Decode(&cfg)
		reg.Add(cfg)
		w.WriteHeader(201)
	})
	mux.HandleFunc("DELETE /api/v1/mcp/servers/{id}", func(w http.ResponseWriter, r *http.Request) {
		reg.Remove(r.PathValue("id"))
		w.WriteHeader(204)
	})
	mux.HandleFunc("POST /api/v1/mcp/servers/{id}/toggle", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Enabled bool `json:"enabled"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		reg.Toggle(r.PathValue("id"), body.Enabled)
		w.WriteHeader(200)
	})
}
