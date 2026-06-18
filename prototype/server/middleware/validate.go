package middleware

import (
	"net/http"
	"strings"
)

const MaxBodySize = 1 << 20 // 1MB

func ValidateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			ct := r.Header.Get("Content-Type")
			if ct != "" && !strings.Contains(ct, "application/json") && !strings.Contains(ct, "text/event-stream") {
				http.Error(w, "unsupported content type", 415)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)
		}
		next.ServeHTTP(w, r)
	})
}
