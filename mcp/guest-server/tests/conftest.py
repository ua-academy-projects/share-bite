import pytest
import respx

from app.config import settings


@pytest.fixture
def base_url() -> str:
    return str(settings.guest_api_base_url).rstrip("/")


@pytest.fixture
def mock_guest_api(base_url: str):
    """
    Activates respx mock router scoped to the Guest API base URL.
    assert_all_mocked=True prevents accidental real HTTP calls during tests.
    """
    with respx.mock(base_url=base_url, assert_all_mocked=True) as router:
        yield router
