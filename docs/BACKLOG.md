# Mitran v0.2.0+ Backlog

Items noted by senior reviewers but deferred from v0.1.0.

## Security (Security Engineer)

- [ ] Full RSA signature verification for JWKS (currently best-effort key fetch only)
- [ ] Token rotation / refresh token flow
- [ ] Rate limiting per-API-key (not just per-tenant)
- [ ] CSP headers with nonce-based script loading for dashboard
- [ ] Secrets management integration (Vault / AWS Secrets Manager)
- [ ] mTLS between engine and worker in production

## Architecture (Architect)

- [ ] Replace in-memory session denylist with Redis-backed store for horizontal scaling
- [ ] Replace in-memory task queue with persistent queue (Redis/NATS/SQS)
- [ ] Database migrations framework (golang-migrate or similar)
- [ ] Graceful shutdown with in-flight request draining
- [ ] OpenTelemetry distributed tracing (replace ad-hoc logging)
- [ ] Event sourcing for audit trail (append-only log beyond HMAC file)
- [ ] Plugin sandboxing via WASM or gRPC isolation

## SDE (Senior SDE)

- [ ] Comprehensive error types (typed errors instead of string messages)
- [ ] Request validation middleware (JSON schema or struct tags)
- [ ] Pagination for all list endpoints
- [ ] Integration test suite with testcontainers
- [ ] Benchmark suite for hot paths (task dispatch, queue ops)
- [ ] Structured logging (slog or zerolog replacing log.Printf)

## Product (PM)

- [ ] Multi-model cost tracking dashboard with spend alerts
- [ ] Agent marketplace / plugin registry
- [ ] Webhook integrations (Slack, Teams, PagerDuty)
- [ ] Custom workflow builder UI
- [ ] Usage analytics and billing per workspace
- [ ] Self-service onboarding wizard in dashboard

## Operations (Project Manager)

- [ ] Helm chart values for all env vars (currently partial)
- [ ] Terraform module for cloud deployment
- [ ] Runbook for common failure modes
- [ ] SLO definitions and error budget tracking
- [ ] Automated canary deployments via Argo Rollouts
- [ ] Backup/restore procedures for persistent state

## DevEx

- [ ] TUI interface for `mitran` CLI (#119, deferred)
- [ ] `mitran upgrade` self-update command
- [ ] Shell completions (bash/zsh/fish)
- [ ] VS Code extension for agent development
- [ ] SDK packages (Go, Python, TypeScript client libraries)
