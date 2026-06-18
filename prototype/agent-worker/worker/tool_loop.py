import logging
from dataclasses import dataclass, field
from typing import Callable

from worker.mcp import MCPRegistry
from worker.llm import invoke as llm_invoke

log = logging.getLogger(__name__)

@dataclass
class ToolCallResult:
    tool_name: str
    server: str
    arguments: dict
    result: str
    success: bool

@dataclass
class LoopConfig:
    max_iterations: int = 20
    approval_mode: str = 'reads'  # interactive, reads, yolo
    timeout_seconds: int = 300

@dataclass
class LoopResult:
    summary: str
    tool_calls: list[ToolCallResult] = field(default_factory=list)
    iterations: int = 0
    final_response: str = ''

def should_auto_approve(tool_name: str, server: str, mode: str) -> bool:
    if mode == 'yolo': return True
    if mode == 'reads':
        read_tools = {'read_file', 'search_files', 'list_tools', 'get', 'list', 'describe', 'search'}
        return any(t in tool_name.lower() for t in read_tools)
    return False  # interactive = never auto

async def run_tool_loop(
    system_prompt: str,
    user_message: str,
    mcp_registry: MCPRegistry,
    config: LoopConfig = LoopConfig(),
    on_approval_needed: Callable = None,
) -> LoopResult:
    messages = [
        {'role': 'system', 'content': system_prompt},
        {'role': 'user', 'content': user_message},
    ]
    
    tools = await mcp_registry.list_all_tools()
    tool_descriptions = [{'name': t.name, 'description': t.description, 'input_schema': t.input_schema} for t in tools]
    
    result = LoopResult(summary='', iterations=0)
    
    for i in range(config.max_iterations):
        result.iterations = i + 1
        
        response = llm_invoke(
            system_prompt=system_prompt,
            user_message=_format_messages(messages[1:]),
            tools=tool_descriptions if tool_descriptions else None,
        )
        
        tool_use = _extract_tool_use(response)
        
        if not tool_use:
            result.final_response = response
            result.summary = response[:500]
            break
        
        tool_name = tool_use['name']
        tool_args = tool_use.get('arguments', {})
        server_name = _find_server_for_tool(tool_name, tools)
        
        approved = should_auto_approve(tool_name, server_name, config.approval_mode)
        if not approved and on_approval_needed:
            approved = await on_approval_needed(tool_name, server_name, tool_args)
        
        if not approved:
            messages.append({'role': 'assistant', 'content': f'Tool {tool_name} was denied by approval policy.'})
            continue
        
        try:
            tool_result = await mcp_registry.call_tool(server_name, tool_name, tool_args)
            success = True
        except Exception as e:
            tool_result = f'Error: {e}'
            success = False
        
        result.tool_calls.append(ToolCallResult(
            tool_name=tool_name, server=server_name,
            arguments=tool_args, result=str(tool_result)[:1000], success=success
        ))
        
        messages.append({'role': 'assistant', 'content': f'Called {tool_name}'})
        messages.append({'role': 'user', 'content': f'Tool result: {tool_result}'})
    
    return result

def _format_messages(messages: list[dict]) -> str:
    return '\n'.join(f"{m['role']}: {m['content']}" for m in messages[-10:])

def _extract_tool_use(response: str) -> dict | None:
    import json, re
    match = re.search(r'\{"name":\s*"(\w+)".*?"arguments":\s*(\{.*?\})\}', response, re.DOTALL)
    if match:
        try:
            return {'name': match.group(1), 'arguments': json.loads(match.group(2))}
        except Exception:
            pass
    return None

def _find_server_for_tool(tool_name: str, tools) -> str:
    for t in tools:
        if t.name == tool_name:
            return getattr(t, 'server', 'unknown')
    return 'unknown'
