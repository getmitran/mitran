# Mitran v0.1.0 — Full Task Backlog

**Team:** 20 people | **Timeline:** 6 months | **Scope:** Full features, no cuts

---

## ✅ COMPLETED (Prior Sessions)

- [x] Vision, HLD, LLD, Pitch Deck, MVP Plan, Competitive Analysis
- [x] Landing page (live at getmitran.vercel.app)
- [x] Go DAG Engine (scheduler, queue, locks, tests)
- [x] Go CLI (`mitran init` wired to engine API)
- [x] Go Core Engine HTTP Server (REST API, scheduler loop, JSON persistence)
- [x] Python Agent Worker (8 agents, AWS Bedrock Claude)
- [x] React Dashboard (queue, kanban, agents, checkpoints, init wizard, settings)
- [x] Multi-env CI/CD (environment CRUD, pipeline templates)
- [x] Project Context Isolation (workspace manager)
- [x] Configurable Settings (paths, modules, per-project)
- [x] Architecture v2 doc (finalized decisions)
- [x] 20-Person Sprint Plan doc
- [x] GitHub repo pushed (getmitran/mitran)

---

## 🔲 REMAINING — Phase 1: Foundation (Sprints 1-4)

### Core Engine (Go)

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 1 | Define gRPC protobuf contracts (proto/ directory) | Go-1 | S1 | P0 |
| 2 | Implement gRPC server in Go engine | Go-1 | S1-S2 | P0 |
| 3 | Replace HTTP REST agent dispatch with gRPC | Go-2 | S2 | P0 |
| 4 | Process supervisor (health checks, auto-restart Python workers) | Go-2 | S2 | P0 |
| 5 | WebSocket server for real-time dashboard updates | Go-3 | S2 | P1 |
| 6 | Session manager (persistent sessions, multi-user) | Go-1 | S3 | P1 |
| 7 | Auth system (JWT tokens, API keys, user roles) | Go-3 | S3 | P1 |
| 8 | Cron scheduler (full cron expressions, timezone, per-agent) | Go-2 | S3 | P1 |
| 9 | Subagent pool (parallel dispatch, result collection) | Go-1 | S4 | P1 |
| 10 | Snapshot/restore (backup, component-level restore) | Go-3 | S4 | P2 |
| 11 | Security audit log (HMAC-signed action log) | Go-2 | S4 | P2 |

### Python Workers (gRPC + LLM)

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 12 | Implement gRPC client in Python (connect to Go server) | Py-1 | S1-S2 | P0 |
| 13 | LLM abstraction layer (Bedrock + OpenAI + local model support) | Py-2 | S1 | P0 |
| 14 | Streaming response support (token-by-token to dashboard) | Py-2 | S2 | P0 |
| 15 | MCP client library (discover tools, call tools, handle results) | Py-3 | S2-S3 | P0 |
| 16 | Agent read-back (agents read prior outputs before generating) | Py-1 | S3 | P1 |
| 17 | Tool calling loop (agent requests tool → approval → execute → return) | Py-4 | S3 | P0 |
| 18 | Skill loader (read SKILL.md files, inject into agent context) | Py-3 | S3 | P1 |
| 19 | Memory write-back (agents persist learnings after task completion) | Py-4 | S4 | P1 |
| 20 | Error handling & retries (LLM timeouts, rate limits, graceful degradation) | Py-1 | S4 | P1 |

### Memory System

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 21 | Semantic memory (JSON key-value facts store) | Py-4 | S2 | P0 |
| 22 | Episodic memory (markdown history per day) | Py-4 | S2 | P1 |
| 23 | Learned corrections (JSON rules, loaded into prompts) | Py-3 | S3 | P1 |
| 24 | Memory API (GET/POST/SEARCH/DELETE via engine) | Go-3 | S3 | P1 |
| 25 | Memory UI page in dashboard | FE-3 | S4 | P2 |

### Dashboard (React)

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 26 | Chat interface (streaming messages, multi-agent selector) | FE-1 | S2-S3 | P0 |
| 27 | WebSocket client (replace polling with real-time) | FE-2 | S2 | P0 |
| 28 | Session management page (list, resume, end sessions) | FE-2 | S3 | P1 |
| 29 | Cron manager page (add, pause, trigger, delete jobs) | FE-3 | S3 | P1 |
| 30 | Artifact library page (list, view, version history) | FE-4 | S4 | P1 |
| 31 | MCP Registry page (view connected servers, available tools) | FE-3 | S4 | P1 |
| 32 | Tool approval queue (pending approvals, approve/reject UI) | FE-1 | S4 | P1 |
| 33 | Agent configuration page (edit agent-spec.json, skills) | FE-4 | S4 | P2 |
| 34 | Security log viewer (filterable audit trail) | FE-2 | S4 | P2 |

