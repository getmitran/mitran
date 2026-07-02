# Mitran — MVP Plan (6 Months)

## Team: 5 People

| Person | Role | Months 1-2 | Months 3-4 | Months 5-6 |
|--------|------|-----------|-----------|-----------|
| P1 (Founder) | Product + Sales | Architecture decisions, specs, first outreach | Design partner calls, demo prep | HN launch, first enterprise conversations |
| P2 (Go) | Core Engine | DAG engine, queue, executor, resource locks | Agent runtime, gRPC, checkpoint store | CLI, install script, process management |
| P3 (Python) | Agents | SDK design, Dev Agent, Docs Agent | CI/CD Agent, Tickets Agent, Wiki Agent | Ops Agent, HR Agent, Dashboard Agent |
| P4 (Full-stack) | Dashboard + Apps | React scaffold, priority queue UI, WebSocket | Checkpoint review UI, ticketing UI | Wiki editor, HR portal, Grafana embeds |
| P5 (DevRel) | Community + Brand | README, docs site, branding, demo script | Blog posts, tutorial videos, community Discord | HN post, launch campaign, design partner support |

## Sprint Breakdown

### Month 1: Foundation
**Goal:** Core engine compiles, one agent runs end-to-end.

| Week | P2 (Go) | P3 (Python) | P4 (Full-stack) | P5 (DevRel) |
|------|---------|-------------|-----------------|-------------|
| 1 | Project scaffold, DAG data structures, SQLite schema | Python SDK scaffold, agent base class, tool interface | React + Vite scaffold, WebSocket client | Logo, domain, landing page wireframe |
| 2 | Priority queue (heap), task CRUD API | Dev Agent v0.1 (read files, generate code, produce diff) | Priority queue UI (task cards, drag reorder) | README draft, architecture diagram |
| 3 | Agent runtime (spawn Python subprocess, gRPC channel) | Integration with core via gRPC, file ACL enforcement | Checkpoint view (show diff, approve/reject buttons) | Docs site (Docusaurus/MkDocs) |
| 4 | Resource lock manager, executor loop | End-to-end: create task → agent runs → checkpoint → approve | Real-time updates via WebSocket | Demo recording (internal) |

**Deliverable:** `mitran serve` starts, create a task, Dev Agent runs, checkpoint appears, human approves.

### Month 2: Multi-Agent + Integrations
**Goal:** 3 agents working, GitHub + Slack integrated.

| Week | P2 (Go) | P3 (Python) | P4 (Full-stack) | P5 (DevRel) |
|------|---------|-------------|-----------------|-------------|
| 5 | GitHub integration (OAuth, repo creation, PR API) | Docs Agent v0.1, CI/CD Agent v0.1 | Agent status panel, progress bars | Landing page live, waitlist form |
| 6 | Slack integration (bot, slash commands, notifications) | Slack tool for agents, GitHub tool for agents | Slack notification for checkpoints | First blog post: "What we're building" |
| 7 | `mitran init` flow (interactive team description) | Orchestrator logic: plan tasks from team description | Init wizard UI in dashboard | Community Discord/Slack setup |
| 8 | Event bus (NATS embedded), webhook handler | Refine all 3 agents based on testing | Event feed in dashboard | Reach out to 10 potential design partners |

**Deliverable:** `mitran init` generates a basic plan, 3 agents execute it, Slack notifications work.

### Month 3: All Agents + Ticketing
**Goal:** All 8 agents functional (basic). Own ticketing system.

| Week | P2 (Go) | P3 (Python) | P4 (Full-stack) | P5 (DevRel) |
|------|---------|-------------|-----------------|-------------|
| 9 | Ticketing API (CRUD, sprints, labels, search) | Tickets Agent, Wiki Agent | Ticketing board UI (kanban + list view) | Tutorial: "Set up CI/CD with Mitran" |
| 10 | Wiki API (pages, tree structure, search via Bleve) | Ops Agent (Prometheus config gen, Grafana dashboards) | Wiki editor (markdown + preview) | Tutorial: "Agent-managed monitoring" |
| 11 | Prometheus + Grafana child process management | HR Agent (PTO CRUD, onboarding checklists) | HR portal (PTO calendar, request form) | Design partner #1 onboarding |
| 12 | Dashboard Agent integration (Grafana API) | Dashboard Agent (auto-create dashboards from metrics) | Grafana embed in Mitran UI | Collect feedback from design partner |

