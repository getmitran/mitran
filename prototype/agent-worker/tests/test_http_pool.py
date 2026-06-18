import pytest
from unittest.mock import patch, MagicMock
from worker.http_pool import get_session, engine_get, engine_post

def test_get_session_singleton():
    s1 = get_session()
    s2 = get_session()
    assert s1 is s2

def test_session_has_retry_adapter():
    session = get_session()
    adapter = session.get_adapter('http://')
    assert adapter.max_retries.total == 3

@patch('worker.http_pool.get_session')
def test_engine_get(mock_session):
    mock_resp = MagicMock(status_code=200)
    mock_session.return_value.get.return_value = mock_resp
    resp = engine_get('/health')
    assert resp.status_code == 200

@patch('worker.http_pool.get_session')
def test_engine_post(mock_session):
    mock_resp = MagicMock(status_code=201)
    mock_session.return_value.post.return_value = mock_resp
    resp = engine_post('/api/v1/tasks', json={'title': 'test'})
    assert resp.status_code == 201