---

## 🔲 REMAINING — Phase 2: Feature Complete (Sprints 5-8)

### Integrations

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 35 | Slack integration (bot, notifications, approve via reaction) | FS-1 | S5 | P1 |
| 36 | GitHub integration (create repos, PRs, branch protection, webhooks) | FS-2 | S5 | P1 |
| 37 | Prometheus bundling (config gen, child process, auto-scrape) | FS-1 | S6 | P1 |
| 38 | Grafana bundling (child process, API for dashboard creation) | FS-2 | S6 | P1 |
| 39 | Slack slash commands (mitran status, approve, reject) | FS-1 | S7 | P2 |
| 40 | GitHub Actions integration (CI/CD agent triggers real workflows) | FS-2 | S7 | P2 |

### Built-in Apps

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 41 | Ticketing app (full kanban + sprints + labels + SLA tracking) | FE-4 + Py-2 | S5-S6 | P1 |
| 42 | Wiki app (markdown editor, page tree, full-text search) | FE-3 + Py-3 | S5-S6 | P1 |
| 43 | HR Portal (PTO requests, onboarding checklists, policies) | FE-2 + Py-4 | S7 | P2 |
| 44 | Monitoring dashboards (Grafana embed, auto-generated panels) | FS-2 | S7 | P2 |

### Agent Enhancement

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 45 | Dev Agent v2 (creates real GitHub PRs, runs tests) | Py-1 | S5 | P1 |
| 46 | CI/CD Agent v2 (triggers real GitHub Actions, monitors runs) | Py-2 | S5 | P1 |
| 47 | Docs Agent v2 (reads codebase, generates API docs from code) | Py-3 | S6 | P1 |
| 48 | Ops Agent v2 (creates real Grafana dashboards, Prometheus rules) | Py-4 | S6 | P1 |
| 49 | Review Agent (automated code review, security scan, style check) | Py-1 | S7 | P1 |
| 50 | Wiki Agent v2 (syncs with docs, generates runbooks from incidents) | Py-3 | S7 | P2 |
| 51 | Tickets Agent v2 (auto-triage, SLA enforcement, sprint planning) | Py-2 | S8 | P2 |
| 52 | HR Agent v2 (PTO approval workflow, onboarding automation) | Py-4 | S8 | P2 |

### Apps Ecosystem

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 53 | App manifest format (app.json spec) | Go-1 | S6 | P1 |
| 54 | App install/enable/disable/uninstall commands | Go-2 | S6 | P1 |
| 55 | App registry API (list, search, install from marketplace) | Go-3 | S7 | P2 |
| 56 | App template generator (`mitran app init`) | Go-2 | S7 | P2 |

---

## 🔲 REMAINING — Phase 3: Polish & Release (Sprints 9-13)

### CLI Completion

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 57 | `mitran serve` (starts engine + workers + prints status) | Go-1 | S9 | P0 |
| 58 | `mitran chat` (REPL mode, talks to agents) | Go-2 | S9 | P1 |
| 59 | `mitran status` (live stats: agents, tasks, memory) | Go-3 | S9 | P1 |
| 60 | `mitran config` (get/set settings from CLI) | Go-2 | S9 | P2 |
| 61 | `mitran agent` (list, create, update, delete agents) | Go-1 | S10 | P2 |
| 62 | `mitran workspace` (list, create, switch projects) | Go-3 | S10 | P2 |
| 63 | `mitran cron` (add, list, pause, trigger from CLI) | Go-2 | S10 | P2 |
| 64 | `mitran snapshot` / `mitran restore` | Go-3 | S11 | P2 |

### DevOps & Infrastructure

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 65 | GitHub Actions CI (Go test + Python test + TS typecheck) | DevOps-1 | S1 | P0 |
| 66 | Docker Compose (full stack: engine + workers + dashboard + monitoring) | DevOps-1 | S5 | P1 |
| 67 | Helm chart for Kubernetes | DevOps-2 | S8 | P2 |
| 68 | Install script (`curl \| sh` — handles Go, Python, Node deps) | DevOps-1 | S9 | P0 |
| 69 | Binary builds (goreleaser for multi-arch Go binary) | DevOps-2 | S10 | P1 |
| 70 | E2E test suite (full init → agent → checkpoint → approve flow) | DevOps-1 | S10 | P1 |
| 71 | Load testing (100 concurrent tasks, measure scheduler throughput) | DevOps-2 | S11 | P2 |
| 72 | Security scan (dependency audit, SAST, secret detection) | DevOps-1 | S11 | P1 |

