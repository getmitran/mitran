package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/getmitran/mitran/server/auth"
	"github.com/getmitran/mitran/server/apierr"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleDeveloper Role = "developer"
	RoleViewer    Role = "viewer"
)

type RBACConfig struct{ Users map[string]Role `json:"users"` }

func LoadRBACConfig() RBACConfig {
	cfg := RBACConfig{Users: make(map[string]Role)}
	if p := os.Getenv("RBAC_CONFIG_PATH"); p != "" {
		if d, err := os.ReadFile(p); err == nil {
			json.Unmarshal(d, &cfg)
		}
	}
	return cfg
}

func CanWrite(r Role) bool { return r == RoleAdmin || r == RoleDeveloper }
func CanAdmin(r Role) bool { return r == RoleAdmin }

func extractUser(req *http.Request) string {
	c, err := req.Cookie("session")
	if err != nil || c.Value == "" {
		return ""
	}
	p := strings.SplitN(c.Value, "|", 3)
	if len(p) != 3 || os.Getenv("MITRAN_SESSION_SECRET") == "" {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(os.Getenv("MITRAN_SESSION_SECRET")))
	mac.Write([]byte(p[0] + "|" + p[1]))
	if !hmac.Equal([]byte(p[2]), []byte(hex.EncodeToString(mac.Sum(nil)))) {
		return ""
	}
	expiry, err := strconv.ParseInt(p[1], 10, 64)
	if err != nil || time.Now().Unix() > expiry {
		return ""
	}
	if auth.IsRevoked(c.Value) {
		return ""
	}
	return p[0]
}

func RBACMiddleware(config RBACConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := extractUser(r)
			if user == "" && os.Getenv("MITRAN_ENV") == "development" {
				user = r.Header.Get("X-User")
			}
			if user == "" {
				apierr.WriteError(w, apierr.Unauthorized("unauthorized"))
				return
			}
			role, ok := config.Users[user]
			if !ok {
				apierr.WriteError(w, apierr.Forbidden("user not found"))
				return
			}
			if r.Method != http.MethodGet && !CanWrite(role) {
				apierr.WriteError(w, apierr.Forbidden("permission denied"))
				return
			}
			if isAdminRoute(r.URL.Path) && !CanAdmin(role) {
				apierr.WriteError(w, apierr.Forbidden("admin access required"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isAdminRoute(path string) bool {
	return strings.HasPrefix(path, "/api/v1/users") || strings.HasPrefix(path, "/api/v1/settings")
}
