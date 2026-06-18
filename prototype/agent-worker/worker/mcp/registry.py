import json
import logging
from pathlib import Path

from .client import MCPClient
from .types import MCPServerConfig, ToolInfo

logger = logging.getLogger(__name__)


class MCPRegistry:
    def __init__(self, config_path: str):
        self.config_path = Path(config_path)
        self._clients: dict[str, MCPClient] = {}
        self._tool_map: dict[str, str] = {}  # tool_name -> server_name

    def _load_config(self) -> list[MCPServerConfig]:
        if not self.config_path.exists():
            logger.warning(f"MCP config not found: {self.config_path}")
            return []
        data = json.loads(self.config_path.read_text())
        servers = data.get("mcpServers", {})
        return [
            MCPServerConfig(
                name=name,
                command=cfg["command"],
                args=cfg.get("args", []),
                env=cfg.get("env", {}),
            )
            for name, cfg in servers.items()
        ]

    async def start_all(self):
        configs = self._load_config()
        for cfg in configs:
            client = MCPClient(name=cfg.name, command=cfg.command, args=cfg.args, env=cfg.env)
            try:
                await client.initialize()
                self._clients[cfg.name] = client
                logger.info(f"MCP server '{cfg.name}' connected")
            except Exception as e:
                logger.error(f"Failed to start MCP server '{cfg.name}': {e}")

    async def list_all_tools(self) -> list[ToolInfo]:
        all_tools: list[ToolInfo] = []
        self._tool_map.clear()
        for name, client in self._clients.items():
            try:
                tools = await client.list_tools()
                for t in tools:
                    self._tool_map[t.name] = name
                all_tools.extend(tools)
            except Exception as e:
                logger.error(f"Failed to list tools from '{name}': {e}")
        return all_tools

    async def call_tool(self, server: str, tool: str, args: dict) -> str:
        client = self._clients.get(server)
        if not client:
            raise ValueError(f"Unknown MCP server: {server}")
        return await client.call_tool(tool, args)

    async def call_tool_auto(self, tool: str, args: dict) -> str:
        """Route tool call to the correct server automatically."""
        server = self._tool_map.get(tool)
        if not server:
            raise ValueError(f"Tool '{tool}' not found in any server")
        return await self.call_tool(server, tool, args)

    async def stop_all(self):
        for name, client in self._clients.items():
            try:
                await client.close()
            except Exception as e:
                logger.error(f"Error closing '{name}': {e}")
        self._clients.clear()
        self._tool_map.clear()
