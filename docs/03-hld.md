# Mitran — High-Level Design (HLD)

## 1. System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    HUMAN OPERATOR                                 │
│  • Priority Queue UI (drag/drop tasks)                          │
│  • Natural Language ("Hey Mitran, write the HLD for auth")      │
│  • Workflow Definitions (on PR → Review → Docs → Deploy)        │
└──────────────────────────────┬──────────────────────────────────┘
                               │ gRPC / WebSocket
┌──────────────────────────────▼──────────────────────────────────┐
│                     MITRAN CORE ENGINE (Go)                       │
│                                                                  │
│  ┌─────────────┐  ┌──────────────┐  ┌───────────────────────┐  │
│  │ DAG Engine  │  │Priority Queue│  │  Checkpoint Manager   │  │
│  │             │  │              │  │                       │  │
│  │ • Topo sort │  │ • Human-     │  │ • State persistence  │  │
│  │ • Dep graph │  │   managed    │  │ • Rollback support   │  │
│  │ • Cycle det │  │ • Per-resource│  │ • Approval gates     │  │
│  └─────────────┘  │   locking    │  └───────────────────────┘  │
│                    └──────────────┘                              │
│  ┌─────────────┐  ┌──────────────┐  ┌───────────────────────┐  │
│  │Agent Runtime│  │ App Registry │  │   Event Bus           │  │
│  │             │  │              │  │                       │  │
│  │ • Lifecycle │  │ • Local apps │  │ • Agent ↔ Agent msgs │  │
│  │ • Sandboxing│  │ • Marketplace│  │ • Webhook triggers   │  │
│  │ • Path ACLs │  │ • Versioning │  │ • Human notifications│  │
│  └─────────────┘  └──────────────┘  └───────────────────────┘  │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                    Integration Layer                        │ │
│  │  GitHub API │ Slack API │ Grafana API │ Prometheus API     │ │
│  └────────────────────────────────────────────────────────────┘ │
└──────────────────────────────┬──────────────────────────────────┘
                               │ Python SDK (gRPC)
┌──────────────────────────────▼──────────────────────────────────┐
│                      AGENT LAYER (Python)                         │
│                                                                  │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐      │
│  │  Dev   │ │  Docs  │ │  CI/CD │ │Tickets │ │  Wiki  │      │
│  │ Agent  │ │ Agent  │ │ Agent  │ │ Agent  │ │ Agent  │      │
│  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘      │
│  ┌────────┐ ┌────────┐ ┌────────┐                             │
│  │  Ops   │ │   HR   │ │Dashboard│                             │
│  │ Agent  │ │ Agent  │ │ Agent  │                             │
│  └────────┘ └────────┘ └────────┘                             │
└─────────────────────────────────────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                     DASHBOARD (TypeScript/React)                  │
│                                                                  │
│  • Priority Queue View    • Agent Activity Feed                 │
│  • Checkpoint Review      • Ticketing UI                        │
│  • Wiki Editor            • Workflow Builder                    │
│  • Grafana Embeds         • HR Portal                          │
└─────────────────────────────────────────────────────────────────┘
```

## 2. Core Components

### 2.1 DAG Engine
**Purpose:** Manages task dependencies and execution order.

**How it works:**
- Tasks are nodes in a directed acyclic graph
- Edges represent dependencies ("Docs Agent can't start until Dev Agent finishes feature X")
- Topological sort determines execution order
- Cycle detection prevents deadlocks
- Human can override order via priority queue

**Comparison with competitors:**
| System | DAG Model |
|--------|-----------|
| Temporal | Workflow-as-code, Go SDK, activity functions |
| Airflow | Python DAG files, static definition |
| MeshClaw | No DAG — sequential tasks in TASK.md, parallel via spawn_run |
| **Mitran** | Dynamic DAG — agents register dependencies, human reorders priority |

### 2.2 Priority Queue
**Purpose:** Human-controlled task backlog. The core primitive of Mitran.

**Properties:**
- One queue per workspace (not per agent)
- Human can: reorder, assign to specific agent, set urgency, pause/resume
- Sequential execution on shared resources (no parallel writes to same file)
- Parallel execution allowed on disjoint resources

**Three interaction modes:**
1. **UI Mode:** Drag-and-drop task cards, assign to agents, approve outputs
2. **NL Mode:** "Mitran, write the HLD for auth service" → creates task, assigns to Docs Agent
3. **Workflow Mode:** "On every PR, run Review Agent → Docs Agent → CI Agent" → template

### 2.3 Checkpoint Manager
**Purpose:** Every agent output is a checkpoint, not final delivery.

**Flow:**
```
Agent works → Produces output → CHECKPOINT → Human reviews
  ├─ Approved → Next task in queue
  ├─ Rejected → Agent re-does with feedback
  └─ Modified → Human edits, agent learns from correction
