# Quickstart: Get Your Agent Team Running in 5 Minutes

## Prerequisites

- Docker and Docker Compose installed
- An LLM API key (AWS Bedrock, OpenAI, or Anthropic)

## 1. Start Mitran

```bash
git clone https://github.com/getmitran/mitran && cd mitran
cp .env.example .env   # Add your LLM API keys
docker compose up -d
```

Dashboard: http://localhost:5173 · API: http://localhost:7780

Wait ~10 seconds for all services to initialize.

## 2. Initialize Your Team

```bash
mitran init --team "My Project" --size 3
```

This provisions 7 agents with your project context:

| Agent | What it does |
|-------|-------------|
| **PM** | Routes tasks, tracks progress, assigns work |
| **Developer** | Writes code, manages branches, creates PRs |
| **Architect** | System design, ADRs, tech debt analysis |
| **DevOps** | CI/CD pipelines, deployments, monitoring |
| **Tester** | Writes and runs tests, coverage reports |
| **Docs** | Documentation, changelogs, API references |
| **Security** | Vulnerability scans, policy checks, audits |

## 3. Verify Health

```bash
curl http://localhost:7780/api/v1/health
```

```json
{"status":"healthy","version":"0.2.0","agents":16,"uptime_seconds":12}
```

## 4. Create a Task

Tasks go to the PM orchestrator, which assigns the right agent:

```bash
curl -X POST http://localhost:7780/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Write unit tests for the auth module","priority":"high"}'
```

```json
{"id":"task_01","title":"Write unit tests for the auth module","agent":"tester","status":"queued"}
```

The PM automatically routed this to the **Tester** agent.

## 5. Watch Progress

```bash
curl http://localhost:7780/api/v1/tasks/task_01
```

```json
{"id":"task_01","status":"in_progress","agent":"tester","steps_completed":2,"steps_total":5}
```

## 6. Check Results

```bash
curl http://localhost:7780/api/v1/tasks/task_01/result
```

```json
{"id":"task_01","status":"completed","output":{"files_created":["tests/test_auth.py"],"coverage":"87%"},"duration_ms":4100}
```

## 7. Use Memory & Crons

```bash
# Search what agents have learned
curl http://localhost:7780/api/v1/memory/search?q=auth+module

# Schedule a recurring security scan
curl -X POST http://localhost:7780/api/v1/crons \
  -H "Content-Type: application/json" \
  -d '{"name":"daily-security-scan","message":"Run dependency audit","cron_expr":"0 9 * * *","agent":"security"}'
```

## Next Steps

- Open the **Dashboard** at http://localhost:5173 for visual Kanban, chat, and monitoring
- Try the **CLI**: `mitran run "Refactor the payment service for better error handling"`
- Browse **Artifacts**: `mitran artifact list` to see versioned agent outputs
- Read the [Architecture Guide](./10-openclaw-architecture-v2.md) for deeper understanding
- Deploy to Kubernetes: `kubectl apply -f deploy/k8s/`
