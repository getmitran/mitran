"""Memory write-back: auto-persist agent learnings after task completion."""

import json
import logging
from dataclasses import dataclass
from typing import Any

from worker.memory.vector import VectorMemory
from worker.llm.client import LLMClient

logger = logging.getLogger(__name__)

EXTRACTION_PROMPT = """Extract structured learnings from this agent task result.
Return JSON with these arrays (empty if none found):
- facts: key-value pairs of factual knowledge (e.g. {"key": "service.port", "value": "8080"})
- corrections: rules learned from mistakes (e.g. {"rule": "always use UTC", "negative": "don't use local time"})
- episodes: notable events worth remembering (e.g. {"summary": "deployed v2 to prod", "context": "..."})

Task ID: {task_id}
Agent: {agent}
Result: {result}
Context: {context}"""


@dataclass
class Learning:
    category: str  # "fact", "correction", "episode"
    content: dict[str, Any]


class MemoryWriteback:
    def __init__(self, memory: VectorMemory, llm: LLMClient):
        self.memory = memory
        self.llm = llm

    async def on_task_complete(
        self, task_id: str, agent: str, result: str, context: str = ""
    ) -> list[Learning]:
        """Extract and persist learnings from completed task output."""
        learnings = await self._extract_learnings(task_id, agent, result, context)
        for learning in learnings:
            await self._persist(learning, task_id, agent)
        logger.info(f"Persisted {len(learnings)} learnings from task {task_id}")
        return learnings

    async def _extract_learnings(
        self, task_id: str, agent: str, result: str, context: str
    ) -> list[Learning]:
        prompt = EXTRACTION_PROMPT.format(
            task_id=task_id, agent=agent, result=result, context=context
        )
        response = await self.llm.complete(prompt)
        try:
            data = json.loads(response)
        except json.JSONDecodeError:
            logger.warning(f"Failed to parse LLM extraction for task {task_id}")
            return []

        learnings = []
        for fact in data.get("facts", []):
            learnings.append(Learning(category="fact", content=fact))
        for correction in data.get("corrections", []):
            learnings.append(Learning(category="correction", content=correction))
        for episode in data.get("episodes", []):
            learnings.append(Learning(category="episode", content=episode))
        return learnings

    async def _persist(self, learning: Learning, task_id: str, agent: str) -> None:
        metadata = {"task_id": task_id, "agent": agent, "category": learning.category}
        text = json.dumps(learning.content)
        await self.memory.store(text=text, metadata=metadata, namespace=learning.category)
