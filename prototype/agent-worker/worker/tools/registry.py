"""Tool registry for Mitran agent worker."""

from dataclasses import dataclass, field
from typing import Any, Callable


class PendingApproval(Exception):
    """Raised when a tool requires interactive approval before execution."""

    def __init__(self, tool_name: str, params: dict):
        self.tool_name = tool_name
        self.params = params
        super().__init__(f"Tool '{tool_name}' requires approval before execution")


@dataclass
class Tool:
    """Definition of an executable tool."""

    name: str
    description: str
    parameters: dict = field(default_factory=dict)
    handler: Callable[..., Any] = field(default=None)
    requires_approval: bool = False


class ToolRegistry:
    """Registry for managing and executing tools."""

    def __init__(self):
        self._tools: dict[str, Tool] = {}

    def register(self, tool: Tool) -> None:
        """Register a tool in the registry."""
        self._tools[tool.name] = tool

    def list_tools(self) -> list[Tool]:
        """List all registered tools."""
        return list(self._tools.values())

    def get(self, name: str) -> Tool | None:
        """Get a tool by name."""
        return self._tools.get(name)

    def execute(self, name: str, params: dict, approval_mode: str = "interactive") -> dict:
        """Execute a tool by name with given parameters.

        Args:
            name: Tool name to execute.
            params: Parameters to pass to the tool handler.
            approval_mode: 'auto' executes immediately, 'interactive' raises
                PendingApproval if the tool requires approval.

        Returns:
            Dict with 'result' key on success, or 'error' key on failure.

        Raises:
            PendingApproval: If approval_mode is 'interactive' and tool requires approval.
            KeyError: If tool is not found.
        """
        tool = self._tools.get(name)
        if not tool:
            raise KeyError(f"Tool '{name}' not found in registry")

        if tool.requires_approval and approval_mode == "interactive":
            raise PendingApproval(tool_name=name, params=params)

        try:
            result = tool.handler(**params)
            return {"result": result}
        except Exception as e:
            return {"error": str(e)}
