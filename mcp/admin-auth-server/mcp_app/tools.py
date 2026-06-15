import json
from fastmcp import FastMCP
from mcp_app.http_client import admin_client
from mcp_app.audit import log_audit_event
from mcp_app.auth import require_role, current_user_ctx

def register_tools(mcp: FastMCP) -> None:

    @mcp.tool(
        description=(
                "Checks connection health between Python MCP server and the Go backend API. "
                "Requires a valid 'auth_token' or active local session. Public/unauthenticated health check is NOT supported by the backend."
        )
    )
    @require_role("moderator")
    async def admin_auth_health_check(auth_token: str | None = None) -> str:
        user = current_user_ctx.get()
        admin_id = user.get("id") if user else "unknown"

        result = await admin_client.get("/mcp/health", auth_token=auth_token)
        if result["is_error"]:
            log_audit_event("admin_auth_health_check", admin_id, "ERROR", result["error_message"])
            return json.dumps({"status": "unhealthy", "error": result["error_message"]})

        log_audit_event("admin_auth_health_check", admin_id, "SUCCESS", "Infrastructure health cleared.")
        return json.dumps(result["data"])


    @mcp.tool(description="Retrieves basic session info for the currently logged-in administrator (ID, role, and status).")
    @require_role("moderator")
    async def get_current_admin_context(auth_token: str | None = None) -> str:
        user = current_user_ctx.get()
        admin_id = user.get("id") if user else "unknown"

        log_details = f"Context matrix returned. Explicit token provided: {bool(auth_token)}"

        log_audit_event("get_current_admin_context", admin_id, "SUCCESS", log_details)
        return json.dumps(user)


    @mcp.tool(description="Verifies if the current administrator has a specific permission for an action. Emits an audit log event.")
    @require_role("admin")
    async def validate_admin_permissions(target_permission: str, auth_token: str | None = None) -> str:
        user = current_user_ctx.get()
        admin_id = user.get("id") if user else "unknown"
        payload = {"permission": target_permission}
        result = await admin_client.post("/mcp/validate-permission", auth_token=auth_token, json=payload)

        if result["is_error"]:
            log_audit_event("validate_admin_permissions", admin_id, "ERROR", result["error_message"])
            return json.dumps({"authorized": False, "reason": result["error_message"]})

        log_audit_event("validate_admin_permissions", admin_id, "SUCCESS", f"Evaluated permission: {target_permission}")
        return json.dumps(result["data"])