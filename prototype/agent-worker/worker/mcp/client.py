import asyncio
import json
import logging
import os
from typing import Any

from .protocol import create_request, parse_response
from .types import MCPResponse, ToolInfo

logger = logging.getLogger(__name__)


class MCPClient:
    def __init__(self, name: str, command: str, args: list[str] = None, env: dict[str, str] = None):
        self.name = name
        self.command = command
        self.args = args or []
        self.env = env
        self._process: asyncio.subprocess.Process | None = None
        self._lock = asyncio.Lock()
        self._capabilities: dict = {}

    async def _start_process(self):
        proc_env = os.environ.copy()
        if self.env:
            proc_env.update(self.env)
        self._process = await asyncio.create_subprocess_exec(
            self.command, *self.args,
            stdin=asyncio.subprocess.PIPE,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            env=proc_env,
        )

    async def _send_request(self, method: str, params: dict | None = None) -> dict:
        if not self._process or self._process.returncode is not None:
            await self._start_process()
        msg = create_request(method, params)
        self._process.stdin.write((msg + "\n").encode())
        await self._process.stdin.drain()
        line = await self._process.stdout.readline()
        if not line:
            raise ConnectionError(f"MCP server '{self.name}' closed stdout")
        return parse_response(line.decode().strip())

    async def _send_with_retry(self, method: str, params: dict | None = None, retries: int = 1) -> dict:
        for attempt in range(retries + 1):
            try:
                return await self._send_request(method, params)
            except (ConnectionError, OSError) as e:
                if attempt == retries:
                    raise
                logger.warning(f"MCP '{self.name}' connection failed, retrying: {e}")
                await self._start_process()

    async def initialize(self) -> dict:
        async with self._lock:
            await self._start_process()
            resp = await self._send_request("initialize", {
                "protocolVersion": "2024-11-05",
                "capabilities": {},
                "clientInfo": {"name": "mitran", "version": "0.1.0"},
            })
            if "error" in resp:
                raise RuntimeError(f"Init failed: {resp['error']}")
            self._capabilities = resp["result"].get("capabilities", {})
            # Send initialized notification (no response expected)
            notif = json.dumps({"jsonrpc": "2.0", "method": "notifications/initialized"})
            self._process.stdin.write((notif + "\n").encode())
            await self._process.stdin.drain()
            return resp["result"]

    async def list_tools(self) -> list[ToolInfo]:
        async with self._lock:
            resp = await self._send_with_retry("tools/list")
        if "error" in resp:
            raise RuntimeError(f"list_tools failed: {resp['error']}")
        tools = resp["result"].get("tools", [])
        return [
            ToolInfo(
                name=t["name"],
                description=t.get("description", ""),
                input_schema=t.get("inputSchema", {}),
            )
            for t in tools
        ]

    async def call_tool(self, tool_name: str, arguments: dict) -> str:
        async with self._lock:
            resp = await self._send_with_retry("tools/call", {
                "name": tool_name,
                "arguments": arguments,
            })
        if "error" in resp:
            raise RuntimeError(f"call_tool failed: {resp['error']}")
        content = resp["result"].get("content", [])
        parts = [c.get("text", "") for c in content if c.get("type") == "text"]
        return "\n".join(parts)

    async def close(self):
        if self._process and self._process.returncode is None:
            self._process.stdin.close()
            try:
                await asyncio.wait_for(self._process.wait(), timeout=5)
            except asyncio.TimeoutError:
                self._process.kill()
        self._process = None
