import json

import httpx
import pytest

from app.constants import OPENAPI_SPECIFICATION_PATH
from app.resources.api import get_api_info, get_openapi_summary
from tests.factories.guest_responses import (
    build_error_response,
    build_info_response,
    build_openapi_response,
)


@pytest.mark.asyncio
async def test_get_api_info_success(mock_guest_api):
    """Returns parsed JSON with API info when Guest API responds with 200."""
    mock_guest_api.get("/info").mock(return_value=build_info_response())

    result = await get_api_info()
    data = json.loads(result)

    assert data["version"] == "1.0.0"
    assert data["commit"] == "abc123"
    assert data["buildTime"] == "2026-01-01T00:00:00Z"
    assert data["environment"] == "dev"


@pytest.mark.asyncio
async def test_get_api_info_unauthorized(mock_guest_api):
    """Raises RuntimeError when Guest API returns 401."""
    mock_guest_api.get("/info").mock(
        return_value=build_error_response(401, "Unauthorized")
    )

    with pytest.raises(RuntimeError, match="Failed to fetch API info"):
        await get_api_info()


@pytest.mark.asyncio
async def test_get_api_info_connection_error(mock_guest_api):
    """Raises RuntimeError when Guest API is unreachable."""
    mock_guest_api.get("/info").mock(
        side_effect=httpx.ConnectError("Connection refused")
    )

    with pytest.raises(RuntimeError, match="Failed to fetch API info"):
        await get_api_info()


@pytest.mark.asyncio
async def test_get_openapi_summary_success(mock_guest_api):
    """Returns valid JSON string with OpenAPI spec when Guest API responds with 200."""
    mock_guest_api.get(OPENAPI_SPECIFICATION_PATH).mock(
        return_value=build_openapi_response()
    )

    result = await get_openapi_summary()
    data = json.loads(result)

    assert data["swagger"] == "2.0"
    assert "info" in data
    assert "paths" in data


@pytest.mark.asyncio
async def test_get_openapi_summary_unauthorized(mock_guest_api):
    """Raises RuntimeError when Guest API returns 401."""
    mock_guest_api.get(OPENAPI_SPECIFICATION_PATH).mock(
        return_value=build_error_response(401, "Unauthorized")
    )

    with pytest.raises(RuntimeError, match="Failed to fetch OpenAPI spec"):
        await get_openapi_summary()


@pytest.mark.asyncio
async def test_get_openapi_summary_connection_error(mock_guest_api):
    """Raises RuntimeError when Guest API is unreachable."""
    mock_guest_api.get(OPENAPI_SPECIFICATION_PATH).mock(
        side_effect=httpx.ConnectError("Connection refused")
    )

    with pytest.raises(RuntimeError, match="Failed to fetch OpenAPI spec"):
        await get_openapi_summary()
