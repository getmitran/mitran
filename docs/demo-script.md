# Mitran Demo Script (5 minutes)

## Prerequisites

- Docker and Docker Compose installed
- AWS credentials configured (for Bedrock Claude access)
- `.env` file populated from `.env.example`

---

## 1. Start Services (~30s)

```bash
cd /path/to/mitran
docker-compose up -d
```

**Expected output:**
```
[+] Running 4/4
 ✔ Container mitran-engine      Started
 ✔ Container mitran-agent-worker Started
 ✔ Container mitran-dashboard    Started
 ✔ Container mitran-prometheus   Started
```

Verify everything is healthy:

```bash
mitran status
```

**Expected output:**
```
┌─────────────────────────────────────┐
│         MITRAN v0.1.0               │
├─────────────────────────────────────┤
│ Engine:       ✅ running (port 7780) │
│ Agent Worker: ✅ running (port 8888) │
│ Dashboard:    ✅ running (port 3000) │
│ Agents:       8 active              │
│ Memory:       JSON store            │
│ Uptime:       12s                   │
└─────────────────────────────────────┘
```

---

## 2. CLI Basics (~45s)

### List available agents

```bash
mitran agent list
```

**Expected output:**
```
NAME            STATUS    DESCRIPTION
dev             active    Development & code generation
docs            active    Documentation writer
ops             active    Operations & monitoring
review          active    Code review & PR management
hr              active    HR & onboarding workflows
cicd            active    CI/CD pipeline management
tickets         active    Ticket triage & management
wiki            active    Knowledge base management
```

### Check configuration

```bash
mitran config get
```

**Expected output:**
```
llm_provider: aws_bedrock
model: us.anthropic.claude-sonnet-4-20250514
memory_backend: json
workspace: ./mitran/projects/default/
```

---

## 3. Chat with Dev Agent (~90s)

```bash
mitran chat --agent dev
```

**Demo conversation:**

```
You: Create a Python FastAPI service with a /health endpoint

Dev Agent: I'll create a FastAPI service for you.

📁 Created: app/main.py
📁 Created: app/requirements.txt
📁 Created: Dockerfile

Here's what I built:
- FastAPI app with /health returning {"status": "ok", "timestamp": ...}
- requirements.txt with fastapi and uvicorn
- Dockerfile for containerized deployment

Want me to add tests or CI configuration?

You: Add pytest tests

Dev Agent: ✅ Created app/tests/test_health.py with 3 test cases:
- test_health_returns_200
- test_health_has_status_field
- test_health_has_timestamp_field

All tests passing. Run with: pytest app/tests/
```

Press `Ctrl+C` to exit chat.

---

## 4. Create a Cron Job (~30s)

```bash
mitran cron add \
  --name "daily-standup" \
  --message "Summarize yesterday's commits and open PRs for the team" \
  --cron "0 9 * * MON-FRI" \
  --agent dev
```

**Expected output:**
```
✅ Cron job created
  ID:       cron_a1b2c3d4
  Name:     daily-standup
  Schedule: 0 9 * * MON-FRI (weekdays at 9:00 AM)
  Agent:    dev
  Next run: Tomorrow 09:00
```

Verify:

```bash
mitran cron list
```

**Expected output:**
```
ID            NAME             SCHEDULE           AGENT   STATUS   NEXT RUN
cron_a1b2c3d4 daily-standup    0 9 * * MON-FRI   dev     active   Tomorrow 09:00
```

---

## 5. Dashboard Walkthrough (~60s)

Open http://localhost:3000 in browser.

### Memory Viewer (click "Memory" in sidebar)

**Show:**
- Semantic memory entries (key-value facts the agents have learned)
- Episodic memory (timestamped conversation history)
- Corrections (user-taught rules that override defaults)

**Demo action:** Click a memory entry to see its full context and source session.

### Cron Manager (click "Crons" in sidebar)

**Show:**
- The `daily-standup` job just created
- Pause/resume toggle
- Manual trigger button
- Execution history with logs

**Demo action:** Click "Trigger Now" → watch the agent execute and produce a standup summary in real-time.

### Artifact Library (click "Artifacts" in sidebar)

**Show:**
- Saved widgets, reports, and generated content
- Version history with diff viewer
- Tags and search/filter

**Demo action:** Open an artifact → show version dropdown → compare v1 vs v2 diff.

---

## 6. GitHub Webhook (~30s)

Simulate a push event:

```bash
curl -X POST http://localhost:7780/api/v1/webhooks/github \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: push" \
  -d '{
    "ref": "refs/heads/main",
    "commits": [
      {
        "id": "abc123",
        "message": "feat: add user authentication",
        "author": {"name": "dev", "email": "dev@example.com"}
      }
    ],
    "repository": {"full_name": "getmitran/mitran"}
  }'
```

**Expected output:**
```json
{
  "status": "accepted",
  "task_id": "task_x7y8z9",
  "agent": "cicd",
  "message": "Processing push event for getmitran/mitran"
}
```

**In dashboard:** Watch the CI/CD agent pick up the event, analyze the commit, and trigger the configured pipeline.

---

## 7. Monitoring & Metrics (~30s)

Open Prometheus at http://localhost:9090.

### Key queries to demo:

```promql
# Agent response latency (p95)
histogram_quantile(0.95, rate(mitran_agent_response_seconds_bucket[5m]))

# Total requests per agent
sum by (agent) (mitran_agent_requests_total)

# Memory store size
mitran_memory_entries_total
```

### From CLI:

```bash
mitran status --metrics
```

**Expected output:**
```
METRICS (last 5m)
─────────────────────────────
Requests:     47 total
Latency p50:  1.2s
Latency p95:  3.8s
Errors:       0
Memory:       234 entries
Cron runs:    3 (all success)
Active tasks: 1
```

---

## Closing (~15s)

> "Mitran gives your team 8 AI agents that handle dev, docs, ops, reviews, CI/CD, tickets, wiki, and HR — all orchestrated through a single CLI, dashboard, or API. Self-hosted, AGPL-licensed, and running on your own infrastructure in under 30 minutes."

```bash
# Cleanup
docker-compose down
```
