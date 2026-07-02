"""Architect Agent - Reviews PRs, enforces patterns, prepares technical backlog."""
import logging
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange
from worker.llm import invoke as llm_invoke

log = logging.getLogger(__name__)


class ArchitectAgent(BaseAgent):
    agent_type = "architect"
    system_prompt = """You are a senior software architect and code reviewer. Your role is to review code changes, enforce architectural patterns, and maintain technical quality.

When reviewing code:
1. Check for architectural consistency and design pattern adherence
2. Identify potential performance issues, race conditions, and security flaws
3. Verify error handling completeness and edge case coverage
4. Assess test coverage and suggest missing test scenarios
5. Evaluate naming conventions and code organization

When preparing technical backlog:
1. Identify technical debt and refactoring opportunities
2. Prioritize items by impact and risk
3. Provide clear acceptance criteria for each item

Output your review as structured markdown with sections:
## Summary
## Issues Found (Critical/Warning/Info)
## Recommendations
## Backlog Items (if applicable)

Be constructive and specific - include line references and code suggestions."""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        project = context.get("project_name", "unknown")
        code_context = context.get("code_diff", "")
        file_list = context.get("files_changed", [])

        prompt = f"""Project: {project}
Files changed: {', '.join(file_list) if file_list else 'N/A'}
Code context:
{code_context if code_context else 'No diff provided'}

Task: {task_description}

Provide a thorough architectural review and recommendations."""

        response = self._call_llm(prompt)

        files = []
        if "backlog" in task_description.lower():
            files.append(FileChange(
                path="TECHNICAL_BACKLOG.md",
                action="create",
                content=response
            ))

        return AgentResult(summary=response, files=files)
