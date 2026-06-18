# Contributing to Mitran

Thank you for your interest in contributing to Mitran! This guide will help you get started.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md). We are committed to providing a welcoming and inclusive experience for everyone.

## Development Setup

### Prerequisites

- Go 1.21+
- Python 3.11+
- Node.js 20+
- Docker (for local services)

### Clone & Install

```bash
git clone https://github.com/getmitran/mitran.git
cd mitran

# Go engine
cd prototype/server
go mod download

# Python agent worker
cd ../agent-worker
python3 -m venv .venv && source .venv/bin/activate
pip install -e ".[dev]"

# Dashboard
cd ../dashboard
npm install
```

### Running Locally

```bash
# Start all services
make dev

# Or individually:
cd prototype/server && go run .
cd prototype/agent-worker && python3 -m worker
cd prototype/dashboard && npm run dev
```

## Branch Naming

| Prefix | Purpose |
|--------|---------|
| `feat/` | New features |
| `fix/` | Bug fixes |
| `docs/` | Documentation changes |
| `refactor/` | Code restructuring |
| `test/` | Adding/fixing tests |

Example: `feat/add-slack-integration`

## Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

[optional body]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `ci`

- Use imperative mood ("add" not "added")
- No period at end of subject
- Subject ≤ 50 characters

## Pull Request Process

1. Fork the repo and create your branch from `main`
2. Make your changes with tests
3. Ensure all tests pass (`make test`)
4. Update documentation if needed
5. Open a PR with a clear title and description
6. Link relevant issues with `Closes #123`
7. Wait for review — maintainers will respond within 48h

## Code Style

| Language | Tool | Command |
|----------|------|---------|
| Go | gofmt + go vet | `make lint-go` |
| Python | black + ruff | `make lint-python` |
| TypeScript | prettier + eslint | `make lint-ts` |

Run all linters: `make lint`

## Testing Requirements

- All new features must include tests
- Bug fixes must include a regression test
- Minimum coverage: maintain existing percentage (no regressions)
- Run the full suite before submitting:

```bash
make test          # All tests
make test-go       # Go unit tests
make test-python   # Python tests
make test-ts       # Dashboard tests
```

## Issue Templates

When creating issues, use the appropriate template:

- **Bug Report** — Include reproduction steps, expected vs actual behavior, environment
- **Feature Request** — Describe the problem, proposed solution, and alternatives considered
- **Documentation** — What's missing or incorrect

## Architecture Overview

```
mitran/
├── prototype/
│   ├── server/          # Go — HTTP engine, DAG scheduler, API
│   ├── agent-worker/    # Python — LLM orchestration, 8 AI agents (Bedrock)
│   ├── dashboard/       # React + Vite + Tailwind — web UI
│   ├── cli/             # Go — CLI tool (`mitran` command)
│   └── proto/           # gRPC protobuf definitions
├── docs/                # Architecture docs, design decisions
└── website/             # Landing page (Vercel)
```

**Key design decisions:**
- Go handles orchestration, scheduling, and API routing
- Python handles all LLM interactions via AWS Bedrock
- gRPC + Protobuf for Go ↔ Python IPC
- React dashboard communicates via REST + WebSocket
- Memory: JSON (now) → SQLite+FTS5 → Postgres+pgvector

## Questions?

Open a [Discussion](https://github.com/getmitran/mitran/discussions) or reach out to the maintainers.
