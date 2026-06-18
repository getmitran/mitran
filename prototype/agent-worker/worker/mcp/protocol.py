import json
from typing import Any

_request_id = 0


def next_id() -> int:
    global _request_id
    _request_id += 1
    return _request_id


def create_request(method: str, params: dict | None = None) -> str:
    msg = {"jsonrpc": "2.0", "id": next_id(), "method": method}
    if params is not None:
        msg["params"] = params
    return json.dumps(msg)


def parse_response(data: str) -> dict:
    msg = json.loads(data)
    if "error" in msg:
        return {"error": msg["error"]}
    return {"result": msg.get("result", {})}
