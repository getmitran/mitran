"""Ollama local LLM provider using /api/generate endpoint."""

import os
import requests

OLLAMA_URL = os.environ.get("OLLAMA_URL", "http://localhost:11434")


def invoke(prompt: str, system: str = None, max_tokens: int = 1024, model: str = "llama3") -> str:
    """Invoke Ollama's generate endpoint. Returns error string if unavailable."""
    payload = {"model": model, "prompt": prompt, "stream": False, "options": {"num_predict": max_tokens}}
    if system:
        payload["system"] = system
    try:
        resp = requests.post(f"{OLLAMA_URL}/api/generate", json=payload, timeout=120)
        resp.raise_for_status()
        return resp.json().get("response", "")
    except requests.ConnectionError:
        return f"[error] Ollama not reachable at {OLLAMA_URL}"
    except Exception as e:
        return f"[error] Ollama request failed: {e}"
