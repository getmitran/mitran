"""Mitran Agent SDK — minimal framework for defining custom agents and tools."""

from __future__ import annotations

import functools
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import Any, Callable, Dict, List, Optional

# ---------- Tool Decorator ----------

_tool_registry: Dict[str, "ToolDef"] = {}


@dataclass
class ToolDef:
    name: str
    description: str
    func: Callable
    parameters: Dict[str, Any] = field(default_factory=dict)


def tool(name: Optional[str] = None, description: str = ""):
    """Decorator to register a function as a callable tool.

    Usage:
        @tool(description="Search codebase for a pattern")
        def search_code(query: str, path: str = ".") -> str:
            ...
    """

    def decorator(func: Callable) -> Callable:
        tool_name = name or func.__name__
        _tool_registry[tool_name] = ToolDef(
            name=tool_name,
            description=description or func.__doc__ or "",
            func=func,
        )

        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            return func(*args, **kwargs)

        wrapper._tool_def = _tool_registry[tool_name]
        return wrapper

    return decorator


def get_tools() -> Dict[str, ToolDef]:
    """Return all registered tools."""
    return dict(_tool_registry)


# ---------- Agent Base Class ----------


class Agent(ABC):
    """Base class for Mitran agents. Subclass and implement handle()."""

    name: str = "unnamed"
    description: str = ""
    model: str = ""

    @abstractmethod
    def handle(self, message: str, context: Optional[Dict[str, Any]] = None) -> str:
        """Process a user message and return a response."""
        ...

    def get_tools(self) -> List[ToolDef]:
        """Override to provide agent-specific tools."""
        return []


# ---------- Agent Registry ----------

_agent_registry: Dict[str, Agent] = {}


def register_agent(name: str, agent: Agent) -> None:
    """Register an agent instance by name."""
    agent.name = name
    _agent_registry[name] = agent


def get_agent(name: str) -> Optional[Agent]:
    """Retrieve a registered agent by name."""
    return _agent_registry.get(name)


def list_agents() -> List[str]:
    """List all registered agent names."""
    return list(_agent_registry.keys())


__all__ = [
    "tool",
    "ToolDef",
    "get_tools",
    "Agent",
    "register_agent",
    "get_agent",
    "list_agents",
]