```

**State persistence:** SQLite (self-hosted) / PostgreSQL (enterprise). Every checkpoint stored with full diff, timestamp, agent ID, human decision.

### 2.4 Agent Runtime
**Purpose:** Manages agent lifecycle, permissions, and execution.

**Permissions model (path ACLs):**
```yaml
# agent-manifest.yaml for Dev Agent
name: dev-agent
permissions:
  read: ["src/**", "docs/**", "config/**"]
  write: ["src/**", "tests/**"]
  execute: ["go build", "go test", "npm run *"]
  integrations: ["github"]
```

Each agent declares its read/write paths. Core engine enforces at runtime. An HR Agent cannot write to `src/`. A Dev Agent cannot read `hr/payroll/`.

**Comparison:**
| System | Isolation |
|--------|-----------|
| Devin | Full Docker VM per task (~2GB overhead) |
| MeshClaw | No isolation — full filesystem access |
| **Mitran** | Path ACLs — lightweight, enforced by core engine |

### 2.5 App Registry
**Purpose:** Manages first-party and third-party agents/apps.

**Two sources:**
1. **Local apps:** `mitran app install ./my-custom-agent/`
2. **Marketplace:** `mitran app install community/security-scanner` (Phase 2)

**App structure:**
```
my-agent/
├── manifest.yaml      # Name, permissions, dependencies
├── agent.py           # Main agent logic (Python SDK)
├── tools/             # Custom tools the agent uses
├── prompts/           # System prompts, templates
└── tests/             # Agent test suite
```

### 2.6 Integration Layer
**Purpose:** Connects Mitran to external systems.

| Integration | Method | What Mitran Does |
|-------------|--------|-----------------|
| GitHub | REST API + Webhooks | Create repos, PRs, branch protection, listen to events |
| Slack | Bot + Slash commands | Notifications, human approvals, NL commands |
| Grafana | HTTP API | Create dashboards, set up data sources, embed in Mitran UI |
| Prometheus | Config generation | Generate prometheus.yml, alerting rules, recording rules |
| Alertmanager | Config generation | Generate alertmanager.yml, route to Slack/PagerDuty |

## 3. Data Flow

### 3.1 The "mitran init" Flow (Killer Demo)
```
User: mitran init
  │
  ▼
Mitran: "What does your company do?"
User: "Mobile fitness app, 8 engineers, TypeScript + Go"
  │
  ▼
┌─────────────────────────────────────────────────┐
│ PLANNING PHASE (Orchestrator Agent)              │
│                                                  │
│ Generates task graph:                            │
│ 1. [Dev Agent] Create GitHub repos + CI         │
│ 2. [Docs Agent] Generate docs structure         │
│ 3. [CI/CD Agent] Set up GitHub Actions          │
│ 4. [Tickets Agent] Create sprint board          │
│ 5. [Wiki Agent] Create onboarding wiki          │
│ 6. [Ops Agent] Install Prometheus + Grafana     │
│ 7. [Ops Agent] Configure alerting               │
│ 8. [HR Agent] Set up PTO tracker                │
└─────────────────────┬───────────────────────────┘
                      │
                      ▼ Human approves plan
┌─────────────────────────────────────────────────┐
│ EXECUTION PHASE (Sequential, checkpointed)       │
│                                                  │
│ Task 1 → Dev Agent → CHECKPOINT → Human ✓       │
│ Task 2 → Docs Agent → CHECKPOINT → Human ✓      │
│ Task 3 → CI/CD Agent → CHECKPOINT → Human ✓     │
│ ... (parallel where resources are disjoint)      │
└─────────────────────────────────────────────────┘
                      │
                      ▼ ~30 minutes later
