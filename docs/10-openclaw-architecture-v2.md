# Mitran — Agent Runtime Architecture (v2 — Finalized)

## 1. Overview

Mitran's agent runtime is an invisible internal module that powers all AI agent behavior. It is **not** a separate product or brand — there is no "OpenClaw" visible to users. When users interact with Mitran agents, they interact with "Mitran agents," period.

The runtime provides:

- Conversational AI chat interface
- MCP (Model Context Protocol) server connectivity
- Persistent memory (phased: JSON → SQLite → PostgreSQL)
- Cron jobs & scheduled tasks
- Subagent orchestration
- Tool approval system
- Skills & knowledge bases
- Artifacts (versioned widgets/content)
- Session management
- Workspace isolation

## 2. Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                     MITRAN PLATFORM                                   │
│                                                                      │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │            MITRAN DASHBOARD (React 18 + Vite + Tailwind)       │  │
│  │                                                               │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │  │
│  │  │ Kanban   │ │ Agents   │ │Settings  │ │  Chat Interface  │ │  │
│  │  │ Board    │ │ Status   │ │  Page    │ │  (Agent Runtime) │ │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘ │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │  │
│  │  │ Memory   │ │ Cron     │ │Artifacts │ │  MCP Registry    │ │  │
│  │  │ Viewer   │ │ Manager  │ │ Library  │ │  & Tool Browser  │ │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘ │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                              │ HTTP/WebSocket                        │
│  ┌───────────────────────────▼───────────────────────────────────┐  │
│  │              MITRAN CORE ENGINE (Go)                           │  │
│  │                                                               │  │
│  │  ┌─────────────┐  ┌──────────────┐  ┌─────────────────────┐  │  │
│  │  │ DAG Engine  │  │ Priority     │  │ Agent Runtime        │  │  │
│  │  │ + Scheduler │  │ Queue        │  │                     │  │  │
│  │  └─────────────┘  └──────────────┘  │ • Session Manager  │  │  │
│  │                                      │ • Memory Router    │  │  │
│  │  ┌─────────────┐  ┌──────────────┐  │ • Cron Scheduler   │  │  │
│  │  │ Project     │  │ Checkpoint   │  │ • Tool Approval    │  │  │
│  │  │ Manager     │  │ Manager      │  │ • Artifact Store   │  │  │
│  │  └─────────────┘  └──────────────┘  │ • Skill Loader     │  │  │
│  │                                      │ • Subagent Pool    │  │  │
│  │  ┌─────────────┐  ┌──────────────┐  │ • Worker Lifecycle │  │  │
│  │  │ Settings    │  │ MCP Client   │  └─────────────────────┘  │  │
│  │  │ Manager     │  │ Connector    │                            │  │
│  │  └─────────────┘  └──────────────┘                            │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                    │ gRPC (protobuf)        │ stdio/SSE              │
│  ┌─────────────────▼──────────────────┐    │                        │
│  │      PYTHON AGENT WORKERS          │    │                        │
│  │                                    │    │                        │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ │    │                        │
│  │  │  Dev   │ │  Docs  │ │  Ops   │ │    │                        │
│  │  │ Worker │ │ Worker │ │ Worker │ │    │                        │
│  │  └────────┘ └────────┘ └────────┘ │    │                        │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ │    │                        │
│  │  │ CI/CD  │ │ Review │ │Custom  │ │    │                        │
│  │  │ Worker │ │ Worker │ │Workers │ │    │                        │
│  │  └────────┘ └────────┘ └────────┘ │    │                        │
│  └────────────────────────────────────┘    │                        │
│                                             │                        │
│  ┌──────────────────────────────────────────▼────────────────────┐  │
│  │              MCP SERVER LAYER (User-Managed)                   │  │
│  │                                                               │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │  │
│  │  │ Slack    │ │ AWS      │ │Playwright│ │  Custom MCP      │ │  │
│  │  │ MCP      │ │ MCP      │ │ MCP      │ │  (user-defined)  │ │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘ │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │  │
│  │  │ GitHub   │ │ Email    │ │ Calendar │ │  File System     │ │  │
│  │  │ MCP      │ │ MCP      │ │ MCP      │ │  MCP             │ │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘ │  │
│  └───────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
```

## 3. Core Separation of Concerns

### Go Engine (Orchestration Only)
- DAG scheduling, task dispatch, priority queuing
- Python worker lifecycle management (start, health-check, restart)
- Session routing, checkpoint persistence
- HTTP/WebSocket API for dashboard
- MCP client connections (discovers and connects to user-started servers)
- gRPC server for worker communication
- **No LLM SDK. No direct model calls. Zero AI logic.**

### Python Workers (All LLM Interaction)
- All LLM API calls (Bedrock, OpenAI, Anthropic, local models)
- Prompt construction, context assembly, response parsing
- Tool call execution (via MCP protocol)
- Memory read/write operations
- Skill loading and application
- **Managed as subprocesses by Go engine**

## 4. Inter-Process Communication (gRPC + Protobuf)

### 4.1 Service Contract

```protobuf
syntax = "proto3";
package mitran.worker.v1;

