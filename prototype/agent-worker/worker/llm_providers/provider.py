from __future__ import annotations
import json
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import AsyncIterator, Optional


@dataclass
class ProviderConfig:
    model: str
    api_key: Optional[str] = None
    region: Optional[str] = None
    base_url: Optional[str] = None
    extra: dict = field(default_factory=dict)


class LLMProvider(ABC):
    def __init__(self, config: ProviderConfig):
        self.config = config

    @abstractmethod
    async def invoke(
        self,
        system: str,
        messages: list[dict],
        max_tokens: int = 4096,
        temperature: float = 0.7,
        stream: bool = False,
    ) -> str:
        ...

    @abstractmethod
    async def stream_invoke(
        self,
        system: str,
        messages: list[dict],
        max_tokens: int = 4096,
        temperature: float = 0.7,
    ) -> AsyncIterator[str]:
        ...


class BedrockProvider(LLMProvider):
    async def invoke(self, system, messages, max_tokens=4096, temperature=0.7, stream=False):
        import boto3

        client = boto3.client("bedrock-runtime", region_name=self.config.region or "us-east-1")
        body = {
            "anthropic_version": "bedrock-2023-05-31",
            "system": system,
            "messages": messages,
            "max_tokens": max_tokens,
            "temperature": temperature,
        }
        if stream:
            resp = client.invoke_model_with_response_stream(
                modelId=self.config.model, body=json.dumps(body)
            )
            chunks = []
            for event in resp["body"]:
                chunk = json.loads(event["chunk"]["bytes"])
                if chunk.get("type") == "content_block_delta":
                    chunks.append(chunk["delta"].get("text", ""))
            return "".join(chunks)

        resp = client.invoke_model(modelId=self.config.model, body=json.dumps(body))
        result = json.loads(resp["body"].read())
        return result["content"][0]["text"]

    async def stream_invoke(self, system, messages, max_tokens=4096, temperature=0.7):
        import boto3

        client = boto3.client("bedrock-runtime", region_name=self.config.region or "us-east-1")
        body = {
            "anthropic_version": "bedrock-2023-05-31",
            "system": system,
            "messages": messages,
            "max_tokens": max_tokens,
            "temperature": temperature,
        }
        resp = client.invoke_model_with_response_stream(
            modelId=self.config.model, body=json.dumps(body)
        )
        for event in resp["body"]:
            chunk = json.loads(event["chunk"]["bytes"])
            if chunk.get("type") == "content_block_delta":
                yield chunk["delta"].get("text", "")


class OpenAIProvider(LLMProvider):
    async def invoke(self, system, messages, max_tokens=4096, temperature=0.7, stream=False):
        import httpx

        headers = {"Authorization": f"Bearer {self.config.api_key}", "Content-Type": "application/json"}
        url = (self.config.base_url or "https://api.openai.com/v1") + "/chat/completions"
        payload = {
            "model": self.config.model,
            "messages": [{"role": "system", "content": system}] + messages,
            "max_tokens": max_tokens,
            "temperature": temperature,
            "stream": stream,
        }
        if stream:
            text = []
            async with httpx.AsyncClient() as client:
                async with client.stream("POST", url, headers=headers, json=payload, timeout=120) as resp:
                    async for line in resp.aiter_lines():
                        if line.startswith("data: ") and line != "data: [DONE]":
                            data = json.loads(line[6:])
                            delta = data["choices"][0].get("delta", {}).get("content", "")
                            if delta:
                                text.append(delta)
            return "".join(text)

        async with httpx.AsyncClient() as client:
            resp = await client.post(url, headers=headers, json=payload, timeout=120)
            resp.raise_for_status()
            return resp.json()["choices"][0]["message"]["content"]

    async def stream_invoke(self, system, messages, max_tokens=4096, temperature=0.7):
        import httpx

        headers = {"Authorization": f"Bearer {self.config.api_key}", "Content-Type": "application/json"}
        url = (self.config.base_url or "https://api.openai.com/v1") + "/chat/completions"
        payload = {
            "model": self.config.model,
            "messages": [{"role": "system", "content": system}] + messages,
            "max_tokens": max_tokens,
            "temperature": temperature,
            "stream": True,
        }
        async with httpx.AsyncClient() as client:
            async with client.stream("POST", url, headers=headers, json=payload, timeout=120) as resp:
                async for line in resp.aiter_lines():
                    if line.startswith("data: ") and line != "data: [DONE]":
                        data = json.loads(line[6:])
                        delta = data["choices"][0].get("delta", {}).get("content", "")
                        if delta:
                            yield delta


