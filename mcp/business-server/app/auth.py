import os
from collections.abc import Mapping

HEADER_AUTH = "Authorization"

def extract_bearer_token(headers: Mapping[str, str] | None = None) -> str | None:
    """
    Extracts Bearer token from HTTP headers
    """
    if not headers:
        return None

    auth_header = next(
        (value for key, value in headers.items() if key.lower() == HEADER_AUTH.lower()),
        None,
    )

    if not auth_header:
        return None

    scheme, _, token = auth_header.partition(" ")
    if scheme.lower() == "bearer" and token.strip():
        return token.strip()

    return None

def resolve_auth_token(headers: Mapping[str, str] | None = None) -> str | None:
    """
    Defines which authorization token to use.
    Priority:
    1. Token from HTTP-headers
    2. Fallback to .env for local stdio mode
    3. Else - return None
    """
    forwarded_token = extract_bearer_token(headers)

    if forwarded_token:
        return forwarded_token

    return os.getenv("BUSINESS_API_AUTH_TOKEN")