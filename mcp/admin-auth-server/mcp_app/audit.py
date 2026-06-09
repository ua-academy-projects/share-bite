import logging
import re
from typing import Any
from mcp_app.config import settings

class RedactingFilter(logging.Filter):
    def filter(self, record: logging.LogRecord) -> bool:
        msg = str(record.msg)
        msg = re.sub(r'(Bearer\s+)[A-Za-z0-9\-._~+/]+=*', r'\1[REDACTED]', msg, flags=re.IGNORECASE)
        msg = re.sub(r'("auth_token"\s*:\s*")[^"]+(")', r'\1[REDACTED]\2', msg, flags=re.IGNORECASE)
        msg = re.sub(r'(token=)[A-Za-z0-9\-._~+/]+=*', r'\1[REDACTED]', msg, flags=re.IGNORECASE)
        msg = re.sub(r'("password"\s*:\s*")[^"]+(")', r'\1[REDACTED]\2', msg, flags=re.IGNORECASE)

        record.msg = msg
        return True

def redact_sensitive_payload(data: Any) -> Any:
    if isinstance(data, dict):
        return {
            k: "[REDACTED]" if k.lower() in ["password", "token", "secret", "auth_token", "authorization", "jwt"]
            else redact_sensitive_payload(v) for k, v in data.items()
        }
    elif isinstance(data, list):
        return [redact_sensitive_payload(item) for item in data]
    return data

logger = logging.getLogger("mcp_server")
logger.setLevel(logging.INFO)

audit_logger = logging.getLogger("mcp_audit")
audit_logger.setLevel(logging.INFO)

if logger.handlers:
    logger.handlers.clear()
if audit_logger.handlers:
    audit_logger.handlers.clear()

security_filter = RedactingFilter()

console_handler = logging.StreamHandler()
console_handler.addFilter(security_filter)
console_formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
console_handler.setFormatter(console_formatter)
logger.addHandler(console_handler)

file_handler = logging.FileHandler(settings.audit_log_destination, encoding="utf-8")
file_handler.addFilter(security_filter)
file_formatter = logging.Formatter('%(asctime)s - [AUDIT_TRACK] - %(message)s')
file_handler.setFormatter(file_formatter)
audit_logger.addHandler(file_handler)


def log_audit_event(action: str, admin_id: str | None, status: str, details: str) -> None:
    safe_details = redact_sensitive_payload(details) if isinstance(details, (dict, list)) else details
    raw_log = f"Action: {action} | AdminID: {admin_id or 'UNAUTHORIZED'} | Status: {status} | Details: {safe_details}"
    audit_logger.info(raw_log)
    logger.info(f"Audit Log Captured -> Action: {action}, Status: {status}")