"""Mitran Agent Registry - All agents (7 new focused + existing legacy)."""
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange

# New focused agents (v0.2.0)
from worker.agents.developer_agent import DeveloperAgent
from worker.agents.architect_agent import ArchitectAgent
from worker.agents.devops_agent import DevOpsAgent
from worker.agents.tester_agent import TesterAgent
from worker.agents.pm_agent import PMAgent
from worker.agents.docs_agent_new import DocumentationAgent
from worker.agents.security_agent import SecurityAgent

# Legacy agents (v0.1.0 - kept for backward compatibility)
from worker.agents.dev_agent import DevAgent
from worker.agents.docs_agent import DocsAgent
from worker.agents.ops_agent import OpsAgent
from worker.agents.cicd_agent import CicdAgent
from worker.agents.tickets_agent import TicketsAgent
from worker.agents.wiki_agent import WikiAgent
from worker.agents.hr_agent import HrAgent
from worker.agents.dashboard_agent import DashboardAgent
from worker.agents.cicd_agent_v2 import CicdAgentV2

AGENTS = {
    # Primary agents (new)
    "developer": DeveloperAgent(),
    "architect": ArchitectAgent(),
    "devops": DevOpsAgent(),
    "tester": TesterAgent(),
    "pm": PMAgent(),
    "docs": DocumentationAgent(),
    "security": SecurityAgent(),
    # Legacy agents (still usable)
    "dev": DevAgent(),
    "docs-legacy": DocsAgent(),
    "ops": OpsAgent(),
    "cicd": CicdAgent(),
    "cicd-v2": CicdAgentV2(),
    "tickets": TicketsAgent(),
    "wiki": WikiAgent(),
    "hr": HrAgent(),
    "dashboard": DashboardAgent(),
}

__all__ = ["AGENTS", "BaseAgent", "AgentResult", "FileChange"]
