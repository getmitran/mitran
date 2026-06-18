import pytest

BASE_URL = "http://localhost:7780"


@pytest.fixture
def base_url():
    return BASE_URL
