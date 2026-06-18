package tenant

import (
	"context"
	"encoding/json"
	"net/http"
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
}

func NewStore() *TenantStore {
	return &TenantStore{tenants: make(map[string]*Tenant)}
}

func (s *TenantStore) Create(id, name string) *Tenant {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := &Tenant{ID: id, Name: name, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	s.tenants[id] = t
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

func (s *TenantStore) Delete(id string) { s.mu.Lock(); delete(s.tenants, id); s.mu.Unlock() }

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
