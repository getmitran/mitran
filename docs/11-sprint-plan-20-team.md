# Mitran v0.1.0 Sprint Plan — 20-Person Team

**Timeline:** 26 weeks (13 sprints × 2 weeks)  
**Target:** Full-featured v0.1.0 release, no scope cuts, no mocking  
**LLM Provider:** AWS Bedrock Claude Sonnet  
**Stack:** Go core · Python agents · React/TS dashboard · gRPC IPC

---

## 1. Team Assignments

### Core Engine (Go) — 3 Engineers: G1, G2, G3

| Engineer | Ownership |
|----------|-----------|
| **G1** | DAG scheduler, priority queue, checkpoint manager, resource locks |
| **G2** | Project manager, workspace isolation, CLI, snapshot/restore |
| **G3** | gRPC server, cron scheduler, subagent pool orchestration, auth |

### Agent Workers (Python) — 4 Engineers: P1, P2, P3, P4

| Engineer | Ownership |
|----------|-----------|
| **P1** | Agent framework, MCP client, tool discovery/calling, skill loader |
| **P2** | Dev Agent, CI/CD Agent, Docs Agent |
| **P3** | Ops Agent, Review Agent, Tickets Agent |
| **P4** | HR Agent, Wiki Agent, memory system (JSON→FTS5), learned corrections |

### Frontend (React/TS) — 4 Engineers: F1, F2, F3, F4

| Engineer | Ownership |
|----------|-----------|
| **F1** | Chat interface (streaming, sessions, multi-agent), approval queue |
| **F2** | Kanban board, ticketing UI, sprint views, SLA tracking |
| **F3** | Agent config, MCP registry, cron manager, security log, settings |
| **F4** | Artifacts store UI, wiki editor, memory browser, sessions, dashboard home |

### Full-Stack / Integrations — 2 Engineers: FS1, FS2

| Engineer | Ownership |
|----------|-----------|
| **FS1** | Slack integration (bot, notifications, approvals, NL commands) |
| **FS2** | GitHub integration (repos, PRs, Actions, webhooks), Grafana/Prometheus bundled setup |

### DevOps — 2 Engineers: D1, D2

| Engineer | Ownership |
|----------|-----------|
| **D1** | Docker multi-service compose, Helm charts, install script (curl\|sh), CI pipelines |
| **D2** | Monitoring infra (Prometheus+Grafana bundled), log aggregation, load testing infra, security hardening |

### QA — 2 Engineers: Q1, Q2

| Engineer | Ownership |
|----------|-----------|
| **Q1** | Integration tests, E2E test suite, API contract tests, regression |
| **Q2** | Security audit, penetration testing, load/stress testing, compliance |

### Designer — 1: DS1

| Person | Ownership |
|--------|-----------|
| **DS1** | UI/UX design system, component library specs, branding, logo, landing page |

### DevRel — 1: DR1

| Person | Ownership |
|--------|-----------|
| **DR1** | Documentation site, README, tutorials, community setup, launch strategy, demo videos |

### Tech Lead / Architect — 1: TL (Founder)

| Person | Ownership |
|--------|-----------|
| **TL** | Architecture decisions, code reviews, proto definitions, cross-team coordination, blocker resolution |

---

## 2. Sprint Timeline

### Phase 1: Foundation (Sprints 1–4, Weeks 1–8)

#### Sprint 1 (Weeks 1–2): Scaffolding & Contracts

| Team | Deliverables |
|------|-------------|
| **TL** | Finalize proto definitions (30+ RPCs), establish coding standards, repo structure |
| **G1** | DAG data model, topological sort, cycle detection, basic scheduler loop |
| **G2** | Project CRUD, workspace directory isolation, CLI scaffold (init, serve, status) |
| **G3** | gRPC server skeleton, health check, reflection; auth token generation/validation |
| **P1** | Python gRPC client stub, agent base class, MCP client connection manager |
| **P2** | Dev Agent skeleton (accepts task, returns result via gRPC) |
| **P3** | Ops Agent skeleton |
| **P4** | Memory module v0: JSON read/write, history append, basic search |
| **F1** | React app scaffold (Vite+Tailwind+TS), router, auth gate, WebSocket client |
| **F2** | Component library setup, design tokens from DS1 |
| **F3** | Settings page shell, API client layer with auth headers |
| **F4** | Dashboard home layout, sidebar navigation |
| **FS1** | Slack app manifest, OAuth flow, event subscription scaffold |
| **FS2** | GitHub App registration, webhook receiver, repo list endpoint |
| **D1** | Monorepo CI (lint+build+test for Go/Python/TS), Docker base images |
| **D2** | docker-compose dev environment (all services), Prometheus scrape config |
| **Q1** | Test framework setup (Go: testify, Python: pytest, TS: vitest), CI integration |
| **Q2** | Threat model document, security requirements checklist |
| **DS1** | Design system v1: colors, typography, spacing, component specs for dashboard |
| **DR1** | docs site scaffold (Docusaurus), README draft, contribution guide |

