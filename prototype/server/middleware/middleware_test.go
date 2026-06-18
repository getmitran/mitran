package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSHeaders(t *testing.T) {
	handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "http://localhost:3000" {
		t.Fatalf("expected origin http://localhost:3000, got %s", v)
	}
	if v := w.Header().Get("Access-Control-Allow-Methods"); v == "" {
		t.Fatal("expected Allow-Methods header")
	}
	if v := w.Header().Get("Access-Control-Allow-Headers"); v == "" {
		t.Fatal("expected Allow-Headers header")
	}
}

func TestCORSPreflight(t *testing.T) {
	handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("OPTIONS", "/api/v1/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for OPTIONS, got %d", w.Code)
	}
}

func TestCORSRejectsNonLocalhost(t *testing.T) {
	handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	req.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "" {
		t.Fatalf("expected no ACAO header for non-localhost, got %s", v)
	}
}

func TestCORSPassesThrough(t *testing.T) {
	called := false
	handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCORSMultipleLocalhostPorts(t *testing.T) {
	handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	origins := []string{"http://localhost:5173", "http://localhost:7780", "http://localhost"}
	for _, origin := range origins {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", origin)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if v := w.Header().Get("Access-Control-Allow-Origin"); v != origin {
			t.Errorf("origin %s: expected ACAO=%s, got %s", origin, origin, v)
		}
	}
}
