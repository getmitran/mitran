# Review V4 — Issues to reach 9/10

Current avg: 8.2/10. Target: 9.0/10 from all 7 reviewers.

## Blocking 9/10 (by reviewer)

### 🔒 Security (8.5 → 9.0) — needs +0.5
- [ ] Fix `csp.go` rand.Read error silently ignored — should fail request
- [ ] Fix `go.mod` declares go 1.25 (doesn't exist) — change to 1.22
- [ ] Add JWKS cache (currently fetches on every SSO callback)

### 🧑‍💻 SDE (7.8 → 9.0) — needs +1.2
- [ ] Migrate remaining handlers to apierr (batch.go, github.go, onboarding.go)
- [ ] Wire ValidateJSON middleware to task create/update routes
- [ ] Add session_test.go for manager (currently only Store tested)
- [ ] Bump randomHex(4) → randomHex(8) for session IDs

### 📦 PM (8.2 → 9.0) — needs +0.8
- [ ] Add quickstart-tutorial to docs/README.md TOC
- [ ] Add requirements.txt to examples/hello-agent
- [ ] Mention dashboard URL in quickstart tutorial

### 📋 EM (8.5 → 9.0) — needs +0.5
- [ ] Fix go.mod version: 1.25 → 1.22
- [ ] Add Trivy filesystem/dependency scan (not just Dockerfile config scan)

### 🏗️ Architect (8.0 → 9.0) — needs +1.0
- [ ] Add f.Sync() before os.Rename in FileQueue rewrite
- [ ] Set Content-Type application/json on /healthz endpoint
- [ ] Add background reaper goroutine for expired MemoryStore entries

### 🧪 QA (7.5 → 9.0) — needs +1.5
- [ ] Add 5+ more dashboard component tests (Settings, Kanban, AgentStatus, Chat, Cron)
- [ ] Add `go test -coverprofile` step to CI with 40% minimum threshold
- [ ] Fix benchmark_test.go — currently 0 runnable benchmarks

### 📣 DevRel (9.0 → 9.0) — already at target ✅
- [ ] Fix quickstart link to v2 arch doc (minor polish)

## Priority Order
1. go.mod fix (blocks EM + Security + release workflow)
2. SDE items (biggest gap: +1.2 needed)
3. QA items (second biggest: +1.5 needed)
4. Architect items (+1.0)
5. PM/DevRel polish
