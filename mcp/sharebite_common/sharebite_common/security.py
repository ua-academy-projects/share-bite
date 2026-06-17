from __future__ import annotations

from collections.abc import Mapping, Sequence
from typing import Any

SECRET_FIELD_NAMES = {
    "password",
    "passwd",
    "token",
    "access_token",
    "refresh_token",
    "auth_token",
    "authorization",
    "secret",
    "api_key",
    "apikey",
    "jwt",
    "cookie",
    "set_cookie",
}

REDACTED = "[REDACTED]"


class SecretLeakError(ValueError):
    pass


class ConfirmationRequiredError(ValueError):
    pass


def redact_secrets(value: Any) -> Any:
    if isinstance(value, Mapping):
        result: dict[Any, Any] = {}
        for key, item in value.items():
            if _is_secret_key(str(key)):
                result[key] = REDACTED
            else:
                result[key] = redact_secrets(item)
        return result

    if isinstance(value, list):
        return [redact_secrets(item) for item in value]

    if isinstance(value, tuple):
        return tuple(redact_secrets(item) for item in value)

    return value


def assert_no_secrets(value: Any) -> None:
    path = _find_secret_path(value)
    if path:
        raise SecretLeakError(f"secret field present in MCP response: {path}")


def require_confirmation(confirm: bool, action: str) -> None:
    if not confirm:
        raise ConfirmationRequiredError(f"confirmation required for action: {action}")


def _find_secret_path(value: Any, prefix: str = "") -> str | None:
    if isinstance(value, Mapping):
        for key, item in value.items():
            key_str = str(key)
            path = f"{prefix}.{key_str}" if prefix else key_str
            if _is_secret_key(key_str):
                return path
            nested = _find_secret_path(item, path)
            if nested:
                return nested

    elif isinstance(value, Sequence) and not isinstance(value, (str, bytes, bytearray)):
        for idx, item in enumerate(value):
            nested = _find_secret_path(item, f"{prefix}[{idx}]")
            if nested:
                return nested

    return None


def _is_secret_key(key: str) -> bool:
    normalized = key.lower().replace("-", "_")
    return normalized in SECRET_FIELD_NAMES or normalized.endswith("_secret") or normalized.endswith("_token")
