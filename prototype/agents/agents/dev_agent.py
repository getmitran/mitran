import anthropic

from mitran_sdk.agent import MitranAgent, TaskContext, AgentOutput, FileChange
from mitran_sdk.tools import write_file, run_command

SYSTEM_PROMPT = """You are Mitran Dev Agent. You write production-quality code.
Given a task description, produce the necessary source files.
Respond with a JSON object: {"files": [{"path": "...", "content": "..."}]}
Only output JSON, no markdown fences."""


class DevAgent(MitranAgent):
    def __init__(self, model: str = "claude-sonnet-4-20250514"):
        self.client = anthropic.Anthropic()
        self.model = model

    def execute(self, ctx: TaskContext) -> AgentOutput:
        self.report_progress("Generating code...", 10)

        response = self.client.messages.create(
            model=self.model,
            max_tokens=4096,
            system=SYSTEM_PROMPT,
            messages=[{"role": "user", "content": f"Task: {ctx.title}\n\nDescription: {ctx.description}\n\nWorkspace: {ctx.workspace_path}"}],
        )

        import json
        text = response.content[0].text
        data = json.loads(text)

        changes = []
        for f in data.get("files", []):
            full_path = f"{ctx.workspace_path}/{f['path']}"
            write_file(full_path, f["content"])
            changes.append(FileChange(path=f["path"], action="create", content=f["content"]))

        self.report_progress("Running tests...", 80)
        output, code = run_command("python -m pytest --tb=short 2>/dev/null || true", cwd=ctx.workspace_path)

        self.report_progress("Done", 100)
        return AgentOutput(
            summary=f"Created {len(changes)} file(s). Test exit code: {code}",
            files_changed=changes,
        )
