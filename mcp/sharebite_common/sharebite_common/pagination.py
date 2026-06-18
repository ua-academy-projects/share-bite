from __future__ import annotations

from dataclasses import dataclass


@dataclass(frozen=True)
class Pagination:
    limit: int
    offset: int

    def to_dict(self) -> dict[str, int]:
        return {"limit": self.limit, "offset": self.offset}


def clamp_pagination(
    limit: int | None,
    offset: int | None = 0,
    *,
    default_limit: int = 20,
    max_limit: int = 100,
) -> Pagination:
    if max_limit < 1:
        raise ValueError("max_limit must be >= 1")

    final_limit = default_limit if limit is None else limit
    final_offset = 0 if offset is None else offset

    if final_limit < 1:
        final_limit = 1
    if final_limit > max_limit:
        final_limit = max_limit
    if final_offset < 0:
        final_offset = 0

    return Pagination(limit=final_limit, offset=final_offset)
