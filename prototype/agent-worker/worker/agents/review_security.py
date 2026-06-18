"""Review Agent - Security Scanner module.

Runs semgrep/bandit on workspace code, parses findings, and auto-creates
tickets for critical vulnerabilities.
"""

import json
import subprocess
import os
from dataclasses import dataclass, field
from typing import Optional


@dataclass
class Finding:
    severity: str
    file: str
    line: int
    rule_id: str
    message: str
    tool: str


class SecurityScanner:
    """Executes static security analysis and manages findings."""

    def __init__(self, ticket_client=None):
        self._findings: list[Finding] = []
        self._ticket_client = ticket_client

    def run_scan(self, workspace_path: str) -> list[Finding]:
        """Run semgrep and bandit on the workspace, return combined findings."""
        self._findings = []
        workspace_path = os.path.abspath(workspace_path)

        semgrep_results = self._run_semgrep(workspace_path)
        bandit_results = self._run_bandit(workspace_path)

        self._findings = semgrep_results + bandit_results
        return self._findings

    def parse_results(self) -> list[dict]:
        """Extract findings as dicts with severity/file/line."""
        return [
            {
                "severity": f.severity,
                "file": f.file,
                "line": f.line,
                "rule_id": f.rule_id,
                "message": f.message,
                "tool": f.tool,
            }
            for f in self._findings
        ]

    def create_tickets(self, findings: Optional[list[Finding]] = None) -> list[dict]:
        """Auto-create tickets for critical/high findings."""
        targets = findings or self._findings
        critical = [f for f in targets if f.severity.lower() in ("critical", "high", "error")]

        tickets = []
        for finding in critical:
            ticket = {
                "title": f"[Security] {finding.rule_id}: {finding.message[:80]}",
                "description": (
                    f"**Tool:** {finding.tool}\n"
                    f"**Severity:** {finding.severity}\n"
                    f"**File:** {finding.file}:{finding.line}\n"
                    f"**Rule:** {finding.rule_id}\n\n"
                    f"{finding.message}"
                ),
                "priority": "High" if finding.severity.lower() == "high" else "Critical",
                "labels": ["security", "auto-scan"],
            }

            if self._ticket_client:
                result = self._ticket_client.create(ticket)
                ticket["id"] = result.get("id")

            tickets.append(ticket)

        return tickets

    def _run_semgrep(self, path: str) -> list[Finding]:
        """Execute semgrep with auto config."""
        try:
            result = subprocess.run(
                ["semgrep", "scan", "--config", "auto", "--json", path],
                capture_output=True, text=True, timeout=300,
            )
        except (FileNotFoundError, subprocess.TimeoutExpired):
            return []

        if not result.stdout.strip():
            return []

        try:
            data = json.loads(result.stdout)
        except json.JSONDecodeError:
            return []

        findings = []
        for r in data.get("results", []):
            findings.append(Finding(
                severity=r.get("extra", {}).get("severity", "WARNING"),
                file=r.get("path", ""),
                line=r.get("start", {}).get("line", 0),
                rule_id=r.get("check_id", "unknown"),
                message=r.get("extra", {}).get("message", ""),
                tool="semgrep",
            ))
        return findings

    def _run_bandit(self, path: str) -> list[Finding]:
        """Execute bandit on Python files."""
        try:
            result = subprocess.run(
                ["bandit", "-r", path, "-f", "json", "-q"],
                capture_output=True, text=True, timeout=300,
            )
        except (FileNotFoundError, subprocess.TimeoutExpired):
            return []

        if not result.stdout.strip():
            return []

        try:
            data = json.loads(result.stdout)
        except json.JSONDecodeError:
            return []

        findings = []
        for r in data.get("results", []):
            findings.append(Finding(
                severity=r.get("issue_severity", "LOW"),
                file=r.get("filename", ""),
                line=r.get("line_number", 0),
                rule_id=r.get("test_id", "unknown"),
                message=r.get("issue_text", ""),
                tool="bandit",
            ))
        return findings
