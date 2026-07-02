"""Docs Agent - Generates and updates documentation from code changes."""
import logging
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange
from worker.llm import invoke as llm_invoke

log = logging.getLogger(__name__)


class DocumentationAgent(BaseAgent):
    agent_type = "docs"
    system_prompt = """You are a technical writer and documentation specialist. Your role is to generate and maintain clear, accurate documentation from code changes.

When generating documentation:
1. Write for the target audience (developers, users, operators)
2. Include code examples, diagrams (mermaid), and API references
3. Structure with clear headings, table of contents, and cross-references
4. Keep docs close to the code they describe
5. Include setup instructions, prerequisites, and troubleshooting

Documentation types you produce:
- README.md (project overview, quickstart)
- API reference (endpoints, params, responses, errors)
- Architecture docs (system design, data flow, decisions)
- Runbooks (operational procedures, debugging guides)
- Changelogs (what changed, migration steps)

Output documentation files as:
=== FILE: docs/filename.md ===
<content>
=== END ===

Principles:
- Accuracy over completeness - never document behavior you can't verify
- Use examples liberally - show, don't just tell
- Keep a consistent voice and formatting style
- Include "last updated" context where relevant"""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        project = context.get("project_name", "unknown")
        project_id = context.get("project_id", "default")
        source_code = context.get("source_code", "")
        existing_docs = context.get("existing_docs", "")

        # Fetch relevant memories
        memories = self._fetch_memory(project_id, task_description)
        memory_context = self._build_memory_context(memories)

        prompt = f"""{memory_context}Project: {project}
Source code:
{source_code if source_code else 'No source provided'}

Existing documentation:
{existing_docs if existing_docs else 'No existing docs'}

Task: {task_description}

Generate or update documentation files."""

        response = self._call_llm(prompt)
        files = self._parse_files(response)

        if not files and response.strip():
            files.append(FileChange(path="docs/README.md", action="create", content=response))

        summary = f"## Documentation Agent Output\n\nGenerated {len(files)} doc files"
        summary += f"\n\n### Files\n" + "\n".join(f"- `{f.path}`" for f in files)

        # Save episode
        self._save_episode(project_id, f"Docs: {task_description} -> generated {len(files)} doc files")

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
