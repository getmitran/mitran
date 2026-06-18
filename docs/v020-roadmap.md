# Mitran v0.2.0 Roadmap — Deferred Architecture Items

Items identified during senior review of v0.1.0, deferred to v0.2.0.

| # | Item | Why It Matters |
|---|------|----------------|
| 1 | Redis/SQS distributed task queue | Enables multi-instance task distribution; eliminates single-process bottleneck for horizontal scaling |
| 2 | WebSocket pub/sub via Redis | Allows dashboard connections to any instance while receiving events from all; required for HA deployments |
| 3 | OpenTelemetry tracing + Prometheus metrics export | Provides end-to-end request tracing across Go↔Python boundary and production-grade observability |
| 4 | SQLite → PostgreSQL migration for tenant/memory stores | Removes single-writer limitation; enables concurrent access, backups, and multi-node shared state |
| 5 | JWT signature verification via JWKS endpoint | Validates tokens cryptographically against issuer public keys instead of shared-secret HMAC; required for SSO/OIDC flows |

## Context

These were flagged as architectural gaps that don't block single-node self-hosted usage (v0.1.0 target) but are prerequisites for any multi-instance or enterprise deployment. Tracked as GitHub issue #23.

## WebSocket Horizontal Scaling

**Issue:** #24 — Current single-process hub architecture limits dashboard connections to ~50 per instance.

**Target architecture:** Redis Pub/Sub as a cross-instance message bus.

1. Each engine instance maintains a local WebSocket hub for its directly-connected dashboards.
2. On event emission (task update, agent output, memory change), the originating instance publishes to a Redis channel (e.g. `mitran:ws:{project_id}`).
3. All engine instances subscribe to relevant Redis channels and fan-out received messages to their local WebSocket connections.
4. Connection routing is stateless — dashboards connect to any instance behind a load balancer; Redis ensures all instances receive all events.
5. Redis Pub/Sub is fire-and-forget (no persistence); missed messages during reconnect are backfilled via REST `/api/v1/events?since=<timestamp>`.

**Reference flow:** `Dashboard → any Engine instance → Redis Pub/Sub → all Engine instances → fan-out to local WS connections`

This unblocks horizontal scaling of the dashboard tier without requiring sticky sessions or shared WebSocket state.