#### Sprint 2 (Weeks 3–4): Core Execution Loop

| Team | Deliverables |
|------|-------------|
| **G1** | Priority queue with preemption, resource lock manager, task state machine |
| **G2** | Checkpoint manager (save/restore task state), workspace file isolation |
| **G3** | Cron scheduler (interval + cron expr), timezone support, job persistence |
| **P1** | MCP tool discovery (list_tools), tool calling with timeout, retry logic |
| **P2** | Dev Agent: code generation, file read/write via MCP tools |
| **P3** | Ops Agent: AWS CLI execution, log analysis; Review Agent: diff parsing |
| **P4** | Memory FTS5: SQLite integration, full-text search, learned corrections CRUD |
| **F1** | Chat page: message list, streaming render, session selector |
| **F2** | Kanban board: columns, drag-drop cards, card detail panel |
| **F3** | Agent config page: list agents, edit config JSON, enable/disable |
| **F4** | Sessions page: list, filter, resume/end buttons |
| **FS1** | Slack message handler, DM routing to correct agent session |
| **FS2** | GitHub PR webhook processing, commit status updates |
| **D1** | Multi-stage Docker builds (Go binary, Python venv, Node bundle) |
| **D2** | Grafana provisioned dashboards (system metrics: CPU, memory, goroutines) |
| **Q1** | gRPC contract tests (Go↔Python), API integration test harness |
| **Q2** | OWASP dependency scan in CI, SAST scanner setup |
| **DS1** | Chat UI mockups, kanban mockups, approval flow wireframes |
| **DR1** | Architecture overview doc, getting started guide draft |

#### Sprint 3 (Weeks 5–6): Agent Intelligence

| Team | Deliverables |
|------|-------------|
| **G1** | DAG execution: parallel node dispatch, failure handling, retry policy |
| **G2** | CLI commands: chat, config, agent, workspace, cron |
| **G3** | Subagent pool: spawn N agents, collect results, timeout handling |
| **P1** | Skill system: SKILL.md parser, per-agent skill loading, skill discovery |
| **P2** | CI/CD Agent: GitHub Actions generation, multi-env deploy configs |
| **P3** | Tickets Agent: kanban CRUD, sprint management, SLA calculation |
| **P4** | History system: conversation storage, decay, episodic memory search |
| **F1** | Chat: multi-agent switching, typing indicators, tool call visualization |
| **F2** | Ticketing UI: create/edit tasks, sprint board, SLA indicators |
| **F3** | Cron manager: list, add, pause/resume, trigger, logs view |
| **F4** | Memory browser: search, view entries, delete; checkpoint viewer |
| **FS1** | Slack approvals (button actions), NL command parsing, thread support |
| **FS2** | GitHub Actions integration: trigger workflows, status polling |
| **D1** | Helm chart v1 (single-node), environment variable management |
| **D2** | Agent-level metrics (tasks/min, LLM tokens, latency p50/p99) |
| **Q1** | Agent behavior tests (given task → expected tool calls), mock LLM layer |
| **Q2** | Auth penetration test, token expiry/refresh validation |
| **DS1** | Artifact viewer mockups, wiki editor mockups, approval queue design |
| **DR1** | Agent authoring guide, skill creation tutorial |

#### Sprint 4 (Weeks 7–8): 🎯 MILESTONE — Internal Demo #1