### QA & Testing

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 73 | Unit test coverage >80% (Go engine) | QA-1 | S6-S10 | P1 |
| 74 | Unit test coverage >80% (Python workers) | QA-2 | S6-S10 | P1 |
| 75 | Integration tests (gRPC contracts, MCP tool calling) | QA-1 | S8 | P1 |
| 76 | Dashboard E2E tests (Playwright or Cypress) | QA-2 | S9 | P2 |
| 77 | Chaos testing (kill workers mid-task, verify recovery) | QA-1 | S11 | P2 |
| 78 | Performance benchmarks (documented, reproducible) | QA-2 | S12 | P2 |

### Design & Branding

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 79 | Logo design (icon + wordmark + favicon) | Designer | S1-S2 | P1 |
| 80 | Dashboard theme (finalize colors, typography, spacing) | Designer | S2-S3 | P1 |
| 81 | Component library (buttons, cards, modals, forms) | Designer | S3-S4 | P1 |
| 82 | Landing page v2 (with logo, updated copy, screenshots) | Designer | S9 | P1 |
| 83 | Social cards (OG images for GitHub, Twitter, HN) | Designer | S12 | P2 |

### DevRel & Launch

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 84 | Documentation site (Docusaurus/MkDocs, hosted) | DevRel | S3-S6 | P1 |
| 85 | API reference (auto-generated from proto + REST) | DevRel | S7 | P1 |
| 86 | Agent SDK developer guide | DevRel | S8 | P1 |
| 87 | Tutorial: "Your first Mitran project in 5 minutes" | DevRel | S9 | P1 |
| 88 | Demo video (2-min, professional) | DevRel | S12 | P0 |
| 89 | HN launch post draft | DevRel | S12 | P0 |
| 90 | Discord community setup (channels, roles, bot) | DevRel | S10 | P2 |
| 91 | Buy getmitran.dev domain + DNS setup | DevRel | S1 | P1 |
| 92 | Twitter/X account + first 10 posts scheduled | DevRel | S11 | P2 |

### Enterprise Prep (Sprints 11-13)

| # | Task | Owner | Sprint | Priority |
|---|------|-------|--------|----------|
| 93 | SSO/SAML integration (OAuth2 + SAML for enterprise) | Go-1 | S11 | P2 |
| 94 | Multi-tenant data isolation (per-org data separation) | Go-2 | S11 | P2 |
| 95 | PostgreSQL backend (replace JSON store for enterprise) | Go-3 | S12 | P2 |
| 96 | Audit log export (CSV/JSON, compliance-ready) | Go-2 | S12 | P2 |
| 97 | Rate limiting + usage tracking | Go-1 | S12 | P2 |
| 98 | Enterprise landing page + pricing | DevRel + Designer | S12 | P2 |
| 99 | 3 design partner onboarding calls | Tech Lead | S10-S12 | P1 |
| 100 | AGPL license file + CLA for contributors | Tech Lead | S1 | P0 |

---

## Summary

| Category | Tasks | Status |
|----------|-------|--------|
| Completed | 13 | ✅ |
| Foundation (S1-S4) | 34 | 🔲 |
| Feature Complete (S5-S8) | 22 | 🔲 |
| Polish & Release (S9-S13) | 31 | 🔲 |
| **Total** | **100** | |

**100 tasks. 20 people. 13 sprints. Zero mocks. v0.1.0.**


## REMAINING v0.1.0 SCOPE (Not Yet Built)

### OpenClaw Runtime Features

| # | Task | Description | Priority |
|---|------|-------------|----------|
| 101 | Vector memory store | Embedded Qdrant/ChromaDB replacing JSON string matching. Semantic similarity search over facts+episodes. | P0 |
| 102 | Skill loader | Read SKILL.md files from workspace, inject into agent system prompts dynamically. Per-agent skill assignment. | P0 |
| 103 | Tool approval system | 3-tier model (interactive/reads/yolo). Per-agent configurable. Dashboard approval queue UI. | P0 |
| 104 | Artifact store | Versioned content store (widgets/HTML/md). CRUD API, version history, slug-based addressing. | P0 |
| 105 | Agent read-back | Before generating, agents read prior outputs. Context window management with sliding window + summarization. | P1 |
| 106 | Memory write-back | After task completion, agents auto-persist learnings to memory store. | P1 |
| 107 | LLM abstraction layer | Support Bedrock + OpenAI + Anthropic direct + Ollama. Provider-agnostic interface with model routing. | P0 |
| 108 | MCP server lifecycle | Auto-start configured MCP servers as child processes. Health monitoring, restart on crash. | P0 |
| 109 | Process supervisor | Monitor Python workers. Auto-restart on crash. Health pings. Multiple instances with load balancing. | P1 |
| 110 | Snapshot/restore system | Full state backup (memory, sessions, config, artifacts, crons). Component-level restore. | P1 |