// Go Engine → Python Worker
service AgentWorker {
  // Execute a task with full context
  rpc ExecuteTask(TaskRequest) returns (stream TaskProgress);
  // Health check
  rpc HealthCheck(HealthRequest) returns (HealthResponse);
  // Graceful shutdown
  rpc Shutdown(ShutdownRequest) returns (ShutdownResponse);
}

message TaskRequest {
  string task_id = 1;
  string agent_name = 2;
  string system_prompt = 3;
  repeated Message conversation_history = 4;
  repeated string skill_contents = 5;
  repeated string available_tools = 6;
  map<string, string> memory_context = 7;
  TaskConfig config = 8;
}

message TaskProgress {
  oneof update {
    string text_chunk = 1;        // Streaming text
    ToolCall tool_request = 2;    // Tool approval needed
    ToolResult tool_result = 3;   // Tool executed
    TaskComplete completion = 4;  // Done
    TaskError error = 5;          // Failed
  }
}

message ToolCall {
  string call_id = 1;
  string server_name = 2;
  string tool_name = 3;
  string arguments_json = 4;
}

message ToolResult {
  string call_id = 1;
  string result_json = 2;
  bool is_error = 3;
}

message TaskComplete {
  string final_response = 1;
  repeated MemoryUpdate memory_updates = 2;
}

message TaskConfig {
  string model = 1;
  string provider = 2;          // bedrock, openai, anthropic, local
  float temperature = 3;
  int32 max_tokens = 4;
  string approval_mode = 5;     // interactive, reads, yolo
}

message Message {
  string role = 1;              // user, assistant, system
  string content = 2;
}

message MemoryUpdate {
  string scope = 1;             // global, project, agent, session
  string key = 2;
  string value = 3;
  string operation = 4;         // set, delete, append
}

message HealthRequest {}
message HealthResponse {
  bool healthy = 1;
  int32 active_tasks = 2;
  int64 uptime_seconds = 3;
}

message ShutdownRequest { int32 timeout_seconds = 1; }
message ShutdownResponse { bool clean = 1; }
```

### 4.2 Worker Lifecycle (Go manages Python)

```
mitran serve
  │
  ├─ Start Go Engine (HTTP :7780, gRPC :7781)
  │
  ├─ For each configured agent worker:
  │    ├─ Spawn: python -m mitran_worker --port <assigned> --agent <name>
  │    ├─ Wait for health check response (timeout: 10s)
  │    ├─ Register in worker pool
  │    └─ Start health check ticker (every 30s)
  │
  ├─ On worker health check failure:
  │    ├─ Log warning
  │    ├─ Attempt restart (max 3 retries with exponential backoff)
  │    └─ If exhausted: mark agent unavailable, notify dashboard
  │
  └─ On shutdown signal (SIGTERM/SIGINT):
       ├─ Send Shutdown RPC to all workers (timeout: 30s)
       ├─ Wait for in-flight tasks to drain
       └─ Force-kill remaining after timeout
```

## 5. Memory Architecture (Phased)

### Phase 1 — v0.1 (JSON Files)

```
.mitran/projects/<project-id>/memory/
├── semantic.json          # Key-value facts
├── lessons.json           # Learned corrections
├── history/
│   └── 2026-06-18.md      # Daily activity log
└── episodes/
    └── <session-id>.json  # Conversation fragments
