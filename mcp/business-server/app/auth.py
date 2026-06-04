from collections.abc import Mapping

HEADER_AUTH = "Authorization"

def extract_bearer_token(headers: Mapping[str, str] | None = None) -> str | None:
    """
    Extracts Bearer token from HTTP headers
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

def resolve_auth_token(
    headers: Mapping[str, str] | None = None,
    explicit_token: str | None = None,
) -> str | None :
    """
    Defines, which authorization token to use.
    Priority:
    1.Explicit_token
    2.Token from HTTP-headers
    3.Else - return None
    """
    if explicit_token:
        return explicit_token

    forwarded_token = extract_bearer_token(headers)

    if forwarded_token:
        return forwarded_token

    return None