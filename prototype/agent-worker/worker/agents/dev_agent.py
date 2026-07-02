import json
import logging
from dataclasses import dataclass, field
from urllib.request import urlopen, Request
from urllib.parse import quote
from urllib.error import URLError
from worker.llm import invoke

_log = logging.getLogger(__name__)

MITRAN_API_BASE = "http://localhost:7780"


@dataclass
class FileChange:
    path: str
    content: str
    action: str = "create"


@dataclass
class AgentResult:
    summary: str
    files: list[FileChange] = field(default_factory=list)


class BaseAgent:
    """Base class for all Mitran agents."""
    agent_type: str = "base"
    system_prompt: str = ""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        raise NotImplementedError

    def _call_llm(self, prompt: str) -> str:
        return invoke(self.system_prompt, prompt)

    def _fetch_memory(self, project_id: str, query: str) -> list:
        """Fetch relevant memories from the memory system."""
        try:
            url = f"{MITRAN_API_BASE}/api/v1/projects/{quote(project_id, safe='')}/memory/search?q={quote(query)}&limit=5"
            req = Request(url, method="GET")
            req.add_header("Content-Type", "application/json")
            with urlopen(req, timeout=5) as resp:
                data = json.loads(resp.read().decode())
                return data if isinstance(data, list) else data.get("results", [])
        except (URLError, json.JSONDecodeError, OSError) as e:
            _log.debug(f"Memory fetch failed (non-fatal): {e}")
            return []

    def _save_episode(self, project_id: str, summary: str) -> None:
        """Save an episode to the memory system after task completion."""
        try:
            url = f"{MITRAN_API_BASE}/api/v1/projects/{quote(project_id, safe='')}/memory/episodes"
            payload = json.dumps({"content": summary, "agent": self.agent_type}).encode()
            req = Request(url, data=payload, method="POST")
            req.add_header("Content-Type", "application/json")
            with urlopen(req, timeout=5) as resp:
                resp.read()
        except (URLError, OSError) as e:
            _log.debug(f"Episode save failed (non-fatal): {e}")

    def _build_memory_context(self, memories: list) -> str:
        """Format memories into a prompt-friendly string."""
        if not memories:
            return ""
        lines = ["Relevant context from memory:"]
        for m in memories[:5]:
            content = m.get("content", "") if isinstance(m, dict) else str(m)
            if content:
                lines.append(f"- {content}")
        return "\n".join(lines) + "\n\n"


SYSTEM_PROMPT = """You are Mitran's Dev Agent — a senior software engineer that generates production-ready project scaffolding.

Given a task description and project context, generate complete, working source files. Output ONLY valid JSON with this structure:
{
  "files": [
    {"path": "relative/path/to/file.py", "content": "full file content here"}
  ],
  "summary": "Brief description of what was generated"
}

Rules:
- Generate real, complete, runnable code — no placeholders or TODOs
- Include proper error handling, type hints, docstrings
- Follow language best practices and conventions
- Include dependency files (requirements.txt, package.json, etc.)
- Include Dockerfile if the project is a service
- Use modern language features appropriate for the stack"""


class DevAgent:
    name = "dev"
    agent_type = "development"

    def execute(self, task: str, context: dict) -> AgentResult:
        user_msg = f"Project: {context.get('team_description', 'Software project')}\nStack: {context.get('tech_stack', 'Python')}\nTeam size: {context.get('team_size', 5)}\n\nTask: {task}"
        raw = invoke(SYSTEM_PROMPT, user_msg)
        try:
            data = json.loads(raw)
        except json.JSONDecodeError:
            return AgentResult(summary=raw, files=[])
        files = [FileChange(path=f["path"], content=f["content"]) for f in data.get("files", [])]
        return AgentResult(summary=data.get("summary", "Dev agent completed"), files=files)
