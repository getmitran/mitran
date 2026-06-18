# Mitran SDK Guide

## Building Custom Agents

All agents extend `BaseAgent` from `worker.agents.dev_agent`:

```python
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange

class MyAgent(BaseAgent):
    agent_type = 'custom'
    system_prompt = """You are a specialist in..."""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        # context contains: project_name, languages, workspace_path, task_id
        response = self._call_llm(task_description)
        files = [FileChange(path='output.md', action='create', content=response)]
        return AgentResult(summary=response, files=files)
```

## Registering Agents

Add to `worker/agents/__init__.py`:
```python
from worker.agents.my_agent import MyAgent
AGENTS['custom'] = MyAgent()
```

## Agent Context

| Key | Type | Description |
|-----|------|-------------|
| project_name | str | Current project name |
| languages | list | Project languages |
| workspace_path | str | Absolute path to workspace |
| task_id | str | Current task ID |
| github_org | str | GitHub org (if configured) |
| github_repo | str | GitHub repo (if configured) |
| team_size | int | Team size |

## REST API Integration

```bash
# Create a task for your agent
curl -X POST http://localhost:7780/api/v1/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title": "My task", "agent": "custom", "priority": 2}'

# Stream chat with agent
curl -X POST http://localhost:8888/stream \
  -H 'Content-Type: application/json' \
  -d '{"message": "Help me with X", "agent": "custom"}'
```

## MCP Tool Servers

Register external tools via the MCP registry:
```bash
curl -X POST http://localhost:7780/api/v1/mcp/servers \
  -H 'Content-Type: application/json' \
  -d '{"id": "my-tool", "name": "My Tool", "command": "npx", "args": ["-y", "my-mcp-server"]}'
```

## WebSocket Events

Connect to `ws://localhost:7780/api/v1/ws` for real-time updates:
```javascript
const ws = new WebSocket('ws://localhost:7780/api/v1/ws')
ws.onmessage = (e) => {
  const event = JSON.parse(e.data)
  // event.type: task.created, task.updated, agent.completed, chat.message
}
```

## Plugin System

```bash
curl -X POST http://localhost:7780/api/v1/plugins \
  -d '{"id": "my-plugin", "name": "My Plugin", "version": "1.0", "entry_point": "plugins/my_plugin.py"}'
```
