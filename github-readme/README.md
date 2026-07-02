<p align="center">
  <h1 align="center">mitran</h1>
  <p align="center"><strong>Your agents build. You decide.</strong></p>
  <p align="center">Open-source platform where AI agents build and maintain your entire team infrastructure.</p>
</p>

<p align="center">
  <a href="https://getmitran.vercel.app">Website</a> ·
  <a href="#quickstart">Quickstart</a> ·
  <a href="#how-it-works">How it works</a> ·
  <a href="docs/">Documentation</a> ·
  <a href="https://discord.gg/mitran">Discord</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/license-AGPL--3.0-blue" alt="License">
  <img src="https://img.shields.io/badge/status-alpha-orange" alt="Status">
  <img src="https://img.shields.io/badge/go-1.22+-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/python-3.11+-3776AB?logo=python" alt="Python">
</p>

---

## What is Mitran?

Run `mitran init`, describe your team, and in **30 minutes** get:

- ⚙️ **CI/CD pipelines** — GitHub Actions, configured and running
- 📝 **Documentation** — API docs, architecture docs, runbooks
- 🎫 **Ticketing** — Sprint boards, triage, SLA tracking (built-in)
- 📊 **Monitoring** — Prometheus + Grafana, dashboards, alerting
- 📚 **Wiki** — Knowledge base, onboarding guides, SOPs
- 👥 **HR portal** — PTO tracking, onboarding checklists
- 📈 **Dashboards** — Team metrics, service health

All built by **8 specialized AI agents**. You control the priority queue and approve every output.

## How it works

```
You: "Set up CI/CD for our Go microservices"
         │
         ▼
┌─ Priority Queue ────────────────────────┐
│  1. [CI/CD Agent] Create GitHub Actions  │
│  2. [Docs Agent] Write deployment docs   │
│  3. [Ops Agent] Add monitoring           │
└──────────────────────────────────────────┘
         │
         ▼  Agent executes
┌─ Checkpoint ────────────────────────────┐
│  "Created 3 workflow files. Tests pass.  │
│   Ready for review."                     │
│                                          │
│  [✓ Approve]  [✗ Reject]  [✏️ Edit]     │
└──────────────────────────────────────────┘
```

**Key design decisions:**
- **Sync execution** — One agent works at a time per resource. Zero conflicts.
- **Human priority queue** — You decide what gets done next, not the AI.
- **Checkpoints** — Every output is reviewed before it takes effect.
- **DAG scheduling** — Tasks respect dependencies automatically.

## Quickstart

```bash
# Install (single binary, no Docker required)
curl -sSL https://install.getmitran.dev | sh

# Set up integrations
mitran setup

# Initialize your infrastructure
mitran init
```

## Architecture

```
┌──────────────────────────────────────────┐
│           Human Operator                  │
│  UI Queue · Natural Language · Workflows  │
└────────────────────┬─────────────────────┘
                     │
┌────────────────────▼─────────────────────┐
│          Core Engine (Go)                 │
│  DAG · Priority Queue · Resource Locks   │
│  Checkpoints · Events · Integrations     │
└────────────────────┬─────────────────────┘
                     │ gRPC
┌────────────────────▼─────────────────────┐
│          Agent Layer (Python SDK)         │
│  Dev · Docs · CI/CD · Tickets · Wiki     │
│  Ops · HR · Dashboard                    │
└────────────────────┬─────────────────────┘
                     │
┌────────────────────▼─────────────────────┐
│         Integrations                      │
│  GitHub · Slack · Grafana · Prometheus   │
└──────────────────────────────────────────┘
```

## Agents

| Agent | Responsibility |
|-------|---------------|
| **Dev** | Code generation, refactoring, PR creation, test writing |
| **Docs** | Technical documentation, API refs, ADRs, runbooks |
| **CI/CD** | Pipeline creation, build optimization, deployment |
| **Tickets** | Issue tracking, sprint planning, triage, SLA tracking |
| **Wiki** | Knowledge base, onboarding, SOPs, cross-referencing |
| **Ops** | Monitoring setup, alert configuration, incident diagnosis |
| **HR** | PTO management, onboarding checklists, policy Q&A |
| **Dashboard** | Grafana dashboards, metric identification, reporting |

## Integrations

| Service | How |
|---------|-----|
| **GitHub** | REST API + Webhooks — repos, PRs, Actions |
| **Slack** | Bot + commands — notifications, approvals, NL interface |
| **Grafana** | HTTP API — dashboards, data sources (Apache 2.0, bundled) |
| **Prometheus** | Config generation — metrics, alerting (Apache 2.0, bundled) |

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Core Engine | Go 1.22+ |
| Agent SDK | Python 3.11+ |
| Dashboard | TypeScript + React |
| Database | SQLite (local) / PostgreSQL (enterprise) |
| Events | NATS (embedded) |
| Monitoring | Prometheus + Grafana (bundled) |
| Search | Bleve (embedded full-text) |

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

```bash
# Clone and build
git clone https://github.com/getmitran/mitran.git
cd mitran
go build ./cmd/mitran
```

## License

AGPL-3.0. See [LICENSE](LICENSE) for details.

Your infrastructure, your data, your rules.
