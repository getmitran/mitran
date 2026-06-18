package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type apiKeyWindow struct {
	mu          sync.Mutex
	count       int
	windowStart time.Time
}

var apiKeyLimits sync.Map

func RateLimitByAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}

		now := time.Now()
		val, _ := apiKeyLimits.LoadOrStore(key, &apiKeyWindow{count: 0, windowStart: now})
		window := val.(*apiKeyWindow)

		window.mu.Lock()
		if now.Sub(window.windowStart) > time.Minute {
			window.count = 0
			window.windowStart = now
		}
		window.count++
		exceeded := window.count > 100
		remaining := time.Minute - now.Sub(window.windowStart)
		window.mu.Unlock()

		if exceeded {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(remaining.Seconds())+1))
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
