# Mitran — Competitive Analysis

## Direct Competitors

### 1. CrewAI (OSS, Python)
**What it is:** Multi-agent Python framework. Role-based agents collaborate on tasks.

**How they did it:**
- Launched early 2024, hit 50K+ GitHub stars
- Pure Python library (`pip install crewai`)
- Agents have roles, goals, backstory
- Sequential or parallel task execution
- No UI, no priority queue, no human checkpoints

**Their model:**
```python
crew = Crew(
    agents=[researcher, writer, editor],
    tasks=[research_task, write_task, edit_task],
    process=Process.sequential
)
result = crew.kickoff()
```

**Where Mitran differs:**
| CrewAI | Mitran |
|--------|--------|
| Library (import in code) | Platform (install and use) |
| Async by default | Sync by default, human controls |
| No UI | Full dashboard + ticketing + wiki |
| No checkpoints | Every output is a checkpoint |
| No resource locking | Path-level locks prevent conflicts |
| Developer builds agents in code | Agents ship ready-to-use |
| No integrations | GitHub, Slack, Grafana built-in |

**Lesson for Mitran:** CrewAI proved demand (50K stars). But it's a developer tool, not a team tool. Mitran is what you build AFTER outgrowing CrewAI.

---

### 2. Devin (Cognition, Closed, $500/mo)
**What it is:** Autonomous AI software engineer. Full VM sandbox per task.

**How they did it:**
- $2B valuation (2024), $175M raised
- Each task gets isolated Docker VM with browser, terminal, editor
- User gives natural language task → Devin plans + implements
- No human-in-the-loop during execution (async model)
- Output is a PR ready for review

**Their model:**
```
User: "Add dark mode to the settings page"
Devin: Plans → Implements → Tests → Creates PR (20-60 min)
Human: Reviews PR, merges or rejects
```

**Where Mitran differs:**
| Devin | Mitran |
|-------|--------|
| Single agent (one software engineer) | 8 specialized agents |
| Full VM isolation (~2GB/task) | Path ACLs (lightweight) |
| Async execution (no checkpoints mid-task) | Sync with human checkpoints |
| One product: code generation | Full team infrastructure |
| Closed source, SaaS only | Open source, self-hosted |
| $500/month per seat | Free (enterprise tier for extras) |
| No ticketing/wiki/HR/monitoring | All included |

**Lesson for Mitran:** Devin proved willingness to pay for AI engineering. But $500/mo is expensive and it only writes code. Mitran is broader (entire infra) and free.

---

### 3. MeshClaw (Personal agent orchestration)
**What it is:** Autonomous agent management layer for individual developers. Memory, crons, subagents, Slack integration.

**How they did it:**
- Single-user personal tool (not a team platform)
- Python-based gateway with kiro-cli backend
- Skills, workspace isolation, persistent memory
- Cron jobs for recurring tasks
- Dashboard + Slack interface
- Apps ecosystem (local install)

**Their model:**
```
User → MeshClaw gateway → kiro-cli agent → tools (MCP)
         ↓                    ↓
      Dashboard           Memory/crons/lessons
```

**Where Mitran differs:**
| MeshClaw | Mitran |
|----------|--------|
| Single-user personal tool | Multi-user team platform |
| One agent (meshclaw) | 8 specialized agents |
| No DAG engine | Full DAG with dependency graph |
| Sequential (TASK.md) or parallel (spawn_run) | DAG-scheduled with human priority |
| Skills are markdown files | Agents are full Python apps |
| Local apps (install from directory) | Local + marketplace |
| No ticketing/wiki/HR | Full internal tools suite |
| No resource locking | Path-level mutex |
| Memory per-user | Shared team knowledge base |

**What Mitran borrows from MeshClaw:**
- ✅ Persistent memory/learning (agents improve over time)
- ✅ Cron/recurring tasks concept (scheduled agent work)
- ✅ Apps ecosystem (install/enable/disable)
- ✅ Checkpoint/approval pattern
- ✅ Dashboard + Slack dual interface
- ✅ CLI-first experience
- ✅ Workspace isolation concept (projects as workspaces)

**What Mitran improves:**
- Multi-user collaboration with roles
- Specialized agents instead of one generalist
- DAG engine for complex task orchestration
- Built-in internal tools (not just developer workflow)
- Marketplace for community extensions

---

### 4. Temporal (OSS, Go)
**What it is:** Durable workflow execution engine. Not AI-native.

