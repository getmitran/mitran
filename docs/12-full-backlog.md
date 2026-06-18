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
