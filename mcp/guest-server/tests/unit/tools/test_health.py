import json

import httpx
import pytest

from app.tools.health import guest_health_check
from tests.factories.guest_responses import build_health_response, build_error_response


@pytest.mark.asyncio
async def test_health_check_success(mock_guest_api):
    """Returns parsed JSON when Guest API responds with 200."""
    mock_guest_api.get("/health").mock(return_value=build_health_response())

    result = await guest_health_check(None)
    data = json.loads(result)

    assert data == {"status": "OK"}


@pytest.mark.asyncio
async def test_health_check_service_unavailable(mock_guest_api):
    """Raises RuntimeError containing status code when Guest API returns 503."""
    mock_guest_api.get("/health").mock(
        return_value=build_error_response(503, "Service Unavailable")
    )

    with pytest.raises(RuntimeError, match="503"):
        await guest_health_check(None)


@pytest.mark.asyncio
async def test_health_check_unauthorized(mock_guest_api):
    """Raises RuntimeError containing status code when Guest API returns 401."""
    mock_guest_api.get("/health").mock(
        return_value=build_error_response(401, "Unauthorized")
    )

    with pytest.raises(RuntimeError, match="401"):
        await guest_health_check(None)


@pytest.mark.asyncio
async def test_health_check_connection_error(mock_guest_api):
    """Raises RuntimeError when Guest API is unreachable."""
    mock_guest_api.get("/health").mock(
        side_effect=httpx.ConnectError("Connection refused")
    )

    with pytest.raises(RuntimeError, match="Failed to connect"):
        await guest_health_check(None)


@pytest.mark.asyncio
async def test_health_check_timeout(mock_guest_api):
    """Raises RuntimeError when the request to Guest API times out."""
    mock_guest_api.get("/health").mock(side_effect=httpx.TimeoutException("Timed out"))

    with pytest.raises(RuntimeError, match="timed out"):
        await guest_health_check(None)
