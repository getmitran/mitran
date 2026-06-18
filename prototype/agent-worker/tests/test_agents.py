"""Tests for Mitran agents: DevAgent, OpsAgent, TicketsAgent SLA, WikiAgent runbooks, memory."""
import json
import sys
import importlib.util
import pytest
from unittest.mock import patch, MagicMock, AsyncMock
from datetime import datetime, timedelta
from pathlib import Path

AGENTS_DIR = Path("/Users/ravitejb/Documents/Mitran/prototype/agent-worker/worker/agents")
WORKER_DIR = Path("/Users/ravitejb/Documents/Mitran/prototype/agent-worker")


def _import_module(name: str, filepath: str):
    """Import a module by file path, bypassing package __init__.py."""
    spec = importlib.util.spec_from_file_location(name, filepath)
    mod = importlib.util.module_from_spec(spec)
    sys.modules[name] = mod
    spec.loader.exec_module(mod)
    return mod


# Pre-load worker.llm to satisfy dev_agent dependency
_llm = _import_module("worker.llm", str(WORKER_DIR / "worker/llm.py"))
_config = _import_module("worker.config", str(WORKER_DIR / "worker/config.py"))

# Now import the modules we need directly
_dev_agent = _import_module("worker.agents.dev_agent", str(AGENTS_DIR / "dev_agent.py"))
DevAgent = _dev_agent.DevAgent
AgentResult = _dev_agent.AgentResult
FileChange = _dev_agent.FileChange
BaseAgent = _dev_agent.BaseAgent

_tickets_sla = _import_module("worker.agents.tickets_sla", str(AGENTS_DIR / "tickets_sla.py"))
SLATracker = _tickets_sla.SLATracker
SLAPolicy = _tickets_sla.SLAPolicy

# wiki_runbooks needs worker.agents.base
_base = _import_module("worker.agents.base", str(AGENTS_DIR / "base.py"))
_wiki_runbooks = _import_module("worker.agents.wiki_runbooks", str(AGENTS_DIR / "wiki_runbooks.py"))
RunbookGenerator = _wiki_runbooks.RunbookGenerator


# --- DevAgent Tests ---

class TestDevAgent:
    def test_handle_task_returns_agent_result(self, dev_agent_context):
        response = json.dumps({
            "files": [{"path": "main.py", "content": "print('hello')"}],
            "summary": "Created main.py"
        })
        with patch.object(_dev_agent, "invoke", return_value=response):
            agent = DevAgent()
            result = agent.execute("Create a hello world script", dev_agent_context)

        assert isinstance(result, AgentResult)
        assert len(result.files) == 1
        assert result.files[0].path == "main.py"
        assert "hello" in result.files[0].content

    def test_handle_task_invalid_json_returns_raw(self, dev_agent_context):
        with patch.object(_dev_agent, "invoke", return_value="This is not JSON"):
            agent = DevAgent()
            result = agent.execute("Do something", dev_agent_context)

        assert isinstance(result, AgentResult)
        assert result.summary == "This is not JSON"
        assert result.files == []

    def test_handle_task_multiple_files(self, dev_agent_context):
        response = json.dumps({
            "files": [
                {"path": "src/main.py", "content": "from .utils import helper"},
                {"path": "src/utils.py", "content": "def helper(): pass"},
                {"path": "requirements.txt", "content": "flask>=3.0"},
            ],
            "summary": "Created project structure"
        })
        with patch.object(_dev_agent, "invoke", return_value=response):
            agent = DevAgent()
            result = agent.execute("Create a Flask app", dev_agent_context)

        assert len(result.files) == 3

    def test_file_change_dataclass(self):
        fc = FileChange(path="test.py", content="pass", action="create")
        assert fc.path == "test.py"
        assert fc.action == "create"

    def test_agent_result_default_files(self):
        r = AgentResult(summary="done")
        assert r.files == []

    def test_base_agent_call_llm(self):
        with patch.object(_dev_agent, "invoke", return_value="test response"):
            agent = BaseAgent()
            agent.system_prompt = "You are helpful"
            result = agent._call_llm("hello")
            assert result == "test response"


# --- OpsAgent Tests ---

