# WebSocket + Memory Integration Changes

## go.mod — Add dependency

Add this line to the `require` block (or create one):

```
require github.com/gorilla/websocket v1.5.3
```

Then run:
```bash
go mod tidy
```

## main.go — Changes to apply

### 1. Add imports

Add these to the import block:

```go
"github.com/getmitran/mitran/server/memory"
"github.com/getmitran/mitran/server/websocket"
```

### 2. Create Hub and MemoryStore (after `store` and `cfg` initialization)

After:
```go
cfg := settings.Load(dataDir)
```

Add:
```go
// WebSocket hub
wsHub := websocket.NewHub()

// Memory store
memStore, err := memory.NewStore(dataDir)
if err != nil {
    log.Fatalf("Failed to initialize memory store: %v", err)
}
```

### 3. Register routes (after existing mux.Handle calls, before health check)

Add:
```go
// WebSocket
mux.HandleFunc("/api/v1/ws", websocket.HandleWS(wsHub))

// Memory
memHandler := &memory.Handler{Store: memStore}
mux.Handle("/api/v1/memory/", memHandler)
mux.Handle("/api/v1/memory/facts", memHandler)
mux.Handle("/api/v1/memory/search", memHandler)
mux.Handle("/api/v1/memory/episodes", memHandler)
mux.Handle("/api/v1/memory/corrections", memHandler)
```

### 4. Broadcast events on task/checkpoint changes

In the `TaskHandler`, after a successful task create, add a broadcast call. The cleanest approach is to pass the hub to handlers that need it.

**Option A (quick):** Make `wsHub` a package-level variable accessible from handlers.

**Option B (clean):** Add `Hub *websocket.Hub` field to `TaskHandler` and `CheckpointHandler`.

For Option B, change in main.go:
```go
mux.Handle("/api/v1/tasks", &handlers.TaskHandler{Store: store, Hub: wsHub})
mux.Handle("/api/v1/tasks/", &handlers.TaskHandler{Store: store, Hub: wsHub})
mux.Handle("/api/v1/checkpoints", &handlers.CheckpointHandler{Store: store, Hub: wsHub})
mux.Handle("/api/v1/checkpoints/", &handlers.CheckpointHandler{Store: store, Hub: wsHub})
```

Then in `handlers/tasks.go`, add the field and broadcast:
```go
type TaskHandler struct {
    Store *db.Store
    Hub   *websocket.Hub  // add this (import websocket package)
}
```

At the end of `create()`:
```go
if h.Hub != nil {
    h.Hub.Broadcast(websocket.Event{Type: "task.created", Data: t})
}
```

At the end of `assign()`:
```go
if h.Hub != nil {
    h.Hub.Broadcast(websocket.Event{Type: "task.updated", Data: map[string]string{"id": id, "status": "running"}})
}
```

In `handlers/checkpoints.go`:
```go
type CheckpointHandler struct {
    Store *db.Store
    Hub   *websocket.Hub
}
```

At the end of `decide()`:
```go
if h.Hub != nil {
    h.Hub.Broadcast(websocket.Event{Type: "task.checkpoint", Data: map[string]string{"id": id, "action": action}})
}
```

### 5. Update banner route list

Add to the printf block:
```go
fmt.Printf("    WS   /api/v1/ws             WebSocket events\n")
fmt.Printf("    GET  /api/v1/memory/facts    Semantic memory\n")
fmt.Printf("    GET  /api/v1/memory/search   Search episodes\n")
```

### 6. Update root JSON endpoint

Add `/api/v1/ws`, `/api/v1/memory` to the endpoints array in the root handler JSON response.
