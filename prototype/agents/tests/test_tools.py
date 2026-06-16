import tempfile
from pathlib import Path

from mitran_sdk.tools import read_file, write_file, search_files, run_command


def test_write_and_read(tmp_path):
    p = str(tmp_path / "test.txt")
    write_file(p, "hello")
    assert read_file(p) == "hello"


def test_write_creates_dirs(tmp_path):
    p = str(tmp_path / "a" / "b" / "c.txt")
    write_file(p, "nested")
    assert read_file(p) == "nested"


def test_search_files(tmp_path):
    (tmp_path / "foo.py").write_text("x")
    (tmp_path / "bar.py").write_text("y")
    (tmp_path / "baz.txt").write_text("z")
    results = search_files("*.py", str(tmp_path))
    assert len(results) == 2


def test_run_command():
    output, code = run_command("echo hello")
    assert code == 0
    assert "hello" in output
