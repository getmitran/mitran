import os


ENGINE_URL = os.environ.get("MITRAN_ENGINE_URL", "http://localhost:7777")
WORKER_PORT = int(os.environ.get("MITRAN_WORKER_PORT", "8888"))
AWS_REGION = os.environ.get("AWS_REGION", "us-east-1")
MODEL_ID = os.environ.get("MITRAN_MODEL_ID", "us.anthropic.claude-sonnet-4-20250514")
WORKER_HOST = os.environ.get("MITRAN_WORKER_HOST", "0.0.0.0")
