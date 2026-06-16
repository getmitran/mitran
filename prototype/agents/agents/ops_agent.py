import json

from mitran_sdk.agent import MitranAgent, TaskContext, AgentOutput, FileChange
from mitran_sdk.tools import write_file


class OpsAgent(MitranAgent):
    def execute(self, ctx: TaskContext) -> AgentOutput:
        self.report_progress("Generating monitoring config...", 30)

        prometheus_yml = """global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'mitran'
    static_targets:
      - targets: ['localhost:8080']
    metrics_path: /metrics
"""

        grafana_dashboard = {
            "dashboard": {
                "title": f"{ctx.title} - Monitoring",
                "panels": [
                    {
                        "title": "Request Rate",
                        "type": "graph",
                        "targets": [{"expr": "rate(http_requests_total[5m])"}],
                    },
                    {
                        "title": "Error Rate",
                        "type": "graph",
                        "targets": [{"expr": "rate(http_requests_total{status=~\"5..\"}[5m])"}],
                    },
                    {
                        "title": "Latency P99",
                        "type": "graph",
                        "targets": [{"expr": "histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))"}],
                    },
                ],
            }
        }

        changes = []
        prom_path = f"{ctx.workspace_path}/monitoring/prometheus.yml"
        write_file(prom_path, prometheus_yml)
        changes.append(FileChange(path="monitoring/prometheus.yml", action="create", content=prometheus_yml))

        dash_content = json.dumps(grafana_dashboard, indent=2)
        dash_path = f"{ctx.workspace_path}/monitoring/grafana-dashboard.json"
        write_file(dash_path, dash_content)
        changes.append(FileChange(path="monitoring/grafana-dashboard.json", action="create", content=dash_content))

        self.report_progress("Done", 100)
        return AgentOutput(
            summary="Generated Prometheus config and Grafana dashboard",
            files_changed=changes,
        )
