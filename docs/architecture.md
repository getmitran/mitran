# Mitran Architecture

## System Overview

Mitran is a multi-process platform with three core components communicating over HTTP/gRPC/WebSocket:

```
┌─────────────────────────────────────────────────────────────────┐
│                        React Dashboard (:5173)                   │
│  Vite + Tailwind + TypeScript + React Router + TanStack Query   │
└───────────────────────────────┬─────────────────────────────────┘
                                │ HTTP REST + WebSocket
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Go Engine (:7780)                           │
│  DAG Scheduler │ Priority Queue │ Resource Lock │ API Gateway   │
└──────┬────────────────────┬─────────────────────────────────────┘
       │ gRPC (protobuf)    │ HTTP callbacks
       ▼                    ▼
┌──────────────────┐  ┌──────────────────────────────────────────┐
│  Python Workers  │  │           OpenClaw Runtime                │
│  (Agent Pool)    │  │  Memory │ Skills │ Crons │ Artifacts     │
│  Claude Bedrock  │  │  MCP │ Snapshots │ Process Supervisor    │
└──────────────────┘  └──────────────────────────────────────────┘
```

## Component Diagram

```
                    ┌─────────┐
                    │  User   │
                    └────┬────┘
                         │ HTTPS
                    ┌────▼────┐
                    │ Reverse │
                    │  Proxy  │
                    └────┬────┘
              ┌──────────┼──────────┐
              │          │          │
         ┌────▼───┐ ┌───▼────┐ ┌──▼───────┐
         │Dashboard│ │  API   │ │WebSocket │
         │  SPA   │ │Gateway │ │  Server  │
         └────────┘ └───┬────┘ └──┬───────┘
                        │         │
                   ┌────▼─────────▼────┐
                   │    Go Engine      │
                   │  ┌─────────────┐  │
                   │  │ DAG Sched.  │  │
                   │  │ Task Queue  │  │
                   │  │ State Mgr   │  │
                   │  └─────────────┘  │
                   └────────┬──────────┘
                   ┌────────┼────────┐
                   │        │        │
              ┌────▼──┐ ┌──▼───┐ ┌──▼──────┐
              │Worker │ │Worker│ │OpenClaw  │
              │ (LLM) │ │(LLM) │ │ Runtime  │
              └───────┘ └──────┘ └──────────┘
```

## Data Flow

### HTTP REST (Engine ↔ Dashboard)
- `GET/POST /api/v1/tasks` — task CRUD
- `GET/POST /api/v1/agents` — agent management
- `GET/PUT /api/v1/settings` — project/module configuration
- `GET /api/v1/metrics` — Prometheus-compatible metrics
- `POST /api/v1/webhooks/{provider}` — GitHub/Slack inbound

### WebSocket (Engine → Dashboard)
- Real-time task status updates
- Agent execution logs streaming
- Approval request notifications
- System health heartbeats

### gRPC (Engine ↔ Python Workers)
- `ExecuteTask` — dispatch work to agent pool
- `StreamOutput` — bidirectional log streaming
- `ReportStatus` — agent health/completion signals

Proto definition: `prototype/proto/mitran.proto`

## OpenClaw Runtime

Invisible internal module providing autonomous agent infrastructure:

### Vector Memory
- **Short-term**: JSON file per session (fast, ephemeral)
- **Mid-term**: SQLite + FTS5 for keyword search
- **Long-term**: ChromaDB for semantic similarity (pgvector in enterprise)
- Phased migration: JSON → SQLite → Postgres/ClickHouse

### Skill Loader
- Skills are markdown files describing capabilities
- Loaded into agent context at session start
- Hot-reload on file change without process restart
- Per-workspace skill isolation

### Tool Approval
- Three modes: `interactive` (prompt every call), `reads` (auto-approve reads), `yolo` (auto-approve all)
- Configurable per-agent and per-workspace
- Audit log of all tool invocations with timestamps

### Artifact Store
- Versioned HTML/markdown/SVG/JSON blobs
- Max 50 versions per artifact (configurable)
- Slug-based addressing for stable URLs
- File-backed or chat-backed storage modes