class AnthropicProvider(LLMProvider):
    async def invoke(self, system, messages, max_tokens=4096, temperature=0.7, stream=False):
        import httpx

        headers = {
            "x-api-key": self.config.api_key,
            "anthropic-version": "2023-06-01",
            "Content-Type": "application/json",
        }
        url = (self.config.base_url or "https://api.anthropic.com/v1") + "/messages"
        payload = {
            "model": self.config.model,
            "system": system,
            "messages": messages,
            "max_tokens": max_tokens,
            "temperature": temperature,
            "stream": stream,
        }
        if stream:
            text = []
            async with httpx.AsyncClient() as client:
                async with client.stream("POST", url, headers=headers, json=payload, timeout=120) as resp:
                    async for line in resp.aiter_lines():
                        if line.startswith("data: "):
                            data = json.loads(line[6:])
                            if data.get("type") == "content_block_delta":
                                text.append(data["delta"].get("text", ""))
            return "".join(text)

        async with httpx.AsyncClient() as client:
            resp = await client.post(url, headers=headers, json=payload, timeout=120)
            resp.raise_for_status()
            return resp.json()["content"][0]["text"]

    async def stream_invoke(self, system, messages, max_tokens=4096, temperature=0.7):
        import httpx

        headers = {
            "x-api-key": self.config.api_key,
            "anthropic-version": "2023-06-01",
            "Content-Type": "application/json",
        }
        url = (self.config.base_url or "https://api.anthropic.com/v1") + "/messages"
        payload = {
            "model": self.config.model,
            "system": system,
            "messages": messages,
            "max_tokens": max_tokens,
            "temperature": temperature,
            "stream": True,
        }
        async with httpx.AsyncClient() as client:
            async with client.stream("POST", url, headers=headers, json=payload, timeout=120) as resp:
                async for line in resp.aiter_lines():
                    if line.startswith("data: "):
                        data = json.loads(line[6:])
                        if data.get("type") == "content_block_delta":
                            yield data["delta"].get("text", "")


class OllamaProvider(LLMProvider):
    async def invoke(self, system, messages, max_tokens=4096, temperature=0.7, stream=False):
        import httpx

        url = (self.config.base_url or "http://localhost:11434") + "/api/chat"
        payload = {
            "model": self.config.model,
            "messages": [{"role": "system", "content": system}] + messages,
            "stream": False,
            "options": {"num_predict": max_tokens, "temperature": temperature},
        }
        async with httpx.AsyncClient() as client:
            resp = await client.post(url, json=payload, timeout=300)
            resp.raise_for_status()
            return resp.json()["message"]["content"]

    async def stream_invoke(self, system, messages, max_tokens=4096, temperature=0.7):
        import httpx

        url = (self.config.base_url or "http://localhost:11434") + "/api/chat"
        payload = {
            "model": self.config.model,
            "messages": [{"role": "system", "content": system}] + messages,
            "stream": True,
            "options": {"num_predict": max_tokens, "temperature": temperature},
        }
        async with httpx.AsyncClient() as client:
            async with client.stream("POST", url, json=payload, timeout=300) as resp:
                async for line in resp.aiter_lines():
                    if line:
                        data = json.loads(line)
                        if not data.get("done") and data.get("message", {}).get("content"):
                            yield data["message"]["content"]


_PROVIDERS: dict[str, type[LLMProvider]] = {
    "bedrock": BedrockProvider,
    "openai": OpenAIProvider,
    "anthropic": AnthropicProvider,
    "ollama": OllamaProvider,
}


def get_provider(name: str, config: ProviderConfig) -> LLMProvider:
    cls = _PROVIDERS.get(name.lower())
    if not cls:
        raise ValueError(f"Unknown provider '{name}'. Available: {list(_PROVIDERS.keys())}")
    return cls(config)
