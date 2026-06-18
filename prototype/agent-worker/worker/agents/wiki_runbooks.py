"""Wiki Agent Runbooks - Generate structured runbooks from incident data."""

import json
from typing import Any

from worker.agents.base import BaseAgent


class RunbookGenerator:
    """Generates structured runbooks from incident data using LLM."""

    def __init__(self, agent: BaseAgent):
        self.agent = agent

    async def generate_from_incident(self, incident_data: dict[str, Any]) -> dict[str, Any]:
        """Generate a structured runbook from incident data.

        Args:
            incident_data: Dict containing incident details (title, description,
                          service, severity, timeline, resolution_notes, dashboards, alerts).

        Returns:
            Structured runbook dict with symptoms, diagnosis, resolution, prevention,
            and linked dashboards/alerts.
        """
        prompt = self._build_prompt(incident_data)
        response = await self.agent.call_llm(prompt)
        runbook = self._parse_response(response)
        runbook["linked_resources"] = self._extract_links(incident_data)
        return runbook

    def _build_prompt(self, incident_data: dict[str, Any]) -> str:
        return f"""Generate a structured runbook from this incident. Return valid JSON only.

Incident:
- Title: {incident_data.get("title", "Unknown")}
- Service: {incident_data.get("service", "Unknown")}
- Severity: {incident_data.get("severity", "Unknown")}
- Description: {incident_data.get("description", "")}
- Timeline: {incident_data.get("timeline", "")}
- Resolution Notes: {incident_data.get("resolution_notes", "")}

Return JSON with these keys:
{{
  "title": "Runbook: <concise title>",
  "symptoms": ["list of observable symptoms that indicate this issue"],
  "diagnosis": {{
    "steps": ["ordered diagnostic steps"],
    "commands": ["relevant CLI commands or queries"],
    "metrics_to_check": ["key metrics/dashboards to inspect"]
  }},
  "resolution": {{
    "steps": ["ordered resolution steps"],
    "rollback": ["rollback steps if resolution fails"],
    "estimated_time": "estimated resolution time"
  }},
  "prevention": {{
    "short_term": ["immediate preventive measures"],
    "long_term": ["architectural or process improvements"]
  }},
  "metadata": {{
    "service": "{incident_data.get("service", "Unknown")}",
    "severity": "{incident_data.get("severity", "Unknown")}",
    "last_occurred": "{incident_data.get("occurred_at", "")}"
  }}
}}"""

    def _parse_response(self, response: str) -> dict[str, Any]:
        """Parse LLM response into structured runbook."""
        try:
            # Strip markdown code fences if present
            text = response.strip()
            if text.startswith("```"):
                text = text.split("\n", 1)[1]
                text = text.rsplit("```", 1)[0]
            return json.loads(text)
        except (json.JSONDecodeError, IndexError):
            return {
                "title": "Runbook (unparsed)",
                "raw_content": response,
                "symptoms": [],
                "diagnosis": {"steps": [], "commands": [], "metrics_to_check": []},
                "resolution": {"steps": [], "rollback": [], "estimated_time": "unknown"},
                "prevention": {"short_term": [], "long_term": []},
                "metadata": {},
            }

    def _extract_links(self, incident_data: dict[str, Any]) -> dict[str, list[str]]:
        """Extract and organize dashboard/alert links from incident data."""
        return {
            "dashboards": incident_data.get("dashboards", []),
            "alerts": incident_data.get("alerts", []),
            "related_runbooks": incident_data.get("related_runbooks", []),
        }
