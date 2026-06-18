package tracing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type contextKey struct{}

// Init initializes the tracing subsystem. Returns a shutdown function.
// TODO: Replace no-op with OTLP exporter (e.g. go.opentelemetry.io/otel/exporters/otlp/otlptrace).
func Init(serviceName string) (func(), error) {
	_ = serviceName
	return func() {}, nil
}

// TraceMiddleware generates a trace ID, sets it on the response, and stores it in context.
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 16)
		rand.Read(b)
		traceID := hex.EncodeToString(b)
		w.Header().Set("X-Trace-ID", traceID)
		ctx := context.WithValue(r.Context(), contextKey{}, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TraceIDFromContext extracts the trace ID from context.
func TraceIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(contextKey{}).(string); ok {
		return id
	}
	return ""
}
