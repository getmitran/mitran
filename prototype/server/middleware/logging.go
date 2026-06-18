package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getmitran/mitran/server/metrics"
)

var M *metrics.Metrics

type contextKey string

const RequestIDKey contextKey = "requestID"

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", reqID)
		sw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(sw, r)
		elapsed := time.Since(start)
		if M != nil {
			M.RecordRequest(elapsed)
			if sw.status >= 500 {
				M.RecordError()
			}
		}
		fmt.Printf("[%s] %s %s %s %d %s\n", start.Format(time.RFC3339), reqID, r.Method, r.URL.Path, sw.status, elapsed)
	})
}
