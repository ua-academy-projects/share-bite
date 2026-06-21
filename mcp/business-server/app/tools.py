from __future__ import annotations

from typing import Any
import base64
import io

from mcp.server.fastmcp import Context, FastMCP

from app.auth import resolve_auth_token
from app.client import BusinessApiClient, BusinessApiError
from app.config import Settings
from app.constants import (
    API_PATH_BUSINESS_PROFILE,
    API_PATH_BUSINESS_VENUES,
    API_PATH_DAILY_SUMMARY,
    API_PATH_RESERVATION_SUMMARY,
    API_PATH_FOOD_BOX_PERFORMANCE,
    API_PATH_ENGAGEMENT_SUMMARY,
    API_PATH_VENUE_ACTIVITY,
    
    TOOL_PREVIEW_VENUE_HOURS_UPDATE,
    API_PATH_NEARBY_BOXES,
    API_PATH_NEARBY_VENUES,
    API_PATH_RECOMMEND_POSTS,
    API_PATH_SEARCH_VENUES,
    API_PATH_UPDATE_VENUE_DETAILS,
    API_PATH_UPDATE_VENUE_HOURS,
    API_PATH_VENUE_DETAILS,
    API_PATH_BUSINESS_FOOD_BOXES,
    API_PATH_FOOD_BOX,
    API_PATH_FOOD_BOX_RESERVATIONS,
)
from app.context_recommender import (
    recommend_venues_by_context as rank_venues_by_context,
)
from app.utils import (
    ForbiddenError,
    changed_fields,
    ensure_venue_owned_by_business,
    validate_date_range,
    validate_discovery_coords,
    validate_pagination,
    validate_profile_update,
    validate_venue_hours,
    validate_venue_update,
    validate_date_range,
    build_venue_hours_preview,
    extract_venue_hours_days,
    ensure_food_box_owned_by_business,
    validate_profile_update,
    validate_venue_hours,
    validate_venue_update,
    validate_food_box_creation,
    validate_food_box_update,
)


def _extract_headers(ctx: Context[Any, Any]) -> dict[str, str]:
    """
    Extract headers from MCP context, regardless of type of object.
    """
    if not ctx.request_context or not ctx.request_context.meta:
        return {}

    meta = ctx.request_context.meta

    if hasattr(meta, "model_dump"):
        meta_dict = meta.model_dump()
    elif isinstance(meta, dict):
        meta_dict = meta
    else:
        try:
            meta_dict = vars(meta)
        except TypeError:
            return {}

    headers = meta_dict.get("headers", {})
    return headers if isinstance(headers, dict) else {}


