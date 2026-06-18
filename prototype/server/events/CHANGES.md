# Events Package - Wiring Guide

## main.go integration

```go
import "mitran/events"

// In main() or server setup:
eventBus := events.New()
events.RegisterRoutes(mux, eventBus)

// Publish events from anywhere with access to the bus:
eventBus.Publish(events.TaskCreated, map[string]string{"id": taskID, "name": taskName}, "engine")
```

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/events?limit=50` | Recent events (JSON array) |
| GET | `/api/v1/events/stream` | SSE stream of real-time events |

## SSE Format

```
event: task.created
data: {"type":"task.created","data":{...},"timestamp":"...","source":"engine"}
```
