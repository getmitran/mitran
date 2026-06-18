# Load Testing

## Prerequisites

```bash
# Install k6
brew install k6

# Install vegeta (optional, for simple throughput testing)
brew install vegeta
```

## k6 (Primary)

Run the full load test suite:

```bash
k6 run k6-script.js
```

Override the base URL:

```bash
k6 run -e BASE_URL=http://localhost:7780 k6-script.js
```

Run a single scenario:

```bash
k6 run --scenario health k6-script.js
```

### Scenarios

| Scenario | Target RPS | p95 Threshold |
|----------|-----------|---------------|
| Health endpoint | 1000 | < 500ms |
| Agent chat | 50 | < 2s |
| Memory operations | 200 | < 500ms |

### Stages

All scenarios follow: ramp-up (30s) → sustain (60s) → ramp-down (10s).

## Vegeta (Quick Throughput Test)

Create request body files:

```bash
echo '{"message":"hello","agent":"dev","project_id":"load-test"}' > chat-body.json
echo '{"key":"test","value":"data","project_id":"load-test"}' > memory-body.json
```

Run against health endpoint at 1000 rps for 30s:

```bash
echo "GET http://localhost:7780/health" | vegeta attack -rate=1000 -duration=30s | vegeta report
```

Run from targets file:

```bash
vegeta attack -targets=vegeta-targets.txt -rate=100 -duration=30s | vegeta report
```

## Interpreting Results

- **p95 < 500ms** for health/memory = PASS
- **p95 < 2s** for chat = PASS (LLM calls are inherently slower)
- **Error rate < 1%** for health/memory, < 5% for chat = PASS
