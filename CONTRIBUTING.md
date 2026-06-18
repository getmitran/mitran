# Contributing to Mitran

Thank you for your interest in contributing to Mitran!

## Development Setup

1. Clone the repo: `git clone https://github.com/getmitran/mitran.git`
2. Copy env: `cp .env.example .env` and fill in values
3. Install deps: `make install-deps`
4. Run: `make run`

## Architecture

- **Engine** (Go): DAG scheduler, REST API, WebSocket hub -- port 7780
- **Agent Worker** (Python): 8 AI agents using AWS Bedrock -- port 8888
- **Dashboard** (React): Vite + Tailwind -- port 5173 (dev) / 3000 (prod)
- **CLI** (Go): `mitran init`, `mitran status`

## Workflow

1. Fork the repo
2. Create a branch: `git checkout -b feat/my-feature`
3. Make changes and test: `make test`
4. Commit (conventional commits): `git commit -m "feat: add X"`
5. Push and open a PR

## Code Style

- Go: `gofmt` + `golangci-lint`
- Python: PEP 8, type hints preferred
- TypeScript: Strict mode, no `any` where avoidable

## License

AGPL-3.0 -- see LICENSE
