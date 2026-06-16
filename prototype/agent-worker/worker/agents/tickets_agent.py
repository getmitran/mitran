import json
from worker.agents.dev_agent import AgentResult, FileChange
from worker.llm import invoke


SYSTEM_PROMPT = """You are Mitran's Tickets Agent — a project manager that creates sprint board structures and backlog items.

Given a project description, generate initial sprint planning artifacts. Output ONLY valid JSON:
{
  "files": [
    {"path": "project/backlog.json", "content": "JSON content"}
  ],
  "summary": "Brief description of what was created"
}

Generate:
- project/backlog.json — initial backlog with epics, stories, tasks derived from the project
- project/sprint-1.json — first sprint with highest-priority items
- project/sprint-template.json — reusable sprint template
- project/labels.json — label taxonomy (priority, type, component)

Rules:
- Create actionable, well-scoped tickets (not vague)
- Include acceptance criteria for each story
- Estimate story points (fibonacci: 1, 2, 3, 5, 8, 13)
- First sprint should cover MVP foundation
- Group by epics that map to major features"""


class TicketsAgent:
    name = "tickets"
    agent_type = "project_management"

    def execute(self, task: str, context: dict) -> AgentResult:
        user_msg = f"Project: {context.get('company_description', 'Software project')}\nStack: {context.get('tech_stack', 'Python')}\nTeam size: {context.get('team_size', 5)}\n\nTask: {task}"
        raw = invoke(SYSTEM_PROMPT, user_msg)
        try:
            data = json.loads(raw)
        except json.JSONDecodeError:
            return AgentResult(summary=raw, files=[])
        files = [FileChange(path=f["path"], content=f["content"]) for f in data.get("files", [])]
        return AgentResult(summary=data.get("summary", "Tickets agent completed"), files=files)
