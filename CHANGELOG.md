# Changelog

All notable changes to this project will be documented in this file.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.2.0] - 2026-07-02

### Added
- **7 Primary Agents** — Developer, Architect, DevOps, Tester, PM, Docs, Security (replacing generic 5-agent model)
- **PM Orchestrator** — Central task routing with intelligent agent assignment
- **Memory System** — Persistent facts, lessons, and episodic recall with LIKE-based search per project
- **Cron Scheduler** — Recurring and one-shot jobs with 5-field cron expressions and interval support
- **Artifacts Store** — Versioned content storage with slug-based identity (50-version retention)
- **`mitran init --team`** — One-command team provisioning with project context setup
- **SQLite persistence** — 13-table schema replacing JSON file store (via modernc.org/sqlite, pure Go)
- **SSE real-time log streaming** — `/api/v1/events` endpoint wired to dispatcher
- **v0.2.0 Roadmap** — 6-phase, 12-week plan covering agent intelligence through community launch
- Kubernetes deployment manifests with ConfigMap-based configuration
- CHANGELOG.md, improved issue templates, community documentation

### Changed
- Agent architecture from 5 generic to 7 focused specialists (legacy agents preserved)
- Brand scope from "company platform" to "team infrastructure" (init team, not init company)
- Dashboard theme to MeshClaw Solarized with dark/light toggle
- Memory backend default from `json` to `sqlite`

### Fixed
- All P0/P1/P2 issues from V5 senior review (108 http.Error migrations, middleware wiring, rate-limit key spoofing)
- Dashboard placeholder pages replaced with functional components

## [0.1.0] - 2026-06-19

### Added
- **Go Engine** — REST API server on port 7780 with DAG task scheduler
- **Python Agent Workers** — 8 agents (Dev, Docs, Ops, Review, HR, CI/CD, Tickets, Wiki) via gRPC
- **React Dashboard** — 32 pages with real-time WebSocket, Kanban boards, chat interface
- **CLI** — `mitran serve`, `chat`, `status`, `config`, `agent`, `workspace`, `cron`, `snapshot`, `app`
- **OpenClaw Runtime** — Subagent orchestration, MCP client, tool dispatch
- **Auth** — HMAC session cookies, JWT API auth, 24h TTL
- **DAG Engine** — Topological sort with cycle detection, parallel TaskGroup execution
- **Apps Ecosystem** — App manifest spec, installer, registry API, template generator
- **HR Portal** — 12-module HRMS (leave, org chart, onboarding, attendance, payroll, performance)
- **Settings API** — Per-project and per-module configuration with path overrides
- **Docker Compose** — Full stack deployment (engine + worker + dashboard)
- **CI/CD** — GitHub Actions for lint, test, build, Docker image push
- **Documentation** — Architecture docs, sprint plans, quickstart tutorial, development SOP

### Infrastructure
- gRPC/protobuf IPC between Go and Python
- WebSocket real-time updates for dashboard
- Gunicorn 4 workers with circuit breaker on Bedrock
- Per-tenant rate limiting
- Atomic JSON writes for data safety
- ASCII startup banner with route table

[0.2.0]: https://github.com/getmitran/mitran/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/getmitran/mitran/releases/tag/v0.1.0
