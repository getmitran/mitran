# Mitran Agent Worker

Python HTTP service that receives task assignments from the Mitran Core Engine, executes them using AWS Bedrock (Claude), and posts results back.

## Architecture

```
┌─────────┐         ┌──────────────┐         ┌─────────────┐
│  CLI    │──init──▶│ Core Engine  │──task──▶│ Agent Worker │
│(port -)│         │ (port 7777)  │◀─result─│ (port 8888)  │
└─────────┘         └──────────────┘         └──────┬──────┘
                                                     │
                                              ┌──────▼──────┐
                                              │ AWS Bedrock  │
                                              │   (Claude)   │
                                              └─────────────┘
```

## Agents

| Agent | Type | Generates |
|-------|------|-----------|
| dev | development | Project scaffolding, source files |
| docs | documentation | README, API docs, architecture docs |
| cicd | ci_cd | GitHub Actions, Jenkinsfile, Makefile, deploy scripts |
| tickets | project_management | Sprint boards, backlog, story templates |
| wiki | knowledge_management | Onboarding guides, runbooks, SOPs |
| ops | operations | Prometheus, AlertManager, Grafana, Docker compose |
| hr | human_resources | PTO policy, onboarding checklist, team directory |
| dashboard | observability | Grafana dashboard JSON |

## Prerequisites

- Python 3.11+
- AWS credentials with Bedrock access (Claude model)
  - Set `AWS_PROFILE` or `AWS_ACCESS_KEY_ID` + `AWS_SECRET_ACCESS_KEY`
  - Region defaults to `us-east-1` (override with `AWS_REGION`)

## Setup

```bash
cd prototype/agent-worker
pip install -e .
```

## Running

**1. Start the Core Engine first:**
```bash
cd ../engine && go run . 
# Runs on port 7777
```

**2. Start the Agent Worker:**
```bash
python run.py
# Runs on port 8888, registers all 8 agents with engine on startup
```

**3. Trigger via CLI or curl:**
```bash
# Via Mitran CLI
cd ../cli && go run . init

# Via curl (direct task execution)
curl -X POST http://localhost:8888/execute \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": "test-1",
    "agent": "dev",
    "task": "Set up a Python Flask REST API with health check and user CRUD endpoints",
    "context": {"tech_stack": "Python/Flask", "team_size": 5}
  }'
```

## Configuration

| Env Variable | Default | Description |
|---|---|---|
| `MITRAN_ENGINE_URL` | `http://localhost:7777` | Core Engine URL |
| `MITRAN_WORKER_PORT` | `8888` | Worker HTTP port |
| `MITRAN_WORKER_HOST` | `0.0.0.0` | Worker bind address |
| `AWS_REGION` | `us-east-1` | AWS region for Bedrock |
| `MITRAN_MODEL_ID` | `us.anthropic.claude-sonnet-4-20250514` | Bedrock model |
| `AWS_PROFILE` | (none) | AWS credential profile |

## Error Handling

If AWS credentials are missing or invalid, the worker returns a clear error:
```json
{"error": "Bedrock invocation failed: Unable to locate credentials"}
```

The worker does NOT mock responses. All output comes from real Bedrock Claude calls.
