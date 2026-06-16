"""GitHub Actions multi-environment workflow template generator."""


def generate_deploy_workflow(project_name: str, environments: list[dict], languages: list[str]) -> str:
    """Generate .github/workflows/deploy.yml content."""
    lang = languages[0] if languages else "python"
    setup_step = _lang_setup(lang)

    env_jobs = ""
    prev_env = None
    for env in environments:
        name = env["name"]
        needs = f"\n    needs: [deploy-{prev_env}]" if prev_env else "\n    needs: [build]"
        approval = ""
        if env.get("approval_required", False) or name == "prod":
            approval = f"\n    environment: {name}"  # GitHub environment with protection rules

        env_jobs += f"""
  deploy-{name}:{needs}{approval}
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - {setup_step}
      - name: Deploy to {name}
        run: make deploy ENV={name}
        env:
          DEPLOY_ENV: {name}
"""
        prev_env = name

    return f"""name: Deploy {project_name}

on:
  push:
    branches: [main]
  workflow_dispatch:
    inputs:
      environment:
        description: Target environment
        required: false
        type: choice
        options:
{_env_options(environments)}

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - {setup_step}
      - name: Build
        run: make build
      - name: Test
        run: make test
      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: build-artifact
          path: dist/
{env_jobs}"""


def generate_rollback_workflow(project_name: str, environments: list[dict]) -> str:
    """Generate .github/workflows/rollback.yml content."""
    return f"""name: Rollback {project_name}

on:
  workflow_dispatch:
    inputs:
      environment:
        description: Environment to rollback
        required: true
        type: choice
        options:
{_env_options(environments)}
      version:
        description: Version to rollback to
        required: true
        type: string

jobs:
  rollback:
    runs-on: ubuntu-latest
    environment: ${{{{ github.event.inputs.environment }}}}
    steps:
      - uses: actions/checkout@v4
        with:
          ref: ${{{{ github.event.inputs.version }}}}
      - name: Rollback deployment
        run: make deploy ENV=${{{{ github.event.inputs.environment }}}} VERSION=${{{{ github.event.inputs.version }}}}
"""


def _env_options(environments: list[dict]) -> str:
    return "\n".join(f"          - {e['name']}" for e in environments)


def _lang_setup(lang: str) -> str:
    setups = {
        "python": "uses: actions/setup-python@v5\n        with:\n          python-version: '3.12'",
        "go": "uses: actions/setup-go@v5\n        with:\n          go-version: '1.22'",
        "node": "uses: actions/setup-node@v4\n        with:\n          node-version: '20'",
        "typescript": "uses: actions/setup-node@v4\n        with:\n          node-version: '20'",
        "java": "uses: actions/setup-java@v4\n        with:\n          distribution: temurin\n          java-version: '21'",
        "rust": "uses: dtolnay/rust-toolchain@stable",
    }
    return setups.get(lang.lower(), setups["python"])