| Team | Deliverables |
|------|-------------|
| **G1** | Checkpoint resume after crash, DAG visualization data endpoint |
| **G2** | Snapshot/restore: create backup, restore by component, CLI commands |
| **G3** | Tool approval: 3-tier system (interactive/reads/yolo), per-agent override |
| **P1** | MCP server lifecycle management, reconnection logic, health monitoring |
| **P2** | Docs Agent: markdown generation, template filling, doc site deployment |
| **P3** | Ops Agent: alarm investigation, runbook execution, incident correlation |
| **P4** | Semantic memory: key-value store, cross-session persistence |
| **F1** | Approval queue page: pending tools, approve/reject, bulk actions |
| **F2** | Sprint planning view: capacity, velocity chart, burndown |
| **F3** | MCP registry page: connected servers, available tools, health status |
| **F4** | Artifact store UI: list, view versions, create, edit, delete |
| **FS1** | Slack cron notifications, channel-specific routing |
| **FS2** | Prometheus auto-discovery, Grafana dashboard generation by Ops Agent |
| **D1** | Install script v1 (curl\|sh): detects OS, installs Go+Python+Node deps |
| **D2** | Load test harness: k6 scripts for gRPC endpoints, baseline benchmarks |
| **Q1** | E2E test: init project → chat → agent executes → result verified |
| **Q2** | Load test execution: 100 concurrent tasks, report latency/throughput |
| **DS1** | Logo final, landing page design, brand guidelines document |
| **DR1** | Demo script, internal presentation deck, video walkthrough |

**Demo #1 scope:** Create project via CLI → assign task via chat → Dev Agent generates code → approve tool calls → see result in dashboard.

---

### Phase 2: Feature Completeness (Sprints 5–8, Weeks 9–16)

#### Sprint 5 (Weeks 9–10): Integrations & Apps

| Team | Deliverables |
|------|-------------|
| **G1** | Apps ecosystem: manifest parser, install/enable/disable lifecycle |
| **G2** | App sandboxing, local app directory, marketplace client stub |
| **G3** | Multi-user auth: roles (admin/operator/viewer), RBAC enforcement on all RPCs |
| **P1** | App runtime: load app agents, app-scoped MCP connections |
| **P2** | Dev Agent: PR review, code quality checks, auto-fix suggestions |
| **P3** | Review Agent: CR analysis, blocking comment detection, approval recommendation |
| **P4** | Wiki Agent: page CRUD, markdown rendering, page tree management, search index |
| **F1** | Chat: file attachments, code blocks with syntax highlight, image preview |
| **F2** | Wiki editor: markdown WYSIWYG, page tree sidebar, version history |
| **F3** | Security log page: audit trail, deny-list management, HMAC status |
| **F4** | Apps page: installed apps, marketplace browse, install/remove buttons |
| **FS1** | Slack: channel tracking, thread context, workspace selector |
| **FS2** | GitHub: branch protection rules, PR template management, Actions logs |
| **D1** | Helm chart v2: multi-replica, HPA, PDB, resource limits |
| **D2** | Security: TLS everywhere, secret rotation, network policies |
| **Q1** | Integration tests for all 8 agents (happy path + error cases) |
| **Q2** | RBAC test matrix: verify all role combinations across all endpoints |
| **DS1** | Wiki editor UX, apps marketplace design, mobile-responsive layouts |
| **DR1** | API reference (auto-generated from proto), integration guides |

#### Sprint 6 (Weeks 11–12): Polish & Depth

| Team | Deliverables |
|------|-------------|
| **G1** | DAG: conditional branches, retry with backoff, timeout per node |
| **G2** | CLI: app, snapshot, restore commands; shell completions (bash/zsh/fish) |
| **G3** | Cron: script-based execution (no LLM), command-based crons, silent mode |
| **P1** | Tool approval hooks: pre/post execution, audit logging, deny patterns |
| **P2** | CI/CD Agent: multi-env (alpha/beta/gamma/prod), rollback on failure |
| **P3** | Tickets Agent: SLA breach notifications, auto-assignment rules |
| **P4** | HR Agent: PTO tracking, onboarding checklists, team directory |
| **F1** | Chat: edit/retry messages, branch conversations, export transcript |
| **F2** | Ticketing: SLA breach indicators, assignment rules config, filters |
| **F3** | Settings: all config exposed, env management, user management |
| **F4** | Dashboard home: activity feed, quick actions, system health widgets |
| **FS1** | Slack: approval buttons in threads, DM-to-channel escalation |
| **FS2** | Grafana: auto-provisioned dashboards per project, alert rules |
| **D1** | Install script v2: offline mode, proxy support, version pinning |
| **D2** | Backup automation: scheduled snapshots, S3-compatible storage |
| **Q1** | Cron execution tests, snapshot/restore round-trip verification |
| **Q2** | Penetration test: injection attacks, SSRF, path traversal |
| **DS1** | HR portal mockups, onboarding flow UX, empty state designs |
| **DR1** | Video tutorials (5 min each): install, first project, agent customization |