```

- No external dependencies
- Simple read/write, grep-searchable
- Sufficient for single-user, <1000 entries

### Phase 2 — v0.2 (SQLite + FTS5)

```sql
-- Full-text search over memory
CREATE VIRTUAL TABLE memory_fts USING fts5(
  key, value, scope, content=memory
);

-- Structured storage
CREATE TABLE memory (
  id TEXT PRIMARY KEY,
  scope TEXT NOT NULL,        -- global|project|agent|session
  key TEXT NOT NULL,
  value TEXT NOT NULL,
  embedding BLOB,            -- via sqlite-vss extension
  created_at TEXT,
  updated_at TEXT
);

CREATE TABLE episodes (
  id TEXT PRIMARY KEY,
  session_id TEXT,
  content TEXT,
  embedding BLOB,
  relevance_score REAL,
  created_at TEXT
);
```

- FTS5 for fast text search
- sqlite-vss for vector similarity (optional)
- Zero-config, embedded, single-file

### Phase 3 — v1.0 (PostgreSQL + pgvector OR ClickHouse)

```sql
-- PostgreSQL + pgvector for production
CREATE EXTENSION vector;

CREATE TABLE memory (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  scope TEXT NOT NULL,
  key TEXT NOT NULL,
  value TEXT NOT NULL,
  embedding vector(384),
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ON memory USING ivfflat (embedding vector_cosine_ops);
```

- Multi-user concurrent access
- Scalable to millions of entries
- Native vector search for semantic retrieval
- OR ClickHouse for analytics-heavy workloads with large history

### Memory Scoping (All Phases)

| Scope | Visibility |
|-------|-----------|
| Global | All agents, all projects |
| Project | All agents within one project |
| Agent | Only that specific agent |
| Session | Only within one conversation |

## 6. MCP Server Lifecycle

### Key Principle: User Starts, Mitran Connects

MCP servers are **not** auto-started by Mitran. Users start them independently (via Docker, systemd, pm2, or manually). Mitran discovers and connects to running servers via configuration.

### 6.1 Configuration

```json
// .mitran/mcp-config.json
{
  "servers": {
    "slack": {
      "transport": "stdio",
      "command": "npx",
      "args": ["-y", "@mitran/mcp-slack"],
      "env": {
        "SLACK_BOT_TOKEN": "${SLACK_BOT_TOKEN}"
      }
    },
    "github": {
      "transport": "sse",
      "url": "http://localhost:3001/mcp",
      "headers": {
        "Authorization": "Bearer ${GITHUB_TOKEN}"
      }
    },
    "internal-api": {
      "transport": "stdio",
      "command": "python",
      "args": ["./mcp-servers/internal-api/server.py"]
    }
  }
}
```

### 6.2 Connection Flow

```
1. User starts MCP servers (their responsibility)
2. mitran serve reads mcp-config.json
3. For each server:
   - stdio: Go engine spawns the process, communicates over stdin/stdout
   - sse: Go engine connects to the running HTTP endpoint
4. Engine discovers available tools via MCP `tools/list`
5. Tools are registered in the tool registry for agent access
6. Dashboard MCP Registry page shows connection status
```

### 6.3 Health & Reconnection

- Go engine pings each MCP connection every 60s
- On disconnect: attempt reconnect with exponential backoff (3 retries)
- On permanent failure: mark server as unavailable, dashboard shows warning
- Users can manually reconnect via Settings > MCP Registry > Reconnect

## 7. Secrets Management

### Now (v0.1): Environment Variables

```bash
# .env or shell exports
export MITRAN_LLM_API_KEY="sk-..."
export SLACK_BOT_TOKEN="xoxb-..."
export GITHUB_TOKEN="ghp_..."
export AWS_PROFILE="default"

# Mitran reads from environment at startup
mitran serve
```

- Simple, works everywhere
- Docker/K8s compatible (inject via secrets)
- `.env` file support with `dotenv` loading

### Later (v1.0): AWS Secrets Manager Integration

```json
// .mitran/secrets-config.json
{
  "provider": "aws-secrets-manager",
  "region": "us-east-1",
  "prefix": "mitran/prod/",
  "secrets": {
    "SLACK_BOT_TOKEN": "mitran/prod/slack-bot-token",
    "GITHUB_TOKEN": "mitran/prod/github-token",
    "LLM_API_KEY": "mitran/prod/llm-api-key"
  },
  "refresh_interval_seconds": 3600
}
```

- Automatic rotation support
- Audit trail via CloudTrail
- Fallback to env vars if Secrets Manager unavailable

## 8. Dashboard Stack

**React 18 + Vite + Tailwind CSS + TypeScript**

Full dependency list (matching MeshClaw's proven stack):

| Package | Purpose |
|---------|---------|
| `react` + `react-dom` | UI framework |
| `vite` | Build tool + dev server |
| `tailwindcss` | Utility-first CSS |
| `typescript` | Type safety |
| `react-router-dom` | Client-side routing |
| `@tanstack/react-query` | Server state management |
| `lucide-react` | Icon library |
| `react-redux` + `@reduxjs/toolkit` | Global state (sessions, preferences) |
| `@monaco-editor/react` | Code editing (agent prompts, skills) |
| `react-markdown` | Markdown rendering (chat, artifacts) |
| `react-virtuoso` | Virtualized lists (large chat history) |

### Dashboard Pages

| Page | Purpose |
|------|---------|
| **Chat** | Conversational interface — talk to any agent |
| **Kanban** | Task management board |
| **Agents** | Status, config, skill editing |
| **Memory** | View/search semantic + episodic memory |
| **Cron** | Manage scheduled jobs |
| **Artifacts** | Browse versioned widgets/content |
| **MCP Registry** | View connected servers + available tools |
| **Sessions** | View/resume past conversations |
| **Approval Queue** | Pending tool approvals |
| **Security Log** | Audit trail of all agent actions |
| **Settings** | Project config, MCP servers, secrets |

## 9. Agent Definition Format

```json
// agents/dev-agent/agent-spec.json
{
  "name": "dev-agent",
  "displayName": "Dev Agent",
  "description": "Writes code, creates PRs, runs tests",
  "model": "claude-sonnet-4",
  "provider": "bedrock",
  "systemPrompt": "./SYSTEM.md",
  "skills": ["./skills/"],
  "tools": {
    "allowed": ["filesystem", "shell", "github"],
    "approval_mode": "reads"
  },
  "workspace": "project",
  "memory": {
    "enabled": true,
    "scope": "project"
  }
}
```

## 10. Deployment Model

### `mitran serve` — Single Command Startup

```
$ mitran serve

  ╔══════════════════════════════════════╗
  ║       MITRAN v0.1.0                  ║
  ╠══════════════════════════════════════╣
  ║  Engine:    http://localhost:7780     ║
  ║  gRPC:      localhost:7781           ║
  ║  Dashboard: http://localhost:7780    ║
  ╠══════════════════════════════════════╣
  ║  Workers:   8 agents ready           ║
  ║  MCP:       5 servers connected      ║
  ║  Memory:    JSON (phase 1)           ║
  ╚══════════════════════════════════════╝

  Routes:
    GET  /api/v1/health
    POST /api/v1/chat
    GET  /api/v1/agents
    ...
```

### Process Tree

```
mitran serve (Go, PID 1)
├── Python Worker: dev-agent (PID 2)
├── Python Worker: docs-agent (PID 3)
├── Python Worker: ops-agent (PID 4)
├── Python Worker: cicd-agent (PID 5)
├── Python Worker: review-agent (PID 6)
├── Python Worker: hr-agent (PID 7)
├── Python Worker: tickets-agent (PID 8)
└── Python Worker: wiki-agent (PID 9)

User-managed (separate processes, NOT started by Mitran):
├── MCP: slack-server (stdio, connected)
├── MCP: github-server (SSE http://localhost:3001)
├── MCP: filesystem (stdio, connected)
└── MCP: browser (stdio, connected)
```

### Health Checks

Go engine monitors each Python worker:
- gRPC `HealthCheck` RPC every 30 seconds
- If 3 consecutive failures: restart worker subprocess
- Max 3 restart attempts with exponential backoff (5s, 15s, 45s)
- After exhaustion: mark agent as `unavailable`, notify via dashboard + Slack

## 11. What's Excluded (Amazon-Internal / Not Applicable)

| Feature | Why Excluded |
|---------|-------------|
| Brazil build system | Amazon proprietary |
| CRUX code reviews | Amazon proprietary |
| Apollo deployments | Amazon proprietary |
| Pipelines (Amazon) | Amazon proprietary (replaced by Mitran CI/CD) |
| AIM (Agent Install Manager) | Amazon proprietary |
| Midway/Kerberos auth | Amazon SSO (replaced by standard OAuth/SAML) |
| ReadInternalWebsites | Amazon internal tool |
| ARCC security governance | Amazon policy (replaced by generic policy engine) |
| Conduit credentials | Amazon credential tool |
| PhoneTool | Amazon employee directory |

## 12. Implementation Priority

| Phase | Components | Timeline |
|-------|-----------|----------|
| **Phase 1** (Month 1-2) | Go engine + gRPC + Python worker scaffold + Chat UI + MCP connector | Core wiring |
| **Phase 2** (Month 2-3) | Memory (JSON→SQLite) + Sessions + Learned corrections | Agent intelligence |
| **Phase 3** (Month 3-4) | Cron scheduler + Artifacts + Skills + Subagents | Automation |
| **Phase 4** (Month 4-5) | Tool approval UI + Security audit + Dashboard polish | Production-ready |
| **Phase 5** (Month 5-6) | Multi-user auth + PostgreSQL memory + Secrets Manager | Enterprise features |

**Team:** 20 engineers. Full scope, no cuts. 6-month delivery.

## 13. Finalized Decisions

These 7 architecture choices are **LOCKED** — not open for debate.

| # | Domain | Decision | Rationale |
|---|--------|----------|-----------|
| 1 | **LLM** | Python workers handle ALL LLM interaction. Go engine only orchestrates/dispatches. No Go LLM SDK. | Python has the richest LLM ecosystem (langchain, anthropic SDK, litellm). Go is for fast I/O, not prompt engineering. |
| 2 | **IPC** | gRPC with protobuf between Go engine and Python workers. Typed service contracts. | Type-safe, high-performance, streaming support, code generation for both languages. |
| 3 | **Memory** | Phased: JSON files (v0.1) → SQLite+FTS5 (v0.2) → PostgreSQL+pgvector or ClickHouse (v1.0) | Start simple (zero deps), scale when needed. Each phase is a clean migration, not a rewrite. |
| 4 | **Deployment** | Multi-process. `mitran serve` starts Go engine which spawns/manages Python workers. MCP servers started separately by user. | Go manages lifecycle with health checks + auto-restart. Clean separation. No monolith. |
| 5 | **Secrets** | Environment variables now. AWS Secrets Manager integration later. | Universally supported. Docker/K8s native. No vendor lock-in at v0.1. |
| 6 | **MCP Lifecycle** | User starts MCP servers independently. Mitran discovers and connects via config. No auto-start. | Keeps Mitran lightweight. Users control their own infra. No magic. |
| 7 | **Dashboard** | React 18 + Vite + Tailwind + TS + react-router-dom + @tanstack/react-query + lucide-react + react-redux + @monaco-editor/react + react-markdown + react-virtuoso | Proven stack (same as MeshClaw). Fast builds, modern DX, extensive ecosystem. |

### Reviewer Feedback Addressed

| Reviewer | Concern | Resolution |
|----------|---------|------------|
| **CTO** | "Who supervises the Python workers? What if one hangs?" | Go engine manages full worker lifecycle: spawn, health-check (30s interval), auto-restart (3x exponential backoff), graceful shutdown via gRPC `Shutdown` RPC. |
| **CEO** | "OpenClaw branding confuses users. Keep it invisible." | Done. No "OpenClaw" anywhere in UI, docs, or user-facing materials. It's just "how Mitran agents work." |
| **Critic** | "Scope is massive. 6 months with what team?" | 20 engineers. Phased delivery (see §12). Each phase is independently shippable. No scope cuts needed. |
