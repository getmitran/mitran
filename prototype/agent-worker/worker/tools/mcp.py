"""MCP (Model Context Protocol) client stub for Mitran agent worker."""

import json
import subprocess
from typing import Any


class MCPClient:
    """Client for connecting to MCP servers via stdio transport.

    This is a stub defining the interface. Full implementation will handle
    JSON-RPC communication over stdin/stdout of a subprocess.
    """

    def __init__(self):
        self._process: subprocess.Popen | None = None
        self._connected = False

    def connect(self, command: list[str]) -> None:
        """Connect to an MCP server by launching it as a subprocess.

        Args:
            command: Command and args to start the MCP server process.
                     e.g. ["npx", "@modelcontextprotocol/server-filesystem"]
        """
        self._process = subprocess.Popen(
            command,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
        )
        self._connected = True
        # TODO: Send initialize request and await response

    def list_tools(self) -> list[dict[str, Any]]:
        """List available tools from the connected MCP server.

        Returns:
            List of tool definitions with name, description, and inputSchema.
        """
        if not self._connected:
            raise RuntimeError("Not connected to an MCP server")
        # TODO: Send tools/list JSON-RPC request
        return []

    def call_tool(self, name: str, params: dict) -> dict:
        """Call a tool on the connected MCP server.

        Args:
            name: Tool name to invoke.
            params: Parameters to pass to the tool.

        Returns:
            Tool execution result.
        """
        if not self._connected:
            raise RuntimeError("Not connected to an MCP server")
        # TODO: Send tools/call JSON-RPC request with name and arguments
        return {"result": None, "error": "MCP client not yet implemented"}

    def disconnect(self) -> None:
        """Disconnect from the MCP server."""
        if self._process:
            self._process.terminate()
            self._process.wait(timeout=5)
            self._process = None
        self._connected = False

    @property
    def is_connected(self) -> bool:
        return self._connected
