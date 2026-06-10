from typing import Any
from fastmcp import FastMCP
from pydantic import BaseModel, Field
from mcp_app.constants import URI_ROLE_PERMISSIONS, URI_AUDIT_SCHEMA, CONTENT_TYPE_JSON

class RolePermissionsResource(BaseModel):
    roles: dict[str, dict[str, Any]]

class AuditSchemaResource(BaseModel):
    schema_url: str = Field(..., alias="$schema")
    title: str
    type: str
    required: list[str]
    properties: dict[str, Any]


ROLE_PERMISSIONS_MATRIX = {
    "roles": {
        "admin": {
            "description": "Адміністратор з повним доступом",
            "permissions": ["admin_auth_health_check", "get_current_admin_context", "validate_admin_permissions"]
        },
        "moderator": {
            "description": "Модератор з обмеженим адмін-доступом",
            "permissions": ["admin_auth_health_check", "get_current_admin_context"]
        }
    }
}

AUDIT_SCHEMA_BLUEPRINT = {
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "title": "AdminAuditEvent",
    "type": "object",
    "required": ["timestamp", "action", "admin_id", "status"],
    "properties": {
        "timestamp": {"type": "string", "format": "date-time"},
        "action": {"type": "string"},
        "admin_id": {"type": "string"},
        "status": {"type": "string", "enum": ["SUCCESS", "DENIED", "ERROR"]},
        "details": {"type": "string"}
    }
}


def register_resources(mcp: FastMCP) -> None:
    @mcp.resource(
        uri=URI_ROLE_PERMISSIONS,
        name="admin_role_permissions",
        title="Admin Role Permissions Matrix",
        description="Returns static matrix of admin and moderator security capabilities inside ShareBite.",
        mime_type=CONTENT_TYPE_JSON
    )
    def get_role_permissions() -> dict[str, Any]:
        return RolePermissionsResource(**ROLE_PERMISSIONS_MATRIX).model_dump(by_alias=True)


    @mcp.resource(
        uri=URI_AUDIT_SCHEMA,
        name="admin_audit_events_schema",
        title="Admin Audit Event JSON Schema",
        description="Returns the official JSON schema blueprint for decoding system audit trail entries.",
        mime_type=CONTENT_TYPE_JSON
    )
    def get_audit_events_schema() -> dict[str, Any]:
        return AuditSchemaResource(**AUDIT_SCHEMA_BLUEPRINT).model_dump(by_alias=True)