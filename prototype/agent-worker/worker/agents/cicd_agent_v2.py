"""Enhanced CI/CD Agent v2 — environment-aware pipeline generation.

Reads project environment config from the Mitran engine API and generates
GitHub Actions workflows with environment-specific jobs, Makefile, docker-compose,
and per-environment .env files.
"""
import json
import requests

from worker.agents.dev_agent import AgentResult, FileChange
from worker.agents.pipeline_templates import (
    generate_deploy_workflow,
    generate_docker_compose,
    generate_makefile,
    generate_rollback_workflow,
)
from worker.llm import invoke

ENGINE_URL = "http://localhost:7777"

SYSTEM_PROMPT = """You are Mitran's CI/CD Agent v2. You customize pipeline configurations based on project context.

Given the base templates and project context, produce ONLY a JSON object with customizations:
{
  "deploy_script": "#!/bin/bash\\nfull deploy script content",
  "env_vars": {
    "alpha": {"KEY": "value"},
    "beta": {"KEY": "value"},
    "prod": {"KEY": "value"}
  },
  "extra_files": [
    {"path": "path/to/file", "content": "file content"}
  ],
  "summary": "Brief description"
}

Rules:
- deploy_script should handle ENV variable to target the right environment
- env_vars per environment for .env files (no real secrets, use placeholders)
- extra_files for any project-specific configs (e.g., nginx.conf, Dockerfile)
- Keep it practical and production-ready"""


DEFAULT_ENVS = [
    {"name": "alpha", "auto_deploy": True, "approval_required": False},
    {"name": "beta", "auto_deploy": True, "approval_required": False},
    {"name": "gamma", "auto_deploy": False, "approval_required": True},
    {"name": "prod", "auto_deploy": False, "approval_required": True},
]


class CicdAgentV2:
    name = "cicd_v2"
    agent_type = "ci_cd"

    def execute(self, task: str, context: dict) -> AgentResult:
        project_id = context.get("project_id", "")
        environments = self._fetch_environments(project_id)
        languages = context.get("languages", context.get("tech_stack", "python").split(","))
        project_name = context.get("project_name", context.get("company_description", "project"))

        # Generate base templates
        deploy_yml = generate_deploy_workflow(project_name, environments, languages)
        rollback_yml = generate_rollback_workflow(project_name, environments)
        makefile = generate_makefile(project_name, environments, languages)
        docker_compose = generate_docker_compose(project_name, languages)

        # Call Bedrock for customization (deploy script, env vars, extras)
        customizations = self._get_customizations(task, context, environments, languages)

        # Assemble output files
        files = [
            FileChange(path=".github/workflows/deploy.yml", content=deploy_yml),
            FileChange(path=".github/workflows/rollback.yml", content=rollback_yml),
            FileChange(path="Makefile", content=makefile),
            FileChange(path="docker-compose.yml", content=docker_compose),
            FileChange(path="scripts/deploy.sh", content=customizations.get("deploy_script", _default_deploy_script())),
        ]

        # Per-environment .env files
        env_vars = customizations.get("env_vars", {})
        for env in environments:
            name = env["name"]
            vars_content = _format_env_file(name, env_vars.get(name, {}))
            files.append(FileChange(path=f"deploy/envs/{name}.env", content=vars_content))

        # Extra files from LLM customization
        for extra in customizations.get("extra_files", []):
            files.append(FileChange(path=extra["path"], content=extra["content"]))

        summary = customizations.get("summary", f"Generated multi-env CI/CD pipeline with {len(environments)} environments")
        return AgentResult(summary=summary, files=files)

    def _fetch_environments(self, project_id: str) -> list[dict]:
        """Fetch environment config from engine API."""
        if not project_id:
            return DEFAULT_ENVS
        try:
            resp = requests.get(f"{ENGINE_URL}/api/v1/projects/{project_id}/environments", timeout=5)
            if resp.status_code == 200:
                envs = resp.json()
                if envs:
                    return envs
        except Exception:
            pass
        return DEFAULT_ENVS

    def _get_customizations(self, task: str, context: dict, environments: list[dict], languages: list[str]) -> dict:
        """Call Bedrock Claude to customize templates for the project."""
        user_msg = (
            f"Project: {context.get('company_description', 'Software project')}\n"
            f"Languages: {', '.join(languages)}\n"
            f"Environments: {json.dumps([e['name'] for e in environments])}\n"
            f"Team size: {context.get('team_size', 5)}\n"
            f"Task: {task}\n\n"
            f"Generate customizations for this project's CI/CD pipeline."
        )
        try:
            raw = invoke(SYSTEM_PROMPT, user_msg)
            return json.loads(raw)
        except (json.JSONDecodeError, RuntimeError):
            return {"deploy_script": _default_deploy_script(), "env_vars": {}, "extra_files": [], "summary": "CI/CD pipeline generated with defaults"}


def _default_deploy_script() -> str:
    return """#!/bin/bash
set -euo pipefail

ENV="${ENV:-alpha}"
VERSION="${VERSION:-latest}"

echo "=== Deploying to $ENV (version: $VERSION) ==="

# Load environment config
if [ -f "deploy/envs/${ENV}.env" ]; then
    export $(grep -v '^#' "deploy/envs/${ENV}.env" | xargs)
fi

# Build artifact
echo "Building..."
make build

# Deploy (customize per your infrastructure)
echo "Deploying artifact to $ENV..."
# docker push / kubectl apply / aws deploy / etc.

echo "=== Deployment to $ENV complete ==="
"""


def _format_env_file(env_name: str, vars: dict) -> str:
    lines = [f"# Environment: {env_name}", f"ENV={env_name}", f"LOG_LEVEL={'debug' if env_name in ('alpha', 'local') else 'info'}"]
    for k, v in vars.items():
        lines.append(f"{k}={v}")
    return "\n".join(lines) + "\n"
