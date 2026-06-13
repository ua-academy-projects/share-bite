from typing import Any

from fastmcp import Context

from .. import auth
from ..http_client import guest_client
from ..server import mcp
from ._utils import unwrap_api_result


@mcp.tool(
    description=(
        "Search guest (UGC) posts. Returns paginated published posts with author info, "
        "images, likes, mentions, and metadata. Use author_id to filter by a specific user."
    )
)
async def search_posts(
    ctx: Context | None,
    limit: int = 20,
    offset: int = 0,
    author_id: str | None = None,
) -> str:
    params: dict[str, Any] = {"limit": limit, "offset": offset}
    if author_id:
        params["customer_id"] = author_id

    result = await guest_client.get(
        "/posts/",
        auth_token=auth.resolve_auth_token(headers=auth.get_headers_from_context(ctx)),
        params=params,
    )
    return unwrap_api_result(result)


@mcp.tool(
    description=(
        "Get a single guest post by its numeric ID. "
        "Returns the post with authors, images, likes, mentions, and venue linkage."
    )
)
async def get_post(ctx: Context | None, post_id: int) -> str:
    result = await guest_client.get(
        f"/posts/{post_id}",
        auth_token=auth.resolve_auth_token(headers=auth.get_headers_from_context(ctx)),
    )
    return unwrap_api_result(result)


@mcp.tool(
    description=(
        "Get all authors (owner and accepted collaborators) of a guest post. "
        "Useful to explain who contributed to a review."
    )
)
async def get_post_authors(ctx: Context | None, post_id: int) -> str:
    result = await guest_client.get(
        f"/posts/{post_id}/authors",
        auth_token=auth.resolve_auth_token(headers=auth.get_headers_from_context(ctx)),
    )
    return unwrap_api_result(result)
