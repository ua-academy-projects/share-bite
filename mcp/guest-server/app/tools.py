import json

from .http_client import guest_client
from .server import mcp


@mcp.tool(description="Checks the basic health status of the Guest API.")
async def guest_health_check() -> str:
    """Perform a health check on the Go backend."""
    result = await guest_client.get("/health")
    if result["is_error"] is True:
        raise RuntimeError(f"Health check failed: {result['error_message']}")
    return json.dumps(result["data"])


@mcp.tool(
    description="Retrieves the detailed operational status of the Guest API, including DB connections."
)
async def get_guest_api_status() -> str:
    """Fetch deep operational status."""
    result = await guest_client.get("/status")
    if result["is_error"] is True:
        raise RuntimeError(f"Status check failed: {result['error_message']}")
    return json.dumps(result["data"])
