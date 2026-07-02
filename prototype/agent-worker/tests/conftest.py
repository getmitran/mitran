"""Shared test fixtures for Mitran agent-worker tests."""
import pytest
from unittest.mock import AsyncMock, MagicMock, patch
from datetime import datetime, timedelta


@pytest.fixture
def mock_llm_response():
    """Returns a context manager that patches the llm invoke function."""
    def _mock(response_text: str):
        # Patch at the source module level since dev_agent imports `invoke` from worker.llm
        return patch("worker.llm.invoke", return_value=response_text)
    return _mock


@pytest.fixture
def mock_bedrock_client():
    """Mock boto3 bedrock-runtime client."""
    import json
    client = MagicMock()
    body_mock = MagicMock()
    body_mock.read.return_value = json.dumps({
        "content": [{"text": "mocked LLM response"}]
    }).encode()
    client.invoke_model.return_value = {"body": body_mock}
    return client


@pytest.fixture
def sample_tickets():
    """Sample tickets for SLA testing."""
    now = datetime.utcnow()
    return [
        {"id": "T-1", "title": "Critical bug", "priority": "critical",
         "created_at": (now - timedelta(hours=2)).isoformat(), "status": "open"},
        {"id": "T-2", "title": "Feature request", "priority": "low",
         "created_at": (now - timedelta(hours=10)).isoformat(), "status": "open"},
        {"id": "T-3", "title": "Medium task", "priority": "medium",
         "created_at": (now - timedelta(hours=5)).isoformat(), "status": "open"},
        {"id": "T-4", "title": "Done ticket", "priority": "critical",
         "created_at": (now - timedelta(hours=5)).isoformat(), "status": "done"},
    ]


@pytest.fixture
def sample_incident():
    """Sample incident data for runbook generation."""
    return {
        "title": "DynamoDB ReadThrottle spike on UserTable",
        "service": "UserService",
        "severity": "SEV-2",
        "description": "Read capacity exceeded causing 5xx errors on /api/users",
        "timeline": "14:00 alert fired, 14:05 on-call paged, 14:20 capacity increased",
        "resolution_notes": "Increased RCU from 500 to 2000. Root cause: batch job scanning full table.",
        "dashboards": ["https://grafana.internal/d/user-service"],
        "alerts": ["UserTable-ReadThrottle", "UserService-5xx-rate"],
        "occurred_at": "2026-06-15T14:00:00Z",
    }


@pytest.fixture
def dev_agent_context():
    """Standard context for DevAgent tests."""
    return {
        "project_name": "test-project",
        "languages": ["Python"],
        "workspace_path": "/tmp/test-workspace",
        "team_description": "Test team",
        "tech_stack": "Python",
        "team_size": 5,
    }
