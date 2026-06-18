# Review V2 Backlog

Issues identified by 7-person senior review panel (2026-06-19).

## Scores

| Reviewer | Score |
|----------|-------|
| Security Engineer | 6.7/10 |
| Sr. SDE (Go) | 5.7/10 |
| Product Manager | 7.4/10 |
| Engineering Manager | 6.6/10 |
| Solutions Architect | 6.5/10 |
| QA/Test Engineer | 5.8/10 |
| DevRel Engineer | 7.8/10 |
| **Average** | **6.6/10** |

---

## P0 — Critical (Security/Correctness)

- [ ] **JWT sig not verified in SSO callback** — `sso.go:62` trusts TLS alone, skips signature check. Wire `VerifyJWT()` from `jwks.go` into callback.
- [ ] **Rate limiter race condition** — `ratelimit_apikey.go:38-45` increments `window.count` on shared pointer without mutex. Add per-key sync.Mutex or use atomic ops.
- [ ] **X-User header bypass when MITRAN_ENV unset** — `rbac.go:58-60` allows impersonation if env var missing. Default to production behavior (deny X-User header).
- [ ] **randomHex() not cryptographically random** — `session/manager.go:102-107` uses `time.Now().UnixNano()%16`. Replace with `crypto/rand`.

## P1 — High (Architecture/Quality)

- [ ] **Wire graceful shutdown into main.go** — `shutdown/graceful.go` exists but `main.go` uses raw `log.Fatal(http.ListenAndServe(...))`. Replace with `shutdown.ListenAndServeGraceful()`.
- [ ] **apierr test/impl API mismatch** — Test references `NewNotFound()`, `Status` field; impl has `NotFound()`, `Code` field. Align test to match implementation.
- [ ] **Unify error response shape** — Handlers use `{"error": msg}` while `apierr` uses `{"error": type, "message": msg}`. Standardize on `apierr.WriteError()` everywhere.
- [ ] **Add CODE_OF_CONDUCT.md** — Referenced in CONTRIBUTING.md but doesn't exist. Use Contributor Covenant.
- [ ] **Check rand.Read error in refresh.go** — `refresh.go:19` ignores error from `rand.Read(b)`. Add error check.
- [ ] **Dual queue abstraction** — `TaskQueue` (channel-based) and `PersistentQueue` (interface) coexist. Unify under single interface.
- [ ] **Wire new modules into main.go** — billing, tracing, session store, onboarding, webhooks, workflows defined but never registered in router.

## P2 — Medium (Operations/DX)

- [ ] **Add docs/README.md index** — No single TOC linking all 28+ docs. Create docs index with descriptions.
- [ ] **Add release workflow (CD)** — CI exists but no image build/push, no tag-based release, no CHANGELOG automation.
- [ ] **Restrict Terraform SG to ALB only** — `deploy/terraform/main.tf` opens 7780/8888 to `0.0.0.0/0`. Should be ALB-only.
- [ ] **Port inconsistency in README** — Says `:7780` in one place, `:7777` in another.
- [ ] **Go version mismatch** — ci.yml uses 1.22, coverage.yml uses 1.21, CONTRIBUTING says 1.21+.
- [ ] **Add .github/ISSUE_TEMPLATE and PR_TEMPLATE** — Missing structured contribution templates.
- [ ] **Add CODEOWNERS** — Unclear who reviews what.
- [ ] **FileQueue.rewrite() not atomic** — Truncate-then-write loses data on crash. Use write-to-temp + rename.
- [ ] **Event log unbounded growth** — `audit/eventsource.go` has no rotation/compaction/size limit.
- [ ] **In-memory session revocation lost on restart** — `sessions.go` uses `sync.Map` only. Persist to file or Redis.
- [ ] **CORS fallback allows any localhost** — `cors.go:19` allows `http://localhost*` when no env var set. Dangerous if deployed without config.

## P3 — Low (Polish/Nice-to-have)

- [ ] **No examples/ directory** — No standalone runnable agent examples for newcomers.
- [ ] **No "hello world" e2e tutorial** — First-time user has no guided 5-minute path.
- [ ] **No GIF/screenshot of dashboard** — README would benefit from visual demo.
- [ ] **Handler tests use fake stubs** — Tests mock the router instead of testing real handlers.
- [ ] **40+ top-level packages** — Excessive for v0.1.0, some should merge (errors+apierr, logging+logger).
- [ ] **Duplicate architecture docs** — `10-openclaw-architecture.md` and `v2` variant — unclear which is canonical.
- [ ] **Plugin sandbox no resource limits** — Only timeout, no memory/CPU/FD limits.
- [ ] **No dashboard React tests** — Zero vitest files despite coverage.yml expecting them.
- [ ] **Tracing is entirely no-op** — `Init()` returns empty function, zero production value today.
- [ ] **EnvSecretStore.Set() uses os.Setenv** — Process-global, visible to child processes.

---

## Fix Strategy

1. **P0 first** — 4 items, ~15 min total. Moves Security 6.7→8.5, SDE 5.7→7.5.
2. **P1 next** — 7 items, ~30 min. Moves Architect 6.5→8, DevRel 7.8→9.
3. **P2 batch** — 11 items, ~1 hour. Moves EM 6.6→8.5, QA 5.8→7.
4. **P3 deferred** — Nice-to-have for v0.2.0.

**Target after P0+P1+P2: ~8.5/10 average.**
