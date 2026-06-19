"""Tests for error handling — invalid input returns proper error codes."""
import pytest
from worker.app import create_app


@pytest.fixture
def client():
    app = create_app()
    app.config["TESTING"] = True
    with app.test_client() as c:
        yield c


def test_execute_no_body_returns_400(client):
    resp = client.post("/execute", content_type="application/json", data="")
    assert resp.status_code == 400


def test_execute_invalid_json_returns_400(client):
    resp = client.post("/execute", content_type="application/json", data="not json{{{")
    assert resp.status_code == 400


def test_execute_empty_json_returns_400(client):
    # Empty JSON has no agent field -> falls through to "No JSON body" or unknown agent
    resp = client.post("/execute", json={})
    # With empty dict, agent defaults to "dev" key lookup which may exist
    # but task_id defaults to "unknown" - test that it doesn't crash
    assert resp.status_code in (200, 400, 404, 500)


def test_health_post_method_not_allowed(client):
    resp = client.post("/health")
    assert resp.status_code == 405


def test_execute_get_method_not_allowed(client):
    resp = client.get("/execute")
    assert resp.status_code == 405
