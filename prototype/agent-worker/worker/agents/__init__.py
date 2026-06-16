from worker.agents.dev_agent import DevAgent
from worker.agents.docs_agent import DocsAgent
from worker.agents.cicd_agent import CicdAgent
from worker.agents.tickets_agent import TicketsAgent
from worker.agents.wiki_agent import WikiAgent
from worker.agents.ops_agent import OpsAgent
from worker.agents.hr_agent import HrAgent
from worker.agents.dashboard_agent import DashboardAgent

AGENTS = {
    "dev": DevAgent(),
    "docs": DocsAgent(),
    "cicd": CicdAgent(),
    "tickets": TicketsAgent(),
    "wiki": WikiAgent(),
    "ops": OpsAgent(),
    "hr": HrAgent(),
    "dashboard": DashboardAgent(),
}
