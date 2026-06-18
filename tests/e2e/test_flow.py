"""End-to-end test: init project -> create task -> agent executes -> result."""
import requests
import time
import pytest

ENGINE = 'http://localhost:7780'
WORKER = 'http://localhost:8888'

@pytest.fixture(scope='module')
def project():
    resp = requests.post(f'{ENGINE}/api/v1/init', json={
        'name': 'e2e-test',
        'description': 'E2E test project',
        'languages': ['Python'],
        'team_size': 2,
    })
    assert resp.status_code == 200
    return resp.json()

def test_health(project):
    resp = requests.get(f'{ENGINE}/health')
    assert resp.status_code == 200
    data = resp.json()
    assert data['status'] == 'healthy'

def test_create_task(project):
    resp = requests.post(f'{ENGINE}/api/v1/tasks', json={
        'title': 'E2E test task',
        'description': 'Automated test',
        'priority': 3,
        'agent': 'dev',
    })
    assert resp.status_code in (200, 201)
    task = resp.json()
    assert task['title'] == 'E2E test task'
    return task

def test_list_tasks(project):
    resp = requests.get(f'{ENGINE}/api/v1/tasks')
    assert resp.status_code == 200
    tasks = resp.json()
    assert isinstance(tasks, list)
    assert len(tasks) >= 1

def test_worker_health():
    try:
        resp = requests.get(f'{WORKER}/health', timeout=5)
        assert resp.status_code == 200
    except requests.ConnectionError:
        pytest.skip('Worker not running')

def test_chat_endpoint():
    try:
        resp = requests.post(f'{WORKER}/chat', json={
            'message': 'Hello',
            'agent': 'dev',
        }, timeout=30)
        assert resp.status_code == 200
    except requests.ConnectionError:
        pytest.skip('Worker not running')
