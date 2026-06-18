"""Retry logic with exponential backoff for LLM and HTTP calls."""
import time
import logging
from functools import wraps
from typing import Callable, TypeVar, Any

log = logging.getLogger(__name__)
T = TypeVar('T')

def retry(
    max_attempts: int = 3,
    base_delay: float = 1.0,
    max_delay: float = 30.0,
    exceptions: tuple = (Exception,),
    on_retry: Callable | None = None,
):
    def decorator(func: Callable[..., T]) -> Callable[..., T]:
        @wraps(func)
        def wrapper(*args: Any, **kwargs: Any) -> T:
            last_exc = None
            for attempt in range(1, max_attempts + 1):
                try:
                    return func(*args, **kwargs)
                except exceptions as e:
                    last_exc = e
                    if attempt == max_attempts:
                        raise
                    delay = min(base_delay * (2 ** (attempt - 1)), max_delay)
                    log.warning(f'Retry {attempt}/{max_attempts} for {func.__name__}: {e} (delay={delay:.1f}s)')
                    if on_retry:
                        on_retry(attempt, e)
                    time.sleep(delay)
            raise last_exc
        return wrapper
    return decorator


def retry_call(
    func: Callable[..., T],
    *args: Any,
    max_attempts: int = 3,
    base_delay: float = 1.0,
    **kwargs: Any,
) -> T:
    last_exc = None
    for attempt in range(1, max_attempts + 1):
        try:
            return func(*args, **kwargs)
        except Exception as e:
            last_exc = e
            if attempt == max_attempts:
                raise
            delay = min(base_delay * (2 ** (attempt - 1)), 30.0)
            log.warning(f'Retry {attempt}/{max_attempts}: {e} (delay={delay:.1f}s)')
            time.sleep(delay)
    raise last_exc
