# Session Manager - Wiring Instructions

## Add to main.go

```go
import "mitran/session"  // adjust import path to match your module

// In main() or server setup:
sessionDir := filepath.Join(dataDir, "sessions")  // e.g. .mitran/data/sessions
sessionMgr, err := session.New(sessionDir)
if err != nil {
    log.Fatalf("failed to init session manager: %v", err)
}
session.RegisterRoutes(mux, sessionMgr)
```

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/sessions | Create session `{"user_id":"...", "agent":"..."}` |
| GET | /api/v1/sessions | List sessions (optional `?user_id=` filter) |
| GET | /api/v1/sessions/:id | Get session with messages |
| POST | /api/v1/sessions/:id/messages | Add message `{"role":"user","content":"..."}` |
| DELETE | /api/v1/sessions/:id | End session (marks inactive) |

## Storage

Sessions persist to `<dir>/sessions.json`. Survives page reloads and server restarts.
