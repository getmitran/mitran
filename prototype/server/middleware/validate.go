package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const maxBodySize = 1 << 20 // 1MB

func ValidateJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			ct := r.Header.Get("Content-Type")
			if !strings.HasPrefix(ct, "application/json") {
				http.Error(w, `{"error":"Content-Type must be application/json"}`, http.StatusUnsupportedMediaType)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		}
		next.ServeHTTP(w, r)
	})
}

func ValidateRequired(fields ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
				return
			}

			var missing []string
			for _, f := range fields {
				val, ok := body[f]
				if !ok {
					missing = append(missing, f)
					continue
				}
				if s, isStr := val.(string); isStr && strings.TrimSpace(s) == "" {
					missing = append(missing, f)
				}
			}

			if len(missing) > 0 {
				msg := fmt.Sprintf(`{"error":"missing required fields","fields":["%s"]}`, strings.Join(missing, `","`))
				http.Error(w, msg, http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
