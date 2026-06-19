package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

type cspContextKey string

const nonceKey cspContextKey = "csp-nonce"

func CSP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		nonce := base64.StdEncoding.EncodeToString(b)

		ctx := context.WithValue(r.Context(), nonceKey, nonce)

		w.Header().Set("Content-Security-Policy", fmt.Sprintf(
			"default-src 'self'; script-src 'self' 'nonce-%s'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'",
			nonce,
		))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NonceFromContext(r *http.Request) string {
	if v, ok := r.Context().Value(nonceKey).(string); ok {
		return v
	}
	return ""
}
