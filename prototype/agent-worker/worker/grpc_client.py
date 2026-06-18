import grpc
import json
from dataclasses import dataclass, field


@dataclass
class TaskAssignment:
    task_id: str
    title: str
    description: str
    agent_type: str
    project_id: str
    workspace_path: str
    context: dict = field(default_factory=dict)
    skill_paths: list[str] = field(default_factory=list)


@dataclass
class TaskResult:
    task_id: str
    summary: str
    files: list[dict] = field(default_factory=list)  # [{path, action, content, diff}]
    questions: list[str] = field(default_factory=list)
    status: str = 'complete'  # 'complete' or 'failed'


class MitranGRPCClient:
    """Connects to Go engine via gRPC. Falls back to HTTP REST if gRPC unavailable."""

    def __init__(self, engine_host='localhost', grpc_port=7781, http_port=7780):
        self.engine_host = engine_host
        self.grpc_port = grpc_port
        self.http_port = http_port
        self._channel = None
        self._use_grpc = False

    def connect(self):
        """Try gRPC connection, fall back to HTTP."""
        try:
            self._channel = grpc.insecure_channel(f'{self.engine_host}:{self.grpc_port}')
            grpc.channel_ready_future(self._channel).result(timeout=3)
            self._use_grpc = True
        except Exception:
            self._use_grpc = False

    def register(self, agent_name: str, agent_type: str, callback_url: str) -> bool:
        """Register agent with engine."""
        import requests
        try:
            resp = requests.post(
                f'http://{self.engine_host}:{self.http_port}/api/v1/agents/register',
                json={'name': agent_name, 'type': agent_type, 'callback_url': callback_url},
                timeout=5,
            )
            return resp.status_code in (200, 201)
        except requests.RequestException:
            return False

    def submit_result(self, result: TaskResult) -> bool:
        """Submit task completion to engine."""
        import requests
        payload = {
            'task_id': result.task_id,
            'checkpoint': {
                'status': result.status,
                'summary': result.summary,
                'files': result.files,
                'questions': result.questions,
            },
        }
        try:
            resp = requests.post(
                f'http://{self.engine_host}:{self.http_port}/api/v1/agents/task-complete',
                json=payload,
                timeout=30,
            )
            return resp.status_code in (200, 201)
        except requests.RequestException:
            return False

    def health_check(self) -> dict:
        """Check engine health."""
        import requests
        try:
            resp = requests.get(f'http://{self.engine_host}:{self.http_port}/health', timeout=5)
            return resp.json()
        except requests.RequestException:
            return {'status': 'unreachable'}

    def close(self):
        if self._channel:
            self._channel.close()
            self._channel = None
