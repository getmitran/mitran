"""Makefile template generator with per-environment deploy targets."""


def generate_makefile(project_name: str, environments: list[dict], languages: list[str]) -> str:
    """Generate Makefile content with per-env targets."""
    lang = languages[0].lower() if languages else "python"
    build_cmd, test_cmd, lint_cmd = _lang_commands(lang)

    env_targets = ""
    for env in environments:
        name = env["name"]
        env_targets += f"""
.PHONY: deploy-{name}
deploy-{name}: build ## Deploy to {name}
\t@echo "Deploying to {name}..."
\tENV={name} ./scripts/deploy.sh
"""

    return f""".DEFAULT_GOAL := help

.PHONY: build test lint deploy clean help

build: ## Build the project
\t{build_cmd}

test: ## Run tests
\t{test_cmd}

lint: ## Run linters
\t{lint_cmd}

deploy: ## Deploy to ENV (usage: make deploy ENV=alpha)
\t@if [ -z "$$ENV" ]; then echo "Usage: make deploy ENV=<environment>"; exit 1; fi
\t@echo "Deploying to $$ENV..."
\t./scripts/deploy.sh

clean: ## Remove build artifacts
\trm -rf dist/ build/ *.egg-info
{env_targets}
help: ## Show this help
\t@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {{FS = ":.*?## "}}; {{printf "\\033[36m%-20s\\033[0m %s\\n", $$1, $$2}}'
"""


def _lang_commands(lang: str) -> tuple[str, str, str]:
    commands = {
        "python": ("python -m build", "pytest", "ruff check ."),
        "go": ("go build ./...", "go test ./...", "golangci-lint run"),
        "node": ("npm run build", "npm test", "npm run lint"),
        "typescript": ("npm run build", "npm test", "npm run lint"),
        "java": ("./gradlew build", "./gradlew test", "./gradlew check"),
        "rust": ("cargo build --release", "cargo test", "cargo clippy"),
    }
    return commands.get(lang, commands["python"])
