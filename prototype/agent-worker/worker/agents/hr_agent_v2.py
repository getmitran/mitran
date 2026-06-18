"""HR Agent v2 - Team management, onboarding docs, org knowledge."""
import logging
from pathlib import Path
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange

log = logging.getLogger(__name__)


class HRAgentV2(BaseAgent):
    agent_type = 'hr'
    system_prompt = """You are an HR and team management specialist. You help with:
    1. Generating onboarding documentation
    2. Creating team handbooks and guidelines
    3. Writing job descriptions and role definitions
    4. Team structure and org chart documentation
    5. Process documentation (standup, retro, planning)
    Output each doc as: === FILE: path/to/file ===\n<content>\n=== END ==="""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        workspace = context.get('workspace_path', '')
        project = context.get('project_name', 'Project')
        team_size = context.get('team_size', 5)

        existing_docs = []
        if workspace:
            for p in Path(workspace).glob('docs/*.md'):
                try:
                    existing_docs.append(f'--- {p.name} ---\n{p.read_text()[:500]}')
                except Exception:
                    pass

        doc_context = '\n'.join(existing_docs[:5]) if existing_docs else 'No existing docs found.'

        prompt = f"""Project: {project}
Team size: {team_size}
Task: {task_description}

Existing documentation:
{doc_context}

Generate the requested HR/team documentation."""

        response = self._call_llm(prompt)
        files = self._parse_files(response)

        if not files:
            files = [FileChange(path='docs/team-handbook.md', action='create', content=f'# {project} Team Handbook\n\n{response}')]

        summary = f'## HR Agent Output\n\nGenerated {len(files)} documents:\n'
        summary += '\n'.join(f'- `{f.path}`' for f in files)

        return AgentResult(summary=summary, files=files)

    def _parse_files(self, response: str) -> list[FileChange]:
        files = []
        parts = response.split('=== FILE: ')
        for part in parts[1:]:
            lines = part.split('\n')
            path = lines[0].strip().rstrip(' =').strip()
            content_lines = []
            for line in lines[1:]:
                if line.strip().startswith('=== END'):
                    break
                content_lines.append(line)
            files.append(FileChange(path=path, action='create', content='\n'.join(content_lines)))
        return files
