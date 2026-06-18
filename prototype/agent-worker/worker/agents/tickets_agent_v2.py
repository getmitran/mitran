"""Tickets Agent v2 - Auto-triages, labels, and routes tickets."""
import logging
import json
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange

log = logging.getLogger(__name__)

class TicketsAgentV2(BaseAgent):
    agent_type = 'tickets'
    system_prompt = """You are a ticket triage specialist. When given a ticket:
    1. Classify priority (critical/high/medium/low)
    2. Assign labels (bug, feature, docs, infra, security)
    3. Suggest assignee based on expertise area
    4. Estimate effort (XS/S/M/L/XL)
    5. Write acceptance criteria
    Output as JSON: {"priority": "", "labels": [], "suggested_assignee": "", "effort": "", "acceptance_criteria": []}"""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        prompt = f"""Triage this ticket:

Title/Description: {task_description}
Project: {context.get('project_name', 'Unknown')}
Languages: {context.get('languages', ['Unknown'])}

Provide triage in JSON format followed by a brief explanation."""
        
        response = self._call_llm(prompt)
        
        # Try to parse JSON from response
        triage_data = None
        try:
            start = response.find('{')
            end = response.rfind('}') + 1
            if start >= 0 and end > start:
                triage_data = json.loads(response[start:end])
        except json.JSONDecodeError:
            pass
        
        if triage_data:
            summary = f'## Ticket Triage\n\n'
            summary += f'**Priority:** {triage_data.get("priority", "medium")}\n'
            summary += f'**Labels:** {", ".join(triage_data.get("labels", []))}\n'
            summary += f'**Effort:** {triage_data.get("effort", "M")}\n'
            summary += f'**Suggested Assignee:** {triage_data.get("suggested_assignee", "unassigned")}\n'
            if triage_data.get('acceptance_criteria'):
                summary += f'\n**Acceptance Criteria:**\n'
                for ac in triage_data['acceptance_criteria']:
                    summary += f'- {ac}\n'
        else:
            summary = f'## Ticket Triage\n\n{response}'
        
        return AgentResult(summary=summary, files=[])
