from __future__ import annotations

from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from typing import Any, Literal

from .security import redact_secrets

AuditStatus = Literal["SUCCESS", "DENIED", "ERROR"]


@dataclass(frozen=True)
class AuditEvent:
    timestamp: str
    action: str
    status: AuditStatus
    actor_id: str | None = None
    request_id: str | None = None
    details: Any | None = None

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


def make_audit_event(
    *,
    action: str,
    status: AuditStatus,
    actor_id: str | None = None,
    request_id: str | None = None,
    details: Any | None = None,
) -> AuditEvent:
    return AuditEvent(
        timestamp=datetime.now(timezone.utc).isoformat(),
        action=action,
        status=status,
        actor_id=actor_id,
        request_id=request_id,
        details=redact_secrets(details),
    )
