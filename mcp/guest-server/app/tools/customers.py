from typing import Any

from fastmcp import Context

from .. import auth
from ..http_client import guest_client
from ..server import mcp
from ._utils import unwrap_api_result


@mcp.tool(
    description="Get a public customer profile by their unique username. No auth required."
)
async def get_customer_by_username(ctx: Context, username: str) -> str:
    result = await guest_client.get(
        f"/customers/{username}",
        auth_token=auth.resolve_auth_token(headers=auth.get_headers_from_context(ctx)),
    )
    return unwrap_api_result(result)


@mcp.tool(
    description=(
        "Get paginated followers for a customer. "
        "If the profile is private, only the owner can view the list."
    )
)
async def get_customer_followers(
    ctx: Context,
    id: str,
    page_size: int = 20,
    page_token: str | None = None,
) -> str:
    params: dict[str, Any] = {"pageSize": page_size}
    if page_token:
        params["pageToken"] = page_token

    result = await guest_client.get(
        f"/customers/{id}/followers",
        auth_token=auth.resolve_auth_token(headers=auth.get_headers_from_context(ctx)),
        params=params,
    )
    return unwrap_api_result(result)


@mcp.tool(
    description=(
        "Get paginated following list for a customer. "
        "If the profile is private, only the owner can view the list."
    )
)
async def get_customer_following(
    ctx: Context,
    id: str,
    page_size: int = 20,
    page_token: str | None = None,
) -> str:
    params: dict[str, Any] = {"pageSize": page_size}
    if page_token:
        params["pageToken"] = page_token

    result = await guest_client.get(
        f"/customers/{id}/following",
        auth_token=auth.resolve_auth_token(headers=auth.get_headers_from_context(ctx)),
        params=params,
    )
    return unwrap_api_result(result)
