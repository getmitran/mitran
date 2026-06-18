"""Agent read-back context for cross-agent output sharing within tasks."""

from dataclasses import dataclass, field
from typing import Callable, Optional


@dataclass
class AgentOutput:
    agent: str
    content: str
    token_estimate: int = 0

    def __post_init__(self):
        if not self.token_estimate:
            self.token_estimate = len(self.content) // 4


class ReadbackContext:
    """Stores prior agent outputs per task_id with sliding window retrieval."""

    def __init__(self, max_tokens: int = 4000):
        self.max_tokens = max_tokens
        self._store: dict[str, list[AgentOutput]] = {}

    def add_output(self, task_id: str, agent: str, content: str) -> None:
        if task_id not in self._store:
            self._store[task_id] = []
        self._store[task_id].append(AgentOutput(agent=agent, content=content))

    def get_context_window(self, current_task_id: str) -> str:
        outputs = self._store.get(current_task_id, [])
        if not outputs:
            return ""

        # Sliding window: most recent first, within token budget
        selected = []
        tokens_used = 0
        for output in reversed(outputs):
            if tokens_used + output.token_estimate > self.max_tokens:
                break
            selected.append(output)
            tokens_used += output.token_estimate

        selected.reverse()
        lines = []
        for o in selected:
            lines.append(f"[{o.agent}]:\n{o.content}")
        return "\n---\n".join(lines)

    def inject_into_prompt(self, base_prompt: str, task_id: str) -> str:
        context = self.get_context_window(task_id)
        if not context:
            return base_prompt
        return (
            f"{base_prompt}\n\n"
            f"## Prior Agent Outputs\n{context}"
        )

    def clear(self, task_id: str) -> None:
        self._store.pop(task_id, None)

    def summarize_if_large(
        self, task_id: str, llm_invoke: Callable[[str], str], threshold_tokens: int = 3000
    ) -> Optional[str]:
        outputs = self._store.get(task_id, [])
        total = sum(o.token_estimate for o in outputs)
        if total <= threshold_tokens:
            return None

        full_text = "\n---\n".join(f"[{o.agent}]: {o.content}" for o in outputs)
        summary = llm_invoke(
            f"Summarize these agent outputs concisely, preserving key decisions and data:\n\n{full_text}"
        )
        self._store[task_id] = [AgentOutput(agent="summary", content=summary)]
        return summary
