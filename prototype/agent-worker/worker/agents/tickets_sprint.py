"""Tickets Agent - Sprint Planning module."""

import json
import math
from dataclasses import dataclass


@dataclass
class Ticket:
    id: str
    title: str
    points: int
    assignee: str = ""
    priority: int = 0


class SprintPlanner:
    """Plans sprints by selecting scope, estimating velocity, and balancing workload."""

    def __init__(self, llm_client=None):
        self.llm = llm_client

    async def suggest_scope(self, backlog: list[dict], velocity: int) -> list[dict]:
        """Use LLM to select tickets from backlog that fit within velocity capacity."""
        if not backlog:
            return []

        prompt = (
            f"Given a sprint velocity of {velocity} points, select tickets from this backlog "
            f"that best fit the sprint. Prioritize by priority field (lower=higher priority), "
            f"then by points fitting within capacity.\n\n"
            f"Backlog:\n{json.dumps(backlog, indent=2)}\n\n"
            f"Return a JSON array of selected ticket IDs that fit within {velocity} points total. "
            f"Format: [\"id1\", \"id2\", ...]"
        )

        if self.llm:
            response = await self.llm.invoke(prompt)
            try:
                selected_ids = json.loads(response)
                return [t for t in backlog if t.get("id") in selected_ids]
            except (json.JSONDecodeError, TypeError):
                pass

        # Fallback: greedy selection by priority then points
        sorted_backlog = sorted(backlog, key=lambda t: (t.get("priority", 99), -t.get("points", 0)))
        selected = []
        remaining = velocity
        for ticket in sorted_backlog:
            pts = ticket.get("points", 0)
            if pts <= remaining:
                selected.append(ticket)
                remaining -= pts
        return selected

    def estimate_velocity(self, history: list[dict]) -> int:
        """Calculate velocity from past sprint completions.

        Args:
            history: list of {"sprint": str, "completed_points": int}
        """
        if not history:
            return 0

        points = [h.get("completed_points", 0) for h in history]

        if len(points) <= 2:
            return math.ceil(sum(points) / len(points))

        # Weighted average: recent sprints count more
        weights = list(range(1, len(points) + 1))
        weighted_sum = sum(p * w for p, w in zip(points, weights))
        return math.ceil(weighted_sum / sum(weights))

    def balance_workload(self, assignees: list[str], tickets: list[dict]) -> dict[str, list[dict]]:
        """Distribute tickets across assignees balancing total points.

        Returns mapping of assignee -> assigned tickets.
        """
        if not assignees:
            return {}
        if not tickets:
            return {a: [] for a in assignees}

        # Sort tickets by points descending for better distribution
        sorted_tickets = sorted(tickets, key=lambda t: t.get("points", 0), reverse=True)

        # Initialize workload tracker
        workload: dict[str, list[dict]] = {a: [] for a in assignees}
        load: dict[str, int] = {a: 0 for a in assignees}

        # Assign pre-assigned tickets first
        unassigned = []
        for ticket in sorted_tickets:
            if ticket.get("assignee") and ticket["assignee"] in assignees:
                workload[ticket["assignee"]].append(ticket)
                load[ticket["assignee"]] += ticket.get("points", 0)
            else:
                unassigned.append(ticket)

        # Distribute remaining tickets to least-loaded assignee
        for ticket in unassigned:
            least_loaded = min(assignees, key=lambda a: load[a])
            workload[least_loaded].append(ticket)
            load[least_loaded] += ticket.get("points", 0)

        return workload
