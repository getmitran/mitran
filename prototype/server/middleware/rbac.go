package middleware

import (
	"encoding/json"
	"net/http"
	"os"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleDeveloper Role = "developer"
	RoleViewer    Role = "viewer"
)

type RBACConfig struct {
	Users map[string]Role `json:"users"`
}

func LoadRBACConfig() RBACConfig {
	cfg := RBACConfig{Users: make(map[string]Role)}
	path := os.Getenv("RBAC_CONFIG_PATH")
	if path == "" {
		return cfg
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	json.Unmarshal(data, &cfg)
	return cfg
}

func CanWrite(role Role) bool  { return role == RoleAdmin || role == RoleDeveloper }
func CanAdmin(role Role) bool  { return role == RoleAdmin }

func RBACMiddleware(config RBACConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Header.Get("X-User")
			if user == "" {
				http.Error(w, `{"error":"missing X-User header"}`, http.StatusUnauthorized)
				return
			}
			role, ok := config.Users[user]
			if !ok {
				http.Error(w, `{"error":"user not found"}`, http.StatusForbidden)
				return
			}
			if r.Method != http.MethodGet && !CanWrite(role) {
				http.Error(w, `{"error":"permission denied"}`, http.StatusForbidden)
				return
			}
			if isAdminRoute(r.URL.Path) && !CanAdmin(role) {
				http.Error(w, `{"error":"admin access required"}`, http.StatusForbidden)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}

func isAdminRoute(path string) bool {
	return len(path) > 10 && (path[:10] == "/api/users" || path[:13] == "/api/settings")
}
