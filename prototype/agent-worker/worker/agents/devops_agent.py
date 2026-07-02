"""DevOps Agent - Debugs failures, manages deployments, infrastructure."""
import logging
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange
from worker.llm import invoke as llm_invoke

log = logging.getLogger(__name__)


class DevOpsAgent(BaseAgent):
    agent_type = "devops"
    system_prompt = """You are a senior DevOps/SRE engineer. Your role is to debug production failures, manage deployments, and maintain infrastructure.

When debugging:
1. Analyze error logs, stack traces, and metrics
2. Identify root cause through systematic elimination
3. Propose immediate mitigation and long-term fix
4. Document the incident timeline

When managing infrastructure:
1. Write IaC (Terraform, CloudFormation, CDK) following best practices
2. Configure CI/CD pipelines with proper gating
3. Set up monitoring, alerting, and observability
4. Implement security hardening and compliance

Output infrastructure files as:
=== FILE: path/to/file ===
<content>
=== END ===

For debugging, output structured analysis:
## Root Cause
## Impact
## Mitigation Steps
## Prevention

Always consider: rollback plans, blast radius, and monitoring gaps."""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        project = context.get("project_name", "unknown")
        project_id = context.get("project_id", "default")
        environment = context.get("environment", "production")
        logs = context.get("error_logs", "")

        # Fetch relevant memories
        memories = self._fetch_memory(project_id, task_description)
        memory_context = self._build_memory_context(memories)

        prompt = f"""{memory_context}Project: {project}
Environment: {environment}
Error context:
{logs if logs else 'No logs provided'}

Task: {task_description}

Provide analysis and/or infrastructure code as needed."""

        response = self._call_llm(prompt)
        files = self._parse_files(response)

        summary = response if not files else f"## DevOps Agent Output\n\nGenerated {len(files)} files\n\n### Files\n" + "\n".join(f"- `{f.path}`" for f in files)

        # Save episode
        self._save_episode(project_id, f"DevOps: {task_description} -> {'generated ' + str(len(files)) + ' files' if files else 'analysis complete'}")

        return AgentResult(summary=summary, files=files)

    def _parse_files(self, response: str) -> list[FileChange]:
        files = []
        parts = response.split("=== FILE: ")
        for part in parts[1:]:
            lines = part.split("\n")
            path = lines[0].strip().rstrip(" =").strip()
            content_lines = []
            for line in lines[1:]:
                if line.strip().startswith("=== END"):
                    break
                content_lines.append(line)
            content = "\n".join(content_lines)
            files.append(FileChange(path=path, action="create", content=content))
        return files
