"""Security Agent - Vulnerability scanning, secrets detection, policy enforcement."""
import logging
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange
from worker.llm import invoke as llm_invoke

log = logging.getLogger(__name__)


class SecurityAgent(BaseAgent):
    agent_type = "security"
    system_prompt = """You are a senior application security engineer. Your role is to scan for vulnerabilities, detect secrets, and enforce security policies.

When scanning code:
1. Check for OWASP Top 10 vulnerabilities
2. Detect hardcoded secrets, API keys, tokens, passwords
3. Identify insecure dependencies and outdated packages
4. Flag improper input validation and output encoding
5. Check authentication/authorization implementation

When enforcing policies:
1. Verify least-privilege access patterns
2. Check encryption at rest and in transit
3. Validate secure configuration defaults
4. Ensure audit logging is present
5. Verify secrets management (no plaintext)

Output security findings as:
## Security Scan Results

### Critical
- [VULN-001] Description | File:Line | Remediation

### High
- [VULN-002] Description | File:Line | Remediation

### Medium / Low / Info

### Secrets Detected
- [SECRET-001] Type | File:Line | Action required

### Policy Violations
- [POLICY-001] Rule violated | File:Line | Required change

Always provide:
- Severity classification (Critical/High/Medium/Low/Info)
- Specific file and line references
- Concrete remediation steps with code examples"""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        project = context.get("project_name", "unknown")
        source_code = context.get("source_code", "")
        dependencies = context.get("dependencies", "")

        prompt = f"""Project: {project}
Source code to scan:
{source_code if source_code else 'No source provided'}

Dependencies:
{dependencies if dependencies else 'No dependency list provided'}

Task: {task_description}

Perform a thorough security analysis and report findings."""

        response = self._call_llm(prompt)

        files = [FileChange(
            path="SECURITY_REPORT.md",
            action="create",
            content=response
        )]

        return AgentResult(summary=response, files=files)
