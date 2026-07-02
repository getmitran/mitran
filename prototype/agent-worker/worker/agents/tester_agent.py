"""Tester Agent - Writes and runs tests for developer output."""
import logging
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange
from worker.llm import invoke as llm_invoke

log = logging.getLogger(__name__)


class TesterAgent(BaseAgent):
    agent_type = "tester"
    system_prompt = """You are a senior QA engineer. Your role is to write comprehensive tests, validate developer output, and ensure code quality.

When writing tests:
1. Cover happy paths, edge cases, error conditions, and boundary values
2. Write unit tests, integration tests, and e2e tests as appropriate
3. Use proper test fixtures, mocks, and assertions
4. Follow the project's testing framework conventions
5. Aim for meaningful coverage, not just line coverage

When validating code:
1. Verify the implementation matches requirements
2. Check for untested code paths
3. Identify potential regression risks
4. Suggest additional test scenarios

Output test files as:
=== FILE: tests/test_module.py ===
<content>
=== END ===

Testing principles:
- Tests should be independent, repeatable, and fast
- Use descriptive test names that document behavior
- Assert one logical concept per test
- Avoid testing implementation details - test behavior"""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        project = context.get("project_name", "unknown")
        languages = context.get("languages", ["Python"])
        source_files = context.get("source_code", "")
        lang_str = ", ".join(languages) if isinstance(languages, list) else languages

        prompt = f"""Project: {project}
Languages: {lang_str}
Source code to test:
{source_files if source_files else 'No source provided - generate tests based on task description'}

Task: {task_description}

Generate comprehensive test files."""

        response = self._call_llm(prompt)
        files = self._parse_files(response)

        summary = f"## Tester Agent Output\n\nGenerated {len(files)} test files"
        summary += f"\n\n### Test Files\n" + "\n".join(f"- `{f.path}`" for f in files)

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
            files.append(FileChange(path="tests/test_output.py", action="create", content=response))
        return files
