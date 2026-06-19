package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimiter_AllowsInitialRequests(t *testing.T) {
	rl := &RateLimiter{buckets: make(map[string]*Bucket), rps: 10, burst: 10}
	for i := 0; i < 10; i++ {
		if !rl.Allow("127.0.0.1") {
			t.Fatalf("request %d should be allowed", i)
		}
	}
}

func TestRateLimiter_BlocksExcess(t *testing.T) {
	rl := &RateLimiter{buckets: make(map[string]*Bucket), rps: 1, burst: 2}
	rl.Allow("127.0.0.1") // 1
	rl.Allow("127.0.0.1") // 2
	if rl.Allow("127.0.0.1") {
		t.Error("third request should be blocked")
	}
}

func TestRateLimiter_Middleware429(t *testing.T) {
	mw := RateLimitMiddleware(1, 1)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	// First request passes
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	// Second request blocked
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req)
	if w2.Code != 429 {
		t.Errorf("expected 429, got %d", w2.Code)
	}
}
