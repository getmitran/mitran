"""Wiki Agent v2 - Maintains project wiki with auto-generated pages."""
import logging
from pathlib import Path
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange

log = logging.getLogger(__name__)


class WikiAgentV2(BaseAgent):
    agent_type = 'wiki'
    system_prompt = """You are a wiki content specialist. You create and maintain project wiki pages.
    When given a wiki task:
    1. Analyze the project structure
    2. Generate well-structured wiki pages with proper linking
    3. Include code examples, diagrams (mermaid), and references
    4. Maintain a table of contents
    Output each page as: === FILE: wiki/page-name.md ===\n<content>\n=== END ==="""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        workspace = context.get('workspace_path', '')
        project = context.get('project_name', 'Project')

        existing_pages = []
        source_summary = []
        if workspace:
            wiki_dir = Path(workspace) / 'wiki'
            if wiki_dir.exists():
                for p in wiki_dir.glob('*.md'):
                    existing_pages.append(p.name)
            for ext in ['*.py', '*.go', '*.ts', '*.tsx']:
                files = list(Path(workspace).rglob(ext))[:5]
                source_summary.extend(str(f.relative_to(workspace)) for f in files)

        prompt = f"""Project: {project}
Existing wiki pages: {existing_pages or 'None'}
Source files: {source_summary[:15] or 'None found'}
Task: {task_description}

Generate wiki page(s) with proper markdown, mermaid diagrams where useful, and cross-links."""

        response = self._call_llm(prompt)
        files = self._parse_files(response)

        if not files:
            files = [FileChange(path='wiki/index.md', action='create',
                               content=f'# {project} Wiki\n\n{response}')]

        summary = f'## Wiki Agent Output\n\nGenerated {len(files)} wiki pages:\n'
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
