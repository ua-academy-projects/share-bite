from __future__ import annotations

from typing import Any

from mcp.server.fastmcp import FastMCP
from pydantic import BaseModel

from app.client import BusinessApiClient, BusinessApiError
from app.config import Settings
from app.constants import (
    URI_PROFILE_SCHEMA,
    URI_VENUE_HOURS_FORMAT,
    URI_VENUE_SCHEMA,
    URI_ANALYTICS_METRICS,
    URI_REPORTING_PERIODS,
)


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


class AnalyticsMetricResource(BaseModel):
    name: str
    description: str
    guidelines: dict[str, str]


class AnalyticsGlossaryResource(BaseModel):
    title: str
    description: str
    metrics: list[AnalyticsMetricResource]


def register_resources(
    mcp: FastMCP, settings: Settings, client: BusinessApiClient
) -> None:
    """Resources registration"""

    @mcp.resource("sharebite://business/api-info")
    def business_api_info() -> dict[str, Any]:
        """Returns information about business API"""
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
                "longitude": {
                    "type": ["number", "null"],
                    "minimum": -180,
                    "maximum": 180,
                },
                "tagIds": {
                    "type": ["array", "null"],
                    "items": {"type": "integer"},
                    "maxItems": 5,
                },
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
                            "openTime": {
                                "type": ["string", "null"],
                                "pattern": "^\\d{2}:\\d{2}$",
                            },
                            "closeTime": {
                                "type": ["string", "null"],
                                "pattern": "^\\d{2}:\\d{2}$",
                            },
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

    @mcp.resource(URI_ANALYTICS_METRICS)
    def business_analytics_metrics() -> dict[str, Any]:
        """Returns the business analytics metrics and interpretation guidelines."""
        return AnalyticsGlossaryResource(
            title="Share-Bite Analytics Metrics",
            description="Guidelines for interpreting business performance metrics and engagement rates.",
            metrics=[
                AnalyticsMetricResource(
                    name="Sell-Through Rate",
                    description="Ratio of reserved items to total created items in boxes.",
                    guidelines={
                        "> 0.8 (80%)": "Excellent performance.",
                        "< 0.5 (50%)": "Requires review of box contents or pricing.",
                    },
                ),
                AnalyticsMetricResource(
                    name="Waste Rate",
                    description="Ratio of expired boxes to total created boxes.",
                    guidelines={
                        "< 0.1 (10%)": "Normal operational waste.",
                        "> 0.2 (20%)": "Critical level, business is losing revenue.",
                    },
                ),
                AnalyticsMetricResource(
                    name="Composite Score",
                    description="Score from 0 to 100 balancing sales and waste. Formula: (SellThroughRate * 100 * 0.7) + ((1 - WasteRate) * 100 * 0.3)",
                    guidelines={
                        "90-100": "Ideal",
                        "70-89": "Good",
                        "50-69": "Satisfactory",
                        "< 50": "Requires immediate intervention.",
                    },
                ),
                AnalyticsMetricResource(
                    name="Average Comments / Likes",
                    description="Average number of interactions per post (total interactions / total posts).",
                    guidelines={
                        "Context": "Used to assess audience loyalty. High engagement usually correlates with a better Sell-Through Rate."
                    },
                ),
                AnalyticsMetricResource(
                    name="Total Boxes / Posts Created",
                    description="Metrics representing the operational activity of the venue.",
                    guidelines={
                        "Context": "Zero values during an active business week may indicate staff negligence or technical issues with their application."
                    },
                ),
            ],
        ).model_dump()

    @mcp.resource(URI_REPORTING_PERIODS)
    def business_reporting_periods() -> dict[str, Any]:
        """Returns information about allowed reporting periods and date constraints."""
        return {
            "title": "Reporting Periods & Constraints",
            "description": "Rules and constraints for querying analytical data.",
            "constraints": {
                "max_days": 90,
                "format": "YYYY-MM-DD",
                "timezone": "UTC",
                "notes": "Queries exceeding 90 days will be rejected to prevent database overload.",
            },
        }
