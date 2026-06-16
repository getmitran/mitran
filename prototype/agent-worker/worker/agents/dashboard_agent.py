import json
from worker.agents.dev_agent import AgentResult, FileChange
from worker.llm import invoke


SYSTEM_PROMPT = """You are Mitran's Dashboard Agent — a data visualization engineer that generates Grafana dashboard configurations.

Given a service context, generate Grafana dashboard JSON. Output ONLY valid JSON:
{
  "files": [
    {"path": "dashboards/service-overview.json", "content": "full Grafana dashboard JSON"}
  ],
  "summary": "Brief description of dashboards generated"
}

Generate:
- dashboards/service-overview.json — main service dashboard with request rate, latency p50/p95/p99, error rate, active connections
- dashboards/infrastructure.json — infra dashboard with CPU, memory, disk, network panels
- dashboards/business-metrics.json — business KPI dashboard template

Rules:
- Use valid Grafana dashboard JSON format (dashboard model v30+)
- Include proper datasource references (use ${DS_PROMETHEUS} variable)
- Create meaningful panel titles and descriptions
- Use appropriate visualization types (time series, stat, gauge, table)
- Include template variables for environment and instance filtering
- Set sensible time ranges and refresh intervals"""


class DashboardAgent:
    name = "dashboard"
    agent_type = "observability"

    def execute(self, task: str, context: dict) -> AgentResult:
        user_msg = f"Service: {context.get('company_description', 'Web service')}\nStack: {context.get('tech_stack', 'Python')}\n\nTask: {task}"
        raw = invoke(SYSTEM_PROMPT, user_msg)
        try:
            data = json.loads(raw)
        except json.JSONDecodeError:
            return AgentResult(summary=raw, files=[])
        files = [FileChange(path=f["path"], content=f["content"]) for f in data.get("files", [])]
        return AgentResult(summary=data.get("summary", "Dashboard agent completed"), files=files)
