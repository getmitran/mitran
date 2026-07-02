"""PM Agent - Sprint planning, roadmap, requirements, assigns work, manages releases."""
import logging
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange
from worker.llm import invoke as llm_invoke

log = logging.getLogger(__name__)


class PMAgent(BaseAgent):
    agent_type = "pm"
    system_prompt = """You are a technical product manager and project orchestrator. Your role is to plan sprints, manage roadmaps, write requirements, assign work, and coordinate releases.

When planning sprints:
1. Break epics into estimable stories with clear acceptance criteria
2. Prioritize by business value, dependencies, and risk
3. Assign work based on team capacity and expertise
4. Identify blockers and dependencies early

When writing requirements:
1. Use clear, unambiguous language
2. Include user stories, acceptance criteria, and edge cases
3. Define success metrics and KPIs
4. Specify non-functional requirements (perf, security, scale)

When managing releases:
1. Track feature completeness and quality gates
2. Coordinate cross-team dependencies
3. Plan rollback strategies
4. Write release notes

Output structured plans as markdown:
## Sprint Plan / Roadmap / Requirements / Release Notes
Include: timeline, assignments, priorities, risks, dependencies."""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        project = context.get("project_name", "unknown")
        team_size = context.get("team_size", 5)
        sprint_duration = context.get("sprint_duration", "2 weeks")

        prompt = f"""Project: {project}
Team size: {team_size}
Sprint duration: {sprint_duration}

Task: {task_description}

Provide a structured plan, requirements doc, or release coordination as needed."""

        response = self._call_llm(prompt)

        files = []
        if any(kw in task_description.lower() for kw in ["sprint", "plan", "roadmap", "requirements", "release"]):
            filename = "SPRINT_PLAN.md"
            if "roadmap" in task_description.lower():
                filename = "ROADMAP.md"
            elif "requirement" in task_description.lower():
                filename = "REQUIREMENTS.md"
            elif "release" in task_description.lower():
                filename = "RELEASE_NOTES.md"
            files.append(FileChange(path=filename, action="create", content=response))

        return AgentResult(summary=response, files=files)
