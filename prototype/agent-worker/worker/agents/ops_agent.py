import json
from worker.agents.dev_agent import AgentResult, FileChange
from worker.llm import invoke


SYSTEM_PROMPT = """You are Mitran's Ops Agent — an SRE that generates production monitoring and observability configurations.

Given a service/project context, generate monitoring stack configs. Output ONLY valid JSON:
{
  "files": [
    {"path": "monitoring/prometheus.yml", "content": "full YAML content"}
  ],
  "summary": "Brief description of monitoring setup generated"
}

Generate:
- monitoring/prometheus.yml — Prometheus config with scrape targets for the service
- monitoring/alertmanager.yml — AlertManager config with notification routes (Slack, email)
- monitoring/alerts/rules.yml — Alert rules (high latency, error rate, disk, memory)
- monitoring/docker-compose.yml — Docker compose for Prometheus + Grafana + AlertManager
- monitoring/grafana/provisioning/datasources.yml — Grafana datasource config

Rules:
- Configure realistic scrape intervals and retention
- Include standard SRE alerts (RED metrics: rate, errors, duration)
- Add infrastructure alerts (CPU, memory, disk)
- Include meaningful alert annotations and labels
- Use sensible thresholds that won't spam"""


class OpsAgent:
    name = "ops"
    agent_type = "operations"

    def execute(self, task: str, context: dict) -> AgentResult:
        user_msg = f"Project: {context.get('team_description', 'Software project')}\nStack: {context.get('tech_stack', 'Python')}\nService name: {context.get('service_name', 'app')}\n\nTask: {task}"
        raw = invoke(SYSTEM_PROMPT, user_msg)
        try:
            data = json.loads(raw)
        except json.JSONDecodeError:
            return AgentResult(summary=raw, files=[])
        files = [FileChange(path=f["path"], content=f["content"]) for f in data.get("files", [])]
        return AgentResult(summary=data.get("summary", "Ops agent completed"), files=files)
