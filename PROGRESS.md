# Mitran — Progress Tracker

## Phase: Ideation ← WE ARE HERE

### Completed
- [x] Idea validation (2026-06-16)
- [x] Brand name decided: **Mitran**
- [x] C-Suite review (CEO, CTO, CFO, Critic)
- [x] Integration strategy decided (GitHub + Slack integrate, ticketing own, OSS tools install)
- [x] Tech stack decided (Go core + Python SDK + TS dashboard)
- [x] Team size decided (start 5, scale to 10)
- [x] Business model decided (AGPL + Enterprise tier)
- [x] Vision document
- [x] High-Level Design (HLD)
- [x] Low-Level Design (LLD)
- [x] Pitch Deck (12 slides)
- [x] MVP Plan (6-month sprint breakdown)
- [x] Competitive Analysis

### Completed (2026-06-16)
- [x] Landing page live (https://getmitran.vercel.app)
- [x] GitHub org + repo (https://github.com/getmitran/mitran)
- [x] Go DAG engine prototype (scheduler, queue, locks — tests pass)
- [x] Go CLI (`mitran init` interactive flow — working)
- [x] Python Agent SDK + 3 standalone agents (tests pass)
- [x] Go Core Engine HTTP server (REST API, scheduler loop, JSON persistence)
- [x] Python Agent Worker (8 agents, AWS Bedrock Claude, real LLM calls)
- [x] React Dashboard (live API polling, InitWizard, approve/reject)
- [x] Multi-environment CI/CD (configurable environments, pipeline templates)
- [x] Running Guide documentation
- [x] Development SOP

### Next Steps
- [ ] Domain registration (getmitran.dev)
- [ ] Logo & branding
- [ ] WebSocket for real-time dashboard updates
- [ ] Test full end-to-end flow with Bedrock credentials
- [ ] GitHub Actions CI for the repo itself
- [ ] Demo video for HN launch
- [ ] Discord community setup
- [ ] Identify first 3 design partners
- [ ] Wire CLI → Engine (CLI calls API instead of simulating)

## Key Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-06-16 | Name: Mitran | Sanskrit "friend/ally", unique in AI space, CLI-friendly |
| 2026-06-16 | All 8 agents in MVP | Founder's conviction — "zero-to-deployed in 30 min" needs all verticals |
| 2026-06-16 | Integrate GitHub + Slack, build own ticketing | Don't compete with best-in-class SCM/chat. Own the workflow layer. |
| 2026-06-16 | Install OSS tools (Grafana, Prometheus) | Apache 2.0 licensed, free to bundle, industry standard |
| 2026-06-16 | AGPL license | Prevents cloud providers from stripping without contributing |
| 2026-06-16 | Go core | Proven for orchestration, single binary, fast enough |
| 2026-06-16 | Start with 5, scale to 10 | CFO recommendation: reduce burn, prove before scaling |
| 2026-06-16 | Self-hosted → Cloud later | GitLab playbook. Build trust with self-hosted, monetize with cloud. |
