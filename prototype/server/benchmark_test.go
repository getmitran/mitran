package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getmitran/mitran/server/middleware"
	"github.com/getmitran/mitran/server/queue"
)

func BenchmarkTaskQueueEnqueue(b *testing.B) {
	q := queue.NewTaskQueue()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			q.Enqueue(&queue.Task{ID: j, Priority: j % 5})
		}
	}
}

func BenchmarkPaginationParse(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?page=3&per_page=25", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queue.ParsePagination(req)
	}
}

func BenchmarkSecurityHeaders(b *testing.B) {
	handler := middleware.SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}
