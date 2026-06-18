"""Review Agent v2 - Analyzes code diffs and provides structured review feedback."""
import logging
import subprocess
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange

log = logging.getLogger(__name__)


class ReviewAgentV2(BaseAgent):
    agent_type = 'review'
    system_prompt = """You are a senior code reviewer. Analyze code changes and provide:
    1. Summary of changes
    2. Potential bugs or issues
    3. Security concerns
    4. Performance implications
    5. Suggestions for improvement
    Be specific with line references. Rate severity: critical/warning/info."""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        workspace = context.get('workspace_path', '')

        # Get git diff if in a git repo
        diff = ''
        if workspace:
            diff = self._run_cmd(f'cd {workspace} && git diff HEAD~1 2>/dev/null || git diff --cached 2>/dev/null || echo "No git changes found"')

        if not diff or diff == 'No git changes found':
            # Fallback: scan recent files
            diff = self._run_cmd(f'cd {workspace} && find . -name "*.py" -o -name "*.go" -o -name "*.ts" | head -5 | xargs head -50 2>/dev/null || echo "No files"')

        prompt = f"""Review request: {task_description}

Code changes/context:
```
{diff[:8000]}
```

Provide a structured code review with:
- Summary
- Issues (critical/warning/info)
- Security concerns
- Suggestions"""

        analysis = self._call_llm(prompt)

        files = []
        if workspace:
            files.append(FileChange(path='review-report.md', action='create', content=f'# Code Review\n\n{analysis}'))

        return AgentResult(summary=f'## Code Review\n\n{analysis}', files=files)

    def _run_cmd(self, cmd: str) -> str:
        try:
            result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=15)
            return result.stdout.strip() or result.stderr.strip() or ''
        except Exception as e:
            return f'Error: {e}'
