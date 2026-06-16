"""Docker Compose template generator for local development."""


def generate_docker_compose(project_name: str, languages: list[str]) -> str:
    """Generate docker-compose.yml for local dev."""
    lang = languages[0].lower() if languages else "python"
    image, port = _lang_defaults(lang)

    return f"""version: "3.9"

services:
  app:
    build: .
    ports:
      - "{port}:{port}"
    environment:
      - ENV=local
      - LOG_LEVEL=debug
    volumes:
      - .:/app
    env_file:
      - deploy/envs/local.env

  db:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: {project_name.lower().replace(" ", "_")}
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: dev
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  pgdata:
"""


def _lang_defaults(lang: str) -> tuple[str, str]:
    defaults = {
        "python": ("python:3.12-slim", "8000"),
        "go": ("golang:1.22-alpine", "8080"),
        "node": ("node:20-alpine", "3000"),
        "typescript": ("node:20-alpine", "3000"),
        "java": ("eclipse-temurin:21-jre", "8080"),
        "rust": ("rust:1.77-slim", "8080"),
    }
    return defaults.get(lang, defaults["python"])
