from typing import Any

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
    limit: int = 20,
    offset: int = 0,
    author_id: str | None = None,
) -> str:
    params: dict[str, Any] = {"limit": limit, "offset": offset}
    if author_id:
        params["customer_id"] = author_id

    result = await guest_client.get(
        "/posts/",
        auth_token=auth.resolve_auth_token(),
        params=params,
    )
    return unwrap_api_result(result)


@mcp.tool(
    description=(
        "Get a single guest post by its numeric ID. "
        "Returns the post with authors, images, likes, mentions, and venue linkage."
    )
)
async def get_post(id: int) -> str:
    result = await guest_client.get(
        f"/posts/{id}",
        auth_token=auth.resolve_auth_token(),
    )
    return unwrap_api_result(result)


@mcp.tool(
    description=(
        "Get all authors (owner and accepted collaborators) of a guest post. "
        "Useful to explain who contributed to a review."
    )
)
async def get_post_authors(id: int) -> str:
    result = await guest_client.get(
        f"/posts/{id}/authors",
        auth_token=auth.resolve_auth_token(),
    )
    return unwrap_api_result(result)
