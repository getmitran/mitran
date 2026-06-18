"""Tests for LLM provider abstraction layer."""
import json
import pytest
from unittest.mock import patch, MagicMock, AsyncMock

from worker.llm_providers import (
    ProviderConfig,
    LLMProvider,
    BedrockProvider,
    OpenAIProvider,
    AnthropicProvider,
    OllamaProvider,
    get_provider,
)


class TestProviderConfig:
    def test_minimal_config(self):
        cfg = ProviderConfig(model="claude-3-sonnet")
        assert cfg.model == "claude-3-sonnet"
        assert cfg.api_key is None
        assert cfg.region is None

    def test_full_config(self):
        cfg = ProviderConfig(
            model="gpt-4", api_key="sk-test", base_url="http://localhost:8080", extra={"org": "test"}
        )
        assert cfg.api_key == "sk-test"
        assert cfg.extra["org"] == "test"


class TestGetProvider:
    def test_returns_bedrock(self):
        p = get_provider("bedrock", ProviderConfig(model="us.anthropic.claude-sonnet-4-20250514"))
        assert isinstance(p, BedrockProvider)

    def test_returns_openai(self):
        p = get_provider("openai", ProviderConfig(model="gpt-4", api_key="sk-test"))
        assert isinstance(p, OpenAIProvider)

    def test_returns_anthropic(self):
        p = get_provider("anthropic", ProviderConfig(model="claude-3-opus", api_key="sk-ant-test"))
        assert isinstance(p, AnthropicProvider)

    def test_returns_ollama(self):
        p = get_provider("ollama", ProviderConfig(model="llama3"))
        assert isinstance(p, OllamaProvider)

    def test_case_insensitive(self):
        p = get_provider("BEDROCK", ProviderConfig(model="test"))
        assert isinstance(p, BedrockProvider)

    def test_unknown_provider_raises(self):
        with pytest.raises(ValueError, match="Unknown provider"):
            get_provider("nonexistent", ProviderConfig(model="test"))


class TestBedrockProvider:
    @pytest.mark.asyncio
    async def test_invoke(self, mock_bedrock_client):
        cfg = ProviderConfig(model="us.anthropic.claude-sonnet-4-20250514", region="us-east-1")
        provider = BedrockProvider(cfg)

        with patch("boto3.client", return_value=mock_bedrock_client):
            result = await provider.invoke("system prompt", [{"role": "user", "content": "hello"}])

        assert result == "mocked LLM response"
        mock_bedrock_client.invoke_model.assert_called_once()
        call_body = json.loads(mock_bedrock_client.invoke_model.call_args[1]["body"])
        assert call_body["system"] == "system prompt"
        assert call_body["messages"][0]["content"] == "hello"

    @pytest.mark.asyncio
    async def test_invoke_with_stream(self, mock_bedrock_client):
        cfg = ProviderConfig(model="us.anthropic.claude-sonnet-4-20250514", region="us-east-1")
        provider = BedrockProvider(cfg)

        # Mock streaming response
        chunk1 = json.dumps({"type": "content_block_delta", "delta": {"text": "Hello"}}).encode()
        chunk2 = json.dumps({"type": "content_block_delta", "delta": {"text": " world"}}).encode()
        mock_bedrock_client.invoke_model_with_response_stream.return_value = {
            "body": [{"chunk": {"bytes": chunk1}}, {"chunk": {"bytes": chunk2}}]
        }

        with patch("boto3.client", return_value=mock_bedrock_client):
            result = await provider.invoke("sys", [{"role": "user", "content": "hi"}], stream=True)

        assert result == "Hello world"


class TestOpenAIProvider:
    @pytest.mark.asyncio
    async def test_invoke(self):
        cfg = ProviderConfig(model="gpt-4", api_key="sk-test")
        provider = OpenAIProvider(cfg)

        mock_resp = MagicMock()
        mock_resp.status_code = 200
        mock_resp.raise_for_status = MagicMock()
        mock_resp.json.return_value = {
            "choices": [{"message": {"content": "OpenAI response"}}]
        }

        with patch("httpx.AsyncClient") as mock_client_cls:
            mock_client = AsyncMock()
            mock_client.__aenter__ = AsyncMock(return_value=mock_client)
            mock_client.__aexit__ = AsyncMock(return_value=None)
            mock_client.post = AsyncMock(return_value=mock_resp)
            mock_client_cls.return_value = mock_client

            result = await provider.invoke("system", [{"role": "user", "content": "test"}])

        assert result == "OpenAI response"


class TestAnthropicProvider:
    @pytest.mark.asyncio
    async def test_invoke(self):
        cfg = ProviderConfig(model="claude-3-opus", api_key="sk-ant-test")
        provider = AnthropicProvider(cfg)

        mock_resp = MagicMock()
        mock_resp.status_code = 200
        mock_resp.raise_for_status = MagicMock()
        mock_resp.json.return_value = {"content": [{"text": "Anthropic response"}]}

        with patch("httpx.AsyncClient") as mock_client_cls:
            mock_client = AsyncMock()
            mock_client.__aenter__ = AsyncMock(return_value=mock_client)
            mock_client.__aexit__ = AsyncMock(return_value=None)
            mock_client.post = AsyncMock(return_value=mock_resp)
            mock_client_cls.return_value = mock_client

            result = await provider.invoke("system", [{"role": "user", "content": "test"}])

        assert result == "Anthropic response"


class TestOllamaProvider:
    @pytest.mark.asyncio
    async def test_invoke(self):
        cfg = ProviderConfig(model="llama3", base_url="http://localhost:11434")
        provider = OllamaProvider(cfg)

        mock_resp = MagicMock()
        mock_resp.status_code = 200
        mock_resp.raise_for_status = MagicMock()
        mock_resp.json.return_value = {"message": {"content": "Ollama response"}}

        with patch("httpx.AsyncClient") as mock_client_cls:
            mock_client = AsyncMock()
            mock_client.__aenter__ = AsyncMock(return_value=mock_client)
            mock_client.__aexit__ = AsyncMock(return_value=None)
            mock_client.post = AsyncMock(return_value=mock_resp)
            mock_client_cls.return_value = mock_client

            result = await provider.invoke("system", [{"role": "user", "content": "test"}])

        assert result == "Ollama response"


class TestProviderInterface:
    """Verify all providers implement the abstract interface."""

    def test_all_providers_subclass_base(self):
        for cls in [BedrockProvider, OpenAIProvider, AnthropicProvider, OllamaProvider]:
            assert issubclass(cls, LLMProvider)

    def test_all_providers_have_invoke(self):
        for cls in [BedrockProvider, OpenAIProvider, AnthropicProvider, OllamaProvider]:
            assert hasattr(cls, "invoke")
            assert hasattr(cls, "stream_invoke")
