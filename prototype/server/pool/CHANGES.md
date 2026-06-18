# Subagent Pool

Parallel agent dispatch with configurable concurrency limits.

## Wiring into main.go

```go
import "mitran/pool"

// In main() after mux creation:
agentPool := pool.New(8) // 8 concurrent workers
agentPool.Start()
defer agentPool.Stop()

pool.RegisterHandlers(mux, agentPool)
```

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/pool/stats | Worker count + queued/running/done/failed |
| GET | /api/v1/pool/items | List all work items with status |
| POST | /api/v1/pool/submit | Submit `{id, agent, payload}` |

## Next Steps

- Wire `Submit` handler to call actual agent worker via gRPC/HTTP
- Add per-agent concurrency limits (rate limiting)
- Add item cancellation endpoint
- Add WebSocket notifications on item completion
