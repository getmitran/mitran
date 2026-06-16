import subprocess
from pathlib import Path


def read_file(path: str) -> str:
    return Path(path).read_text()


def write_file(path: str, content: str) -> None:
    p = Path(path)
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(content)


def search_files(pattern: str, root: str = ".") -> list[str]:
    return [str(p) for p in Path(root).rglob(pattern)]


def run_command(cmd: str, cwd: str | None = None) -> tuple[str, int]:
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True, cwd=cwd)
    return result.stdout + result.stderr, result.returncode
