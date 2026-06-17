from collections.abc import Mapping
from typing import TYPE_CHECKING
if TYPE_CHECKING:
    from fastmcp import Context

from .config import settings
from .constants import HEADER_AUTH


def extract_bearer_token(
    headers: Mapping[str, str] | None = None,
) -> str | None:
    """
    Extract bearer token from incoming HTTP headers.
    """
    if not headers:
        return None

    auth_header = headers.get(HEADER_AUTH) or headers.get(HEADER_AUTH.lower())

    if not auth_header:
        return None

    scheme, _, token = auth_header.partition(" ")
    if scheme.lower() == "bearer" and token.strip():
        return token.strip()

    return None


def get_headers_from_context(ctx: "Context | None") -> Mapping[str, str] | None:
    """Safely extract HTTP headers from FastMCP Context if available."""
    if not ctx or not getattr(ctx, "request_context", None):
        return None
    request = getattr(ctx.request_context, "request", None)
    if not request:
        return None
    return getattr(request, "headers", None)


def resolve_auth_token(
    headers: Mapping[str, str] | None = None,
    explicit_token: str | None = None,
) -> str | None:
    """
    Resolve authentication token.

    Priority:
    1. Explicit token override
    2. Forwarded user token (HTTP mode)
    3. Configured Guest API token
    4. Anonymous access
    """

    if explicit_token:
        return explicit_token

    forwarded_token = extract_bearer_token(headers)

    if forwarded_token:
        return forwarded_token

    return settings.guest_api_auth_token
