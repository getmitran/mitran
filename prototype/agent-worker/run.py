#!/usr/bin/env python3
"""Mitran Agent Worker entry point."""
from worker.app import app, register_agents
from worker.config import WORKER_PORT, WORKER_HOST

if __name__ == "__main__":
    register_agents()
    app.run(host=WORKER_HOST, port=WORKER_PORT, debug=False)