### LLM Abstraction
- Provider interface: AWS Bedrock (primary), OpenAI, Anthropic direct
- Model selection per agent/task
- Token budget tracking and automatic context compaction
- Retry with exponential backoff on transient failures

### MCP Lifecycle
- User starts MCP servers independently
- Mitran discovers and connects via stdio/SSE transport
- Health monitoring with automatic reconnection
- Tool schema introspection at connect time

### Process Supervisor
- `mitran` CLI manages all processes (engine, workers, dashboard)
- Graceful shutdown with drain timeout
- Crash restart with backoff
- PID file management under `.mitran/`

### Snapshot/Restore
- Portable tar.gz of full runtime state
- Components: memory, crons, artifacts, config, lessons
- Merge or replace restore modes
- Automatic daily snapshots (configurable retention)

## Agent Types

| Agent | Capabilities |
|-------|-------------|
| **Dev** | Code generation, PR creation, build verification, test writing |
| **Docs** | Technical writing, API docs, README updates, changelog |
| **Ops** | Grafana dashboards, alerting rules, runbook generation |
| **Review** | Security scanning (semgrep/bandit), code review, style checks |
| **CI/CD** | Pipeline config, GitHub Actions, multi-env deployments |
| **Tickets** | SLA enforcement, sprint planning, status tracking |
| **Wiki** | Knowledge base, incident runbooks, onboarding guides |
| **HR** | PTO workflows, onboarding checklists, team roster |

All agents share: LLM access (Claude Sonnet via Bedrock), filesystem I/O, MCP tool invocation, memory read/write.

## Storage Architecture

```
.mitran/
├── projects/
│   └── <project-id>/
│       ├── settings.json          # Project config
│       ├── memory/
│       │   ├── sessions/          # Per-session JSON (short-term)
│       │   ├── index.db           # SQLite + FTS5 (mid-term)
│       │   └── vectors/           # ChromaDB collection (long-term)
│       ├── artifacts/
│       │   └── <slug>/
│       │       ├── current.html   # Latest version
│       │       └── versions/      # Historical snapshots
│       ├── crons.json             # Scheduled jobs
│       └── lessons.json           # Learned corrections
├── config.json                    # Global settings
└── snapshots/                     # Backup archives
```

- **SQLite**: Task metadata, agent state, execution history, FTS5 indexes
- **ChromaDB**: Embedding vectors for semantic memory search
- **Filesystem**: Artifacts, snapshots, skill files, proto definitions

## Deployment Topology

### Development (Single Machine)
```
mitran start        # Launches all processes
  ├── Go Engine     (port 7780)
  ├── Python Workers (port 8888, pool of N)
  └── Dashboard     (port 5173, Vite dev server)
```

### Production (Docker Compose)
```yaml
services:
  engine:      # Go binary, stateless
  worker:      # Python, scales horizontally (replicas: N)
  dashboard:   # Nginx serving built React SPA
  postgres:    # Replaces SQLite for metadata
  chromadb:    # Vector store
  prometheus:  # Metrics collection
  grafana:     # Dashboards
```

### Enterprise (Kubernetes)
- Engine: Deployment with HPA on CPU/request count
- Workers: Deployment with HPA on queue depth
- Dashboard: Static file serving via CDN
- Storage: Managed Postgres (RDS) + managed vector DB
- Secrets: AWS Secrets Manager (env vars in dev)

## API Gateway Pattern

The Go engine acts as the unified API gateway:

1. **Authentication**: Token-based (JWT) with configurable providers
2. **Rate limiting**: Per-client token bucket (in-memory, Redis for distributed)
3. **Request routing**: Path-prefix dispatch to internal handlers
4. **WebSocket upgrade**: `/ws` endpoint for real-time streaming
5. **Middleware chain**: Logging → Auth → RateLimit → CORS → Handler
6. **Health endpoints**: `/health` (liveness), `/ready` (readiness with dependency checks)
7. **API versioning**: `/api/v1/` prefix, additive-only changes within version

All external traffic enters through the engine. Workers are never exposed directly — they communicate exclusively via gRPC on the internal network.
