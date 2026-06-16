import anthropic

from mitran_sdk.agent import MitranAgent, TaskContext, AgentOutput, FileChange
from mitran_sdk.tools import read_file, write_file, search_files

SYSTEM_PROMPT = """You are Mitran Docs Agent. You generate clear, concise documentation.
Given source files, produce markdown documentation (README.md and API docs).
Respond with JSON: {"files": [{"path": "...", "content": "..."}]}
Only output JSON, no markdown fences."""


class DocsAgent(MitranAgent):
    def __init__(self, model: str = "claude-sonnet-4-20250514"):
        self.client = anthropic.Anthropic()
        self.model = model

    def execute(self, ctx: TaskContext) -> AgentOutput:
        self.report_progress("Scanning codebase...", 10)

        source_files = search_files("*.py", ctx.workspace_path)[:20]
        code_context = ""
        for f in source_files:
            try:
                content = read_file(f)
                code_context += f"\n--- {f} ---\n{content}\n"
            except Exception:
                continue

        self.report_progress("Generating docs...", 50)

        response = self.client.messages.create(
            model=self.model,
            max_tokens=4096,
            system=SYSTEM_PROMPT,
            messages=[{"role": "user", "content": f"Generate documentation for this project:\n\n{code_context[:8000]}"}],
        )

        import json
        data = json.loads(response.content[0].text)

        changes = []
        for f in data.get("files", []):
            full_path = f"{ctx.workspace_path}/{f['path']}"
            write_file(full_path, f["content"])
            changes.append(FileChange(path=f["path"], action="create", content=f["content"]))

        self.report_progress("Done", 100)
        return AgentOutput(
            summary=f"Generated {len(changes)} documentation file(s)",
            files_changed=changes,
        )
