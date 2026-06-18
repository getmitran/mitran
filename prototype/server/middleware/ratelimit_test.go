package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimiter_AllowsInitialRequests(t *testing.T) {
	rl := NewRateLimiter(10, 10)
	for i := 0; i < 10; i++ {
		if !rl.Allow("127.0.0.1") { t.Fatalf("request %d should be allowed", i) }
	}
}

func TestRateLimiter_BlocksExcess(t *testing.T) {
	rl := NewRateLimiter(1, 2)
	rl.Allow("127.0.0.1") // 1
	rl.Allow("127.0.0.1") // 2
	if rl.Allow("127.0.0.1") { t.Error("third request should be blocked") }
}

func TestRateLimiter_Middleware429(t *testing.T) {
	rl := NewRateLimiter(1, 1)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))
	// First request passes
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 200 { t.Errorf("expected 200, got %d", w.Code) }
	// Second request blocked
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req)
	if w2.Code != 429 { t.Errorf("expected 429, got %d", w2.Code) }
}
