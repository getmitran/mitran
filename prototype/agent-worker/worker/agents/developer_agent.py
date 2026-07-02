"""Developer Agent - Writes code, fixes bugs, implements features."""
import logging
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange
from worker.llm import invoke as llm_invoke

log = logging.getLogger(__name__)


class DeveloperAgent(BaseAgent):
    agent_type = "developer"
    system_prompt = """You are a senior software developer. Your role is to write production-quality code, fix bugs, and implement features.

When given a task:
1. Analyze requirements and plan the implementation
2. Write clean, well-structured code with proper error handling
3. Include type hints, docstrings, and inline comments where helpful
4. Follow the project's existing conventions and patterns
5. Consider edge cases and failure modes

Output each file as:
=== FILE: path/to/file ===
<content>
=== END ===

Rules:
- Generate real, complete, runnable code - no placeholders or TODOs
- Include proper imports and dependency declarations
- Follow language-specific best practices (PEP 8 for Python, etc.)
- Handle errors gracefully with meaningful messages
- Write self-documenting code with clear variable/function names"""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        project = context.get("project_name", "unknown")
        project_id = context.get("project_id", "default")
        languages = context.get("languages", ["Python"])
        lang_str = ", ".join(languages) if isinstance(languages, list) else languages

        # Fetch relevant memories
        memories = self._fetch_memory(project_id, task_description)
        memory_context = self._build_memory_context(memories)

        prompt = f"""{memory_context}Project: {project}
Languages: {lang_str}
Task: {task_description}

Generate production-quality code files to complete this task."""

        response = self._call_llm(prompt)
        files = self._parse_files(response)

        summary = f"## Developer Agent Output\n\nGenerated {len(files)} files"
        summary += f"\n\n### Files\n" + "\n".join(f"- `{f.path}`" for f in files)

        # Save episode
        self._save_episode(project_id, f"Developer: {task_description} -> generated {len(files)} files")

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
        if not files and response.strip():
            files.append(FileChange(path="output.txt", action="create", content=response))
        return files
