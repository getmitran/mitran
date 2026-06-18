// NOTE: In-memory + JSON persistence is suitable for <50 tenants.
// For larger deployments, migrate to SQLite (see docs/v020-roadmap.md).
package tenant

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Tenant struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type TenantStore struct {
	mu      sync.RWMutex
	tenants map[string]*Tenant
	path    string
}

func dataPath() string {
	if p := os.Getenv("TENANT_DATA_PATH"); p != "" {
		return p
	}
	return "./data/tenants.json"
}

func NewStore() *TenantStore {
	s := &TenantStore{tenants: make(map[string]*Tenant), path: dataPath()}
	s.load()
	if len(s.tenants) > 40 {
		log.Println("WARNING: tenant count approaching file-based limit, consider migrating to SQLite")
	}
	return s
}

func (s *TenantStore) load() {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	var tenants []*Tenant
	if json.Unmarshal(data, &tenants) == nil {
		for _, t := range tenants {
			s.tenants[t.ID] = t
		}
	}
}

func (s *TenantStore) persist() {
	out := make([]*Tenant, 0, len(s.tenants))
	for _, t := range s.tenants {
		out = append(out, t)
	}
	writeJSONAtomic(s.path, out)
}

func (s *TenantStore) Create(id, name string) *Tenant {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := &Tenant{ID: id, Name: name, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	s.tenants[id] = t
	s.persist()
	return t
}

func (s *TenantStore) Get(id string) *Tenant {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tenants[id]
}

func (s *TenantStore) List() []*Tenant {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Tenant, 0, len(s.tenants))
	for _, t := range s.tenants {
		out = append(out, t)
	}
	return out
}

func (s *TenantStore) Delete(id string) {
	s.mu.Lock()
	delete(s.tenants, id)
	s.persist()
	s.mu.Unlock()
}

type ctxKey struct{}

func FromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKey{}).(string); ok {
		return v
	}
	return ""
}

func Middleware(store *TenantStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tid := r.Header.Get("X-Tenant-ID")
			if tid == "" {
				next.ServeHTTP(w, r)
				return
			}
			if store.Get(tid) == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"error": "unknown tenant"})
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, tid)))
		})
	}
}

func writeJSONAtomic(path string, data any) error {
	os.MkdirAll(filepath.Dir(path), 0o755)
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(data); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, path)
}
