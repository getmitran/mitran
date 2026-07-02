"""Built-in tools for Mitran agent worker."""

import subprocess
import urllib.request
import json

from worker.tools.registry import Tool


def read_file(path: str) -> str:
    """Read file contents."""
    with open(path, "r") as f:
        return f.read()


def write_file(path: str, content: str) -> bool:
    """Write content to a file."""
    with open(path, "w") as f:
        f.write(content)
    return True


def run_command(cmd: str, cwd: str = ".") -> dict:
    """Run a shell command and return stdout, stderr, exit_code."""
    result = subprocess.run(
        cmd, shell=True, cwd=cwd, capture_output=True, text=True, timeout=300
    )
    return {
        "stdout": result.stdout,
        "stderr": result.stderr,
        "exit_code": result.returncode,
    }


def http_request(method: str, url: str, body: str = "") -> dict:
    """Make an HTTP request."""
    data = body.encode("utf-8") if body else None
    req = urllib.request.Request(url, data=data, method=method.upper())
    req.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            return {
                "status": resp.status,
                "headers": dict(resp.headers),
                "body": resp.read().decode("utf-8"),
            }
    except urllib.error.HTTPError as e:
        return {
            "status": e.code,
            "headers": dict(e.headers),
            "body": e.read().decode("utf-8"),
        }


def git_diff(cwd: str = ".") -> str:
    """Get git diff output."""
    result = subprocess.run(
        ["git", "diff"], cwd=cwd, capture_output=True, text=True, timeout=30
    )
    return result.stdout


def git_commit(cwd: str = ".", message: str = "") -> bool:
    """Stage all changes and commit."""
    subprocess.run(["git", "add", "-A"], cwd=cwd, check=True, timeout=30)
    result = subprocess.run(
        ["git", "commit", "-m", message], cwd=cwd, capture_output=True, text=True, timeout=30
    )
    return result.returncode == 0


BUILTIN_TOOLS = [
    Tool(
        name="read_file",
        description="Read the contents of a file at the given path.",
        parameters={"path": {"type": "string", "description": "File path to read"}},
        handler=read_file,
        requires_approval=False,
    ),
    Tool(
        name="write_file",
        description="Write content to a file, creating or overwriting it.",
        parameters={
            "path": {"type": "string", "description": "File path to write"},
            "content": {"type": "string", "description": "Content to write"},
        },
        handler=write_file,
        requires_approval=True,
    ),
    Tool(
        name="run_command",
        description="Execute a shell command and return stdout, stderr, and exit code.",
        parameters={
            "cmd": {"type": "string", "description": "Shell command to execute"},
            "cwd": {"type": "string", "description": "Working directory", "default": "."},
        },
        handler=run_command,
        requires_approval=True,
    ),
    Tool(
        name="http_request",
        description="Make an HTTP request to a URL.",
        parameters={
            "method": {"type": "string", "description": "HTTP method (GET, POST, etc.)"},
            "url": {"type": "string", "description": "Request URL"},
            "body": {"type": "string", "description": "Request body", "default": ""},
        },
        handler=http_request,
        requires_approval=False,
    ),
    Tool(
        name="git_diff",
        description="Get the current git diff in a directory.",
        parameters={
            "cwd": {"type": "string", "description": "Repository directory", "default": "."},
        },
        handler=git_diff,
        requires_approval=False,
    ),
    Tool(
        name="git_commit",
        description="Stage all changes and create a git commit.",
        parameters={
            "cwd": {"type": "string", "description": "Repository directory", "default": "."},
            "message": {"type": "string", "description": "Commit message"},
        },
        handler=git_commit,
        requires_approval=True,
    ),
]
