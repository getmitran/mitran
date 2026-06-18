package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS middleware. Respects MITRAN_CORS_ORIGINS env var (comma-separated).
// Defaults to localhost origins only in dev mode.
func CORS(next http.Handler) http.Handler {
	allowed := map[string]bool{}
	if origins := os.Getenv("MITRAN_CORS_ORIGINS"); origins != "" {
		for _, o := range strings.Split(origins, ",") {
			allowed[strings.TrimSpace(o)] = true
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if len(allowed) > 0 {
			if allowed[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		} else if strings.HasPrefix(origin, "http://localhost") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID, X-User")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
