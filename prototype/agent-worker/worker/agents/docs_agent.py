import json
from worker.agents.dev_agent import AgentResult, FileChange
from worker.llm import invoke


SYSTEM_PROMPT = """You are Mitran's Docs Agent — a technical writer that generates comprehensive project documentation.

Given a project description and context, generate documentation files. Output ONLY valid JSON:
{
  "files": [
    {"path": "docs/relative/path.md", "content": "full markdown content"}
  ],
  "summary": "Brief description of docs generated"
}

Generate:
- README.md with project overview, quickstart, architecture summary
- docs/api.md with API reference (endpoints, params, responses)
- docs/architecture.md with system design, data flow, component diagram descriptions
- docs/deployment.md with deployment instructions per environment
- CONTRIBUTING.md with development setup, coding standards, PR process

Rules:
- Write clear, concise, developer-friendly documentation
- Include code examples and command snippets
- Use proper markdown formatting with headers, tables, code blocks
- Tailor content to the specific tech stack and project type"""


class DocsAgent:
    name = "docs"
    agent_type = "documentation"

    def execute(self, task: str, context: dict) -> AgentResult:
        user_msg = f"Project: {context.get('team_description', 'Software project')}\nStack: {context.get('tech_stack', 'Python')}\nTeam size: {context.get('team_size', 5)}\n\nTask: {task}"
        raw = invoke(SYSTEM_PROMPT, user_msg)
        try:
            data = json.loads(raw)
        except json.JSONDecodeError:
            return AgentResult(summary=raw, files=[])
        files = [FileChange(path=f["path"], content=f["content"]) for f in data.get("files", [])]
        return AgentResult(summary=data.get("summary", "Docs agent completed"), files=files)