class TestOpsAgent:
    def test_ops_v2_class_structure(self):
        """Verify OpsAgentV2 has expected methods via AST."""
        import ast
        with open(AGENTS_DIR / "ops_agent_v2.py") as f:
            tree = ast.parse(f.read())
        classes = {n.name: n for n in ast.walk(tree) if isinstance(n, ast.ClassDef)}
        assert "OpsAgentV2" in classes
        methods = [n.name for n in ast.walk(classes["OpsAgentV2"]) if isinstance(n, ast.FunctionDef)]
        assert "execute" in methods
        assert "_run_cmd" in methods

    def test_ops_v2_system_prompt(self):
        """Verify system prompt contains ops keywords."""
        with open(AGENTS_DIR / "ops_agent_v2.py") as f:
            content = f.read()
        assert "DevOps" in content or "SRE" in content
        assert "monitor" in content.lower()


# --- SLATracker Tests ---

class TestSLATracker:
    def test_default_policies(self):
        tracker = SLATracker()
        assert "critical" in tracker.policies
        assert tracker.policies["critical"].max_hours == 1
        assert tracker.policies["low"].max_hours == 72

    def test_get_deadline(self):
        tracker = SLATracker()
        created = datetime(2026, 6, 1, 10, 0, 0)
        deadline = tracker.get_deadline("high", created)
        assert deadline == created + timedelta(hours=4)

    def test_is_breached_true(self):
        tracker = SLATracker()
        created = datetime(2026, 6, 1, 10, 0, 0)
        resolved = datetime(2026, 6, 1, 15, 0, 0)
        assert tracker.is_breached("high", created, resolved) is True

    def test_is_breached_false(self):
        tracker = SLATracker()
        created = datetime(2026, 6, 1, 10, 0, 0)
        resolved = datetime(2026, 6, 1, 12, 0, 0)
        assert tracker.is_breached("high", created, resolved) is False

    def test_is_breached_unknown_priority_defaults_medium(self):
        tracker = SLATracker()
        created = datetime(2026, 6, 1, 10, 0, 0)
        resolved = datetime(2026, 6, 2, 15, 0, 0)
        assert tracker.is_breached("unknown_priority", created, resolved) is True

    def test_check_breaches(self, sample_tickets):
        tracker = SLATracker()
        breached = tracker.check_breaches(sample_tickets)
        ids = [t["id"] for t in breached]
        assert "T-1" in ids
        assert "T-4" not in ids

    def test_check_breaches_adds_metadata(self, sample_tickets):
        tracker = SLATracker()
        breached = tracker.check_breaches(sample_tickets)
        for t in breached:
            assert t["sla_breached"] is True
            assert "sla_overdue_hours" in t

    def test_auto_escalate_calls_notify(self, sample_tickets):
        tracker = SLATracker()
        notifications = []
        tracker.auto_escalate(sample_tickets, notify_fn=notifications.append)
        assert len(notifications) > 0
        assert "SLA BREACH" in notifications[0]

    def test_auto_escalate_compliant_tickets(self):
        tracker = SLATracker()
        now = datetime.utcnow()
        compliant = [{"id": "T-1", "priority": "low", "created_at": now.isoformat(), "status": "open"}]
        result = tracker.auto_escalate(compliant)
        assert result == []

    def test_sla_report(self, sample_tickets):
        tracker = SLATracker()
        report = tracker.get_sla_report(sample_tickets)
        assert report["total"] == 4
        assert "compliance_rate" in report
        assert report["compliance_rate"] <= 100.0

    def test_sla_report_empty_list(self):
        tracker = SLATracker()
        report = tracker.get_sla_report([])
        assert report["total"] == 0
        assert report["compliance_rate"] == 100.0

    def test_custom_policies(self):
        custom = [SLAPolicy("urgent", 2)]
        tracker = SLATracker(policies=custom)
        assert "urgent" in tracker.policies
        assert "critical" not in tracker.policies


# --- WikiAgent RunbookGenerator Tests ---

