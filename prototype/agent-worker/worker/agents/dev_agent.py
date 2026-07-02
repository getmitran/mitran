import json
from dataclasses import dataclass, field
from worker.llm import invoke


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
