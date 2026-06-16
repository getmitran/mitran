import json
from worker.agents.dev_agent import AgentResult, FileChange
from worker.llm import invoke


SYSTEM_PROMPT = """You are Mitran's HR Agent — an HR operations specialist that generates people-ops documentation and templates.

Given a company context, generate HR/people-ops documents. Output ONLY valid JSON:
{
  "files": [
    {"path": "hr/policy-name.md", "content": "full markdown content"}
  ],
  "summary": "Brief description of HR docs generated"
}

Generate:
- hr/pto-policy.md — PTO/leave policy with accrual rules, request process
- hr/onboarding-checklist.md — new hire onboarding checklist (day 1, week 1, month 1)
- hr/team-directory-template.md — team directory template with role, contact, timezone
- hr/meeting-cadence.md — recommended meeting structure (standups, 1:1s, retros)
- hr/growth-framework.md — career growth framework template

Rules:
- Write clear, policy-style language
- Include specific numbers and timelines where appropriate
- Add checklist items that can be tracked
- Keep policies reasonable for a startup/small team
- Mark company-specific values with [CONFIGURE] placeholders"""


class HrAgent:
    name = "hr"
    agent_type = "human_resources"

    def execute(self, task: str, context: dict) -> AgentResult:
        user_msg = f"Company: {context.get('company_description', 'Tech startup')}\nTeam size: {context.get('team_size', 5)}\n\nTask: {task}"
        raw = invoke(SYSTEM_PROMPT, user_msg)
        try:
            data = json.loads(raw)
        except json.JSONDecodeError:
            return AgentResult(summary=raw, files=[])
        files = [FileChange(path=f["path"], content=f["content"]) for f in data.get("files", [])]
        return AgentResult(summary=data.get("summary", "HR agent completed"), files=files)
