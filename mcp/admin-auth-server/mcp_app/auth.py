from contextvars import ContextVar
from functools import wraps
from typing import Any, Callable
from mcp_app.http_client import admin_client
from mcp_app.audit import log_audit_event
from mcp_app.config import settings

current_user_ctx: ContextVar[dict[str, Any] | None] = ContextVar("current_user_ctx", default=None)

async def resolve_admin_context(auth_token: str | None = None) -> dict[str, Any]:
    result = await admin_client.get("/mcp/context", auth_token=auth_token)
    if result["is_error"]:
        raise ValueError(result["error_message"])
    return result["data"]

def require_role(required_role: str):
    def decorator(func: Callable):
        @wraps(func)
        async def wrapper(*args, **kwargs):
            auth_token = kwargs.get("auth_token")

            if not settings.enforce_authentication:
                mock_user = {"id": "local_dev", "role": "admin", "roles": ["admin"]}
                token_context = current_user_ctx.set(mock_user)
                try:
                    return await func(*args, **kwargs)
                finally:
                    current_user_ctx.reset(token_context)
            if not auth_token and not settings.local_refresh_token:
                log_audit_event(func.__name__, None, "DENIED", "No credentials available.")
                return "ERROR: No active session. Please provide an 'auth_token' or set a refresh token."

            try:
                user = await resolve_admin_context(auth_token)
            except ValueError as err:
                log_audit_event(func.__name__, None, "DENIED", f"Auth failed: {str(err)}")
                return f"ERROR: Critical Authentication Failure: {str(err)}. Please re-login."

            user_role = user.get("role")
            user_roles = user.get("roles", [])
            all_roles = [user_role] if user_role else []
            if isinstance(user_roles, list):
                all_roles.extend(user_roles)

            if "admin" not in all_roles and required_role not in all_roles:
                log_audit_event(func.__name__, user.get("id"), "DENIED", f"Requires role: {required_role}")
                return f"ERROR: Forbidden. This tool requires '{required_role}' role."

            token_context = current_user_ctx.set(user)
            try:
                return await func(*args, **kwargs)
            finally:
                current_user_ctx.reset(token_context)
        return wrapper
    return decorator