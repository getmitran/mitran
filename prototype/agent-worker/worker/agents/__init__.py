from worker.agents.dev_agent import DevAgent
from worker.agents.dev_agent_v2 import DevAgentV2
from worker.agents.docs_agent import DocsAgent
from worker.agents.docs_agent_v2 import DocsAgentV2
from worker.agents.cicd_agent import CicdAgent
from worker.agents.cicd_agent_v2 import CicdAgentV2
from worker.agents.tickets_agent import TicketsAgent
from worker.agents.tickets_agent_v2 import TicketsAgentV2
from worker.agents.wiki_agent import WikiAgent
from worker.agents.wiki_agent_v2 import WikiAgentV2
from worker.agents.ops_agent import OpsAgent
from worker.agents.ops_agent_v2 import OpsAgentV2
from worker.agents.review_agent_v2 import ReviewAgentV2
from worker.agents.hr_agent import HrAgent
from worker.agents.hr_agent_v2 import HRAgentV2
from worker.agents.dashboard_agent import DashboardAgent

AGENTS = {
    "dev": DevAgentV2(),  # v2 with GitHub PR creation
    "dev_v1": DevAgent(),  # legacy dev agent without GitHub integration
    "docs": DocsAgentV2(),  # v2 with real code-based doc generation
    "docs_v1": DocsAgent(),  # legacy docs agent
    "cicd": CicdAgentV2(),  # v2 with multi-environment support
    "cicd_v1": CicdAgent(),  # legacy single-env agent
    "tickets": TicketsAgentV2(),  # v2 with auto-triage and routing
    "tickets_v1": TicketsAgent(),  # legacy tickets agent
    "wiki": WikiAgentV2(),  # v2 with auto-generated pages
    "wiki_v1": WikiAgent(),  # legacy wiki agent
    "ops": OpsAgentV2(),  # v2 with real infra monitoring
    "ops_v1": OpsAgent(),  # legacy ops agent
    "review": ReviewAgentV2(),  # v2 with diff analysis and structured feedback
    "hr": HRAgentV2(),  # v2 with onboarding docs and org knowledge
    "hr_v1": HrAgent(),  # legacy hr agent
    "dashboard": DashboardAgent(),
}
