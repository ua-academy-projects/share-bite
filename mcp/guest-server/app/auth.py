from collections.abc import Mapping

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

    if auth_header.startswith("Bearer "):
        return auth_header.removeprefix("Bearer ").strip()

    return None


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
