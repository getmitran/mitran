import logging
import requests
from flask import Flask, jsonify, request as flask_request
from worker.config import ENGINE_URL, WORKER_PORT, WORKER_HOST
from worker.agents import AGENTS
from worker.agents.dev_agent import AgentResult

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
log = logging.getLogger(__name__)

app = Flask(__name__)


def register_agents():
    """Register all 8 agents with the Core Engine on startup."""
    callback_base = f"http://localhost:{WORKER_PORT}"
    for name, agent in AGENTS.items():
        payload = {
            "name": name,
            "type": agent.agent_type,
            "callback_url": f"{callback_base}/execute",
        }
        try:
            resp = requests.post(f"{ENGINE_URL}/api/v1/agents/register", json=payload, timeout=5)
            if resp.status_code in (200, 201):
                log.info(f"Registered agent: {name}")
            else:
                log.warning(f"Failed to register {name}: {resp.status_code} {resp.text}")
        except requests.ConnectionError:
            log.warning(f"Engine not reachable at {ENGINE_URL} — will retry on next request")
        except Exception as e:
            log.error(f"Error registering {name}: {e}")


@app.route("/health", methods=["GET"])
def health():
    return jsonify({"status": "ok", "agents": list(AGENTS.keys())})


@app.route("/execute", methods=["POST"])
def execute():
    """Receive a task assignment from the Engine, execute it, post result back."""
    data = flask_request.get_json()
    if not data:
        return jsonify({"error": "No JSON body"}), 400

    task_id = data.get("task_id", "unknown")
    agent_name = data.get("agent", "dev")
    task_description = data.get("task", "")
    context = data.get("context", {})

    agent = AGENTS.get(agent_name)
    if not agent:
        return jsonify({"error": f"Unknown agent: {agent_name}"}), 404

    log.info(f"Executing task {task_id} with agent '{agent_name}': {task_description[:80]}")

    try:
        result: AgentResult = agent.execute(task_description, context)
    except RuntimeError as e:
        error_msg = str(e)
        log.error(f"Agent {agent_name} failed: {error_msg}")
        # Post failure back to engine
        _post_result(task_id, {"status": "failed", "error": error_msg})
        return jsonify({"error": error_msg}), 500

    # Post success result back to engine
    checkpoint = {
        "status": "completed",
        "summary": result.summary,
        "files": [{"path": f.path, "content": f.content, "action": f.action} for f in result.files],
    }
    _post_result(task_id, checkpoint)

    return jsonify({"status": "completed", "summary": result.summary, "file_count": len(result.files)})


def _post_result(task_id: str, checkpoint: dict):
    """Post task result back to the Core Engine."""
    payload = {"task_id": task_id, "checkpoint": checkpoint}
    try:
        resp = requests.post(f"{ENGINE_URL}/api/v1/agents/task-complete", json=payload, timeout=30)
        if resp.status_code not in (200, 201):
            log.warning(f"Engine rejected result for {task_id}: {resp.status_code}")
    except requests.ConnectionError:
        log.warning(f"Could not reach engine to post result for task {task_id}")
    except Exception as e:
        log.error(f"Error posting result: {e}")


def create_app():
    return app
