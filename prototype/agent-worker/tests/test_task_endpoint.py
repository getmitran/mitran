"""Tests for the /execute task endpoint."""
import pytest
from unittest.mock import patch, MagicMock
from worker.app import create_app
from worker.agents.dev_agent import AgentResult


@pytest.fixture
def client():
    app = create_app()
    app.config["TESTING"] = True
    with app.test_client() as c:
        yield c


def test_execute_accepts_valid_json(client):
    mock_result = AgentResult(summary="done", files=[])
    with patch("worker.app.AGENTS", {"dev": MagicMock(execute=MagicMock(return_value=mock_result))}):
        with patch("worker.app._post_result"):
            resp = client.post("/execute", json={
                "task_id": "t-1",
                "agent": "dev",
                "task": "hello",
                "context": {},
            })
    assert resp.status_code == 200


def test_execute_returns_completed_status(client):
    mock_result = AgentResult(summary="built it", files=[])
    with patch("worker.app.AGENTS", {"dev": MagicMock(execute=MagicMock(return_value=mock_result))}):
        with patch("worker.app._post_result"):
            data = client.post("/execute", json={
                "task_id": "t-2",
                "agent": "dev",
                "task": "test task",
                "context": {},
            }).get_json()
    assert data["status"] == "completed"
    assert data["summary"] == "built it"


def test_execute_unknown_agent_returns_404(client):
    resp = client.post("/execute", json={
        "task_id": "t-3",
        "agent": "nonexistent",
        "task": "x",
    })
    assert resp.status_code == 404
