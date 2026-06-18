"""Ops Agent - Alerting subsystem.

Defines alert rules, notification channels, escalation policies,
and generates Prometheus alerting rules YAML.
"""

import yaml
from dataclasses import dataclass, field
from enum import Enum
from typing import Optional


class Severity(str, Enum):
    INFO = "info"
    WARNING = "warning"
    CRITICAL = "critical"


class ChannelType(str, Enum):
    SLACK = "slack"
    EMAIL = "email"


@dataclass
class AlertRule:
    name: str
    expr: str  # PromQL expression
    threshold: float
    duration: str  # e.g. "5m", "1h"
    severity: Severity = Severity.WARNING
    summary: str = ""
    description: str = ""


@dataclass
class NotificationChannel:
    name: str
    channel_type: ChannelType
    target: str  # webhook URL or email address
    severity_filter: list[Severity] = field(default_factory=lambda: list(Severity))


@dataclass
class EscalationPolicy:
    name: str
    stages: list[dict]  # [{"wait": "5m", "channels": ["slack-ops"]}, ...]


@dataclass
class AlertingConfig:
    rules: list[AlertRule] = field(default_factory=list)
    channels: list[NotificationChannel] = field(default_factory=list)
    escalations: list[EscalationPolicy] = field(default_factory=list)

    def add_rule(self, rule: AlertRule):
        self.rules.append(rule)

    def add_channel(self, channel: NotificationChannel):
        self.channels.append(channel)

    def add_escalation(self, policy: EscalationPolicy):
        self.escalations.append(policy)

    def to_prometheus_rules_yaml(self) -> str:
        """Generate Prometheus alerting rules YAML."""
        groups = [{
            "name": "mitran_alerts",
            "rules": [
                {
                    "alert": rule.name,
                    "expr": f"{rule.expr} > {rule.threshold}" if ">" not in rule.expr else rule.expr,
                    "for": rule.duration,
                    "labels": {"severity": rule.severity.value},
                    "annotations": {
                        "summary": rule.summary or rule.name,
                        "description": rule.description or f"{rule.name} triggered",
                    },
                }
                for rule in self.rules
            ],
        }]
        return yaml.dump({"groups": groups}, default_flow_style=False, sort_keys=False)


# Default alert rules for Mitran services
DEFAULT_RULES = [
    AlertRule(
        name="HighCPUUsage",
        expr="process_cpu_seconds_total",
        threshold=0.8,
        duration="5m",
        severity=Severity.WARNING,
        summary="High CPU usage detected",
        description="CPU usage above 80% for 5 minutes",
    ),
    AlertRule(
        name="HighMemoryUsage",
        expr="process_resident_memory_bytes / node_memory_MemTotal_bytes",
        threshold=0.85,
        duration="5m",
        severity=Severity.CRITICAL,
        summary="High memory usage",
        description="Memory usage above 85% for 5 minutes",
    ),
    AlertRule(
        name="HighErrorRate",
        expr="rate(http_requests_total{status=~'5..'}[5m]) / rate(http_requests_total[5m])",
        threshold=0.05,
        duration="2m",
        severity=Severity.CRITICAL,
        summary="High error rate",
        description="Error rate above 5% for 2 minutes",
    ),
    AlertRule(
        name="HighLatency",
        expr="histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
        threshold=2.0,
        duration="5m",
        severity=Severity.WARNING,
        summary="High p95 latency",
        description="p95 latency above 2s for 5 minutes",
    ),
    AlertRule(
        name="AgentTaskFailures",
        expr="rate(agent_task_failures_total[10m])",
        threshold=0.1,
        duration="10m",
        severity=Severity.WARNING,
        summary="Agent task failure rate elevated",
        description="Agent task failures above threshold for 10 minutes",
    ),
]

DEFAULT_CHANNELS = [
    NotificationChannel(
        name="slack-ops",
        channel_type=ChannelType.SLACK,
        target="${SLACK_OPS_WEBHOOK_URL}",
        severity_filter=[Severity.WARNING, Severity.CRITICAL],
    ),
    NotificationChannel(
        name="email-oncall",
        channel_type=ChannelType.EMAIL,
        target="${ONCALL_EMAIL}",
        severity_filter=[Severity.CRITICAL],
    ),
]

DEFAULT_ESCALATION = EscalationPolicy(
    name="default",
    stages=[
        {"wait": "0m", "channels": ["slack-ops"]},
        {"wait": "15m", "channels": ["slack-ops", "email-oncall"]},
        {"wait": "30m", "channels": ["email-oncall"]},
    ],
)


def create_default_config() -> AlertingConfig:
    """Create alerting config with sensible defaults."""
    config = AlertingConfig()
    for rule in DEFAULT_RULES:
        config.add_rule(rule)
    for ch in DEFAULT_CHANNELS:
        config.add_channel(ch)
    config.add_escalation(DEFAULT_ESCALATION)
    return config


def generate_alerting_rules_file(output_path: Optional[str] = None) -> str:
    """Generate Prometheus alerting rules and optionally write to file."""
    config = create_default_config()
    rules_yaml = config.to_prometheus_rules_yaml()
    if output_path:
        with open(output_path, "w") as f:
            f.write(rules_yaml)
    return rules_yaml