def register_tools(
    mcp: FastMCP,
    settings: Settings,
    client: BusinessApiClient | None = None,
) -> None:
    client = client or BusinessApiClient(
        base_url=settings.business_api_base_url,
        timeout_seconds=settings.request_timeout_seconds,
        api_token=settings.business_api_token,
    )

    @mcp.tool()
    async def business_health_check(
        ctx: Context[Any, Any],
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
                "/business/healthz",
                auth_token=final_token,
                request_id=request_id,
            )
        except BusinessApiError as exc:
            raise RuntimeError(str(exc)) from exc

        return {
            "ok": True,
            "service": "business-api",
            "base_url": settings.business_api_base_url,
        }

    @mcp.tool()
    async def get_business_api_status(
        ctx: Context[Any, Any],
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
        
    # BUSINESS-ORG TOOLS
    @mcp.tool()
    async def recommend_venues_by_context(
        context: dict[str, Any],
        venues: list[dict[str, Any]],
        limit: int = 10,
    ) -> dict[str, Any]:
        """
        Rank venue candidates for context such as a date, meetup, or budget plan.
        This deterministic AI-style ranker extracts user intent and returns
        explainable scores without calling external AI APIs.
        """
        ranked = rank_venues_by_context(context=context, venues=venues, limit=limit)
        return _tool_success(result={"items": ranked, "total": len(ranked)})

    @mcp.tool()
    async def get_business_profile(
        business_id: int,
        auth_token: str,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Retrieve a business profile.

        business_id must be provided explicitly by the caller.
        auth_token is forwarded to business-api for authorization.
        """
        try:
            data = await client.get(
                API_PATH_BUSINESS_PROFILE.format(business_id=business_id),
                auth_token=auth_token,
                request_id=request_id,
            )
            return _tool_success(
                business_id=business_id,
                result=_unwrap(data),
            )
        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def update_business_profile(
        business_id: int,
        payload: dict[str, Any],
        auth_token: str,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Update a business profile after validating allowed fields.

        Side effects: updates business profile fields in Business API.
        Required permission: business:profile:update.
        Auth: forwards caller bearer token to Business API.
        """
        validation_errors = validate_profile_update(payload)
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        try:
            before = _unwrap(
                await client.get(
                    API_PATH_BUSINESS_PROFILE.format(business_id=business_id),
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )
            after = _unwrap(
                await client.patch(
                    API_PATH_BUSINESS_PROFILE.format(business_id=business_id),
                    json_data=payload,
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )

            return _tool_success(
                business_id=business_id,
                changed_fields=changed_fields(
                    before,
                    after,
                    ("name", "avatar", "banner", "description"),
                ),
                result=after,
            )
        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def list_business_venues(
        business_id: int,
        skip: int = 0,
        limit: int = 10,
        auth_token: str | None = None,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        List venues for a business.

        business_id must be provided explicitly by the caller.
        """
        try:
            data = await client.get(
                API_PATH_BUSINESS_VENUES.format(business_id=business_id),
                auth_token=auth_token,
                request_id=request_id,
                params={
                    "skip": max(skip, 0),
                    "limit": max(1, min(limit, 100)),
                },
            )
            return _tool_success(
                business_id=business_id,
                result=_unwrap(data),
            )
        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def get_venue_details(
        business_id: int,
        venue_id: int,
        auth_token: str,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """Get venue details after verifying that the venue belongs to business_id."""
        try:
            venue = _unwrap(
                await client.get(
                    API_PATH_VENUE_DETAILS.format(venue_id=venue_id),
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )
            ensure_venue_owned_by_business(venue, business_id)

            return _tool_success(
                business_id=business_id,
                venue_id=venue_id,
                result=venue,
            )
        except (BusinessApiError, ForbiddenError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def update_venue_details(
        business_id: int,
        venue_id: int,
        payload: dict[str, Any],
        auth_token: str,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Update venue details after validation and ownership check.

        Side effects: updates venue fields in Business API.
        Required permission: business:venue:update.
        Auth: forwards caller bearer token to Business API.
        """
        validation_errors = validate_venue_update(payload)
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        try:
            before = _unwrap(
                await client.get(
                    API_PATH_VENUE_DETAILS.format(venue_id=venue_id),
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )
            ensure_venue_owned_by_business(before, business_id)

            after = _unwrap(
                await client.patch(
                    API_PATH_UPDATE_VENUE_DETAILS.format(venue_id=venue_id),
                    json_data=payload,
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )

            return _tool_success(
                business_id=business_id,
                venue_id=venue_id,
                changed_fields=changed_fields(
                    before,
                    after,
                    (
                        "name",
                        "avatar",
                        "banner",
                        "description",
                        "latitude",
                        "longitude",
                        "tagIds",
                    ),
                ),
                result=after,
            )
        except (BusinessApiError, ForbiddenError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def update_venue_hours(
        business_id: int,
        venue_id: int,
        payload: dict[str, Any],
        auth_token: str,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Update venue hours after validation and ownership check.

        Side effects: replaces venue weekly opening hours in Business API.
        Required permission: business:venue-hours:update.
        Auth: forwards caller bearer token to Business API.
        """
        validation_errors = validate_venue_hours(payload)
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        try:
            venue = _unwrap(
                await client.get(
                    API_PATH_VENUE_DETAILS.format(venue_id=venue_id),
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )
            ensure_venue_owned_by_business(venue, business_id)

            after = _unwrap(
                await client.patch(
                    API_PATH_UPDATE_VENUE_HOURS.format(venue_id=venue_id),
                    json_data=payload,
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )

            return _tool_success(
                business_id=business_id,
                venue_id=venue_id,
                changed_fields=["days"],
                result=after,
            )
        except (BusinessApiError, ForbiddenError, RuntimeError) as exc:
            return _tool_error(str(exc))
    
    @mcp.tool()
    async def search_venues(
        ctx: Context[Any, Any],
        q: str | None = None,
        tags: str | None = None,
        skip: int = 0,
        limit: int = 10,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Search venues by keyword and/or location tags.
        Anonymous access allowed; Bearer token forwarded if present.
        """
        headers = _extract_headers(ctx)
        final_token = resolve_auth_token(headers=headers)

        validation_errors = validate_pagination(skip, limit)
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        query = (q or "").strip()
        tags = (tags or "").strip()
        if not query and not tags:
            return _tool_error("at least one search filter is required: q or tags")

        params: dict[str, Any] = {
            "skip": max(skip, 0),
            "limit": max(1, min(limit, 100)),
        }
        if query:
            params["q"] = query
        if tags:
            params["tags"] = tags

        try:
            data = await client.get(
                API_PATH_SEARCH_VENUES,
                auth_token=final_token,
                request_id=request_id,
                params=params,
            )
            return _tool_success(result=_unwrap(data))
        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def get_recommended_venues(
        ctx: Context[Any, Any],
        lat: float,
        lon: float,
        skip: int = 0,
        limit: int = 10,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        List nearby venues sorted by distance from user coordinates.
        Anonymous access allowed; Bearer token forwarded if present.
        """
        headers = _extract_headers(ctx)
        final_token = resolve_auth_token(headers=headers)

        validation_errors = validate_discovery_coords(lat, lon) + validate_pagination(
            skip, limit
        )
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        params: dict[str, Any] = {
            "lat": lat,
            "lon": lon,
            "skip": max(skip, 0),
            "limit": max(1, min(limit, 100)),
        }

        try:
            data = await client.get(
                API_PATH_NEARBY_VENUES,
                auth_token=final_token,
                request_id=request_id,
                params=params,
            )
            return _tool_success(result=_unwrap(data))
        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool(name=TOOL_PREVIEW_VENUE_HOURS_UPDATE)
    async def preview_venue_hours_update(
        business_id: int,
        venue_id: int,
        payload: dict[str, Any],
        auth_token: str,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """Preview venue-hours changes without performing a write."""
        validation_errors = validate_venue_hours(payload)
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        try:
            venue = _unwrap(
                await client.get(
                    API_PATH_VENUE_DETAILS.format(venue_id=venue_id),
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )
            ensure_venue_owned_by_business(venue, business_id)

            current_days = extract_venue_hours_days(venue)
            if current_days is None:
                return _tool_error(
                    "current venue hours are unavailable for preview"
                )

            preview = build_venue_hours_preview(
                current_days=current_days,
                proposed_days=payload.get("days"),
            )

            return _tool_success(
                business_id=business_id,
                venue_id=venue_id,
                changed_fields=preview["changed_fields"],
                result={
                    "current_days": preview["current_days"],
                    "preview_days": preview["preview_days"],
                    "day_changes": preview["day_changes"],
                    "summary": preview["summary"],
                },
            )
        except (BusinessApiError, ForbiddenError, RuntimeError) as exc:
            return _tool_error(str(exc))

        
    @mcp.tool()
    async def get_feed_items(
        ctx: Context[Any, Any],
        lat: float,
        lon: float,
        skip: int = 0,
        limit: int = 10,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Get personalized post recommendations based on user behavior and location.
        Requires Bearer authentication.
        """
        headers = _extract_headers(ctx)
        final_token = resolve_auth_token(headers=headers)

        if not final_token:
            return _tool_error("Unauthorized: Missing authentication token")

        validation_errors = validate_discovery_coords(lat, lon) + validate_pagination(
            skip, limit
        )
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        params: dict[str, Any] = {
            "lat": lat,
            "lon": lon,
            "skip": max(skip, 0),
            "limit": max(1, min(limit, 100)),
        }

        try:
            data = await client.get(
                API_PATH_RECOMMEND_POSTS,
                auth_token=final_token,
                request_id=request_id,
                params=params,
            )
            return _tool_success(result=_unwrap(data))
        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def search_boxes(
        ctx: Context[Any, Any],
        lat: float,
        lon: float,
        skip: int = 0,
        limit: int = 10,
        org_id: int | None = None,
        category_id: int | None = None,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        List available food boxes sorted by distance from user coordinates.
        Anonymous access allowed; Bearer token forwarded if present.
        Optional filters: org_id, category_id.
        """
        headers = _extract_headers(ctx)
        final_token = resolve_auth_token(headers=headers)

        validation_errors = validate_discovery_coords(lat, lon) + validate_pagination(
            skip, limit
        )
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        params: dict[str, Any] = {
            "lat": lat,
            "lon": lon,
            "skip": max(skip, 0),
            "limit": max(1, min(limit, 100)),
        }
        if org_id is not None:
            params["org_id"] = org_id
        if category_id is not None:
            params["category_id"] = category_id

        try:
            data = await client.get(
                API_PATH_NEARBY_BOXES,
                auth_token=final_token,
                request_id=request_id,
                params=params,
            )
            return _tool_success(result=_unwrap(data))
        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    # BUSINESS ANALYTICS TOOLS
    @mcp.tool()
    async def get_business_daily_summary(
        ctx: Context[Any, Any],
        start_date: str,
        end_date: str,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Retrieve daily summary analytics for the brand.
        Returns the number of created boxes, posts, and total venues.

        Args:
            start_date: Start date for the filter in YYYY-MM-DD format.
            end_date: End date for the filter in YYYY-MM-DD format.
        """

        date_error = validate_date_range(start_date, end_date)
        if date_error:
            return _tool_error(f"Validation failed: {date_error}")

        headers = _extract_headers(ctx)
        auth_token = resolve_auth_token(headers=headers)

        if not auth_token:
            return _tool_error("Unauthorized: Missing authentication token")

        try:
            data = await client.get(
                API_PATH_DAILY_SUMMARY,
                auth_token=auth_token,
                request_id=request_id,
                params={
                    "start_date": start_date,
                    "end_date": end_date,
                },
            )
            return _tool_success(result=_unwrap(data))
        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def get_reservation_summary(
        ctx: Context[Any, Any],
        start_date: str,
        end_date: str,
        venue_id: int | None = None,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Retrieve reservation summary analytics (items sold, reserved, available, and potential revenue).

        Args:
            start_date: Start date for the filter in YYYY-MM-DD format.
            end_date: End date for the filter in YYYY-MM-DD format.
            venue_id: Optional ID of a specific venue. If omitted, returns stats for the entire brand.
        """

        date_error = validate_date_range(start_date, end_date)
        if date_error:
            return _tool_error(f"Validation failed: {date_error}")

        headers = _extract_headers(ctx)
        auth_token = resolve_auth_token(headers=headers)

        if not auth_token:
            return _tool_error("Unauthorized: Missing authentication token")

        params: dict[str, Any] = {
            "start_date": start_date,
            "end_date": end_date,
        }

        if venue_id is not None:
            params["venue_id"] = venue_id

        try:
            data = await client.get(
                API_PATH_RESERVATION_SUMMARY,
                auth_token=auth_token,
                request_id=request_id,
                params=params,
            )

            return _tool_success(result=_unwrap(data), venue_id=venue_id)
        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def get_food_box_performance(
        ctx: Context[Any, Any],
        start_date: str,
        end_date: str,
        venue_id: int | None = None,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Retrieve food box performance metrics for the business.
        Returns data on created boxes, expired boxes, average discount, sell-through rate, waste rate, and a composite performance score.

        Args:
            start_date: Start date for the filter in YYYY-MM-DD format.
            end_date: End date for the filter in YYYY-MM-DD format.
            venue_id: Optional ID of a specific venue. If omitted, returns stats for the entire brand.
        """

        date_error = validate_date_range(start_date, end_date)
        if date_error:
            return _tool_error(f"Validation failed: {date_error}")

        headers = _extract_headers(ctx)
        auth_token = resolve_auth_token(headers=headers)

        if not auth_token:
            return _tool_error("Unauthorized: Missing authentication token")

        params: dict[str, Any] = {
            "start_date": start_date,
            "end_date": end_date,
        }
        if venue_id is not None:
            params["venue_id"] = venue_id

        try:
            data = await client.get(
                API_PATH_FOOD_BOX_PERFORMANCE,
                auth_token=auth_token,
                request_id=request_id,
                params=params,
            )

            return _tool_success(result=_unwrap(data), venue_id=venue_id)

        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def get_engagement_summary(
        ctx: Context[Any, Any],
        start_date: str,
        end_date: str,
        venue_id: int | None = None,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Retrieve social engagement metrics for business posts.
        Returns total posts created, total comments, total likes, and average engagement rates.

        Args:
            start_date: Start date for the filter in YYYY-MM-DD format.
            end_date: End date for the filter in YYYY-MM-DD format.
            venue_id: Optional ID of a specific venue. If omitted, returns stats for the entire brand.
        """

        date_error = validate_date_range(start_date, end_date)
        if date_error:
            return _tool_error(f"Validation failed: {date_error}")

        headers = _extract_headers(ctx)
        auth_token = resolve_auth_token(headers=headers)

        if not auth_token:
            return _tool_error("Unauthorized: Missing authentication token")

        params: dict[str, Any] = {
            "start_date": start_date,
            "end_date": end_date,
        }

        if venue_id is not None:
            params["venue_id"] = venue_id

        try:
            data = await client.get(
                API_PATH_ENGAGEMENT_SUMMARY,
                auth_token=auth_token,
                request_id=request_id,
                params=params,
            )

            return _tool_success(result=_unwrap(data), venue_id=venue_id)

        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def get_venue_activity_summary(
        ctx: Context[Any, Any],
        start_date: str,
        end_date: str,
        venue_id: int,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Retrieve operational activity summary for a specific venue.
        Returns the number of created boxes, posts, and the venue's name.

        Args:
            start_date: Start date for the filter in YYYY-MM-DD format.
            end_date: End date for the filter in YYYY-MM-DD format.
            venue_id: Required ID of the specific venue to query.
        """

        date_error = validate_date_range(start_date, end_date)
        if date_error:
            return _tool_error(f"Validation failed: {date_error}")

        headers = _extract_headers(ctx)
        auth_token = resolve_auth_token(headers=headers)

        if not auth_token:
            return _tool_error("Unauthorized: Missing authentication token")

        try:
            data = await client.get(
                API_PATH_VENUE_ACTIVITY.format(venue_id=venue_id),
                auth_token=auth_token,
                request_id=request_id,
                params={
                    "start_date": start_date,
                    "end_date": end_date,
                },
            )

            return _tool_success(result=_unwrap(data), venue_id=venue_id)

        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def list_business_food_boxes(
        business_id: int,
        skip: int = 0,
        limit: int = 10,
        auth_token: str | None = None,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        List food boxes for a business.
        business_id must be provided explicitly by the caller.
        """
        try:
            data = await client.get(
                API_PATH_BUSINESS_FOOD_BOXES.format(business_id=business_id),
                auth_token=auth_token,
                request_id=request_id,
                params={
                    "skip": max(skip, 0),
                    "limit": max(1, min(limit, 100)),
                },
            )
            return _tool_success(
                business_id=business_id,
                result=_unwrap(data),
            )
        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))
        
    @mcp.tool()
    async def get_food_box(
        business_id: int,
        food_box_id: int,
        auth_token: str,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Get food box details after verifying that the food box belongs to business_id.
        """
        try:
            food_box = _unwrap(
                await client.get(
                    API_PATH_FOOD_BOX.format(food_box_id=food_box_id),
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )
            ensure_food_box_owned_by_business(food_box, business_id)

            return _tool_success(
                business_id=business_id,
                food_box_id=food_box_id,
                result=food_box,
            )
        except (BusinessApiError, ForbiddenError, RuntimeError) as exc:
            return _tool_error(str(exc))


    @mcp.tool()
    async def create_food_box(
        business_id: int,
        payload: dict[str, Any],
        auth_token: str,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Create a new food box with image upload.
        
        Required fields: venue_id, price_full, expires_at, quantity, image_base64
        Optional fields: category_id, price_discount
        
        image_base64 should be a base64-encoded image file (PNG, JPG, etc).
        """
        validation_errors = validate_food_box_creation(payload)
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        try:
            image_base64 = payload.pop("image_base64")
            
            if "," in image_base64:
                image_base64 = image_base64.split(",", 1)[1]
                
            try:
                image_bytes = base64.b64decode(image_base64, validate=True)
            except Exception as e:
                return _tool_error(f"invalid base64 encoding for image: {str(e)}")

            form_data = {}
            for key in ["venue_id", "price_full", "expires_at", "quantity", "category_id", "price_discount"]:
                if key in payload and payload[key] is not None:
                    form_data[key] = str(payload[key])

            form_data["image"] = ("image.jpg", image_bytes, "image/jpeg")

            result = _unwrap(
                await client.post_form(
                    API_PATH_BUSINESS_FOOD_BOXES.format(business_id=business_id),
                    form_data=form_data,
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )

            return _tool_success(
                business_id=business_id,
                result=result,
            )
        except (BusinessApiError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def update_food_box(
        business_id: int,
        food_box_id: int,
        payload: dict[str, Any],
        auth_token: str,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Update food box details with validation.
        
        Supported fields: category_id, price_full, price_discount, expires_at
        """
        validation_errors = validate_food_box_update(payload)
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        try:
            before = _unwrap(
                await client.get(
                    API_PATH_FOOD_BOX.format(food_box_id=food_box_id),
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )
            
            after = _unwrap(
                await client.patch(
                    API_PATH_FOOD_BOX.format(food_box_id=food_box_id),
                    json_data=payload,
                    auth_token=auth_token,
                    request_id=request_id,
                )
            )

            return _tool_success(
                business_id=business_id,
                food_box_id=food_box_id,
                changed_fields=changed_fields(
                    before,
                    after,
                    ("categoryId", "fullPrice", "discountPrice", "expiresAt"),
                ),
                result=after,
            )
        except (BusinessApiError, ForbiddenError, RuntimeError) as exc:
            return _tool_error(str(exc))

    @mcp.tool()
    async def get_food_box_reservations(
        business_id: int,
        food_box_id: int,
        skip: int = 0,
        limit: int = 10,
        auth_token: str | None = None,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Get reservations for a food box.
        """
        try:
            reservations = _unwrap(
                await client.get(
                    API_PATH_FOOD_BOX_RESERVATIONS.format(food_box_id=food_box_id),
                    auth_token=auth_token,
                    request_id=request_id,
                    params={
                        "skip": max(skip, 0),
                        "limit": max(1, min(limit, 100)),
                    },
                )
            )

            return _tool_success(
                business_id=business_id,
                food_box_id=food_box_id,
                result=reservations,
            )
        except (BusinessApiError, ForbiddenError, RuntimeError) as exc:
            return _tool_error(str(exc))
    
def _unwrap(result: dict[str, Any]) -> dict[str, Any]:
    if result.get("is_error") is True:
        message = result.get("error_message", "unknown business api error")
        raise RuntimeError(f"Business API error: {message}")

    data = result.get("data")
    if data is None:
        return result

    if isinstance(data, dict):
        return data

    return {"value": data}


def _tool_success(
    *,
    result: dict[str, Any],
    business_id: int | None = None,
    venue_id: int | None = None,
    food_box_id: int | None = None,
    changed_fields: list[str] | None = None,
) -> dict[str, Any]:
    response: dict[str, Any] = {
        "ok": True,
        "error": None,
        "validation_errors": [],
        "changed_fields": changed_fields or [],
        "result": result,
    }

    if business_id is not None:
        response["business_id"] = business_id

    if venue_id is not None:
        response["venue_id"] = venue_id

    if food_box_id is not None:
        response["food_box_id"] = food_box_id

    return response


def _tool_error(
    message: str,
    validation_errors: list[dict[str, str]] | None = None,
) -> dict[str, Any]:
    return {
        "ok": False,
        "error": message,
        "validation_errors": validation_errors or [],
        "changed_fields": [],
        "result": None,
    }
