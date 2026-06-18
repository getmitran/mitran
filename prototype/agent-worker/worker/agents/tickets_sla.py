"""Tickets Agent - SLA Enforcement module."""
from dataclasses import dataclass
from datetime import datetime, timedelta
from typing import Optional


@dataclass
class SLAPolicy:
    priority: str
    max_hours: int


DEFAULT_SLAS = [
    SLAPolicy("critical", 1),
    SLAPolicy("high", 4),
    SLAPolicy("medium", 24),
    SLAPolicy("low", 72),
]


class SLATracker:
    def __init__(self, policies: list[SLAPolicy] = None):
        self.policies = {p.priority: p for p in (policies or DEFAULT_SLAS)}

    def get_deadline(self, priority: str, created_at: datetime) -> datetime:
        policy = self.policies.get(priority, SLAPolicy("medium", 24))
        return created_at + timedelta(hours=policy.max_hours)

    def is_breached(self, priority: str, created_at: datetime, resolved_at: Optional[datetime] = None) -> bool:
        deadline = self.get_deadline(priority, created_at)
        check_time = resolved_at or datetime.utcnow()
        return check_time > deadline

    def check_breaches(self, tickets: list[dict]) -> list[dict]:
        breached = []
        now = datetime.utcnow()
        for t in tickets:
            if t.get("status") in ("done", "closed"):
                continue
            created = datetime.fromisoformat(t["created_at"]) if isinstance(t.get("created_at"), str) else t.get("created_at", now)
            priority = t.get("priority", "medium")
            if self.is_breached(priority, created):
                deadline = self.get_deadline(priority, created)
                t["sla_breached"] = True
                t["sla_overdue_hours"] = round((now - deadline).total_seconds() / 3600, 1)
                breached.append(t)
        return breached

    def auto_escalate(self, tickets: list[dict], notify_fn=None) -> list[dict]:
        breached = self.check_breaches(tickets)
        if notify_fn and breached:
            for t in breached:
                notify_fn(f"SLA BREACH: [{t.get('priority','?').upper()}] {t.get('title','')} - overdue by {t.get('sla_overdue_hours',0)}h")
        return breached

    def get_sla_report(self, tickets: list[dict]) -> dict:
        total = len(tickets)
        breached = len(self.check_breaches(tickets))
        compliant = total - breached
        return {
            "total": total,
            "compliant": compliant,
            "breached": breached,
            "compliance_rate": round(compliant / total * 100, 1) if total > 0 else 100.0,
            "by_priority": {p: {"max_hours": pol.max_hours} for p, pol in self.policies.items()},
        }
