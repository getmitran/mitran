import asyncio
import json
import tempfile
from pathlib import Path
from unittest.mock import AsyncMock, patch, MagicMock

import pytest

from worker.mcp.registry import MCPRegistry
from worker.mcp.types import ToolInfo


@pytest.fixture
def mcp_config(tmp_path):
    config = {
        "mcpServers": {
            "filesystem": {
                "command": "node",
                "args": ["server.js"],
                "env": {"HOME": "/tmp"}
            },
            "github": {
                "command": "gh-mcp",
                "args": ["--token", "xxx"]
            }
        }
    }
    config_path = tmp_path / "mcp-config.json"
    config_path.write_text(json.dumps(config))
    return str(config_path)


@pytest.fixture
def empty_config(tmp_path):
    config_path = tmp_path / "mcp-config.json"
    config_path.write_text(json.dumps({"mcpServers": {}}))
    return str(config_path)


def test_load_config(mcp_config):
    registry = MCPRegistry(mcp_config)
    configs = registry._load_config()
    assert len(configs) == 2
    assert configs[0].name == "filesystem"
    assert configs[0].command == "node"
    assert configs[0].args == ["server.js"]
    assert configs[1].name == "github"


def test_load_missing_config(tmp_path):
    registry = MCPRegistry(str(tmp_path / "nonexistent.json"))
    configs = registry._load_config()
    assert configs == []


@pytest.mark.asyncio
async def test_start_all_and_list_tools(mcp_config):
    mock_tools = [
        ToolInfo(name="read_file", description="Read a file", input_schema={"type": "object"}),
        ToolInfo(name="write_file", description="Write a file", input_schema={"type": "object"}),
    ]

    with patch("worker.mcp.registry.MCPClient") as MockClient:
        instance = AsyncMock()
        instance.initialize = AsyncMock(return_value={})
        instance.list_tools = AsyncMock(return_value=mock_tools)
        instance.close = AsyncMock()
        MockClient.return_value = instance

        registry = MCPRegistry(mcp_config)
        await registry.start_all()
        assert len(registry._clients) == 2

        tools = await registry.list_all_tools()
        # Each server returns 2 tools, but same names so tool_map has 2 entries
        assert len(tools) == 4  # 2 servers * 2 tools each

        await registry.stop_all()
        assert len(registry._clients) == 0


@pytest.mark.asyncio
async def test_call_tool_auto(mcp_config):
    with patch("worker.mcp.registry.MCPClient") as MockClient:
        instance = AsyncMock()
        instance.initialize = AsyncMock(return_value={})
        instance.list_tools = AsyncMock(return_value=[
            ToolInfo(name="read_file", description="Read", input_schema={})
        ])
        instance.call_tool = AsyncMock(return_value="file content")
        instance.close = AsyncMock()
        MockClient.return_value = instance

        registry = MCPRegistry(mcp_config)
        await registry.start_all()
        await registry.list_all_tools()

        result = await registry.call_tool_auto("read_file", {"path": "/tmp/x"})
        assert result == "file content"


@pytest.mark.asyncio
async def test_call_tool_unknown_server(mcp_config):
    registry = MCPRegistry(mcp_config)
    with pytest.raises(ValueError, match="Unknown MCP server"):
        await registry.call_tool("nonexistent", "tool", {})


@pytest.mark.asyncio
async def test_call_tool_auto_unknown_tool(empty_config):
    registry = MCPRegistry(empty_config)
    await registry.start_all()
    with pytest.raises(ValueError, match="not found"):
        await registry.call_tool_auto("unknown_tool", {})
