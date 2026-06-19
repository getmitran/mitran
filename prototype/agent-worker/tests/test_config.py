"""Tests for config loading and defaults."""
import os
from unittest.mock import patch


def test_default_engine_url():
    with patch.dict(os.environ, {}, clear=True):
        # Re-import to pick up defaults
        import importlib
        import worker.config as cfg
        importlib.reload(cfg)
        assert cfg.ENGINE_URL == "http://localhost:7780"


def test_default_worker_port():
    with patch.dict(os.environ, {}, clear=True):
        import importlib
        import worker.config as cfg
        importlib.reload(cfg)
        assert cfg.WORKER_PORT == 8888


def test_custom_engine_url():
    with patch.dict(os.environ, {"MITRAN_ENGINE_URL": "http://custom:9999"}):
        import importlib
        import worker.config as cfg
        importlib.reload(cfg)
        assert cfg.ENGINE_URL == "http://custom:9999"


def test_default_aws_region():
    with patch.dict(os.environ, {}, clear=True):
        import importlib
        import worker.config as cfg
        importlib.reload(cfg)
        assert cfg.AWS_REGION == "us-east-1"


def test_default_model_id():
    with patch.dict(os.environ, {}, clear=True):
        import importlib
        import worker.config as cfg
        importlib.reload(cfg)
        assert "claude" in cfg.MODEL_ID.lower() or "anthropic" in cfg.MODEL_ID.lower()
