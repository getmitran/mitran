from worker.agents.dev_agent import DevAgent
from worker.agents.dev_agent_v2 import DevAgentV2
from worker.agents.docs_agent import DocsAgent
from worker.agents.cicd_agent import CicdAgent
from worker.agents.cicd_agent_v2 import CicdAgentV2
from worker.agents.tickets_agent import TicketsAgent
from worker.agents.wiki_agent import WikiAgent
from worker.agents.ops_agent import OpsAgent
from worker.agents.hr_agent import HrAgent
from worker.agents.dashboard_agent import DashboardAgent

AGENTS = {
    "dev": DevAgentV2(),  # v2 with GitHub PR creation
    "dev_v1": DevAgent(),  # legacy dev agent without GitHub integration
    "docs": DocsAgent(),
    "cicd": CicdAgentV2(),  # v2 with multi-environment support
    "cicd_v1": CicdAgent(),  # legacy single-env agent
    "tickets": TicketsAgent(),
    "wiki": WikiAgent(),
    "ops": OpsAgent(),
    "hr": HrAgent(),
    "dashboard": DashboardAgent(),
}
