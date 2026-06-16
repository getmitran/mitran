from mitran_sdk.agent import TaskContext, AgentOutput, FileChange


def test_task_context_creation():
    ctx = TaskContext(task_id="t1", title="Test", description="A test task", workspace_path="/tmp")
    assert ctx.task_id == "t1"
    assert ctx.history == []


def test_agent_output_defaults():
    out = AgentOutput(summary="done", files_changed=[])
    assert out.questions is None
    assert out.next_tasks is None


def test_file_change():
    fc = FileChange(path="main.py", action="create", content="print('hi')")
    assert fc.diff is None
    assert fc.action == "create"


def test_agent_output_with_changes():
    changes = [
        FileChange(path="a.py", action="create", content="x=1"),
        FileChange(path="b.py", action="modify", content="y=2", diff="+y=2"),
    ]
    out = AgentOutput(summary="2 files", files_changed=changes, questions=["confirm?"])
    assert len(out.files_changed) == 2
    assert out.questions == ["confirm?"]