┌─────────────────────────────────────────────────┐
│ COMPLETE: Your infrastructure is live            │
│ • 3 GitHub repos with CI/CD                     │
│ • Documentation site                            │
│ • Sprint board with templates                   │
│ • Monitoring dashboards                         │
│ • Alerting configured                           │
│ • Wiki with onboarding guide                    │
│ • PTO tracker                                   │
└─────────────────────────────────────────────────┘
```

### 3.2 Steady-State Operation
After init, agents operate continuously:
- **Dev Agent:** Watches PRs, writes code when assigned tasks, runs tests
- **Docs Agent:** Auto-updates docs when code changes, writes ADRs
- **CI/CD Agent:** Monitors pipelines, fixes flaky tests, optimizes build times
- **Tickets Agent:** Triages incoming tickets, assigns priority, tracks SLAs
- **Ops Agent:** Monitors alerts, performs initial diagnosis, escalates to human
- **Wiki Agent:** Keeps wiki in sync with docs, generates runbooks from incidents
- **HR Agent:** Manages PTO requests, sends reminders, tracks compliance
- **Dashboard Agent:** Creates/updates Grafana dashboards based on new services

## 4. Technology Stack

| Layer | Technology | Why |
|-------|-----------|-----|
| Core Engine | Go 1.22+ | Single binary, goroutines, proven for orchestration |
| Agent SDK | Python 3.11+ | LLM ecosystem, easy for agent authors |
| Dashboard | TypeScript + React | Rich UI, component ecosystem |
| Database | SQLite (local) / PostgreSQL (enterprise) | Zero-config local, scalable enterprise |
| Message Queue | NATS (embedded) | Lightweight, Go-native, event bus |
| Cache | In-process (local) / Redis (enterprise) | Speed |
| LLM | Claude/GPT-4/local models via SDK | Provider-agnostic |
| Monitoring | Prometheus + Grafana (bundled) | Apache 2.0, industry standard |
| Search | Bleve (embedded) | Full-text search for tickets/wiki, Go-native |

## 5. Deployment Architecture

### Self-Hosted (Phase 1)
```
Single binary: mitran serve
├── Core Engine (Go)
├── Agent Runtime (Python subprocess pool)
├── Dashboard (embedded static files)
├── SQLite database
├── NATS (embedded)
├── Prometheus (child process)
└── Grafana (child process)
```

Everything runs from one `mitran` binary + Python for agents. Install:
```bash
curl -sSL https://install.mitran.dev | sh
mitran init
```

### Managed Cloud (Phase 2)
```
Kubernetes cluster:
├── Core Engine (Deployment, 3 replicas)
├── Agent Workers (StatefulSet, auto-scaling)
├── Dashboard (Deployment + CDN)
├── PostgreSQL (managed RDS/CloudSQL)
├── NATS (cluster mode)
├── Redis (ElastiCache)
├── Prometheus + Grafana (managed)
└── S3/GCS (artifact storage)
```

## 6. Security Model

| Layer | Mechanism |
|-------|-----------|
| Agent permissions | Path ACLs declared in manifest, enforced by core |
| Human auth | Local: token-based. Enterprise: SSO/SAML |
| API auth | mTLS between core and agents |
| Secrets | Encrypted at rest (AES-256-GCM), injected at runtime |
| Audit | Every action logged: who, what, when, approved-by |
| LLM data | No training data sent to LLM providers (configurable) |
| Network | All services bind to localhost by default. Enterprise: internal network only |

## 7. Non-Functional Requirements

| Requirement | Target |
|-------------|--------|
| Startup time | < 10 seconds (including all agents) |
| Agent response time | < 5 seconds for simple tasks, < 5 minutes for complex |
| Checkpoint latency | < 500ms from agent output to human notification |
| Concurrent agents | Up to 8 (one per type), parallel on disjoint resources |
| Data retention | All checkpoints retained for 90 days (configurable) |
| Uptime (self-hosted) | Depends on host. No external dependencies at runtime. |
| Uptime (managed) | 99.9% SLA |
