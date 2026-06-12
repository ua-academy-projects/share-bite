import json

from fastmcp import Context

from .. import auth

from ..http_client import guest_client, APIResponse, APIErrorResponse
from ..server import mcp


def _unwrap_guest_result(action: str, result: APIResponse | APIErrorResponse) -> str:
    if result.get("is_error"):
        err_msg = result.get("error_message", "Unknown error")
        status_code = result.get("status", "N/A")

        raise RuntimeError(f"{action} failed (HTTP {status_code}): {err_msg}")

    return json.dumps(result.get("data", {}))


@mcp.tool(description="Checks the basic health status of the Guest API")
async def guest_health_check(ctx: Context) -> str:
    """
    Checks the basic health status of the Guest API.
    Use this to verify if the service container is up and running.
    """
    result = await guest_client.get(
        "/health",
        auth_token=auth.resolve_auth_token(headers=auth.get_headers_from_context(ctx)),
    )
    return _unwrap_guest_result("Health check", result)


@mcp.tool(description="Fetches the deep operational status of the Guest API")
async def get_guest_api_status(ctx: Context) -> str:
    """
    Fetches the deep operational status of the Guest API.
    Returns the connection status for internal components like PostgreSQL and Redis.
    """
    result = await guest_client.get(
        "/status",
        auth_token=auth.resolve_auth_token(headers=auth.get_headers_from_context(ctx)),
    )
    return _unwrap_guest_result("Status check", result)
