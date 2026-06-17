from __future__ import annotations

import os
from collections.abc import Mapping

AUTH_HEADER = "authorization"


def normalize_bearer_token(token: str | None) -> str | None:
    if token is None:
        return None
    token = token.strip()
    if not token:
        return None
    if token.lower().startswith("bearer "):
        return token
    return f"Bearer {token}"


def extract_auth_token(
    headers: Mapping[str, str] | None = None,
    *,
    fallback_env: str | None = None,
) -> str | None:
    headers = headers or {}
    for key, value in headers.items():
        if key.lower() == AUTH_HEADER:
            return _strip_bearer(value)

    if fallback_env:
        return _strip_bearer(os.getenv(fallback_env))

    return None


def _strip_bearer(value: str | None) -> str | None:
    if value is None:
        return None
    value = value.strip()
    if not value:
        return None
    if value.lower().startswith("bearer "):
        return value[7:].strip() or None
    return value
