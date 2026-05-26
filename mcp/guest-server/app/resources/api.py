import json

from app.auth import resolve_auth_token

from ..constants import (
    CONTENT_TYPE_JSON,
    OPENAPI_SPECIFICATION_PATH,
    URI_API_INFO,
    URI_OPENAPI_SUMMARY,
)
from ..http_client import guest_client
from ..server import mcp


@mcp.resource(
    uri=URI_API_INFO,
    name="guest_api_info",
    title="Guest API Info",
    description="General information about Guest API, version and configuration.",
    mime_type=CONTENT_TYPE_JSON,
)
async def get_api_info() -> str:
    """Fetch and return general API information."""
    result = await guest_client.get("/info", auth_token=resolve_auth_token())
    if result["is_error"] is True:
        raise RuntimeError(f"Failed to fetch API info: {result['error_message']}")
    return json.dumps(result["data"])


@mcp.resource(
    uri=URI_OPENAPI_SUMMARY,
    name="guest_openapi_summary",
    title="Guest API OpenAPI Summary",
    description="Swagger/OpenAPI documentation for understanding Guest API structure.",
    mime_type=CONTENT_TYPE_JSON,
)
async def get_openapi_summary() -> str:
    """Fetch and return the OpenAPI JSON specification."""
    result = await guest_client.get(
        OPENAPI_SPECIFICATION_PATH, auth_token=resolve_auth_token()
    )
    if result["is_error"] is True:
        raise RuntimeError(f"Failed to fetch OpenAPI spec: {result['error_message']}")
    return json.dumps(result["data"])
