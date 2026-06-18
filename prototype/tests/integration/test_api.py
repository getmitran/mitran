import pytest
import httpx
import uuid

pytestmark = pytest.mark.integration

BASE_URL = "http://localhost:7780"


@pytest.fixture
def client():
    with httpx.Client(base_url=BASE_URL, timeout=30) as c:
        yield c


@pytest.fixture
def session_id():
    return str(uuid.uuid4())


class TestHealthCheck:
    def test_health_endpoint(self, client):
        resp = client.get("/api/v1/health")
        assert resp.status_code == 200
        data = resp.json()
        assert data.get("status") == "ok"


class TestSessionFlow:
    def test_create_session_send_message_get_response(self, client, session_id):
        # Create session
        resp = client.post("/api/v1/sessions", json={"id": session_id, "name": "test-session"})
        assert resp.status_code in (200, 201)

        # Send message
        resp = client.post(
            f"/api/v1/sessions/{session_id}/messages",
            json={"content": "Hello, Mitran", "role": "user"},
        )
        assert resp.status_code in (200, 201)
        data = resp.json()
        assert "id" in data or "content" in data

        # Get messages
        resp = client.get(f"/api/v1/sessions/{session_id}/messages")
        assert resp.status_code == 200
        messages = resp.json()
        assert isinstance(messages, list)
        assert len(messages) >= 1


class TestMemoryFlow:
    def test_store_and_recall(self, client):
        key = f"test_key_{uuid.uuid4().hex[:8]}"
        value = "integration test value"

        # Store
        resp = client.post("/api/v1/memory", json={"key": key, "value": value})
        assert resp.status_code in (200, 201)

        # Recall
        resp = client.get(f"/api/v1/memory/{key}")
        assert resp.status_code == 200
        data = resp.json()
        assert data.get("value") == value

        # Cleanup
        client.delete(f"/api/v1/memory/{key}")


class TestCronFlow:
    def test_create_list_delete(self, client):
        cron_name = f"test-cron-{uuid.uuid4().hex[:8]}"

        # Create
        resp = client.post(
            "/api/v1/crons",
            json={"name": cron_name, "schedule": "0 * * * *", "task": "echo test"},
        )
        assert resp.status_code in (200, 201)
        cron_id = resp.json().get("id")
        assert cron_id

        # List
        resp = client.get("/api/v1/crons")
        assert resp.status_code == 200
        crons = resp.json()
        assert any(c.get("id") == cron_id or c.get("name") == cron_name for c in crons)

        # Delete
        resp = client.delete(f"/api/v1/crons/{cron_id}")
        assert resp.status_code in (200, 204)


class TestArtifactFlow:
    def test_create_get_delete(self, client):
        artifact_name = f"test-artifact-{uuid.uuid4().hex[:8]}"

        # Create
        resp = client.post(
            "/api/v1/artifacts",
            json={"name": artifact_name, "content": "<h1>Test</h1>", "kind": "widget"},
        )
        assert resp.status_code in (200, 201)
        data = resp.json()
        slug = data.get("slug") or data.get("id")
        assert slug

        # Get
        resp = client.get(f"/api/v1/artifacts/{slug}")
        assert resp.status_code == 200
        data = resp.json()
        assert artifact_name in (data.get("name", ""), data.get("title", ""))

        # Delete
        resp = client.delete(f"/api/v1/artifacts/{slug}")
        assert resp.status_code in (200, 204)
