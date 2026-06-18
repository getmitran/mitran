# Mitran

> Agent-native company infrastructure platform. Your AI team runs dev, ops, tickets, wiki, and more — so humans focus on decisions.

[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)](https://go.dev)
[![Python](https://img.shields.io/badge/Python-3.11+-3776AB.svg)](https://python.org)

## What is Mitran?

Mitran (Sanskrit/Tamil: "friend, ally") is an open-source platform where AI agents handle your company's internal tooling — development workflows, operations, ticketing, documentation, and CI/CD — as a coordinated team. You describe your company, run `mitran init`, and get a fully operational agent workforce in 30 minutes.

## Features

- **5 Core Agents** — Dev, Ops, Tickets, Wiki, HR (with CI/CD, Review, and Chat agents)
- **OpenClaw Runtime** — Persistent memory, scheduled jobs, subagent orchestration, self-learning
- **MCP Support** — Connect any Model Context Protocol server for extensible tool use
- **Multi-LLM** — AWS Bedrock (Claude), OpenAI, Anthropic, local models via configurable providers
- **DAG Engine** — Topological task scheduling with dependency resolution and parallel execution
- **gRPC IPC** — Go orchestrator ↔ Python workers communicate via Protocol Buffers
- **React Dashboard** — Real-time WebSocket updates, Kanban boards, agent status, chat interface
- **Self-Hosted First** — AGPL-3.0, runs on your infrastructure, no data leaves your network

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    React Dashboard (:5173)               │
│              WebSocket + REST API Consumer               │
└───────────────────────────┬─────────────────────────────┘
                            │ HTTP/WS
┌───────────────────────────▼─────────────────────────────┐
│                   Go Engine (:7780)                      │
│  REST API · DAG Scheduler · Memory · Auth · WebSocket   │
└───────────────────────────┬─────────────────────────────┘
                            │ gRPC (protobuf)
┌───────────────────────────▼─────────────────────────────┐
│              Python Agent Workers (:8888)                │
│  8 Agents · LLM Orchestration · MCP Client · Tools      │
└───────────────────────────┬─────────────────────────────┘
                            │
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
         AWS Bedrock    MCP Servers    Integrations
         (Claude)       (user-managed) (Slack/GitHub)
```

## Quick Start (Docker Compose)

```bash
git clone https://github.com/getmitran/mitran.git
cd mitran
cp .env.example .env  # Add your LLM API keys
docker compose up
```

Dashboard: http://localhost:5173 · API: http://localhost:7780 · Agents: http://localhost:8888

## Manual Setup

### Prerequisites

- Go 1.22+
- Python 3.11+
- Node.js 18+
- protoc (Protocol Buffers compiler)

### Install & Run

```bash
# 1. Engine (Go)
cd prototype/server
go mod download
go run main.go

# 2. Agent Workers (Python)
cd prototype/agent-worker
pip install -r requirements.txt
python -m worker.main

# 3. Dashboard (React)
cd prototype/dashboard
npm install
npm run dev
```

## CLI Usage

```bash
# Initialize a new project
mitran init --name "Acme Corp"

# Start all services
mitran up

# Check agent status
mitran status

# Run a task
mitran run "Set up CI/CD for the payments service"

# View logs
mitran logs --agent dev
```

## Configuration

Environment variables (`.env`):

| Variable | Description | Default |
|----------|-------------|---------|
| `MITRAN_ENGINE_PORT` | Go engine port | `7780` |
| `MITRAN_WORKER_PORT` | Python worker port | `8888` |
| `MITRAN_LLM_PROVIDER` | LLM backend | `bedrock` |
| `AWS_REGION` | AWS region for Bedrock | `us-east-1` |
| `MITRAN_MEMORY_BACKEND` | Storage backend | `json` |
| `MITRAN_LOG_LEVEL` | Log verbosity | `info` |

Per-project settings live in `.mitran/projects/<id>/settings.json`.

## Project Structure

```
prototype/
├── server/          # Go engine (REST, DAG, memory, auth)
├── agent-worker/    # Python agents (LLM, MCP, tools)
├── dashboard/       # React + Vite + Tailwind
├── proto/           # gRPC/protobuf definitions
├── cli/             # Go CLI binary
└── agents/          # Agent configurations
docs/                # Architecture, design, sprint plans
docker-compose.yml   # Full stack deployment
Makefile             # Build shortcuts
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, coding standards, and PR guidelines.

We welcome contributions of all kinds — bug fixes, new agents, integrations, docs, and tests.

## License

[AGPL-3.0](LICENSE) — Free to use, modify, and self-host. Enterprise licensing available for managed hosting and custom SLA.

---

Built with 🤝 by the Mitran community · [getmitran.dev](https://getmitran.vercel.app)
