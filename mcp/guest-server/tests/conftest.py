import pytest
import respx

from unittest.mock import patch

from app.http_client import guest_client


@pytest.fixture
def mock_guest_api():
    with respx.mock(base_url=guest_client.base_url) as respx_mock:
        yield respx_mock


@pytest.fixture
def fake_auth_token():
    """Patch resolve_auth_token without hardcoding its signature."""
    with patch("app.auth.resolve_auth_token", return_value="fake-jwt"):
        yield "fake-jwt"
