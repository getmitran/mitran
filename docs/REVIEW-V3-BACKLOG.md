# Review V3 Backlog — Final Polish

## P3 Polish (6)
- [ ] 1. Rewrite handler tests to use real router (not fake stubs)
- [ ] 2. Merge duplicate packages (errors→apierr, logging→logger)
- [ ] 3. Plugin sandbox resource limits (memory/CPU/FD via ulimit)
- [ ] 4. Expand dashboard React tests (add tests for key pages)
- [ ] 5. Wire OTLP SDK in tracing when OTEL_EXPORTER_OTLP_ENDPOINT is set
- [ ] 6. Fix EnvSecretStore.Set() — don't use os.Setenv

## Architectural Wiring (3)
- [ ] 7. Wire billing/tracing/session/onboarding/webhooks/workflows into main.go
- [ ] 8. Unify error responses — all handlers use apierr.WriteError()
- [ ] 9. Unify queue abstraction (remove TaskQueue, use PersistentQueue)

## Frontend Polish (3)
- [ ] 10. Group sidebar nav items into collapsible sections
- [ ] 11. Add dashboard screenshot to README
- [ ] 12. Create "Hello World" quickstart tutorial in docs

## CI/CD (2)
- [ ] 13. Add container scanning (Trivy) to CI workflow
- [ ] 14. Document branch protection rules

## Testing Gaps (5)
- [ ] 15. Tests for auth package (SSO, JWT verify, refresh)
- [ ] 16. Tests for session/manager.go
- [ ] 17. Tests for queue (persistent + task queue)
- [ ] 18. Tests for audit (HMAC chain + event sourcing)
- [ ] 19. Tests for graceful shutdown

## Deferred
- [ ] 20. TUI interface (mitran tui) — deferred to v0.2.0