**How they did it:**
- Fork of Uber's Cadence (Go)
- Workflow-as-code paradigm
- Durable execution: survives process crashes
- SDKs in Go, Java, Python, TypeScript
- Temporal Cloud: managed hosting
- $1.5B valuation

**Their model:**
```go
// Temporal workflow (developer writes this)
func MyWorkflow(ctx workflow.Context) error {
    result := workflow.ExecuteActivity(ctx, ProcessOrder, order)
    // ...durable execution, retries, timeouts handled
}
```

**Where Mitran differs:**
| Temporal | Mitran |
|----------|--------|
| Generic workflow engine | AI-agent-native platform |
| Developer codes every workflow | Agents generate + execute workflows |
| No LLM integration | LLM is the core of every agent |
| No UI for task management | Full dashboard with priority queue |
| Steep learning curve | `mitran init` → running in 30 min |
| Infrastructure for developers | Tool for entire team |

**Lesson for Mitran:** Temporal proved Go is right for durable orchestration. Their workflow-as-code pattern is what Mitran's DAG engine should feel like internally. But Temporal requires developers to build everything — Mitran ships the building blocks ready.

---

### 5. Atlassian (Jira + Confluence + Bitbucket)
**What it is:** Legacy suite for project management, documentation, and source control.

**How they did it:**
- Started with Jira (2002), added Confluence, Bitbucket
- $60B market cap
- Millions of teams use it
- Server → Cloud migration (killed self-hosted in 2024)
- AI features added (Atlassian Intelligence, 2023)

**Where Mitran differs:**
| Atlassian | Mitran |
|-----------|--------|
| Human-operated (people file tickets, write docs) | Agent-operated (agents do the work) |
| $50K-500K/year for enterprise | Free (AGPL) + enterprise tier |
| Decades of tech debt | Built from scratch for AI |
| Adding AI as afterthought | AI-native from day 1 |
| Cloud-only (killed self-hosted) | Self-hosted first |
| Separate products (Jira ≠ Confluence) | Unified platform |

**Lesson for Mitran:** Atlassian killed their self-hosted option in 2024. Millions of "Atlassian Server refugees" need alternatives. Mitran + self-hosted + free is attractive to this audience.

---

## Feature Comparison Matrix

| Feature | Mitran | CrewAI | Devin | MeshClaw | Temporal | Atlassian |
|---------|--------|--------|-------|----------|----------|-----------|
| Multi-agent orchestration | ✅ | ✅ | ❌ | ⚠️ (spawn) | ❌ | ❌ |
| Human priority queue | ✅ | ❌ | ❌ | ❌ | ❌ | ⚠️ (backlog) |
| Checkpoints/approvals | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ |
| DAG scheduling | ✅ | ❌ | ❌ | ❌ | ✅ | ❌ |
| Built-in ticketing | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ (Jira) |
| Built-in wiki | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ (Confluence) |
| Built-in CI/CD | ✅ | ❌ | ❌ | ❌ | ❌ | ⚠️ (Bitbucket) |
| Monitoring/alerting | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| HR portal | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Self-hosted | ✅ | ✅ | ❌ | ✅ | ✅ | ❌ (killed) |
| Open source | ✅ (AGPL) | ✅ (MIT) | ❌ | ❌ | ✅ (MIT) | ❌ |
| Single binary install | ✅ | ❌ (pip) | N/A | ❌ (pip) | ✅ | N/A |
| Natural language interface | ✅ | ❌ | ✅ | ✅ | ❌ | ❌ |
| Agent marketplace | ✅ (Phase 2) | ❌ | ❌ | ⚠️ (local) | ❌ | ✅ (marketplace) |
| Multi-user/team | ✅ | ❌ | ✅ | ❌ | ✅ | ✅ |
| Resource locking | ✅ | ❌ | ✅ (VM) | ❌ | ✅ | N/A |
| Persistent memory | ✅ | ❌ | ❌ | ✅ | ❌ | N/A |

## Positioning Summary

```
                    Full Team Platform
                           ▲
                           │
         Atlassian ●       │      ● MITRAN
         (legacy,          │        (AI-native,
          expensive)       │         open-source,
                           │         self-hosted)
                           │
  Framework ───────────────┼──────────── Product
                           │
         CrewAI ●          │      ● Devin
         Temporal ●        │        (single agent,
         LangGraph ●       │         closed, expensive)
                           │
                           ▼
                     Single Purpose
```

Mitran occupies the **top-right quadrant** — full team platform + product (not framework). Nobody else is here.
