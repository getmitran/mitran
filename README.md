# Mitran

> Agent-native team infrastructure. AI agents run your dev, ops, and docs — humans focus on decisions.

[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)](https://go.dev)
[![Python](https://img.shields.io/badge/Python-3.11+-3776AB.svg)](https://python.org)

## What is Mitran?

Mitran (Sanskrit/Tamil: "friend") is an open-source platform where AI agents operate as your team's internal tooling layer — development, operations, testing, documentation, and project management — coordinated by a PM orchestrator. Run `mitran init --team` and get a working agent workforce in 30 minutes.

## Agents

**7 Primary Agents** (PM orchestrator + 6 specialists):

| Agent | Role |
|-------|------|
| **PM** | Orchestrates tasks, assigns work, tracks progress |
| **Developer** | Writes code, creates PRs, manages branches |
| **Architect** | System design, ADRs, dependency analysis |
| **DevOps** | CI/CD, deployments, infrastructure, monitoring |
| **Tester** | Unit/integration/E2E tests, coverage analysis |
| **Docs** | Documentation, changelogs, API references |
| **Security** | Vulnerability scanning, policy enforcement, audits |

Plus 9 legacy agents (Dev, Ops, Tickets, Wiki, HR, CI/CD, Review, Chat, Support) for backward compatibility.

## Key Features

- **Memory** — Persistent facts, lessons, and episodic recall per project
- **Crons** — Scheduled agent jobs with 5-field cron expressions and intervals
- **Artifacts** — Versioned content store (50-version retention, slug-based)
- **DAG Engine** — Topological task scheduling with parallel execution
- **MCP Support** — Connect any Model Context Protocol server
- **Multi-LLM** — Bedrock, OpenAI, Anthropic, Ollama (local)
- **gRPC IPC** — Go orchestrator ↔ Python workers via protobuf
- **React Dashboard** — Real-time WebSocket, Kanban, chat, monitoring

## Architecture

```
┌───────────────────────────────────────────────────┐
│         React Dashboard (:5173)                   │
│         WebSocket + REST Consumer                 │
└────────────────────┬──────────────────────────────┘
                     │ HTTP/WS
┌────────────────────▼──────────────────────────────┐
│              Go Engine (:7780)                     │
│  REST API · DAG · Memory · Crons · Artifacts      │
│  Auth · WebSocket · PM Orchestrator               │
└────────────────────┬──────────────────────────────┘
                     │ gRPC
┌────────────────────▼──────────────────────────────┐
│          Python Agent Workers (:8888)             │
│  7 Primary + 9 Legacy Agents · LLM · MCP · Tools │
└────────────────────┬──────────────────────────────┘
                     │
         ┌───────────┼───────────┐
         ▼           ▼           ▼
    AWS Bedrock   MCP Servers  Integrations
    (Claude)      (extensible) (Slack/GitHub)
```

## Quick Start

```bash
git clone https://github.com/getmitran/mitran.git && cd mitran
cp .env.example .env   # Add your LLM API keys
docker compose up
```

Dashboard: http://localhost:5173 · API: http://localhost:7780

### Init Team Flow

```bash
mitran init --team "Acme Corp" --size 5
# → Provisions all 7 agents with project context
# → Creates .mitran/ workspace with memory, cron configs
# → Ready in ~30 minutes
```

## CLI

```bash
mitran init --team "Name"     # Initialize agent team
mitran up                      # Start all services
mitran status                  # Agent health + task queue
mitran run "task description"  # Dispatch to PM orchestrator
mitran logs --agent developer  # Stream agent logs
mitran cron list               # View scheduled jobs
mitran memory search "query"   # Search agent memory
```

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `MITRAN_ENGINE_PORT` | Go engine port | `7780` |
| `MITRAN_WORKER_PORT` | Python worker port | `8888` |
| `MITRAN_LLM_PROVIDER` | LLM backend | `bedrock` |
| `AWS_REGION` | AWS region for Bedrock | `us-east-1` |
| `MITRAN_MEMORY_BACKEND` | Storage (json/sqlite) | `sqlite` |
| `MITRAN_LOG_LEVEL` | Log verbosity | `info` |

## Project Structure

```
prototype/
├── server/          # Go engine (API, DAG, memory, crons, artifacts, auth)
├── agent-worker/    # Python agents (LLM, MCP, tools)
├── dashboard/       # React + Vite + Tailwind
├── proto/           # gRPC/protobuf definitions
├── cli/             # Go CLI binary
└── agents/          # Agent configurations (7 primary + 9 legacy)
deploy/
├── k8s/             # Kubernetes manifests
├── helm/            # Helm chart
└── terraform/       # Infrastructure as code
docs/                # Architecture, design, roadmap
```

## Deployment

```bash
# Kubernetes
kubectl apply -f deploy/k8s/

# Docker Compose (production)
docker compose -f docker-compose.prod.yml up -d
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and PR guidelines.

## License

[AGPL-3.0](LICENSE) — Free to use, modify, and self-host.

---

Built with 🤝 by the Mitran community · [getmitran.dev](https://getmitran.vercel.app)
