package tracing

import (
	"log"
	"os"
)

// OTLPConfigured indicates whether the OTLP exporter endpoint is set.
// Middleware can use this to conditionally add trace context propagation headers.
var OTLPConfigured bool

// InitOTLP checks for OTLP exporter configuration and logs readiness.
// Actual OpenTelemetry SDK integration deferred to reduce dependencies.
func InitOTLP(serviceName string) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint != "" {
		log.Printf("tracing: OTLP endpoint detected: %s (service=%s). Full OTLP integration requires the otel SDK dependency.", endpoint, serviceName)
		OTLPConfigured = true
	} else {
		log.Printf("tracing: no-op mode (service=%s)", serviceName)
	}
}
