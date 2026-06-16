import json
from worker.agents.dev_agent import AgentResult, FileChange
from worker.llm import invoke


SYSTEM_PROMPT = """You are Mitran's Wiki Agent — a knowledge management specialist that generates internal wiki structures.

Given a project and team context, generate wiki pages. Output ONLY valid JSON:
{
  "files": [
    {"path": "wiki/page-name.md", "content": "full markdown content"}
  ],
  "summary": "Brief description of wiki pages generated"
}

Generate:
- wiki/onboarding.md — new engineer onboarding guide (setup, access, first PR)
- wiki/architecture.md — system architecture overview with component descriptions
- wiki/runbooks/incident-response.md — incident response runbook template
- wiki/runbooks/deployment.md — deployment runbook with rollback steps
- wiki/sops/on-call.md — on-call SOP with escalation paths
- wiki/sops/code-review.md — code review standards and process

Rules:
- Write actionable, step-by-step content
- Include checklists where appropriate
- Add placeholder sections for team-specific info (marked with [TODO])
- Use consistent formatting and navigation links between pages"""


class WikiAgent:
    name = "wiki"
    agent_type = "knowledge_management"

    def execute(self, task: str, context: dict) -> AgentResult:
        user_msg = f"Project: {context.get('company_description', 'Software project')}\nStack: {context.get('tech_stack', 'Python')}\nTeam size: {context.get('team_size', 5)}\n\nTask: {task}"
        raw = invoke(SYSTEM_PROMPT, user_msg)
        try:
            data = json.loads(raw)
        except json.JSONDecodeError:
            return AgentResult(summary=raw, files=[])
        files = [FileChange(path=f["path"], content=f["content"]) for f in data.get("files", [])]
        return AgentResult(summary=data.get("summary", "Wiki agent completed"), files=files)