### HR Portal (Full App)

| # | Task | Description | Priority |
|---|------|-------------|----------|
| 111 | HR Portal dashboard page | Full React page: employee directory, PTO calendar, onboarding tracker, policy library. | P1 |
| 112 | PTO request system | Submit requests with date picker, type selection, auto-calculate balance, manager approval workflow. | P1 |
| 113 | Onboarding checklists | Template-based flows. Assign to new hire. Track completion. Auto-create tasks in ticketing. | P1 |
| 114 | Policy library | CRUD for company policies (markdown). Versioned. Searchable. HR agent uses as context. | P2 |
| 115 | Employee directory | Team members with role, department, start date, manager. Org chart view. | P2 |
| 116 | PTO balance tracking | Per-employee accrual, usage history, remaining balance. Year-end rollover rules. | P2 |

### CLI Commands

| # | Task | Description | Priority |
|---|------|-------------|----------|
| 117 | mitran serve | Start engine + all workers + print status table. Watch for crashes. Single command. | P0 |
| 118 | mitran chat | REPL-mode conversational interface. Agent selector. Streaming responses. History. | P0 |
| 119 | mitran tui | Terminal UI (Bubble Tea). Split panes: chat + task list + agent status. | P1 |
| 120 | mitran status | Live stats: running agents, queued tasks, memory entries, cron jobs, uptime, worker health. | P1 |
| 121 | mitran config get/set | Read/write settings from CLI. Dot-notation keys. JSON output. | P2 |
| 122 | mitran agent list/create/delete | Manage agent configs from CLI. Hot-reload on change. | P2 |
| 123 | mitran workspace list/create/switch | Multi-project. Each workspace = isolated memory + sessions + config. | P1 |
| 124 | mitran cron add/list/pause/trigger | Full cron management from CLI. | P2 |
| 125 | mitran snapshot/restore | CLI interface to snapshot system. Component selection. Dry-run. | P2 |
| 126 | mitran app install/enable/disable | Install MCP apps from local dirs or marketplace. | P2 |

### Apps Ecosystem

| # | Task | Description | Priority |
|---|------|-------------|----------|
| 127 | App manifest spec (app.json) | Define format: name, version, permissions, UI components, crons, skills. | P1 |
| 128 | App install from directory | Read app.json, validate, copy to store, register routes/crons/skills. | P1 |
| 129 | App enable/disable per workspace | Workspace-scoped activation. | P2 |
| 130 | App registry API | List installed, search marketplace, version checking, updates. | P2 |
| 131 | App template generator | mitran app init scaffolds app.json + entry + sample skill + cron. | P2 |
| 132 | App marketplace scaffold | Remote registry API structure for community apps. | P2 |

### Dashboard Pages (Missing)

| # | Task | Description | Priority |
|---|------|-------------|----------|
| 133 | Memory viewer page | Browse facts, episodes, corrections. Search. Add/edit/delete. Timeline. | P1 |
| 134 | Cron manager page | List jobs, pause/resume, trigger, run history, create new with form. | P1 |
| 135 | Artifact library page | Grid/list of artifacts. Preview. Version dropdown. Tag filtering. | P1 |
| 136 | MCP Registry page | Connected servers with status. Tool browser. Enable/disable. Add form. | P1 |
| 137 | Agent config editor | Monaco editor for agent-spec. Skills assignment. Model select. Test prompt. | P2 |
| 138 | Security log viewer | Filterable audit trail. Action type, agent, timestamp. Export CSV. | P2 |

### Real Integrations

| # | Task | Description | Priority |
|---|------|-------------|----------|
| 139 | Prometheus bundling | Ship binary. Auto-generate scrape config. Child process. /metrics endpoints. | P1 |
| 140 | Grafana bundling | Ship binary. Auto-provision datasource. Pre-built agent dashboards. | P1 |
| 141 | Monitoring dashboard page | Embed Grafana iframes. Auto-generated per-agent panels. | P2 |
| 142 | Slack slash commands | /mitran status, approve, reject, chat. Real Slack app manifest. | P1 |
| 143 | GitHub webhooks inbound | Receive PR/issue/push events. Route to appropriate agents. Auto-review PRs. | P1 |

### Agent Enhancements (Real Capabilities)

