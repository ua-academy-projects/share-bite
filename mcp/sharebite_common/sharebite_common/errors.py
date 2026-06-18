from __future__ import annotations

from dataclasses import dataclass
from typing import Any


@dataclass(frozen=True)
class MCPError:
    ok: bool
    error: str
    status: int | None = None
    details: Any | None = None

    def to_dict(self) -> dict[str, Any]:
        payload: dict[str, Any] = {"ok": self.ok, "error": self.error}
        if self.status is not None:
            payload["status"] = self.status
        if self.details is not None:
            payload["details"] = self.details
        return payload


def map_http_error(status: int | None, message: str | None = None) -> MCPError:
    safe_message = _safe_message(status)
    details = message if status is not None and 400 <= status < 500 else None
    return MCPError(ok=False, error=safe_message, status=status, details=details)


def _safe_message(status: int | None) -> str:
    if status == 400:
        return "validation failed"
    if status == 401:
        return "unauthorized"
    if status == 403:
        return "forbidden"
    if status == 404:
        return "not found"
    if status == 408:
        return "upstream request timed out"
    if status == 409:
        return "conflict"
    if status == 429:
        return "rate limited"
    if status is not None and status >= 500:
        return "upstream service error"
    return "upstream request failed"
