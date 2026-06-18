"""Docs Agent v2 - Generates real documentation from source code."""
import os
import logging
from pathlib import Path
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange

log = logging.getLogger(__name__)

class DocsAgentV2(BaseAgent):
    agent_type = 'docs'
    system_prompt = """You are a technical documentation specialist. You generate comprehensive docs from source code.
    When given code files:
    1. Analyze the code structure and purpose
    2. Generate API documentation
    3. Create usage examples
    4. Write a clear README with installation, usage, and architecture sections
    Output each doc file as: === FILE: path/to/file ===\n<content>\n=== END ==="""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        workspace = context.get('workspace_path', '')
        languages = context.get('languages', ['Python'])
        project = context.get('project_name', 'project')
        
        # Scan workspace for source files
        source_files = []
        if workspace and os.path.isdir(workspace):
            for ext in self._get_extensions(languages):
                for p in Path(workspace).rglob(f'*{ext}'):
                    if '.git' not in str(p) and 'node_modules' not in str(p):
                        try:
                            content = p.read_text()[:2000]  # First 2000 chars
                            source_files.append(f'--- {p.name} ---\n{content}')
                        except: pass
                    if len(source_files) >= 10: break
        
        # Generate docs via LLM
        source_context = '\n\n'.join(source_files[:10]) if source_files else 'No source files found.'
        prompt = f"""Project: {project}
Languages: {', '.join(languages) if isinstance(languages, list) else languages}
Task: {task_description}

Source files:
{source_context}

Generate comprehensive documentation files including README.md, API docs, and usage guide."""
        
        response = self._call_llm(prompt)
        files = self._parse_files(response)
        
        # Fallback: generate a basic README
        if not files:
            readme = f'# {project}\n\n{response}'
            files = [FileChange(path='docs/README.md', action='create', content=readme)]
        
        summary = f'## Docs Agent Output\n\nGenerated {len(files)} documentation files:\n'
        summary += '\n'.join(f'- `{f.path}`' for f in files)
        
        return AgentResult(summary=summary, files=files)
    
    def _get_extensions(self, languages):
        ext_map = {'Python': '.py', 'Go': '.go', 'JavaScript': '.js', 'TypeScript': '.ts', 'Rust': '.rs', 'Java': '.java'}
        if isinstance(languages, str): languages = [languages]
        return [ext_map.get(l, '.py') for l in languages]
    
    def _parse_files(self, response: str) -> list[FileChange]:
        files = []
        parts = response.split('=== FILE: ')
        for part in parts[1:]:
            lines = part.split('\n')
            path = lines[0].strip().rstrip(' =').strip()
            content_lines = []
            for line in lines[1:]:
                if line.strip().startswith('=== END'): break
                content_lines.append(line)
            files.append(FileChange(path=path, action='create', content='\n'.join(content_lines)))
        return files