**Deliverable:** All 8 agents operational. Ticketing, wiki, HR built-in. Monitoring installed.

### Month 4: Polish + Reliability
**Goal:** Production-quality for design partners. Install experience smooth.

| Week | P2 (Go) | P3 (Python) | P4 (Full-stack) | P5 (DevRel) |
|------|---------|-------------|-----------------|-------------|
| 13 | Install script (`curl \| sh`), auto-detect OS/arch | Agent error handling, retry logic, graceful degradation | Error states in UI, loading skeletons, mobile responsive | Onboard design partner #2 and #3 |
| 14 | Config file (YAML), `mitran config` command | Agent test suite (unit + integration) | Settings page, agent configuration UI | Write comparison blog: "Mitran vs CrewAI vs Temporal" |
| 15 | Backup/restore (`mitran backup`, `mitran restore`) | LLM provider abstraction (Claude/GPT-4/local) | Onboarding wizard (first-time setup in UI) | Demo video for website |
| 16 | Security: audit log, path ACL enforcement tests | Agent quality pass: reduce hallucinations, improve output | Audit log viewer in dashboard | Press kit, screenshots, demo GIF |

**Deliverable:** 3 design partners running Mitran daily. No critical bugs. Install takes < 5 min.

### Month 5: Enterprise Features + Scale
**Goal:** Enterprise tier ready for beta.

| Week | P2 (Go) | P3 (Python) | P4 (Full-stack) | P5 (DevRel) |
|------|---------|-------------|-----------------|-------------|
| 17 | SSO/SAML integration (enterprise auth) | Workflow mode: define templates ("on PR → review → docs") | SSO login flow, team management UI | Enterprise landing page |
| 18 | PostgreSQL backend (enterprise storage) | Agent memory: learn from past corrections | Workflow builder UI (visual) | Enterprise pricing page |
| 19 | Multi-user support (roles: admin, operator, viewer) | Custom agent template: SDK + docs + test harness | Role-based UI (hide admin panels for viewers) | Start enterprise outreach |
| 20 | Rate limiting, usage tracking, billing hooks | Local app install (`mitran app install ./`) | App management page | App developer documentation |

**Deliverable:** Enterprise-ready beta. SSO works. Multi-user tested. Billing framework in place.

### Month 6: Launch
**Goal:** Public launch. HN front page. First revenue.

| Week | P2 (Go) | P3 (Python) | P4 (Full-stack) | P5 (DevRel) |
|------|---------|-------------|-----------------|-------------|
| 21 | Performance testing, load testing with 8 concurrent agents | Final agent quality pass, edge case fixes | Performance optimization, bundle size | Launch plan: HN post draft, tweets, blog |
| 22 | CLI polish (`mitran --help` experience) | Marketplace scaffold (publish, discover, install) | Marketplace browse UI (Phase 2 preview) | Pre-launch: tell design partners, get testimonials |
| 23 | Final bug fixes, security audit | Documentation for all agents | Final UI polish, dark mode, accessibility | **LAUNCH DAY**: HN, Twitter, Reddit, Discord |
| 24 | Post-launch bug fixes, scale support | Community bug fixes, first external contributions | Post-launch UX improvements | Support incoming users, close enterprise leads |

**Deliverable:** Public on GitHub. 3K+ stars. 3 enterprise design partners. First consulting revenue.

## Success Metrics (Month 6)

| Metric | Target |
|--------|--------|
| GitHub stars | 3,000+ |
| Self-hosted installs | 500+ |
| Enterprise design partners | 3-5 |
| Revenue (consulting) | $50K |
| Active community (Discord) | 200+ |
| Agents operational | 8/8 |
| Install success rate | > 90% |
| `mitran init` completion rate | > 80% |

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|-----------|
| 8 agents too ambitious for 1 Python engineer in 6mo | HIGH | Ship half-baked agents | Prioritize Dev + CI/CD + Ops first. Others can be minimal viable. |
| LLM reliability for complex tasks | MEDIUM | Users don't trust output | Conservative checkpointing. Agent asks human when confidence < 80%. |
| Install fails on diverse environments | MEDIUM | Bad first impression | Test on Ubuntu 22/24, macOS 13/14/15, Fedora. Docker fallback. |
| No enterprise customers in 6 months | MEDIUM | Revenue delayed | Consulting revenue fills gap. Lower enterprise price for early adopters. |
| Competitor launches similar product | LOW | Positioning challenged | Speed + open-source + community moat. Ship first, improve forever. |
