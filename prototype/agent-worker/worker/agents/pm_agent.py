"""PM Agent - Orchestrator that decomposes goals into tasks and assigns them to specialist agents."""
import json
import logging
import httpx
from dataclasses import dataclass, field
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange
from worker.llm import invoke as llm_invoke

log = logging.getLogger(__name__)

ENGINE_API_BASE = "http://localhost:7780/api/v1"

SYSTEM_PROMPT = """You are a Product Manager AI orchestrator for Mitran. Given a goal, decompose it into actionable tasks for specialist agents:

- developer: writing code, implementing features, fixing bugs
- architect: system design, architecture review, API design, tech decisions
- tester: writing tests, test plans, integration/load testing
- devops: deployment, CI/CD, infrastructure, debugging production issues
- docs: documentation, READMEs, API docs, user guides
- security: security scanning, vulnerability assessment, access control review

Return ONLY a valid JSON array of task objects. Each object must have:
- "title": concise task title (imperative mood, max 80 chars)
- "agent": one of [developer, architect, tester, devops, docs, security]
- "priority": integer 1-5 (1=critical, 5=nice-to-have)
- "depends_on": array of task indices (0-based) this task depends on, empty if none
- "description": 1-2 sentence description of what needs to be done

Example output for goal "Add user authentication":
[
  {"title": "Design auth architecture and token flow", "agent": "architect", "priority": 1, "depends_on": [], "description": "Define JWT vs session strategy, token refresh flow, and storage approach."},
  {"title": "Implement auth middleware and login endpoint", "agent": "developer", "priority": 1, "depends_on": [0], "description": "Build login/register endpoints with password hashing and JWT issuance."},
  {"title": "Write auth integration tests", "agent": "tester", "priority": 2, "depends_on": [1], "description": "Test login, token refresh, invalid credentials, and session expiry."},
  {"title": "Review auth for OWASP top 10 vulnerabilities", "agent": "security", "priority": 2, "depends_on": [1], "description": "Check for injection, broken auth, sensitive data exposure."},
  {"title": "Document auth API endpoints and usage", "agent": "docs", "priority": 3, "depends_on": [1], "description": "Write OpenAPI spec and usage guide for auth endpoints."},
  {"title": "Add auth service to deployment pipeline", "agent": "devops", "priority": 3, "depends_on": [1], "description": "Configure secrets management and add auth env vars to CI/CD."}
]

Rules:
- Order tasks by dependency (earlier tasks first)
- Ensure depends_on indices reference earlier tasks only (no circular deps)
- Include at least one task for testing and security when applicable
- Keep tasks atomic — one clear deliverable per task
- Priority 1-2 for blocking work, 3 for important, 4-5 for polish"""


class PMAgent(BaseAgent):
    """PM Agent acts as orchestrator: decomposes goals into subtasks and assigns them to specialist agents."""

    agent_type = "pm"
    SYSTEM_PROMPT = SYSTEM_PROMPT

    # Keep backward-compat with BaseAgent
    system_prompt = SYSTEM_PROMPT

    def decompose_goal(self, goal: str) -> list[dict]:
        """Parse a natural language goal into structured subtasks using LLM.

        Returns:
            List of dicts with keys: title, agent, priority, depends_on, description
        """
        prompt = f"Decompose this goal into tasks:\n\n{goal}"
        raw = llm_invoke(SYSTEM_PROMPT, prompt)

        # Extract JSON from response (handle markdown code blocks)
        text = raw.strip()
        if text.startswith("```"):
            lines = text.split("\n")
            # Remove first and last lines (```json and ```)
            lines = [l for l in lines[1:] if not l.strip().startswith("```")]
            text = "\n".join(lines)

        try:
            tasks = json.loads(text)
        except json.JSONDecodeError:
            # Try to find JSON array in the response
            start = text.find("[")
            end = text.rfind("]") + 1
            if start >= 0 and end > start:
                tasks = json.loads(text[start:end])
            else:
                log.error("Failed to parse LLM response as JSON: %s", text[:200])
                raise ValueError(f"LLM did not return valid JSON task list")

        # Validate structure
        valid_agents = {"developer", "architect", "tester", "devops", "docs", "security"}
        validated = []
        for i, task in enumerate(tasks):
            validated.append({
                "title": str(task.get("title", f"Task {i+1}")),
                "agent": task.get("agent", "developer") if task.get("agent") in valid_agents else "developer",
                "priority": max(1, min(5, int(task.get("priority", 3)))),
                "depends_on": [d for d in task.get("depends_on", []) if isinstance(d, int) and 0 <= d < i],
                "description": str(task.get("description", task.get("title", ""))),
            })

        return validated

    def assign_tasks(self, subtasks: list[dict], context: dict) -> list[dict]:
        """POST each subtask to the engine API to create real tasks in the system.

        Args:
            subtasks: List from decompose_goal()
            context: Project context with project_id, etc.

        Returns:
            List of created task responses from the API
        """
        project_id = context.get("project_id", "default")
        results = []

        for task in subtasks:
            payload = {
                "title": task["title"],
                "description": task["description"],
                "agent_type": task["agent"],
                "priority": task["priority"],
                "depends_on": task["depends_on"],
                "project_id": project_id,
                "status": "pending",
            }

            try:
                resp = httpx.post(
                    f"{ENGINE_API_BASE}/tasks",
                    json=payload,
                    timeout=10.0,
                )
                if resp.status_code in (200, 201):
                    results.append({"status": "created", "task": task["title"], "response": resp.json()})
                else:
                    results.append({"status": "failed", "task": task["title"], "error": resp.text})
                    log.warning("Failed to create task '%s': %s", task["title"], resp.status_code)
            except httpx.RequestError as e:
                results.append({"status": "error", "task": task["title"], "error": str(e)})
                log.error("Network error creating task '%s': %s", task["title"], e)

        return results

    def execute(self, task_description: str, context: dict) -> AgentResult:
        """Main execution: parse goal, decompose into subtasks, assign to agents.

        The task_description is treated as a goal to decompose and orchestrate.
        """
        project_id = context.get("project_id", "default")
        log.info("PM Agent orchestrating goal: %s", task_description[:100])

        # Fetch relevant memories
        memories = self._fetch_memory(project_id, task_description)
        memory_context = self._build_memory_context(memories)

        # Step 1: Decompose the goal into subtasks
        try:
            subtasks = self.decompose_goal(task_description)
        except (ValueError, json.JSONDecodeError) as e:
            return AgentResult(
                summary=f"Failed to decompose goal: {e}",
                files=[],
            )

        # Step 2: Assign tasks to the engine
        assignment_results = self.assign_tasks(subtasks, context)

        # Step 3: Build summary
        created = sum(1 for r in assignment_results if r["status"] == "created")
        failed = sum(1 for r in assignment_results if r["status"] != "created")

        task_table = "| # | Title | Agent | Priority | Depends On |\n|---|-------|-------|----------|------------|\n"
        for i, t in enumerate(subtasks):
            deps = ", ".join(str(d) for d in t["depends_on"]) or "—"
            task_table += f"| {i} | {t['title']} | {t['agent']} | {t['priority']} | {deps} |\n"

        summary = f"""## Goal Decomposition Complete

**Goal:** {task_description}
**Tasks created:** {created}/{len(subtasks)} ({failed} failed)

{task_table}
"""

        # Optionally write a plan file
        files = [FileChange(path="ORCHESTRATION_PLAN.md", action="create", content=summary)]

        # Save episode
        self._save_episode(project_id, f"PM: decomposed '{task_description}' into {len(subtasks)} tasks, {created} created")

        return AgentResult(summary=summary, files=files)
