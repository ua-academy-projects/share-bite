import json

import httpx
import pytest

from app.tools.health import get_guest_api_status
from tests.factories.guest_responses import build_status_response, build_error_response


@pytest.mark.asyncio
async def test_get_guest_api_status_all_connected(mock_guest_api):
    """Returns parsed JSON with all components connected when Guest API responds with 200."""
    mock_guest_api.get("/status").mock(return_value=build_status_response())

    result = await get_guest_api_status(None)
    data = json.loads(result)

    assert data == {"app": "share-bite", "database": "connected", "redis": "connected"}


@pytest.mark.asyncio
async def test_get_guest_api_status_db_disconnected(mock_guest_api):
    """Raises RuntimeError when database is disconnected and Guest API returns 503."""
    mock_guest_api.get("/status").mock(
        return_value=build_status_response(db="disconnected", status=503)
    )

    with pytest.raises(RuntimeError, match="503"):
        await get_guest_api_status(None)


@pytest.mark.asyncio
async def test_get_guest_api_status_redis_disconnected(mock_guest_api):
    """Raises RuntimeError when Redis is disconnected and Guest API returns 503."""
    mock_guest_api.get("/status").mock(
        return_value=build_status_response(redis="disconnected", status=503)
    )

    with pytest.raises(RuntimeError, match="503"):
        await get_guest_api_status(None)


@pytest.mark.asyncio
async def test_get_guest_api_status_unauthorized(mock_guest_api):
    """Raises RuntimeError containing status code when Guest API returns 401."""
    mock_guest_api.get("/status").mock(
        return_value=build_error_response(401, "Unauthorized")
    )

    with pytest.raises(RuntimeError, match="401"):
        await get_guest_api_status(None)


@pytest.mark.asyncio
async def test_get_guest_api_status_connection_error(mock_guest_api):
    """Raises RuntimeError when Guest API is unreachable."""
    mock_guest_api.get("/status").mock(
        side_effect=httpx.ConnectError("Connection refused")
    )

    with pytest.raises(RuntimeError, match="Failed to connect"):
        await get_guest_api_status(None)


@pytest.mark.asyncio
async def test_get_guest_api_status_timeout(mock_guest_api):
    """Raises RuntimeError when the request to Guest API times out."""
    mock_guest_api.get("/status").mock(side_effect=httpx.TimeoutException("Timed out"))

    with pytest.raises(RuntimeError, match="timed out"):
        await get_guest_api_status(None)
