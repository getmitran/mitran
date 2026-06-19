package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkPaginationParse(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?page=3&page_size=25", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ParsePagination(req)
	}
}

func BenchmarkSecurityHeaders(b *testing.B) {
	handler := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}
