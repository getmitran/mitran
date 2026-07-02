import json
from worker.agents.dev_agent import AgentResult, FileChange
from worker.llm import invoke


SYSTEM_PROMPT = """You are Mitran's CI/CD Agent — a DevOps engineer that generates production-grade pipeline configurations.

Given project context and deployment requirements, generate CI/CD pipeline configs. Output ONLY valid JSON:
{
  "files": [
    {"path": ".github/workflows/ci.yml", "content": "full YAML content"}
  ],
  "summary": "Brief description of pipelines generated"
}

Generate pipeline configs that support multi-environment deployment. The environments are specified in the project context (defaults: dev, staging, prod).

Generate:
- .github/workflows/ci.yml — lint, test, build on every PR
- .github/workflows/deploy.yml — deploy to each environment with approval gates
- Makefile — common dev commands (build, test, lint, deploy)
- scripts/deploy.sh — deployment script supporting environment parameter
- docker-compose.yml — local development setup

Rules:
- Support configurable environments (read from context, default: dev/staging/prod)
- Include proper secret references (not hardcoded values)
- Add approval gates for production deployments
- Include rollback mechanisms
- Use matrix builds for multi-version testing where appropriate
- Generate Jenkinsfile alternative if context specifies Jenkins"""


class CicdAgent:
    name = "cicd"
    agent_type = "ci_cd"

    def execute(self, task: str, context: dict) -> AgentResult:
        envs = context.get("environments", ["dev", "staging", "prod"])
        user_msg = (
            f"Project: {context.get('team_description', 'Software project')}\n"
            f"Stack: {context.get('tech_stack', 'Python')}\n"
            f"Environments: {', '.join(envs)}\n"
            f"CI system: {context.get('ci_system', 'GitHub Actions')}\n\n"
            f"Task: {task}"
        )
        raw = invoke(SYSTEM_PROMPT, user_msg)
        try:
            data = json.loads(raw)
        except json.JSONDecodeError:
            return AgentResult(summary=raw, files=[])
        files = [FileChange(path=f["path"], content=f["content"]) for f in data.get("files", [])]
        return AgentResult(summary=data.get("summary", "CI/CD agent completed"), files=files)
