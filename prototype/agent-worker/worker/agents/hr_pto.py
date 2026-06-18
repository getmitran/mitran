"""HR Agent - PTO Workflow"""
import json
import os
import uuid
from datetime import datetime

DATA_DIR = os.path.join(os.path.dirname(__file__), "..", "..", "data", "pto")


def _ensure_dir():
    os.makedirs(DATA_DIR, exist_ok=True)


def _load():
    _ensure_dir()
    path = os.path.join(DATA_DIR, "requests.json")
    if os.path.exists(path):
        with open(path) as f:
            return json.load(f)
    return {"requests": [], "balances": {}}


def _save(data):
    _ensure_dir()
    with open(os.path.join(DATA_DIR, "requests.json"), "w") as f:
        json.dump(data, f, indent=2)


class PTOWorkflow:
    def submit_request(self, employee: str, dates: list[str], pto_type: str) -> dict:
        data = _load()
        request = {
            "id": str(uuid.uuid4())[:8],
            "employee": employee,
            "dates": dates,
            "type": pto_type,
            "status": "pending",
            "submitted_at": datetime.utcnow().isoformat(),
        }
        data["requests"].append(request)
        _save(data)
        return request

    def check_balance(self, employee: str) -> dict:
        data = _load()
        default = {"vacation": 20, "sick": 10, "personal": 5}
        balance = data["balances"].get(employee, default)
        used = {}
        for r in data["requests"]:
            if r["employee"] == employee and r["status"] == "approved":
                t = r["type"]
                used[t] = used.get(t, 0) + len(r["dates"])
        return {
            "employee": employee,
            "allocated": balance,
            "used": used,
            "remaining": {k: balance.get(k, 0) - used.get(k, 0) for k in balance},
        }

    def approve(self, request_id: str, manager: str) -> dict:
        data = _load()
        for r in data["requests"]:
            if r["id"] == request_id:
                r["status"] = "approved"
                r["approved_by"] = manager
                r["decided_at"] = datetime.utcnow().isoformat()
                _save(data)
                return r
        return {"error": f"Request {request_id} not found"}

    def deny(self, request_id: str, reason: str) -> dict:
        data = _load()
        for r in data["requests"]:
            if r["id"] == request_id:
                r["status"] = "denied"
                r["deny_reason"] = reason
                r["decided_at"] = datetime.utcnow().isoformat()
                _save(data)
                return r
        return {"error": f"Request {request_id} not found"}

    def get_calendar(self, team: list[str]) -> dict:
        data = _load()
        calendar = {}
        for r in data["requests"]:
            if r["employee"] in team and r["status"] == "approved":
                for d in r["dates"]:
                    calendar.setdefault(d, []).append(r["employee"])
        return {"team": team, "calendar": calendar}
