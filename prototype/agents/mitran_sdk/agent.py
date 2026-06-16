from abc import ABC, abstractmethod
from dataclasses import dataclass, field


@dataclass
class FileChange:
    path: str
    action: str  # "create", "modify", "delete"
    content: str | None = None
    diff: str | None = None


@dataclass
class AgentOutput:
    summary: str
    files_changed: list[FileChange] = field(default_factory=list)
    questions: list[str] | None = None
    next_tasks: list[dict] | None = None


@dataclass
class TaskContext:
    task_id: str
    title: str
    description: str
    workspace_path: str
    history: list[dict] = field(default_factory=list)


class MitranAgent(ABC):
    @abstractmethod
    def execute(self, ctx: TaskContext) -> AgentOutput:
        ...

    def ask_human(self, question: str) -> str:
        return input(f"[MITRAN] {question}\n> ")

    def report_progress(self, msg: str, pct: float = 0.0) -> None:
        print(f"[{pct:.0f}%] {msg}")