#### Sprint 7 (Weeks 13–14): Hardening

| Team | Deliverables |
|------|-------------|
| **G1** | Performance: benchmark DAG with 1000 nodes, optimize hot paths |
| **G2** | Workspace: quota enforcement, garbage collection, disk usage tracking |
| **G3** | Security: audit log (append-only), HMAC verification, deny-list engine |
| **P1** | MCP: reconnection resilience, server crash recovery, health probes |
| **P2** | Dev Agent: workspace-aware (respects project boundaries), git safety |
| **P3** | Ops Agent: incident response playbook, automated escalation |
| **P4** | Memory: garbage collection, size limits, export/import JSON |
| **F1** | Chat: keyboard shortcuts, accessibility (WCAG 2.1 AA), dark mode |
| **F2** | Wiki: search across all pages, backlinks, table of contents |
| **F3** | Cron: execution history, error logs, manual trigger with params |
| **F4** | Artifact: diff viewer between versions, publish/share URLs |
| **FS1** | Slack: rate limiting, retry queue, graceful degradation |
| **FS2** | GitHub: webhook retry, signature verification, event deduplication |
| **D1** | Production docker-compose: restart policies, health checks, volumes |
| **D2** | Chaos testing: kill agents mid-task, verify checkpoint recovery |
| **Q1** | Performance test suite: 50 concurrent users, 500 tasks/hour |
| **Q2** | Security audit report, CVE scan, dependency license compliance |
| **DS1** | Dark mode design tokens, accessibility audit, icon set finalization |
| **DR1** | Troubleshooting guide, FAQ, community forum setup (GitHub Discussions) |

#### Sprint 8 (Weeks 15–16): 🎯 MILESTONE — Internal Demo #2

| Team | Deliverables |
|------|-------------|
| **G1** | DAG visualization in dashboard (live execution view) |
| **G2** | `mitran init` wizard: interactive project setup, template selection |
| **G3** | Multi-project support: switch contexts, cross-project agent sharing |
| **P1** | Agent marketplace: publish custom agents, version management |
| **P2** | All 3 agents feature-complete with comprehensive tool coverage |
| **P3** | All 3 agents feature-complete, SLA tracking operational |
| **P4** | HR Agent + Wiki Agent feature-complete, memory system stable |
| **F1** | All chat features complete, streaming stable, zero dropped messages |
| **F2** | Ticketing + Wiki fully functional end-to-end |
| **F3** | All admin pages complete and tested |
| **F4** | Artifacts + sessions + memory pages polished |
| **FS1** | Slack integration: full feature parity with dashboard interactions |
| **FS2** | GitHub + monitoring fully wired, dashboards auto-generated |
| **D1** | One-command deploy to any cloud (AWS/GCP/Azure) via Helm |
| **D2** | Production monitoring: PagerDuty/OpsGenie alerting, SLO dashboards |
| **Q1** | Full regression suite (200+ tests), <5 min CI run |
| **Q2** | Load test: 200 concurrent users, identify and fix bottlenecks |
| **DS1** | All UI pages designed and handed off, design system documented |
| **DR1** | Full documentation site live, 10+ tutorial articles |

**Demo #2 scope:** Full workflow — init project → Slack notification → multiple agents collaborate → kanban updates → wiki page generated → monitoring dashboard auto-created → snapshot taken.

---

### Phase 3: Release Readiness (Sprints 9–13, Weeks 17–26)

#### Sprint 9 (Weeks 17–18): Integration Testing & Bug Fixing