class TestRunbookGenerator:
    @pytest.mark.asyncio
    async def test_generate_from_incident(self, sample_incident):
        mock_agent = MagicMock()
        mock_agent.call_llm = AsyncMock(return_value=json.dumps({
            "title": "Runbook: DynamoDB ReadThrottle",
            "symptoms": ["5xx errors", "ReadThrottle alarm"],
            "diagnosis": {"steps": ["Check RCU"], "commands": ["aws cw"], "metrics_to_check": ["ReadThrottle"]},
            "resolution": {"steps": ["Increase RCU"], "rollback": ["Revert"], "estimated_time": "15min"},
            "prevention": {"short_term": ["Auto-scaling"], "long_term": ["Redesign"]},
            "metadata": {"service": "UserService", "severity": "SEV-2"},
        }))

        gen = RunbookGenerator(mock_agent)
        runbook = await gen.generate_from_incident(sample_incident)

        assert runbook["title"] == "Runbook: DynamoDB ReadThrottle"
        assert len(runbook["symptoms"]) == 2
        assert runbook["linked_resources"]["dashboards"] == sample_incident["dashboards"]
        mock_agent.call_llm.assert_called_once()

    @pytest.mark.asyncio
    async def test_generate_handles_invalid_json(self, sample_incident):
        mock_agent = MagicMock()
        mock_agent.call_llm = AsyncMock(return_value="Not valid JSON")

        gen = RunbookGenerator(mock_agent)
        runbook = await gen.generate_from_incident(sample_incident)

        assert runbook["title"] == "Runbook (unparsed)"
        assert runbook["raw_content"] == "Not valid JSON"

    @pytest.mark.asyncio
    async def test_generate_handles_markdown_fenced_json(self, sample_incident):
        mock_agent = MagicMock()
        fenced = '```json\n{"title": "Runbook: Test", "symptoms": [], "diagnosis": {"steps": [], "commands": [], "metrics_to_check": []}, "resolution": {"steps": [], "rollback": [], "estimated_time": "5min"}, "prevention": {"short_term": [], "long_term": []}, "metadata": {}}\n```'
        mock_agent.call_llm = AsyncMock(return_value=fenced)

        gen = RunbookGenerator(mock_agent)
        runbook = await gen.generate_from_incident(sample_incident)
        assert runbook["title"] == "Runbook: Test"

    def test_extract_links(self, sample_incident):
        gen = RunbookGenerator(MagicMock())
        links = gen._extract_links(sample_incident)
        assert links["dashboards"] == ["https://grafana.internal/d/user-service"]
        assert len(links["alerts"]) == 2

    def test_extract_links_empty(self):
        gen = RunbookGenerator(MagicMock())
        assert gen._extract_links({}) == {"dashboards": [], "alerts": [], "related_runbooks": []}

    def test_build_prompt(self, sample_incident):
        gen = RunbookGenerator(MagicMock())
        prompt = gen._build_prompt(sample_incident)
        assert "DynamoDB ReadThrottle" in prompt
        assert "UserService" in prompt


# --- Memory Store/Recall Tests ---

try:
    import chromadb as _chromadb
    HAS_CHROMADB = True
except ImportError:
    HAS_CHROMADB = False


@pytest.mark.skipif(not HAS_CHROMADB, reason="chromadb not installed")
class TestVectorMemory:
    @pytest.fixture
    def memory_store(self, tmp_path):
        from worker.memory.vector_store import VectorMemory
        return VectorMemory(str(tmp_path / "test_chroma"))

    def test_add_and_list_facts(self, memory_store):
        memory_store.add_fact("user.name", "Ravi")
        memory_store.add_fact("user.team", "AFT Support")
        facts = memory_store.list_facts()
        assert len(facts) == 2

    def test_fact_upsert_overwrites(self, memory_store):
        memory_store.add_fact("key1", "value1")
        memory_store.add_fact("key1", "value2")
        facts = memory_store.list_facts()
        assert len(facts) == 1
        assert facts[0]["value"] == "value2"

    def test_add_episode(self, memory_store):
        memory_store.add_episode("Deployed v2.0", "deployment")
        episodes = memory_store.list_episodes()
        assert len(episodes) == 1
        assert episodes[0]["category"] == "deployment"

    def test_search(self, memory_store):
        memory_store.add_fact("db.host", "postgres.internal:5432")
        memory_store.add_episode("Database migration completed")
        results = memory_store.search("database")
        assert len(results) > 0

    def test_add_correction(self, memory_store):
        memory_store.add_correction("git push -f", "git push", "never force push")
        assert memory_store.stats()["corrections"] == 1

    def test_delete_fact(self, memory_store):
        memory_store.add_fact("temp", "val")
        assert memory_store.delete_fact("temp") is True

    def test_stats(self, memory_store):
        memory_store.add_fact("k1", "v1")
        memory_store.add_episode("ep1")
        s = memory_store.stats()
        assert s == {"facts": 1, "episodes": 1, "corrections": 0}

    def test_search_specific_collection(self, memory_store):
        memory_store.add_fact("api.url", "https://api.example.com")
        memory_store.add_episode("API was slow")
        results = memory_store.search("api", collection="facts")
        assert all(r["collection"] == "facts" for r in results)
