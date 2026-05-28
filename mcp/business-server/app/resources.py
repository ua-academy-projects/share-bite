from __future__ import annotations

from typing import Any

from mcp.server.fastmcp import FastMCP

from app.client import BusinessApiClient, BusinessApiError
from app.config import Settings


def register_resources(mcp: FastMCP, settings: Settings) -> None:
    client = BusinessApiClient(
        base_url=settings.business_api_base_url,
        timeout_seconds=settings.request_timeout_seconds,
    )

    @mcp.resource("sharebite://business/api-info")
    def business_api_info() -> dict[str, Any]:
        return {
            "service": "business-api",
            "base_url": settings.business_api_base_url,
            "timeout_seconds": settings.request_timeout_seconds,
            "auth": {
                "authorization_forwarding": True,
                "business_id_policy": "Business ID is never guessed; pass it explicitly or derive it from authenticated context in tools that require it.",
            },
            "request_id": {
                "header": "X-Request-ID",
                "propagation": True,
            },
        }

    @mcp.resource("sharebite://business/openapi-summary")
    async def business_openapi_summary() -> dict[str, Any]:
        try:
            data = await client.get("/swagger/doc.json")
        except BusinessApiError as exc:
            raise RuntimeError(str(exc)) from exc

        info = data.get("info", {})
        paths = data.get("paths", {})

        return {
            "title": info.get("title"),
            "version": info.get("version"),
            "description": info.get("description"),
            "base_url": settings.business_api_base_url,
            "path_count": len(paths) if isinstance(paths, dict) else None,
            "paths": sorted(paths.keys()) if isinstance(paths, dict) else [],
        }