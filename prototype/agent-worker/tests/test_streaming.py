import pytest
from unittest.mock import patch, MagicMock

def test_stream_endpoint_import():
    from worker.stream_endpoint import register_stream_routes
    assert callable(register_stream_routes)

def test_llm_streaming_import():
    from worker.llm_streaming import stream_invoke, get_client
    assert callable(stream_invoke)
    assert callable(get_client)
