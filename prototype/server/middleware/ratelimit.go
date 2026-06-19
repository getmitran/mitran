package middleware

import (
	"math"
	"net/http"
	"sync"
	"time"
	"github.com/getmitran/mitran/server/apierr"
)

type Bucket struct {
	tokens   float64
	max      float64
	refill   float64
	lastTime time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*Bucket
	rps     float64
	burst   int
	calls   int
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.calls++
	if rl.calls%1000 == 0 {
		cutoff := time.Now().Add(-10 * time.Minute)
		for k, b := range rl.buckets {
			if b.lastTime.Before(cutoff) {
				delete(rl.buckets, k)
			}
		}
	}

	b, ok := rl.buckets[key]
	if !ok {
		b = &Bucket{tokens: float64(rl.burst), max: float64(rl.burst), refill: rl.rps, lastTime: time.Now()}
		rl.buckets[key] = b
	}
	now := time.Now()
	b.tokens = math.Min(b.max, b.tokens+now.Sub(b.lastTime).Seconds()*b.refill)
	b.lastTime = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func RateLimitMiddleware(rps float64, burst int) func(http.Handler) http.Handler {
	if rps == 0 {
		rps = 100
	}
	if burst == 0 {
		burst = 200
	}
	rl := &RateLimiter{buckets: make(map[string]*Bucket), rps: rps, burst: burst}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr
			if tenant := r.Header.Get("X-Tenant-ID"); tenant != "" {
				key = tenant + ":" + key
			}
			if !rl.Allow(key) {
				w.Header().Set("Retry-After", "1")
				apierr.WriteError(w, apierr.Internal("rate limit exceeded"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
