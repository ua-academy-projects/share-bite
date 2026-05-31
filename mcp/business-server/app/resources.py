from __future__ import annotations

from typing import Any

from mcp.server.fastmcp import FastMCP
from pydantic import BaseModel

from app.client import BusinessApiClient, BusinessApiError
from app.config import Settings


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
