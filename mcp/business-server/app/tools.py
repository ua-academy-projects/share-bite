from __future__ import annotations

from typing import Any

from mcp.server.fastmcp import FastMCP, Context
from app.auth import resolve_auth_token

from app.client import BusinessApiClient, BusinessApiError
from app.config import Settings


def register_tools(mcp: FastMCP, settings: Settings, client: BusinessApiClient) -> None:

    @mcp.tool()
    async def business_health_check(
        ctx: Context,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Check whether the ShareBite business-api is reachable.

        auth_token is forwarded to business-api when provided.
        request_id is propagated as X-Request-ID when provided.
        """
        headers = _extract_headers(ctx)
        final_token = resolve_auth_token(headers=headers)

        if not final_token:
            raise RuntimeError("Unauthorized: Missing authentication token")

        try:
            data = await client.get(
                "/swagger/doc.json",
                auth_token=final_token,
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
        ctx: Context,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Return basic business-api status and OpenAPI metadata.

        This tool does not infer or guess a business ID.
        """

        headers = _extract_headers(ctx)

        final_token = resolve_auth_token(headers=headers)

        if not final_token:
            raise RuntimeError("Unauthorized: Missing authentication token")

        try:
            data = await client.get(
                "/swagger/doc.json",
                auth_token=final_token,
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


def _extract_headers(ctx: Context) -> dict[str, str]:
    """
    Extract headers from MCP context, regardless of type of object.
    """
    if not ctx.request_context or not ctx.request_context.meta:
        return {}

    meta = ctx.request_context.meta

    if hasattr(meta, "model_dump"):
        meta_dict = meta.model_dump()
    elif hasattr(meta, "dict"):
        meta_dict = meta.dict()
    elif isinstance(meta, dict):
        meta_dict = meta 
    else:
        try:
            meta_dict = vars(meta)
        except TypeError:
            return {}

    headers = meta_dict.get("headers", {})
    return headers if isinstance(headers, dict) else {}