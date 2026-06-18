# Changelog

All notable changes to Mitran will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-06-18

### Added
- Go core engine with DAG scheduler, gRPC IPC, and WebSocket real-time updates
- Python agent worker with 5 agent types (Dev, Docs, Ops, Review, CI/CD) on AWS Bedrock
- React + Vite + Tailwind dashboard with 12 pages (Kanban, Memory, Crons, Artifacts, Settings, etc.)
- OpenClaw autonomous runtime (memory, crons, subagents, MCP orchestration)
- CLI with 10 commands (`init`, `start`, `stop`, `status`, `agent`, `task`, `config`, `env`, `logs`, `version`)
- 12 HRMS modules (employees, leave, attendance, payroll, recruitment, onboarding, offboarding, training, performance, analytics, org chart, documents)
- Slack slash commands, GitHub webhook ingestion, Prometheus /metrics endpoint
- Grafana provisioning with pre-built dashboards and alerting rules
- Docker Compose and Helm chart for deployment
- GitHub Actions CI pipeline with lint, test, and build stages
- 180 tasks completed across platform, agents, integrations, enterprise, and DevRel

### Security
- HMAC session cookies with 24h TTL and MITRAN_SESSION_SECRET
- Tamper-proof audit log with HMAC integrity verification
- Role-based access control (RBAC) with per-module permissions
- Per-tenant rate limiting with configurable thresholds
- OIDC/SSO support for enterprise authentication
- Circuit breaker on external LLM calls (Bedrock, Ollama fallback)

### Fixed
- 28 review findings resolved: 3 build-breaking Go import issues, 5 TypeScript type errors, 5 cosmetic issues (naming conflicts, missing imports), and 15 code quality improvements across all modules
