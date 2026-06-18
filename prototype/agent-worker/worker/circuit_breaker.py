import time

class CircuitOpenError(Exception):
    pass

class CircuitBreaker:
    CLOSED, OPEN, HALF_OPEN = 0, 1, 2

    def __init__(self, threshold=5, cooldown=30.0):
        self.state, self.failures, self.threshold = self.CLOSED, 0, threshold
        self.cooldown, self.opened_at = cooldown, 0

    def call(self, fn, *args, **kwargs):
        if self.state == self.OPEN:
            if time.time() - self.opened_at >= self.cooldown:
                self.state = self.HALF_OPEN
            else:
                raise CircuitOpenError("Circuit is open")
        try:
            result = fn(*args, **kwargs)
            if self.state == self.HALF_OPEN:
                self.state, self.failures = self.CLOSED, 0
            return result
        except Exception:
            self.failures += 1
            if self.failures >= self.threshold:
                self.state, self.opened_at = self.OPEN, time.time()
            raise

def with_retry(fn, max_attempts=3, base_delay=1.0):
    for attempt in range(max_attempts):
        try:
            return fn()
        except Exception:
            if attempt == max_attempts - 1:
                raise
            time.sleep((2 ** attempt) * base_delay)
