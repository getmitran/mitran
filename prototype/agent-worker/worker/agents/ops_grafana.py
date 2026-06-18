"""Ops Agent Grafana Integration - Creates dashboards via Grafana HTTP API."""

import os
import json
import requests
from typing import Any


GRAFANA_URL = os.getenv("GRAFANA_URL", "http://localhost:3000")
GRAFANA_API_KEY = os.getenv("GRAFANA_API_KEY", "")


def _headers() -> dict[str, str]:
    h = {"Content-Type": "application/json"}
    if GRAFANA_API_KEY:
        h["Authorization"] = f"Bearer {GRAFANA_API_KEY}"
    return h


def create_dashboard(title: str, panels: list[dict[str, Any]]) -> dict[str, Any]:
    """Create a Grafana dashboard via POST /api/dashboards/db.

    Args:
        title: Dashboard title.
        panels: List of panel dicts. Each panel should have at minimum
                'title', 'type' (e.g. 'graph', 'stat', 'timeseries'),
                and 'targets' (list of query dicts).

    Returns:
        Grafana API response dict with id, uid, url, status, version.
    """
    positioned_panels = []
    for i, panel in enumerate(panels):
        p = {
            "id": i + 1,
            "gridPos": {"h": 8, "w": 12, "x": (i % 2) * 12, "y": (i // 2) * 8},
            **panel,
        }
        positioned_panels.append(p)

    payload = {
        "dashboard": {
            "id": None,
            "uid": None,
            "title": title,
            "panels": positioned_panels,
            "schemaVersion": 39,
            "version": 0,
            "refresh": "5s",
        },
        "message": f"Created by Mitran Ops Agent: {title}",
        "overwrite": False,
    }

    resp = requests.post(
        f"{GRAFANA_URL}/api/dashboards/db",
        headers=_headers(),
        json=payload,
        timeout=30,
    )
    resp.raise_for_status()
    return resp.json()


class OpsGrafanaIntegration:
    """Mixin for OpsAgentV2 providing Grafana dashboard creation."""

    def handle_grafana_request(self, action: str, params: dict[str, Any]) -> dict[str, Any]:
        if action == "create_dashboard":
            return create_dashboard(
                title=params["title"],
                panels=params.get("panels", []),
            )
        return {"error": f"Unknown grafana action: {action}"}
