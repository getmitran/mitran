# Mitran SLO Definitions

## Service Level Objectives

| SLO | Target | Measurement |
|-----|--------|-------------|
| API Availability | 99.9% | `sum(rate(http_requests_total{code!~"5.."}[5m])) / sum(rate(http_requests_total[5m]))` |
| P50 Latency | < 100ms | `histogram_quantile(0.50, rate(http_request_duration_seconds_bucket[5m]))` |
| P99 Latency | < 500ms | `histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))` |
| Task Queue Depth | < 1000 | `mitran_task_queue_depth` |
| Error Rate | < 1% | `sum(rate(http_requests_total{code=~"5.."}[5m])) / sum(rate(http_requests_total[5m])) * 100` |

## Measurement Method

All SLOs are measured via Prometheus metrics scraped from the `/metrics` endpoint at 15s intervals. Grafana dashboards visualize compliance over rolling 30-day windows.

**Key metrics:**
- `http_requests_total` — counter, labels: `method`, `path`, `code`
- `http_request_duration_seconds` — histogram, labels: `method`, `path`
- `mitran_task_queue_depth` — gauge, current pending tasks

## Error Budget Calculation

```
Error Budget = 1 - SLO Target

Monthly budget (minutes) = 30 days × 24h × 60min × (1 - target)

Example (99.9% availability):
  Budget = 0.1% = 43.2 minutes/month

Burn rate = actual_error_rate / allowed_error_rate
Budget consumed (%) = (downtime_minutes / 43.2) × 100
```

## Alerting Thresholds

| Alert | Condition | Severity |
|-------|-----------|----------|
| HighErrorRate | error_rate > 1% for 5m | warning |
| CriticalErrorRate | error_rate > 5% for 2m | critical |
| LatencyP99High | P99 > 500ms for 5m | warning |
| LatencyP99Critical | P99 > 2s for 2m | critical |
| QueueDepthHigh | queue_depth > 800 for 5m | warning |
| QueueDepthCritical | queue_depth > 1000 for 2m | critical |
| BudgetBurnFast | >50% budget burned in <24h | critical |

## Escalation Policy

1. **Warning** — Slack notification to `#mitran-ops`
2. **Critical** — Page on-call engineer via PagerDuty
3. **Budget burn >50% in <24h** — Page on-call immediately, notify engineering lead
4. **Budget burn >80% in <48h** — Incident declared, all hands on deck

### On-Call Response SLA

| Severity | Acknowledge | Mitigate |
|----------|-------------|----------|
| Warning | 30 min | 2 hours |
| Critical | 5 min | 30 min |
| Incident | Immediate | 15 min |
