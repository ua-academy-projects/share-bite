import logging
import re
from typing import Any

class RedactingFilter(logging.Filter):
    def filter(self, record: logging.LogRecord) -> bool:
        msg = str(record.getMessage())
        msg = re.sub(r'(Bearer\s+)[A-Za-z0-9\-._~+/]+=*', r'\1[REDACTED]', msg, flags=re.IGNORECASE)
        msg = re.sub(r'("password"\s*:\s*")[^"]+(")', r'\1[REDACTED]\2', msg, flags=re.IGNORECASE)
        msg = re.sub(r'(token=)[A-Za-z0-9\-._~+/]+=*', r'\1[REDACTED]', msg, flags=re.IGNORECASE)
        msg = re.sub(r'("?[A-Za-z_]*token"?\s*[:=]\s*")[^"]+(")', r'\1[REDACTED]\2', msg, flags=re.IGNORECASE)

        record.msg = msg
        record.args = ()
        return True

def redact_sensitive_payload(data: Any) -> Any:
    if isinstance(data, dict):
        sensitive_markers = ("password", "token", "secret", "authorization", "jwt", "api_key")
        return {
            k: "[REDACTED]" if any(marker in k.lower() for marker in sensitive_markers)
            else redact_sensitive_payload(v) for k, v in data.items()
        }
    elif isinstance(data, list):
        return [redact_sensitive_payload(item) for item in data]
    elif isinstance(data, str):
        data = re.sub(r'(Bearer\s+)[A-Za-z0-9\-._~+/]+=*', r'\1[REDACTED]', data, flags=re.IGNORECASE)
        data = re.sub(r'(token=)[A-Za-z0-9\-._~+/]+=*', r'\1[REDACTED]', data, flags=re.IGNORECASE)
        data = re.sub(r'("?[A-Za-z_]*token"?\s*[:=]\s*")[^"]+(")', r'\1[REDACTED]\2', data, flags=re.IGNORECASE)
        return data
    return data

def setup_logger(name: str = "admin_auth") -> logging.Logger:
    app_logger = logging.getLogger(name)
    app_logger.setLevel(logging.INFO)

    if not app_logger.handlers:
        handler = logging.StreamHandler()
        handler.addFilter(RedactingFilter())
        formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
        handler.setFormatter(formatter)
        app_logger.addHandler(handler)

    return app_logger

logger = setup_logger()