"""Agent-to-Agent Delegation Manager for Mitran.

Handles inter-agent task handoffs, delegation chains, and provides
standard handoff patterns for common workflows.
"""

import time
import uuid
from typing import Optional

import requests

SERVER_BASE = "http://localhost:7780"


class DelegationManager:
    """Manages inter-agent task delegation and result polling."""

    def __init__(self, base_url: str = SERVER_BASE):
        self.base_url = base_url.rstrip("/")
        self._chains: dict[str, list[dict]] = {}

    def delegate(self, from_agent: str, to_agent: str, task: str, context: dict) -> str:
        """Delegate a task from one agent to another.

        Posts to the server's task API with delegation metadata.

        Args:
            from_agent: Name of the delegating agent.
            to_agent: Name of the target agent.
            task: Task description or instruction.
            context: Additional context dict (may contain 'task_id' for chaining).

        Returns:
            The task_id of the newly created delegated task.
        """
        payload = {
            "title": task,
            "agent": to_agent,
            "priority": context.get("priority", "medium"),
            "context": {
                "delegated_by": from_agent,
                "parent_task_id": context.get("task_id"),
                **{k: v for k, v in context.items() if k not in ("task_id", "priority")},
            },
        }

        resp = requests.post(
            f"{self.base_url}/api/v1/tasks",
            json=payload,
            timeout=10,
        )
        resp.raise_for_status()
        data = resp.json()
        task_id = data.get("id") or data.get("task_id") or str(uuid.uuid4())

        # Track delegation chain
        parent_id = context.get("task_id")
        chain_entry = {
            "task_id": task_id,
            "from_agent": from_agent,
            "to_agent": to_agent,
            "task": task,
            "parent_task_id": parent_id,
        }
        self._chains[task_id] = self._chains.get(parent_id, []) + [chain_entry]

        return task_id

    def wait_for_result(self, task_id: str, timeout: int = 60) -> dict:
        """Poll for task completion.

        Checks task status every 2 seconds until completed, failed, or timeout.

        Args:
            task_id: The task identifier to poll.
            timeout: Maximum seconds to wait (default 60).

        Returns:
            Dict with 'status', 'result', and 'task_id' keys.

        Raises:
            TimeoutError: If task does not complete within timeout.
        """
        deadline = time.time() + timeout
        terminal_states = {"completed", "failed", "cancelled"}

        while time.time() < deadline:
            try:
                resp = requests.get(
                    f"{self.base_url}/api/v1/tasks/{task_id}",
                    timeout=10,
                )
                resp.raise_for_status()
                data = resp.json()
                status = data.get("status", "unknown")

                if status in terminal_states:
                    return {
                        "status": status,
                        "result": data.get("result"),
                        "task_id": task_id,
                        "data": data,
                    }
            except requests.RequestException:
                pass  # Retry on transient errors

            time.sleep(2)

        raise TimeoutError(f"Task {task_id} did not complete within {timeout}s")

    def get_chain(self, task_id: str) -> list:
        """Return the delegation chain for a task.

        Args:
            task_id: Any task_id in the chain.

        Returns:
            List of delegation entries showing the handoff path.
        """
        return self._chains.get(task_id, [])


# ---------------------------------------------------------------------------
# Standard handoff patterns
# ---------------------------------------------------------------------------

_default_manager: Optional[DelegationManager] = None


def _get_manager() -> DelegationManager:
    global _default_manager
    if _default_manager is None:
        _default_manager = DelegationManager()
    return _default_manager


def handoff_to_tester(code_files: list, from_agent: str = "developer") -> str:
    """Hand off code files to the tester agent for validation.

    Args:
        code_files: List of file paths to test.
        from_agent: Originating agent name.

    Returns:
        task_id of the created testing task.
    """
    mgr = _get_manager()
    task = f"Run tests and validate changes in: {', '.join(code_files)}"
    context = {"code_files": code_files, "type": "test_validation"}
    return mgr.delegate(from_agent, "tester", task, context)


def handoff_to_architect(description: str, from_agent: str = "developer") -> str:
    """Hand off an architectural question or review to the architect agent.

    Args:
        description: Description of what needs architectural review.
        from_agent: Originating agent name.

    Returns:
        task_id of the created architecture task.
    """
    mgr = _get_manager()
    task = f"Review architecture: {description}"
    context = {"description": description, "type": "architecture_review"}
    return mgr.delegate(from_agent, "architect", task, context)


def handoff_to_docs(changes: list, from_agent: str = "developer") -> str:
    """Hand off a list of changes to the docs agent for documentation updates.

    Args:
        changes: List of change descriptions needing documentation.
        from_agent: Originating agent name.

    Returns:
        task_id of the created documentation task.
    """
    mgr = _get_manager()
    task = f"Update documentation for: {', '.join(changes)}"
    context = {"changes": changes, "type": "documentation_update"}
    return mgr.delegate(from_agent, "docs", task, context)