| Team | Deliverables |
|------|-------------|
| **G1** | Fix all P0/P1 bugs from Demo #2 feedback, DAG edge cases |
| **G2** | CLI polish: help text, error messages, shell completions verified |
| **G3** | Auth: session management, token refresh, multi-device support |
| **P1** | MCP stability: 72-hour soak test without disconnects |
| **P2** | Agent accuracy: benchmark against 50 real-world tasks, >90% success |
| **P3** | Agent accuracy: benchmark, SLA calculation verified against edge cases |
| **P4** | Memory stability: no data loss under concurrent writes |
| **F1** | Cross-browser testing (Chrome, Firefox, Safari), mobile responsive |
| **F2** | UI performance: <100ms interaction latency, virtualized long lists |
| **F3** | Form validation, error states, loading states for all pages |
| **F4** | Artifact sharing: public URLs, embed codes |
| **FS1** | Slack: edge cases (rate limits, large messages, file uploads) |
| **FS2** | GitHub: large repo handling, monorepo support |
| **D1** | Install script: tested on Ubuntu 22/24, Debian 12, macOS 13/14, AL2023 |
| **D2** | Disaster recovery runbook, backup verification automated |
| **Q1** | Bug bash: all team members test cross-functionally for 2 days |
| **Q2** | Final security audit, fix all HIGH/CRITICAL findings |
| **DS1** | Landing page implementation, marketing assets |
| **DR1** | Launch blog post draft, HN submission prep, ProductHunt listing |

#### Sprint 10 (Weeks 19–20): Performance & Scale

| Team | Deliverables |
|------|-------------|
| **G1** | Optimize: <50ms task dispatch latency, <10ms queue operations |
| **G2** | Large project support: 10,000+ tasks without degradation |
| **G3** | Concurrent auth: 100+ simultaneous users, session pooling |
| **P1** | LLM token optimization: caching, context compression, cost tracking |
| **P2-P4** | Agent response time: <30s for simple tasks, <120s for complex |
| **F1-F4** | Bundle size <500KB gzipped, Lighthouse score >90 |
| **FS1** | Slack: 1000+ messages/min throughput |
| **FS2** | GitHub: webhook processing <2s, no dropped events |
| **D1** | Horizontal scaling guide: multi-node deployment tested |
| **D2** | Performance dashboards: latency, throughput, error rate SLOs |
| **Q1** | Performance regression tests added to CI |
| **Q2** | Load test: 500 concurrent users, 72-hour soak test |
| **DS1** | Performance perception UX: skeleton screens, optimistic updates |
| **DR1** | Performance tuning guide, capacity planning doc |

#### Sprint 11 (Weeks 21–22): Documentation & Developer Experience

| Team | Deliverables |
|------|-------------|
| **G1** | Plugin system: Go plugin interface for custom schedulers |
| **G2** | Migration tools: import from Jira/Linear/Trello (Tickets module) |
| **G3** | API versioning strategy, deprecation headers |
| **P1** | Agent SDK: public Python package for custom agent development |
| **P2-P4** | Agent fine-tuning: prompt optimization based on QA feedback |
| **F1-F4** | Accessibility audit fixes, i18n framework (strings extracted) |
| **FS1** | Slack app directory listing preparation |
| **FS2** | GitHub Marketplace listing preparation |
| **D1** | Terraform modules for AWS/GCP deployment |
| **D2** | Upgrade path: v0.0.x → v0.1.0 migration script |
| **Q1** | User acceptance testing with 5 external beta testers |
| **Q2** | Compliance checklist: GDPR data handling, SOC2 prep |
| **DS1** | Sticker pack, social media assets, conference booth design |
| **DR1** | SDK documentation, agent development tutorial, video series (8 eps) |

#### Sprint 12 (Weeks 23–24): 🎯 MILESTONE — Internal Demo #3 (Release Candidate)

| Team | Deliverables |
|------|-------------|
| **G1** | All engine features frozen, only bug fixes |
| **G2** | CLI feature freeze, man pages generated |
| **G3** | Auth + cron + subagent pool: feature freeze |
| **P1-P4** | All 8 agents frozen, prompt versions pinned |
| **F1-F4** | UI feature freeze, all pages complete |
| **FS1-FS2** | Integration freeze, webhook handlers stable |
| **D1** | Release pipeline: build → test → package → publish (automated) |
| **D2** | Production environment provisioned, DNS configured |
| **Q1** | Full regression pass, zero P0/P1 bugs |
| **Q2** | Final penetration test, sign-off report |
| **DS1** | All assets delivered, landing page live |
| **DR1** | All docs reviewed, launch checklist complete |

**Demo #3 scope:** Full production-ready demo to stakeholders. Install from scratch on clean machine → full workflow → scale test → recovery from failure.

#### Sprint 13 (Weeks 25–26): 🚀 RELEASE — v0.1.0

