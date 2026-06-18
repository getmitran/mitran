import os
import requests
from flask import Blueprint, jsonify

health_bp = Blueprint("health", __name__)
ENGINE_URL = os.getenv("ENGINE_URL", "http://localhost:7780")

@health_bp.route("/health")
def health():
    try:
        r = requests.get(f"{ENGINE_URL}/api/v1/health", timeout=3)
        connected = r.status_code == 200
    except Exception:
        connected = False
    return jsonify({"status": "ok", "engine_connected": connected})
