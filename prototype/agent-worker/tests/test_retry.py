import pytest
from unittest.mock import patch
from worker.retry import retry, retry_call

def test_retry_succeeds_first_try():
    @retry(max_attempts=3)
    def succeed(): return 42
    assert succeed() == 42

def test_retry_succeeds_after_failure():
    calls = [0]
    @retry(max_attempts=3, base_delay=0.01)
    def flaky():
        calls[0] += 1
        if calls[0] < 3: raise ValueError('fail')
        return 'ok'
    assert flaky() == 'ok'
    assert calls[0] == 3

def test_retry_exhausts_attempts():
    @retry(max_attempts=2, base_delay=0.01)
    def always_fail(): raise RuntimeError('boom')
    with pytest.raises(RuntimeError, match='boom'):
        always_fail()

def test_retry_call_function():
    calls = [0]
    def flaky():
        calls[0] += 1
        if calls[0] < 2: raise IOError('tmp')
        return 'done'
    result = retry_call(flaky, max_attempts=3, base_delay=0.01)
    assert result == 'done'
