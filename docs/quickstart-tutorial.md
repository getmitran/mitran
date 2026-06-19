# Build Your First Agent in 5 Minutes

## Prerequisites

- Docker and Docker Compose installed
- `curl` available in your terminal

## 1. Start Mitran

```bash
git clone https://github.com/getmitran/mitran && cd mitran
docker compose up -d
```

Open http://localhost:5173 for the dashboard UI.

Wait ~10 seconds for all services to initialize.

## 2. Verify Health

```bash
curl http://localhost:7780/api/v1/health
```

```json
{"status":"healthy","version":"0.1.0","agents":8,"uptime_seconds":12}
```

## 3. Create a Task

```bash
curl -X POST http://localhost:7780/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Write a README","agent":"dev","priority":"high"}'
```

```json
{"id":"task_01","title":"Write a README","agent":"dev","status":"queued","created_at":"2026-06-19T12:00:00Z"}
```

## 4. Watch the Agent Process It

```bash
curl http://localhost:7780/api/v1/tasks/task_01
```

```json
{"id":"task_01","title":"Write a README","agent":"dev","status":"in_progress","steps_completed":2,"steps_total":4}
```

The Dev agent analyzes context, generates the file, and validates output.

## 5. Check Results

```bash
curl http://localhost:7780/api/v1/tasks/task_01/result
```

```json
{"id":"task_01","status":"completed","output":{"files_created":["README.md"],"summary":"Generated README with project overview, setup instructions, and usage examples."},"duration_ms":3200}
```

## Next Steps

- Open the dashboard at `http://localhost:7780` to see all agents in action
- Try other agents: `docs`, `ops`, `review`, `cicd`, `tickets`, `wiki`, `hr`
- Read the [Architecture Guide](./10-openclaw-architecture-v2.md) for deeper understanding
