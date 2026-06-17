from __future__ import annotations

from uuid import uuid4

REQUEST_ID_HEADER = "X-Request-ID"


def get_or_create_request_id(request_id: str | None = None) -> str:
    if request_id and request_id.strip():
        return request_id.strip()
    return str(uuid4())
