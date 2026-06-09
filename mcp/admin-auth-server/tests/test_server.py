from typing import Any
import pytest
from unittest.mock import AsyncMock, MagicMock, patch
from mcp_app.http_client import admin_client
from mcp_app.tools import register_tools


@pytest.fixture(scope="module")
def extracted_tools():
    class ToolInterceptor:
        def __init__(self):
            self.funcs = {}
        def tool(self, *_args: Any, **_kwargs: Any):
            def decorator(func: Any):
                self.funcs[func.__name__] = func
                return func
            return decorator

    interceptor = ToolInterceptor()
    register_tools(interceptor)
    return interceptor.funcs


@pytest.fixture(autouse=True)
def mock_dependencies():
    admin_client.get = AsyncMock()
    admin_client.post = AsyncMock()
    mock_audit = MagicMock()
    with patch("mcp_app.auth.log_audit_event", mock_audit), \
            patch("mcp_app.tools.log_audit_event", mock_audit):
        yield mock_audit


@pytest.mark.asyncio
async def test_successful_moderator_access(extracted_tools: dict, mock_dependencies: Any):
    admin_auth_health_check = extracted_tools["admin_auth_health_check"]

    admin_client.get.side_effect = [
        {"is_error": False, "data": {"id": "mod_123", "role": "moderator"}},
        {"is_error": False, "data": {"status": "healthy", "database": "connected"}}
    ]

    response = await admin_auth_health_check(auth_token="valid_moderator_jwt")
    assert "healthy" in response
    assert admin_client.get.call_count == 2
    mock_dependencies.assert_called_with("admin_auth_health_check", "mod_123", "SUCCESS", "Infrastructure health cleared.")


@pytest.mark.asyncio
async def test_unauthorized_access_rejection(extracted_tools: dict, mock_dependencies: Any):
    admin_auth_health_check = extracted_tools["admin_auth_health_check"]

    admin_client.get.return_value = {
        "is_error": True,
        "error_message": "Go API Error (401): Unauthorized token payload"
    }

    response = await admin_auth_health_check(auth_token="invalid_or_expired_token")

    assert "ERROR" in response
    assert "Unauthorized" in response
    mock_dependencies.assert_called_once()
    assert mock_dependencies.call_args.args[2] == "DENIED"


@pytest.mark.asyncio
async def test_insufficient_role_forbidden(extracted_tools: dict, mock_dependencies: Any):
    validate_admin_permissions = extracted_tools["validate_admin_permissions"]

    admin_client.get.return_value = {
        "is_error": False,
        "data": {"id": "mod_123", "role": "moderator"}
    }
    response = await validate_admin_permissions(target_permission="root_write", auth_token="moderator_token")

    assert "ERROR" in response
    assert "Forbidden" in response
    assert mock_dependencies.call_args.args[2] == "DENIED"


@pytest.mark.asyncio
async def test_completely_missing_auth_rejection(extracted_tools: dict, mock_dependencies: Any):
    admin_auth_health_check = extracted_tools["admin_auth_health_check"]
    with patch("mcp_app.auth.settings") as mock_settings:
        mock_settings.enforce_authentication = True
        mock_settings.local_refresh_token = None

        response = await admin_auth_health_check(auth_token=None)

        assert "ERROR" in response
        assert "Missing" in response or "No active session" in response
        assert mock_dependencies.call_args.args[2] == "DENIED"


@pytest.mark.asyncio
async def test_successful_admin_write_action(extracted_tools: dict, mock_dependencies: Any):
    validate_admin_permissions = extracted_tools["validate_admin_permissions"]
    admin_client.get.return_value = {"is_error": False, "data": {"id": "admin_999", "role": "admin"}}
    admin_client.post.return_value = {"is_error": False, "data": {"authorized": True}}

    response = await validate_admin_permissions(target_permission="delete_orders", auth_token="super_admin_jwt")
    assert "authorized" in response
    assert mock_dependencies.call_count == 1
    assert mock_dependencies.call_args.args[0] == "validate_admin_permissions"
    assert mock_dependencies.call_args.args[1] == "admin_999"
    assert mock_dependencies.call_args.args[2] == "SUCCESS"

@pytest.mark.asyncio
async def test_health_check_connection_failure(extracted_tools: dict, mock_dependencies: Any):
    admin_auth_health_check = extracted_tools["admin_auth_health_check"]
    admin_client.get.side_effect = [
        {"is_error": False, "data": {"id": "mod_123", "role": "moderator"}},
        {"is_error": True, "error_message": "Network error: Connection refused by remote host"}
    ]

    response = await admin_auth_health_check(auth_token="any_token")
    assert "unhealthy" in response
    assert "Connection refused" in response
    assert mock_dependencies.call_args.args[2] == "ERROR"