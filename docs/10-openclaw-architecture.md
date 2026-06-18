> **SUPERSEDED**: See [Architecture v2](10-openclaw-architecture-v2.md) for the current design.

# Mitran — OpenClaw Integration Architecture

## 1. Overview

**OpenClaw** is Mitran's built-in agent management layer — an open-source clone of the MeshClaw architecture, stripped of all Amazon-internal dependencies. It provides:

- Conversational AI chat interface
- MCP (Model Context Protocol) server management
- Persistent memory (vector-store based)
- Cron jobs & scheduled tasks
- Subagent orchestration
- Tool approval system
- Skills & knowledge bases
- Artifacts (versioned widgets/content)
- Session management
- Workspace isolation

OpenClaw is NOT a separate product — it's a **module within Mitran** that powers the agent runtime. Think of it as "Mitran's brain."

## 2. Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                     MITRAN PLATFORM                                   │
│                                                                      │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │                    MITRAN DASHBOARD (React)                    │  │
│  │                                                               │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │  │
│  │  │ Kanban   │ │ Agents   │ │Settings  │ │  Chat (OpenClaw) │ │  │
│  │  │ Board    │ │ Status   │ │  Page    │ │  Interface       │ │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘ │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │  │
│  │  │ Memory   │ │ Cron     │ │Artifacts │ │  MCP Registry    │ │  │
│  │  │ Viewer   │ │ Manager  │ │ Library  │ │  & Tool Browser  │ │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘ │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                              │ HTTP/WS                               │
│  ┌───────────────────────────▼───────────────────────────────────┐  │
│  │                  MITRAN CORE ENGINE (Go)                       │  │
│  │                                                               │  │
│  │  ┌─────────────┐  ┌──────────────┐  ┌─────────────────────┐  │  │
│  │  │ DAG Engine  │  │ Priority     │  │ OpenClaw Runtime     │  │  │
│  │  │ + Scheduler │  │ Queue        │  │                     │  │  │
│  │  └─────────────┘  └──────────────┘  │ • Session Manager  │  │  │
│  │                                      │ • Memory (Vector)  │  │  │
│  │  ┌─────────────┐  ┌──────────────┐  │ • Cron Scheduler   │  │  │
│  │  │ Project     │  │ Checkpoint   │  │ • Tool Approval    │  │  │
│  │  │ Manager     │  │ Manager      │  │ • Artifact Store   │  │  │
│  │  └─────────────┘  └──────────────┘  │ • Skill Loader     │  │  │
│  │                                      │ • Subagent Pool    │  │  │
│  │  ┌─────────────┐  ┌──────────────┐  └─────────────────────┘  │  │
│  │  │ Settings    │  │ MCP Server   │                            │  │
│  │  │ Manager     │  │ Registry     │                            │  │
│  │  └─────────────┘  └──────────────┘                            │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                              │ gRPC/HTTP                             │
│  ┌───────────────────────────▼───────────────────────────────────┐  │
│  │                    MCP SERVER LAYER                            │  │
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
│                              │                                       │
│  ┌───────────────────────────▼───────────────────────────────────┐  │
│  │                    AGENT WORKER LAYER (Python)                 │  │
│  │                                                               │  │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────────┐ │  │
│  │  │  Dev   │ │  Docs  │ │  Ops   │ │  CI/CD │ │  Custom    │ │  │
│  │  │ Agent  │ │ Agent  │ │ Agent  │ │ Agent  │ │  Agents    │ │  │
│  │  └────────┘ └────────┘ └────────┘ └────────┘ └────────────┘ │  │
│  └───────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
```

## 3. OpenClaw Module Design

### 3.1 What OpenClaw Provides (MeshClaw Features → Mitran)

| MeshClaw Feature | OpenClaw Equivalent | Implementation |
|-----------------|-------------------|----------------|
| Chat interface (conversational agent) | ✅ Chat tab in Mitran Dashboard | React component + WebSocket |
| Session management | ✅ Persistent sessions per user | Go session store |
| Memory (semantic + episodic) | ✅ Vector-store memory | Embedded Qdrant or ChromaDB |
| Learned corrections | ✅ Same format | Stored in memory subsystem |
| Cron jobs | ✅ Same API | Go scheduler (already exists in engine) |
| Subagent orchestration | ✅ spawn_run equivalent | Go worker pool |
| Tool approval (interactive/reads/yolo) | ✅ Same 3-tier model | Per-agent config |
| Artifacts (versioned widgets) | ✅ Same concept | SQLite/JSON artifact store |
| Skills (SKILL.md files) | ✅ Per-agent skill loading | Filesystem-based |
| Workspaces | ✅ Project-scoped | Maps to Mitran projects |
| Apps ecosystem | ✅ Same manifest format | App registry |
| Slack integration | ✅ Bot + slash commands | Slack MCP server |
| CLI (chat, tui, commands) | ✅ `mitran chat`, `mitran tui` | Go CLI (Cobra) |
| Dashboard (web UI) | ✅ Merged into Mitran dashboard | React |
| Snapshot/restore | ✅ Project-level backup | Go implementation |
| Security audit | ✅ Action logging | Audit log table |

### 3.2 What's EXCLUDED (Amazon-Internal)

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

## 4. MCP Server Registry

### 4.1 How MCP Servers Work in Mitran

When a company sets up Mitran, they configure which MCP servers are available — similar to writing `mcp.json` for Kiro/Cline, but managed via the Mitran Settings UI.

```json
// .mitran/mcp-config.json (equivalent to mcp.json)
{
  "servers": {
    "slack": {
      "command": "npx",
      "args": ["-y", "@mitran/mcp-slack"],
      "env": {
        "SLACK_BOT_TOKEN": "${SLACK_BOT_TOKEN}",
        "SLACK_TEAM_ID": "${SLACK_TEAM_ID}"
      }
    },
    "github": {
      "command": "npx",
      "args": ["-y", "@mitran/mcp-github"],
      "env": {
        "GITHUB_TOKEN": "${GITHUB_TOKEN}"
      }
    },
    "aws": {
      "command": "npx",
      "args": ["-y", "@mitran/mcp-aws"],
      "env": {
        "AWS_PROFILE": "${AWS_PROFILE}"
      }
    },
    "browser": {
      "command": "npx",
      "args": ["-y", "@playwright/mcp-server"]
    },
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@mitran/mcp-filesystem"],
      "env": {
        "ALLOWED_DIRS": "${PROJECT_WORKSPACE}"
      }
    },
    "custom-internal-api": {
      "command": "python",
      "args": ["./mcp-servers/internal-api/server.py"],
      "env": {
        "API_BASE_URL": "https://api.company.internal",
        "API_KEY": "${INTERNAL_API_KEY}"
      }
    }
  }
}
```

### 4.2 MCP Registry Dashboard UI

The Settings page includes an "MCP Servers" section where users can:
- Add/remove MCP server configurations
- Test connectivity (health check)
- View available tools per server
- Set which agents can access which servers
- Configure environment variables (secrets stored encrypted)

### 4.3 Built-in MCP Servers

| Server | Purpose | Bundled? |
|--------|---------|----------|
| `@mitran/mcp-filesystem` | Read/write files in project workspace | Yes |
| `@mitran/mcp-shell` | Execute shell commands (with approval) | Yes |
| `@mitran/mcp-slack` | Slack API (messages, channels, reactions) | Optional |
| `@mitran/mcp-github` | GitHub API (repos, PRs, Actions) | Optional |
| `@mitran/mcp-aws` | AWS CLI wrapper | Optional |
| `@mitran/mcp-browser` | Playwright browser automation | Optional |
| `@mitran/mcp-email` | SMTP/IMAP email | Optional |
| `@mitran/mcp-calendar` | Google/Outlook calendar | Optional |
| `@mitran/mcp-database` | SQL query execution | Optional |

## 5. Agent Definition Format

Same as MeshClaw's agent-spec pattern:

```json
// agents/dev-agent/agent-spec.json
{
  "name": "dev-agent",
  "displayName": "Dev Agent",
  "description": "Writes code, creates PRs, runs tests",
  "model": "claude-sonnet-4",
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

```markdown
<!-- agents/dev-agent/SYSTEM.md -->
You are a senior software engineer working within the Mitran platform.
You write production-quality code, tests, and documentation.
...
```

```markdown
<!-- agents/dev-agent/skills/code-review/SKILL.md -->
# Code Review Skill
When reviewing code, check for:
- Security issues
- Performance problems
- Test coverage
...
```

## 6. Memory Architecture

### 6.1 Vector Store

OpenClaw's memory system uses an **embedded vector database** (no external service needed):

```
Memory Store
├── Semantic Memory (key-value facts)
│   └── Stored as plain JSON (fast lookup by key)
├── Episodic Memory (conversation fragments)
│   └── Embedded via sentence-transformers, stored in vector DB
├── Learned Corrections (rules)
│   └── Loaded into agent system prompt
└── History (daily activity log)
    └── Markdown files per day
```

**Implementation:** Use `pgvector` extension with SQLite (via `sqlite-vss`) for self-hosted, or Qdrant embedded mode.

### 6.2 Memory Scoping

| Scope | Visibility |
|-------|-----------|
| Global | All agents, all projects |
| Project | All agents within one project |
| Agent | Only that specific agent |
| Session | Only within one conversation |

## 7. Tool Approval Model

Identical to MeshClaw's 3-tier:

| Mode | Behavior |
|------|----------|
| `interactive` | Prompt user for every tool call (default) |
| `reads` | Auto-approve read-only tools (fs read, list, search) |
| `yolo` | Auto-approve everything (requires isolated workspace) |

Configured per-agent in `agent-spec.json` and overridable per-session via the dashboard.

## 8. Dashboard UI Additions

New pages/sections to add to Mitran's React dashboard:

| Page | Purpose |
|------|---------|
| **Chat** | Conversational interface (talk to any agent) |
| **Memory** | View/search semantic + episodic memory |
| **Cron** | Manage scheduled jobs (add, pause, trigger) |
| **Artifacts** | Browse versioned widgets/content |
| **MCP Registry** | Configure MCP servers + view tools |
| **Sessions** | View/resume past conversations |
| **Agent Config** | Edit agent specs, skills, prompts |
| **Approval Queue** | Pending tool approvals (if interactive mode) |
| **Security Log** | Audit trail of all agent actions |

## 9. Settings Updates

### New settings in `.mitran/settings.json`:

```json
{
  "openclaw": {
    "default_model": "claude-sonnet-4",
    "model_provider": "bedrock",
    "bedrock_region": "us-east-1",
    "memory": {
      "enabled": true,
      "vector_db": "sqlite-vss",
      "embedding_model": "all-MiniLM-L6-v2"
    },
    "approval_mode": "reads",
    "max_subagents": 5,
    "session_ttl_hours": 24,
    "artifact_max_versions": 50
  },
  "mcp_servers": {},
  "integrations": {
    "slack": { "bot_token": "", "team_id": "" },
    "github": { "token": "", "org": "" },
    "email": { "smtp_host": "", "smtp_port": 587 }
  }
}
```

## 10. Migration Path

For existing MeshClaw users who want to adopt Mitran:

```bash
# Export MeshClaw state
meshclaw memory export -o memory.json
meshclaw snapshot ./meshclaw-backup

# Import into Mitran
mitran import --from-meshclaw ./meshclaw-backup
# Converts: memory, lessons, crons, artifacts, skills
```

## 11. What's Different from MeshClaw

| Aspect | MeshClaw | Mitran + OpenClaw |
|--------|----------|-------------------|
| Runtime | Python (kiro-cli wrapper) | Go native |
| Single user | Yes | Multi-user with roles |
| LLM access | via kiro-cli | Direct Bedrock/OpenAI/local |
| MCP config | Via AIM + agent-spec | Via settings UI + mcp-config.json |
| Auth | Midway token | OAuth2 / SAML / API key |
| Dashboard | Separate SPA | Unified with platform |
| Billing | N/A (internal tool) | Stripe for enterprise |
| Distribution | pip install | Single binary (`mitran serve`) |

## 12. Implementation Priority

| Phase | Components |
|-------|-----------|
| **Phase 1** | MCP registry + Chat UI + Session manager |
| **Phase 2** | Memory (vector store) + Learned corrections |
| **Phase 3** | Cron scheduler + Artifacts + Skills |
| **Phase 4** | Subagents + Tool approval UI + Security audit |
| **Phase 5** | Snapshot/restore + Import from MeshClaw |
