from __future__ import annotations

from typing import Any

from mcp.server.fastmcp import FastMCP
from pydantic import BaseModel

from app.client import BusinessApiClient, BusinessApiError
from app.config import Settings
from app.constants import URI_PROFILE_SCHEMA, URI_VENUE_HOURS_FORMAT, URI_VENUE_SCHEMA

class AuthInfo(BaseModel):
    authorization_forwarding: bool
    business_id_policy: str


class RequestIdInfo(BaseModel):
    header: str
    propagation: bool


class BusinessApiInfoResource(BaseModel):
    service: str
    base_url: str
    timeout_seconds: float
    auth: AuthInfo
    request_id: RequestIdInfo


class BusinessOpenApiSummaryResource(BaseModel):
    title: str | None
    version: str | None
    description: str | None
    base_url: str
    path_count: int | None
    paths: list[str]


def register_resources(mcp: FastMCP, settings: Settings) -> None:
    client = BusinessApiClient(
        base_url=settings.business_api_base_url,
        timeout_seconds=settings.request_timeout_seconds,
    )

    @mcp.resource("sharebite://business/api-info")
    def business_api_info() -> dict[str, Any]:
        return BusinessApiInfoResource(
            service="business-api",
            base_url=settings.business_api_base_url,
            timeout_seconds=settings.request_timeout_seconds,
            auth=AuthInfo(
                authorization_forwarding=True,
                business_id_policy="Business ID is never guessed; pass it explicitly or derive it from authenticated context in tools that require it.",
            ),
            request_id=RequestIdInfo(
                header="X-Request-ID",
                propagation=True,
            ),
        ).model_dump()

    @mcp.resource("sharebite://business/openapi-summary")
    async def business_openapi_summary() -> dict[str, Any]:
        try:
            data = await client.get("/swagger/doc.json")
        except BusinessApiError as exc:
            raise RuntimeError(str(exc)) from exc

        info_raw = data.get("info")
        info = info_raw if isinstance(info_raw, dict) else {}
        paths = data.get("paths", {})

        return BusinessOpenApiSummaryResource(
            title=info.get("title"),
            version=info.get("version"),
            description=info.get("description"),
            base_url=settings.business_api_base_url,
            path_count=len(paths) if isinstance(paths, dict) else None,
            paths=sorted(paths.keys()) if isinstance(paths, dict) else [],
        ).model_dump()

    @mcp.resource(URI_PROFILE_SCHEMA)
    def business_profile_schema() -> dict[str, Any]:
        return {
            "type": "object",
            "properties": {
                "name": {"type": "string", "minLength": 3, "maxLength": 40},
                "avatar": {"type": ["string", "null"]},
                "banner": {"type": ["string", "null"]},
                "description": {"type": ["string", "null"]},
            },
            "additionalProperties": False,
        }


    @mcp.resource(URI_VENUE_SCHEMA)
    def business_venue_schema() -> dict[str, Any]:
        return {
            "type": "object",
            "properties": {
                "name": {"type": ["string", "null"]},
                "avatar": {"type": ["string", "null"]},
                "banner": {"type": ["string", "null"]},
                "description": {"type": ["string", "null"]},
                "latitude": {"type": ["number", "null"], "minimum": -90, "maximum": 90},
                "longitude": {"type": ["number", "null"], "minimum": -180, "maximum": 180},
                "tagIds": {"type": ["array", "null"], "items": {"type": "integer"}, "maxItems": 5},
            },
            "additionalProperties": False,
        }


    @mcp.resource(URI_VENUE_HOURS_FORMAT)
    def business_venue_hours_format() -> dict[str, Any]:
        return {
            "type": "object",
            "required": ["days"],
            "properties": {
                "days": {
                    "type": "array",
                    "minItems": 1,
                    "maxItems": 7,
                    "items": {
                        "type": "object",
                        "required": ["weekday", "openTime", "closeTime"],
                        "properties": {
                            "weekday": {"type": "integer", "minimum": 1, "maximum": 7},
                            "openTime": {"type": ["string", "null"], "pattern": "^\\d{2}:\\d{2}$"},
                            "closeTime": {"type": ["string", "null"], "pattern": "^\\d{2}:\\d{2}$"},
                        },
                        "additionalProperties": False,
                    },
                }
            },
            "additionalProperties": False,
            "example": {
                "days": [
                    {"weekday": 1, "openTime": "09:00", "closeTime": "18:00"},
                    {"weekday": 7, "openTime": None, "closeTime": None},
                ]
            },
        }