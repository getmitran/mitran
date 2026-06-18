"""HTTP connection pool for persistent engine communication."""
import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry
from worker.config import ENGINE_URL

_session = None

def get_session() -> requests.Session:
    global _session
    if _session is None:
        _session = requests.Session()
        retry = Retry(total=3, backoff_factor=0.5, status_forcelist=[502, 503, 504])
        adapter = HTTPAdapter(pool_connections=10, pool_maxsize=20, max_retries=retry)
        _session.mount('http://', adapter)
        _session.mount('https://', adapter)
    return _session

def engine_get(path: str, **kwargs) -> requests.Response:
    return get_session().get(f'{ENGINE_URL}{path}', timeout=30, **kwargs)

def engine_post(path: str, json=None, **kwargs) -> requests.Response:
    return get_session().post(f'{ENGINE_URL}{path}', json=json, timeout=30, **kwargs)

def engine_patch(path: str, json=None, **kwargs) -> requests.Response:
    return get_session().patch(f'{ENGINE_URL}{path}', json=json, timeout=30, **kwargs)
