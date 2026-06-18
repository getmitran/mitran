# Cron Scheduler - main.go Wiring

Add the following to `main.go`:

## Import

```go
import "mitran/cron"
```

## Initialize (after dataDir setup)

```go
cronDir := filepath.Join(dataDir, "cron")
scheduler, err := cron.New(cronDir, func(job cron.Job) {
    log.Printf("[cron] executing job %s (%s) for agent %s", job.Name, job.ID, job.Agent)
    // TODO: dispatch job.Payload to job.Agent via agent worker
})
if err != nil {
    log.Fatalf("cron init: %v", err)
}
scheduler.Start()
defer scheduler.Stop()
```

## Register Routes (with existing mux)

```go
cron.RegisterRoutes(mux, scheduler)
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/cron | Create a cron job |
| GET | /api/v1/cron | List all jobs |
| DELETE | /api/v1/cron/:id | Remove a job |
| POST | /api/v1/cron/:id/pause | Pause a job |
| POST | /api/v1/cron/:id/resume | Resume a paused job |
| POST | /api/v1/cron/:id/trigger | Run job immediately |

## Create Job Request Body

```json
{
  "name": "daily-standup",
  "agent": "dev",
  "schedule": "24h",
  "payload": "Generate standup summary"
}
```

Valid schedules: `1m`, `5m`, `15m`, `1h`, `24h`
