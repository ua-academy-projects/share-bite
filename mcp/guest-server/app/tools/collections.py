import json
from typing import Any

from fastmcp import Context

from .. import auth
from ..http_client import guest_client
from ..server import mcp
from ._utils import unwrap_api_result


@mcp.tool(
    description="List current user's collections. Requires authentication (Bearer token)."
)
async def list_my_collections(
    ctx: Context,
    page_size: int = 20,
    page_token: str | None = None,
) -> str:
    token = auth.resolve_auth_token(headers=auth.get_headers_from_context(ctx))
    if not token:
        return json.dumps(
            {
                "error": "unauthorized",
                "message": "Authentication required to list your collections.",
            }
        )

    params: dict[str, Any] = {"pageSize": page_size}
    if page_token:
        params["pageToken"] = page_token

    result = await guest_client.get(
        "/collections/me",
        auth_token=token,
        params=params,
    )
    return unwrap_api_result(result)


@mcp.tool(
    description=(
        "Get a collection by ID. Public collections are accessible without auth; "
        "private collections require ownership."
    )
)
async def get_collection(ctx: Context, collection_id: str) -> str:
    result = await guest_client.get(
        f"/collections/{collection_id}",
        auth_token=auth.resolve_auth_token(headers=auth.get_headers_from_context(ctx)),
    )
    return unwrap_api_result(result)


@mcp.tool(
    description=(
        "List venues in a collection, ordered by sort order. "
        "Public collections accessible without auth."
    )
)
async def get_collection_venues(ctx: Context, collection_id: str) -> str:
    result = await guest_client.get(
        f"/collections/{collection_id}/venues",
        auth_token=auth.resolve_auth_token(headers=auth.get_headers_from_context(ctx)),
    )
    return unwrap_api_result(result)