| Team | Deliverables |
|------|-------------|
| **G1-G3** | Critical bug fixes only, on-call rotation for launch |
| **P1-P4** | Prompt fixes only, monitor agent success rates |
| **F1-F4** | Hotfix capability, error tracking (Sentry) configured |
| **FS1-FS2** | Integration monitoring, webhook health alerts |
| **D1** | Release cut: Docker images tagged, Helm chart published, PyPI/npm packages |
| **D2** | Production deployment, monitoring verified, runbook tested |
| **Q1** | Smoke test on production, deployment verification |
| **Q2** | Post-launch security monitoring active |
| **DS1** | Launch day social assets, product screenshots |
| **DR1** | HN post, ProductHunt launch, Twitter/X announcement, blog live |

**Release day checklist:**
- [ ] All CI green on release branch
- [ ] Docker images published to GHCR
- [ ] Helm chart published
- [ ] Install script points to v0.1.0
- [ ] Documentation site updated
- [ ] Landing page live with download links
- [ ] GitHub release created with changelog
- [ ] Slack/Discord community channels open

---

## 3. Dependency Map

```
┌─────────────────────────────────────────────────────────────┐
│                    CRITICAL PATH                             │
│                                                             │
│  Proto Definitions (TL, Sprint 1)                           │
│       │                                                     │
│       ├──→ gRPC Server (G3) ──→ gRPC Client (P1)           │
│       │         │                     │                     │
│       │         ▼                     ▼                     │
│       │    Auth (G3) ──→ Dashboard Auth (F1)                │
│       │                                                     │
│       ├──→ DAG Scheduler (G1) ──→ Subagent Pool (G3)       │
│       │         │                     │                     │
│       │         ▼                     ▼                     │
│       │    Checkpoint (G1) ──→ Agent Workers (P1-P4)        │
│       │                                                     │
│       └──→ CLI Scaffold (G2) ──→ Install Script (D1)       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Blocking Dependencies (must complete before dependent starts)

| Blocker | Blocked By | Sprint |
|---------|-----------|--------|
| Proto definitions (TL) | All gRPC work | Sprint 1, Week 1 |
| gRPC server (G3) | Python agent framework (P1) | Sprint 1 |
| DAG scheduler (G1) | Subagent pool (G3) | Sprint 2 |
| Auth system (G3) | Dashboard auth gate (F1) | Sprint 1–2 |
| Agent base class (P1) | All 8 agents (P2-P4) | Sprint 1 |
| MCP client (P1) | Tool-using agents (P2-P4) | Sprint 2 |
| Design system (DS1) | All frontend pages (F1-F4) | Sprint 1 |
| Memory module (P4) | Learned corrections, history (P4) | Sprint 2 |
| Skill system (P1) | Per-agent skills (P2-P4) | Sprint 3 |
| Docker images (D1) | Helm charts (D1), install script (D1) | Sprint 2 |
| Ticketing backend (P3) | Ticketing UI (F2) | Sprint 3 |
| Wiki backend (P4) | Wiki editor (F4) | Sprint 5 |
| App manifest system (G1) | App UI (F4), marketplace (FS2) | Sprint 5 |

### Parallel Tracks (no blocking dependencies between them)

- Frontend pages can develop against mock API contracts in parallel with backend
- Slack integration (FS1) and GitHub integration (FS2) are independent
- QA test writing parallels feature development (1 sprint behind)
- DevOps infra (D1, D2) runs independently until integration points
- Designer (DS1) works 1–2 sprints ahead of frontend implementation

---

## 4. Milestones

| Milestone | Sprint | Date (relative) | Success Criteria |
|-----------|--------|-----------------|------------------|
| **Foundation Complete** | 4 | Week 8 | CLI init → chat → agent executes task → dashboard shows result |
| **Feature Complete** | 8 | Week 16 | All 8 agents working, all dashboard pages functional, integrations live |
| **Release Candidate** | 12 | Week 24 | Zero P0 bugs, security audit passed, docs complete, install tested on 4 OS |
| **v0.1.0 Release** | 13 | Week 26 | Public launch, community channels open, landing page live |

### Internal Demo Schedule

| Demo | Audience | Format | Duration |
|------|----------|--------|----------|
| Demo #1 (Sprint 4) | Core team | Live walkthrough | 30 min |
| Demo #2 (Sprint 8) | Extended team + advisors | Full workflow demo | 45 min |
| Demo #3 (Sprint 12) | Stakeholders + beta testers | Production-ready showcase | 60 min |

---

## 5. Risk Register

| # | Risk | Probability | Impact | Mitigation |
|---|------|-------------|--------|------------|
| 1 | **LLM rate limits / costs exceed budget** | High | High | Token caching, prompt compression, cost alerts at 50%/80%/100% of monthly budget. Fallback to smaller models for simple tasks. |
| 2 | **gRPC contract churn delays Python team** | Medium | High | Freeze proto definitions by end of Sprint 1 Week 1. Use backward-compatible additions only after freeze. |
| 3 | **Agent accuracy below 80% on real tasks** | Medium | High | Continuous prompt benchmarking from Sprint 4. Dedicate P2-P4 Sprint 10 entirely to prompt optimization. Skill system allows rapid iteration without code changes. |
| 4 | **MCP server ecosystem immaturity** | Medium | Medium | Build thin MCP adapters in-house for critical tools (filesystem, git, shell). Don't depend on external MCP servers for core functionality. |
| 5 | **Frontend scope creep (15 dashboard pages)** | High | Medium | Design system components reduce per-page effort. F1-F4 pair on complex pages. Feature freeze at Sprint 12 non-negotiable. |
| 6 | **Single point of failure: TL/Architect** | Low | Critical | Document all architecture decisions in ADRs. G1 is backup architect. Weekly architecture sync ensures knowledge sharing. |
| 7 | **Security vulnerability in agent execution** | Medium | Critical | Sandboxed execution from Sprint 1. Q2 continuous security testing. Tool deny-list from Sprint 4. No network access without explicit approval. |
| 8 | **Team member attrition (20% over 6 months)** | Medium | High | Cross-training: each component has primary + backup owner. Comprehensive docs from Sprint 1. |
| 9 | **Docker/Helm complexity blocks self-hosted users** | Medium | Medium | Install script (curl\|sh) as primary path. Docker-compose for advanced users. Test on 4 OS variants every sprint from Sprint 9. |
| 10 | **Slack/GitHub API breaking changes** | Low | Medium | Pin API versions. Abstract behind interface layer. FS1/FS2 monitor changelogs weekly. |
| 11 | **Memory system data loss** | Low | High | WAL mode for SQLite, fsync after writes. Automated backup every hour. Snapshot/restore tested weekly from Sprint 4. |
| 12 | **Launch fails to get traction** | Medium | High | DR1 starts community building from Sprint 7. Beta testers from Sprint 11. Launch on HN + ProductHunt same day. |

---

## 6. Sprint Cadence & Ceremonies

| Ceremony | Frequency | Duration | Attendees |
|----------|-----------|----------|-----------|
| Sprint Planning | Bi-weekly (Monday, Sprint start) | 90 min | All |
| Daily Standup | Daily | 15 min | All (async option via Slack) |
| Demo Day | Bi-weekly (Friday, Sprint end) | 60 min | All |
| Retro | Bi-weekly (Friday, after Demo) | 45 min | All |
| Architecture Sync | Weekly (Wednesday) | 30 min | TL, G1-G3, P1 |
| Frontend Sync | Weekly (Tuesday) | 30 min | F1-F4, DS1 |
| Security Review | Monthly | 60 min | TL, Q2, D2 |

---

## 7. Definition of Done (per Sprint)

- [ ] All code committed and reviewed (≥1 reviewer)
- [ ] Unit tests passing (>80% coverage for new code)
- [ ] Integration tests passing for cross-service features
- [ ] No P0 or P1 bugs open at sprint end
- [ ] Documentation updated for new features
- [ ] Docker build succeeds
- [ ] Demo-able to stakeholders

---

## 8. Budget Estimate (Monthly)

| Category | Cost/Month | Notes |
|----------|-----------|-------|
| AWS Bedrock (Claude Sonnet) | $2,000–5,000 | Development + testing tokens |
| Infrastructure (dev) | $500 | EKS/EC2 for staging |
| CI/CD (GitHub Actions) | $200 | Build minutes |
| Tools & SaaS | $300 | Monitoring, error tracking |
| **Total** | **~$3,000–6,000/mo** | Excluding salary |

---

*Generated: 2026-06-18 | Mitran v0.1.0 Sprint Plan | 20-person team | 26 weeks*
