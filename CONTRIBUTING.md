# Contributing to Mitran

Thank you for your interest in contributing! Here's how to get started.

## Development Setup

```bash
# Prerequisites
go 1.22+
python 3.11+
node 18+ (for dashboard)

# Clone
git clone https://github.com/getmitran/mitran.git
cd mitran

# Build core engine
go build ./cmd/mitran

# Run tests
go test ./...

# Run agents (Python)
cd agents
pip install -e .
```

## Project Structure

```
mitran/
├── cmd/mitran/          # CLI entrypoint
├── internal/            # Core engine (Go)
│   ├── engine/          # DAG, queue, executor
│   ├── agent/           # Agent runtime, ACLs
│   ├── checkpoint/      # State persistence
│   ├── integration/     # GitHub, Slack, Grafana
│   └── api/             # gRPC + REST + WebSocket
├── agents/              # Agent implementations (Python)
├── dashboard/           # Web UI (TypeScript/React)
└── docs/                # Documentation
```

## How to Contribute

### Bug Reports
Open an issue with:
- Steps to reproduce
- Expected vs actual behavior
- Mitran version (`mitran --version`)

### Feature Requests
Open a discussion first. Describe the problem you're solving, not just the solution.

### Pull Requests
1. Fork the repo
2. Create a branch (`git checkout -b feat/my-feature`)
3. Make changes
4. Run tests (`go test ./...`)
5. Submit PR with clear description

### Writing Agents
See the [Agent SDK docs](docs/agent-sdk.md) for creating custom agents.

## Code Style

- **Go:** `gofmt` + `golangci-lint`
- **Python:** `ruff` + `black`
- **TypeScript:** `eslint` + `prettier`

## License

By contributing, you agree that your contributions will be licensed under AGPL-3.0.
