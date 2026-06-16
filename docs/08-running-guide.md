# Mitran — Running the System

## Quick Start

```bash
# Terminal 1: Core Engine
cd /Users/ravitejb/Documents/Mitran/prototype/server
go run .
# → Mitran engine running on :7777

# Terminal 2: Agent Worker
cd /Users/ravitejb/Documents/Mitran/prototype/agent-worker
pip install -e .
python run.py
# → 8 agents registered, Worker running on :8888

# Terminal 3: Dashboard (optional)
cd /Users/ravitejb/Documents/Mitran/prototype/dashboard
npm install && npm run dev
# → http://localhost:5173

# Terminal 4: Trigger
curl -X POST http://localhost:7777/api/v1/init \
  -H "Content-Type: application/json" \
  -d '{"name":"MyApp","description":"SaaS analytics platform","languages":"Python, React","team_size":5}'
```

## Prerequisites

| Requirement | Version | Check |
|-------------|---------|-------|
| Go | 1.22+ | `go version` |
| Python | 3.11+ | `python3 --version` |
| Node.js | 18+ | `node --version` |
| AWS credentials | Bedrock access in us-east-1 | `aws sts get-caller-identity` |

### Go Proxy (if behind firewall)
```bash
go env -w GOPROXY=https://goproxy.io,direct
```

### AWS Bedrock Access
The agent worker calls `us.anthropic.claude-sonnet-4-20250514` via Bedrock. Ensure credentials are available:
```bash
# Option 1: Environment variables
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AWS_REGION=us-east-1

# Option 2: AWS profile
export AWS_PROFILE=my-profile

# Option 3: Ada (Amazon internal)
ada credentials update --account <id> --role Admin --provider conduit --once
```

## Architecture

```
┌─────────┐     HTTP      ┌──────────────┐     HTTP      ┌──────────────┐
│   CLI   │──────────────▶│ Core Engine  │──────────────▶│ Agent Worker │
│  :stdin │               │    :7777     │               │    :8888     │
└─────────┘               │              │◀──────────────│              │
                          │ • REST API   │  task-complete │ • 8 agents  │
┌─────────┐     HTTP      │ • Scheduler  │               │ • Bedrock   │
│Dashboard│──────────────▶│ • JSON store │               │ • Templates │
│  :5173  │  (poll 3-5s)  └──────────────┘               └──────────────┘
└─────────┘
```

## API Reference

### Projects
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/init` | Create project + generate task plan |
| POST | `/api/v1/projects` | Create project manually |
| GET | `/api/v1/projects` | List all projects |
| GET | `/api/v1/projects/:id` | Get project with tasks |

### Tasks
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tasks` | List tasks (filter: ?status=queued) |
| POST | `/api/v1/tasks` | Create task |
| PUT | `/api/v1/tasks/:id/priority` | Reorder priority |
| POST | `/api/v1/tasks/:id/assign` | Assign to agent |

### Checkpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/checkpoints` | List pending checkpoints |
| POST | `/api/v1/checkpoints/:id/approve` | Approve with feedback |
| POST | `/api/v1/checkpoints/:id/reject` | Reject with feedback |

### Agents
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/agents` | List agents + status |
| POST | `/api/v1/agents/register` | Agent self-registration |
| POST | `/api/v1/agents/task-complete` | Submit work result |

### Environments (Multi-env CI/CD)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/projects/:id/environments` | Set environments |
| GET | `/api/v1/projects/:id/environments` | Get environments |
| PUT | `/api/v1/projects/:id/environments/:name` | Update env config |
| POST | `/api/v1/projects/:id/environments/:name/deploy` | Trigger deploy |
| GET | `/api/v1/projects/:id/deployments` | Deployment history |

### Pipelines
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/projects/:id/pipelines` | Create pipeline |
| GET | `/api/v1/projects/:id/pipelines` | List pipelines |
| POST | `/api/v1/projects/:id/pipelines/:pid/trigger` | Trigger run |
| GET | `/api/v1/projects/:id/pipelines/:pid/runs` | Run history |

## Configuration

All config stored in `.mitran/config.json`:

```json
{
  "engine_port": 7777,
  "worker_port": 8888,
  "aws_region": "us-east-1",
  "model_id": "us.anthropic.claude-sonnet-4-20250514",
  "scheduler_interval_seconds": 2,
  "default_environments": ["alpha", "beta", "gamma", "prod"]
}
```

Environment variables override:
| Variable | Default | Description |
|----------|---------|-------------|
| `MITRAN_PORT` | 7777 | Engine HTTP port |
| `MITRAN_WORKER_PORT` | 8888 | Worker HTTP port |
| `MITRAN_ENGINE_URL` | http://localhost:7777 | Engine URL (worker uses) |
| `MITRAN_MODEL_ID` | us.anthropic.claude-sonnet-4-20250514 | Bedrock model |
| `AWS_REGION` | us-east-1 | AWS region for Bedrock |

## Data Persistence

Prototype uses JSON file storage in `.mitran/` directory:
```
.mitran/
├── projects.json      # Project definitions
├── tasks.json         # Task queue
├── checkpoints.json   # Agent outputs pending review
├── agents.json        # Registered agents
├── deployments.json   # Deployment history
└── config.json        # Runtime configuration
```

## Troubleshooting

| Problem | Fix |
|---------|-----|
| `go: proxy.golang.org: i/o timeout` | `go env -w GOPROXY=https://goproxy.io,direct` |
| `Bedrock invocation failed` | Check AWS credentials + Bedrock model access |
| `Connection refused :7777` | Start engine first: `cd server && go run .` |
| `No agents registered` | Start worker: `cd agent-worker && python run.py` |
| Dashboard shows "Engine offline" | Engine not running or CORS issue |
| Tasks stuck in "queued" | Check agent worker is running + registered |