| # | Task | Description | Priority |
|---|------|-------------|----------|
| 144 | Ops Agent real Grafana | Creates actual Grafana dashboard JSON via API. | P1 |
| 145 | Ops Agent real alerting | Define alert rules, notification channels, escalation policies. | P2 |
| 146 | Tickets Agent SLA enforcement | Track SLA per priority. Auto-escalate overdue. Breach notifications. | P1 |
| 147 | Tickets Agent sprint planning | AI-suggested scope from backlog, velocity estimation, workload balance. | P2 |
| 148 | Wiki Agent incident runbooks | Auto-generate runbooks from incident history. Link to dashboards/alerts. | P2 |
| 149 | HR Agent PTO approval workflow | Receive requests, check balance, route to manager, update calendar. | P1 |
| 150 | Review Agent security scan | Run actual semgrep/gosec, parse results, create tickets for findings. | P1 |

### Enterprise Prep

| # | Task | Description | Priority |
|---|------|-------------|----------|
| 151 | SSO/SAML integration | OAuth2 + SAML provider. Enterprise IdP connection. Role mapping. | P2 |
| 152 | Multi-tenant isolation | Per-org data separation. Org-scoped agents, memory, sessions. | P2 |
| 153 | PostgreSQL backend | Replace SQLite for production. Migration tool. Connection pooling. | P2 |
| 154 | Audit log export | CSV/JSON export. Compliance format. Retention policies. Signed logs. | P2 |
| 155 | Usage metering | Track LLM tokens per user/agent, API calls, storage. Billing-ready. | P2 |

### DevRel and Launch

| # | Task | Description | Priority |
|---|------|-------------|----------|
| 156 | Documentation site | Docusaurus/MkDocs on getmitran.dev. Auto-deploy. Search. Versioned. | P1 |
| 157 | Install script | curl pipe sh. Detects OS, installs Go/Python/Node deps, downloads binary. | P0 |
| 158 | Binary builds (goreleaser) | Multi-arch Go binaries. GitHub Release artifacts. Homebrew formula. | P1 |
| 159 | Tutorial 5 min quickstart | Install to first agent output. With screenshots. | P1 |
| 160 | Demo video 2 min | Professional recording: init, agents working, PR created, dashboard. | P0 |
| 161 | HN launch post | Show HN draft + timing strategy. | P1 |
| 162 | Discord community | Server, channels, roles, welcome bot, contributor guidelines. | P2 |
| 163 | Buy getmitran.dev | Domain + DNS + email forwarding. | P1 |
| 164 | Social presence | Twitter @getmitran. First 10 posts. GitHub social preview. Badges. | P2 |

### Testing and Quality

| # | Task | Description | Priority |
|---|------|-------------|----------|
| 165 | Go unit tests 80 pct coverage | Cover all packages fully. | P1 |
| 166 | Python unit tests 80 pct coverage | Cover all agents and utilities fully. | P1 |
| 167 | gRPC integration tests | Go to Python contract testing. Mock LLM for deterministic results. | P1 |
| 168 | Dashboard E2E tests (Playwright) | Init wizard, task creation, kanban, chat, settings, navigation. | P2 |
| 169 | Chaos testing | Kill workers mid-task, network partition, OOM. Verify recovery. | P2 |
| 170 | Performance benchmarks | 100 concurrent tasks. Scheduler throughput. Published results. | P2 |

### Design and Branding

| # | Task | Description | Priority |
|---|------|-------------|----------|
| 171 | Logo design | Icon + wordmark + favicon. SVG. Dark/light variants. | P1 |
| 172 | Dashboard theme finalization | Consistent colors, typography, spacing. Design tokens. | P1 |
| 173 | Component library | Reusable: Button, Card, Modal, Form, Select, Toast. Storybook. | P2 |
| 174 | Landing page v2 | With logo, screenshots, feature grid, comparison, CTA. | P1 |
| 175 | Social cards | OG images for GitHub, Twitter, HN. Template. | P2 |

## Updated Summary

| Category | Tasks | Status |
|----------|-------|--------|
| Completed (v0.1.0 built) | 1-100 | Done |
| OpenClaw Runtime | 101-110 | Remaining |
| HR Portal | 111-116 | Remaining |
| CLI Commands | 117-126 | Remaining |
| Apps Ecosystem | 127-132 | Remaining |
| Dashboard Pages | 133-138 | Remaining |
| Integrations | 139-143 | Remaining |
| Agent Enhancements | 144-150 | Remaining |
| Enterprise | 151-155 | Remaining |
| DevRel | 156-164 | Remaining |
| Testing | 165-170 | Remaining |
| Design | 171-175 | Remaining |
| **TOTAL** | **175** | 100 done, 75 remaining |
