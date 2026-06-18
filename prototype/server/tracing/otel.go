package tracing

import (
	"log"
	"os"
)

// InitOTLP checks for OTLP exporter configuration and logs readiness.
// Actual OpenTelemetry SDK integration deferred to reduce dependencies.
func InitOTLP(serviceName string) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint != "" {
		log.Printf("tracing: OTLP exporter configured: %s (service=%s)", endpoint, serviceName)
	} else {
		log.Printf("tracing: no-op mode (service=%s)", serviceName)
	}
}
