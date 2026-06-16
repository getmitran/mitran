# Mitran — Development SOP

## Development Workflow

### Adding a New Agent

1. Create `worker/agents/my_agent.py`:
```python
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange

class MyAgent(BaseAgent):
    agent_type = "my-agent"
    system_prompt = "You are a specialist in..."

    def execute(self, task_id: str, title: str, description: str, project: dict) -> AgentResult:
        prompt = f"Given project: {project['name']}. Task: {description}"
        response = self._call_llm(prompt)
        return AgentResult(
            summary="## What I did\n...",
            files=[FileChange(path="...", action="create", content=response)]
        )
```

2. Register in `worker/agents/__init__.py`:
```python
from worker.agents.my_agent import MyAgent
AGENTS["my-agent"] = MyAgent()
```

3. Add task generation in `server/handlers/init.go` within the plan builder.

### Adding a New API Endpoint

1. Create handler in `server/handlers/my_handler.go`
2. Register route in `server/main.go`
3. Add store methods in `server/db/sqlite.go` if persistence needed
4. Update dashboard API client in `dashboard/src/api.ts`

### Adding a Pipeline Template

1. Create template in `worker/agents/pipeline_templates/my_template.py`:
```python
def generate_my_config(project: dict, environments: list[dict]) -> str:
    # Return file content as string
    ...
```

2. Import in `worker/agents/pipeline_templates/__init__.py`
3. Call from `cicd_agent_v2.py`

## Testing

```bash
# Go engine tests
cd prototype/engine && go test ./...

# Python agent SDK tests
cd prototype/agents && pip install -e ".[dev]" && pytest

# Dashboard type check
cd prototype/dashboard && npx tsc --noEmit

# Integration test (requires engine + worker running)
curl -X POST http://localhost:7777/api/v1/init \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","description":"Test project","languages":"Go","team_size":3}'
```

## Release Process

1. Run all tests
2. Update version in relevant `go.mod` / `pyproject.toml` / `package.json`
3. Commit: `git commit -m "chore: bump to vX.Y.Z"`
4. Tag: `git tag vX.Y.Z`
5. Push: owner pushes tag + main branch
6. GitHub Actions builds binaries (future)

## Code Style

| Language | Formatter | Linter |
|----------|-----------|--------|
| Go | `gofmt` | `golangci-lint` |
| Python | `black` | `ruff` |
| TypeScript | `prettier` | `eslint` |

## Directory Structure

```
/Users/ravitejb/Documents/Mitran/
├── docs/                    # All documentation
├── website/                 # Landing page (Vercel)
├── github-readme/           # Public-facing README + CONTRIBUTING
└── prototype/
    ├── engine/              # Go: DAG, queue, locks, scheduler
    ├── server/              # Go: HTTP API, persistence, routing
    ├── cli/                 # Go: mitran init/serve/status
    ├── agents/              # Python: SDK + standalone agents
    ├── agent-worker/        # Python: Flask worker + 8 agents + templates
    └── dashboard/           # TypeScript: React + Tailwind UI
```
