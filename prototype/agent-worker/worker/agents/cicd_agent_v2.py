"""CI/CD Agent v2 - Triggers real GitHub Actions workflows."""
import logging
from worker.agents.dev_agent import BaseAgent, AgentResult, FileChange
from worker.integrations.github import (
    create_workflow_file, trigger_workflow, get_workflow_runs, GitHubAPIError
)
from worker.llm import invoke as llm_invoke

log = logging.getLogger(__name__)

DEFAULT_WORKFLOW = '''name: Mitran CI
on:
  workflow_dispatch:
    inputs:
      task_id:
        description: Mitran task ID
        required: true
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run task
        run: echo "Running Mitran task ${{ github.event.inputs.task_id }}"
'''

class CicdAgentV2(BaseAgent):
    agent_type = 'cicd'
    system_prompt = """You are a CI/CD specialist. You create and trigger GitHub Actions workflows.
    When given a deployment or CI task:
    1. Determine the workflow needed
    2. Create or update the workflow YAML
    3. Trigger the workflow
    4. Report status"""

    def execute(self, task_description: str, context: dict) -> AgentResult:
        github_org = context.get('github_org', '')
        github_repo = context.get('github_repo', '')
        task_id = context.get('task_id', 'unknown')

        if not github_org or not github_repo:
            return AgentResult(
                summary='## CI/CD Agent\n\n⚠️ GitHub not configured. Set github_org and github_repo in project settings.',
                files=[]
            )

        # Generate workflow via LLM
        prompt = f'Generate a GitHub Actions workflow YAML for: {task_description}\nOutput ONLY the YAML, no markdown.'
        workflow_yaml = self._call_llm(prompt)

        # Clean LLM output
        if '```' in workflow_yaml:
            lines = workflow_yaml.split('\n')
            workflow_yaml = '\n'.join(l for l in lines if not l.strip().startswith('```'))

        if not workflow_yaml.strip().startswith('name:'):
            workflow_yaml = DEFAULT_WORKFLOW

        workflow_path = f'.github/workflows/mitran-{task_id[:8]}.yml'

        try:
            create_workflow_file(github_org, github_repo, workflow_yaml, workflow_path)
            log.info(f'Created workflow: {workflow_path}')

            # Trigger it
            trigger_workflow(github_org, github_repo, workflow_path.split('/')[-1], {'task_id': task_id})
            status = '✅ Workflow created and triggered'
        except GitHubAPIError as e:
            log.warning(f'CI/CD action failed: {e}')
            status = f'⚠️ GitHub API error: {e}'

        return AgentResult(
            summary=f'## CI/CD Agent\n\n{status}\n\n**Workflow:** `{workflow_path}`\n**Repo:** {github_org}/{github_repo}',
            files=[FileChange(path=workflow_path, action='create', content=workflow_yaml)]
        )
