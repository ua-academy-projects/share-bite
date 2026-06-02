from __future__ import annotations

from typing import Any

from mcp.server.fastmcp import FastMCP

from app.client import BusinessApiClient, BusinessApiError
from app.config import Settings


def register_tools(mcp: FastMCP, settings: Settings) -> None:
    client = BusinessApiClient(
        base_url=settings.business_api_base_url,
        timeout_seconds=settings.request_timeout_seconds,
    )

    @mcp.tool()
    async def business_health_check(
        auth_token: str | None = None,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Check whether the ShareBite business-api is reachable.

        auth_token is forwarded to business-api when provided.
        request_id is propagated as X-Request-ID when provided.
        """
        try:
            data = await client.get(
                "/swagger/doc.json",
                auth_token=auth_token,
                request_id=request_id,
            )
        except BusinessApiError as exc:
            raise RuntimeError(str(exc)) from exc

        info_raw = data.get("info")
        info = info_raw if isinstance(info_raw, dict) else {}

        return {
            "ok": True,
            "service": "business-api",
            "base_url": settings.business_api_base_url,
            "title": info.get("title"),
            "version": info.get("version"),
        }

    @mcp.tool()
    async def get_business_api_status(
        auth_token: str | None = None,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Return basic business-api status and OpenAPI metadata.

        This tool does not infer or guess a business ID.
        """
        try:
            data = await client.get(
                "/swagger/doc.json",
                auth_token=auth_token,
                request_id=request_id,
            )
        except BusinessApiError as exc:
            raise RuntimeError(str(exc)) from exc

        info_raw = data.get("info")
        info = info_raw if isinstance(info_raw, dict) else {}
        paths = data.get("paths", {})

        return {
            "ok": True,
            "service": "business-api",
            "base_url": settings.business_api_base_url,
            "api": {
                "title": info.get("title"),
                "version": info.get("version"),
                "description": info.get("description"),
                "path_count": len(paths) if isinstance(paths, dict) else None,
            },
        }
