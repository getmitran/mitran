"""Ops Agent v2 - Real infrastructure monitoring and health checks."""
import logging
import subprocess
import json
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange

log = logging.getLogger(__name__)

class OpsAgentV2(BaseAgent):
    agent_type = 'ops'
    system_prompt = """You are a DevOps/SRE specialist. You monitor infrastructure, check service health,
    analyze logs, and recommend fixes. When given an ops task:
    1. Check service health endpoints
    2. Analyze system metrics
    3. Identify issues and root causes
    4. Recommend remediation steps
    Output findings in a structured report."""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        workspace = context.get('workspace_path', '')
        
        # Gather system health data
        checks = {}
        checks['disk'] = self._run_cmd('df -h / | tail -1')
        checks['memory'] = self._run_cmd('free -h 2>/dev/null || vm_stat 2>/dev/null | head -5')
        checks['load'] = self._run_cmd('uptime')
        checks['docker'] = self._run_cmd('docker ps --format "{{.Names}}: {{.Status}}" 2>/dev/null || echo "Docker not running"')
        checks['ports'] = self._run_cmd('ss -tlnp 2>/dev/null || netstat -tlnp 2>/dev/null | head -20')
        
        # Ask LLM to analyze
        health_report = json.dumps(checks, indent=2)
        prompt = f"""System health data:
{health_report}

User request: {task_description}

Analyze the system health, identify any issues, and provide recommendations.
Format as a structured ops report with sections: Status, Issues, Recommendations."""
        
        analysis = self._call_llm(prompt)
        
        # Write report to workspace
        report_content = f"# Ops Report\n\n{analysis}\n\n## Raw Data\n```json\n{health_report}\n```"
        files = []
        if workspace:
            files.append(FileChange(path='ops-report.md', action='create', content=report_content))
        
        return AgentResult(summary=f'## Ops Agent Report\n\n{analysis}', files=files)
    
    def _run_cmd(self, cmd: str) -> str:
        try:
            result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=10)
            return result.stdout.strip() or result.stderr.strip() or 'No output'
        except Exception as e:
            return f'Error: {e}'
