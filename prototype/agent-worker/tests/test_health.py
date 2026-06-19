"""Tests for the /health endpoint."""
import pytest
from worker.app import create_app


@pytest.fixture
def client():
    app = create_app()
    app.config["TESTING"] = True
    with app.test_client() as c:
        yield c


def test_health_returns_200(client):
    resp = client.get("/health")
    assert resp.status_code == 200


def test_health_returns_status_ok(client):
    data = client.get("/health").get_json()
    assert data["status"] == "ok"


def test_health_includes_agents_list(client):
    data = client.get("/health").get_json()
    assert "agents" in data
    assert isinstance(data["agents"], list)
    assert len(data["agents"]) > 0
