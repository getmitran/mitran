# Plugin Authoring Guide

Write custom agents for Mitran by extending the `BaseAgent` class.

## Quick Start

```python
from worker.agents.base import BaseAgent

class SummarizerAgent(BaseAgent):
    """Agent that summarizes text content."""

    def get_capabilities(self) -> dict:
        return {
            "name": "summarizer",
            "description": "Summarizes documents and text",
            "skills": ["summarization", "extraction"],
            "tools": ["read_file"],
        }

    async def handle_task(self, task: dict) -> dict:
        content = task.get("content", "")
        if not content and task.get("file_path"):
            content = self.tools.read_file(task["file_path"])

        summary = await self.call_llm(
            prompt=f"Summarize this concisely:\n\n{content}",
            max_tokens=500,
        )

        self.memory.store("last_summary", summary)
        return {"status": "complete", "result": summary}
```

## BaseAgent Class

Import from `worker.agents.base`:

```python
from worker.agents.base import BaseAgent
```

Your agent must implement two methods:

| Method | Purpose |
|--------|---------|
| `get_capabilities()` | Returns metadata dict (name, description, skills, tools) |
| `handle_task(task)` | Processes a task dict, returns result dict |

## LLM Integration

Call the LLM from any agent method:

```python
response = await self.call_llm(
    prompt="Your prompt here",
    system="Optional system message",
    max_tokens=1024,
    temperature=0.7,
)
```

The LLM provider is configured globally (AWS Bedrock Claude by default). Your agent doesn't need to manage credentials.

## Tool Registration

Declare tools in `get_capabilities()` and use them via `self.tools`:

```python
def get_capabilities(self) -> dict:
    return {
        "name": "my_agent",
        "tools": ["read_file", "write_file", "run_command"],
    }

async def handle_task(self, task: dict) -> dict:
    content = self.tools.read_file("/path/to/file.txt")
    self.tools.write_file("/path/to/output.txt", result)
    return {"status": "complete"}
```

## Skill Loading

Load reusable skills (prompt templates, workflows) from the skills directory:

```python
skill = self.load_skill("code_review")
prompt = skill.render(context={"code": source_code})
```

## Memory Access

Persist and recall data across tasks:

```python
# Store a value
self.memory.store("project_context", {"repo": "mitran", "branch": "main"})

# Recall a value
context = self.memory.recall("project_context")

# Search semantic memory
results = self.memory.search("deployment configuration", limit=5)
```

## Configuration

Agents read configuration from environment variables:

```bash
# .env or shell export
MITRAN_AGENT_SUMMARIZER_MAX_TOKENS=2000
MITRAN_AGENT_SUMMARIZER_MODEL=claude-sonnet
```

Access in code:

```python
import os

max_tokens = int(os.environ.get("MITRAN_AGENT_SUMMARIZER_MAX_TOKENS", "1024"))
```

## Testing Locally

Run your agent in isolation:

```bash
cd prototype/agent-worker

# Run agent directly
python -m worker.agents.test_runner --agent summarizer --task '{"content": "Hello world"}'

# Run with pytest
pytest tests/agents/test_summarizer.py -v
```

Write tests using the mock LLM client:

```python
import pytest
from worker.agents.summarizer import SummarizerAgent
from worker.testing import MockLLMClient, MockMemory

@pytest.fixture
def agent():
    return SummarizerAgent(llm=MockLLMClient(), memory=MockMemory())

async def test_summarize(agent):
    result = await agent.handle_task({"content": "Long text here..."})
    assert result["status"] == "complete"
    assert "result" in result
```

## Packaging and Registration

1. Place your agent module in `prototype/agent-worker/worker/agents/`:

```
worker/agents/
├── base.py
├── summarizer.py    ← your agent
└── __init__.py
```

2. Register in `worker/agents/__init__.py`:

```python
from .summarizer import SummarizerAgent

AGENT_REGISTRY = {
    "summarizer": SummarizerAgent,
    # ... other agents
}
```

3. The engine discovers registered agents on startup. No additional configuration needed.

4. Verify registration:

```bash
curl http://localhost:8888/agents | python -m json.tool
```

Your agent should appear in the list with its capabilities.
