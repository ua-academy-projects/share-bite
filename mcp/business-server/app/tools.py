from typing import Any, Protocol

from .constants import (
    API_PATH_BUSINESS_PROFILE,
    API_PATH_BUSINESS_VENUES,
    API_PATH_UPDATE_VENUE_DETAILS,
    API_PATH_UPDATE_VENUE_HOURS,
    API_PATH_VENUE_DETAILS,
    TOOL_GET_BUSINESS_PROFILE,
    TOOL_GET_VENUE_DETAILS,
    TOOL_LIST_BUSINESS_VENUES,
    TOOL_UPDATE_BUSINESS_PROFILE,
    TOOL_UPDATE_VENUE_DETAILS,
    TOOL_UPDATE_VENUE_HOURS,
)
from .utils import (
    AccessError,
    ForbiddenError,
    changed_fields,
    ensure_venue_owned_by_business,
    resolve_business_access,
    validate_profile_update,
    validate_venue_hours,
    validate_venue_update,
)


class BusinessAPIClient(Protocol):
    async def get(
        self,
        path: str,
        auth_token: str | None = None,
        request_id: str | None = None,
        params: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        ...

    async def patch(
        self,
        path: str,
        json_data: dict[str, Any],
        auth_token: str | None = None,
        request_id: str | None = None,
    ) -> dict[str, Any]:
        ...


def register_tools(mcp, api_client: BusinessAPIClient) -> None:
    @mcp.tool(
        name=TOOL_GET_BUSINESS_PROFILE,
        description="Retrieve authenticated business profile.",
        exclude_args=["auth_context"],
    )
    async def get_business_profile(
        business_id: int | None = None,
        auth_context: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        try:
            access = resolve_business_access(auth_context, business_id)
            data = await _api_get(
                api_client,
                API_PATH_BUSINESS_PROFILE.format(business_id=access["business_id"]),
                access,
            )
            return {
                "ok": True,
                "error": None,
                "validation_errors": [],
                "changed_fields": [],
                "business_id": access["business_id"],
                "result": data,
            }
        except (AccessError, ForbiddenError, RuntimeError) as e:
            return _tool_error(str(e))

    @mcp.tool(
        name=TOOL_UPDATE_BUSINESS_PROFILE,
        description="Update authenticated business profile.",
        exclude_args=["auth_context"],
    )
    async def update_business_profile(
        payload: dict[str, Any],
        business_id: int | None = None,
        auth_context: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        validation_errors = validate_profile_update(payload)
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        try:
            access = resolve_business_access(auth_context, business_id)
            before = await _api_get(
                api_client,
                API_PATH_BUSINESS_PROFILE.format(business_id=access["business_id"]),
                access,
            )
            after = await _api_patch(
                api_client,
                API_PATH_BUSINESS_PROFILE.format(business_id=access["business_id"]),
                payload,
                access,
            )

            fields = ("name", "avatar", "banner", "description")
            return {
                "ok": True,
                "error": None,
                "validation_errors": [],
                "changed_fields": changed_fields(before, after, fields),
                "business_id": access["business_id"],
                "result": after,
            }
        except (AccessError, ForbiddenError, RuntimeError) as e:
            return _tool_error(str(e))

    @mcp.tool(
        name=TOOL_LIST_BUSINESS_VENUES,
        description="List venues for authenticated business.",
        exclude_args=["auth_context"],
    )
    async def list_business_venues(
        business_id: int | None = None,
        skip: int = 0,
        limit: int = 10,
        auth_context: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        try:
            access = resolve_business_access(auth_context, business_id)
            data = await _api_get(
                api_client,
                API_PATH_BUSINESS_VENUES.format(business_id=access["business_id"]),
                access,
                params={"skip": max(skip, 0), "limit": max(1, min(limit, 100))},
            )
            return {
                "ok": True,
                "error": None,
                "validation_errors": [],
                "changed_fields": [],
                "business_id": access["business_id"],
                "result": data,
            }
        except (AccessError, ForbiddenError, RuntimeError) as e:
            return _tool_error(str(e))

    @mcp.tool(
        name=TOOL_GET_VENUE_DETAILS,
        description="Get venue details with ownership check.",
        exclude_args=["auth_context"],
    )
    async def get_venue_details(
        venue_id: int,
        business_id: int | None = None,
        auth_context: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        try:
            access = resolve_business_access(auth_context, business_id)
            venue = await _api_get(
                api_client,
                API_PATH_VENUE_DETAILS.format(venue_id=venue_id),
                access,
            )
            ensure_venue_owned_by_business(venue, access["business_id"])
            return {
                "ok": True,
                "error": None,
                "validation_errors": [],
                "changed_fields": [],
                "business_id": access["business_id"],
                "venue_id": venue_id,
                "result": venue,
            }
        except (AccessError, ForbiddenError, RuntimeError) as e:
            return _tool_error(str(e))

    @mcp.tool(
        name=TOOL_UPDATE_VENUE_DETAILS,
        description="Update venue details with ownership check.",
        exclude_args=["auth_context"],
    )
    async def update_venue_details(
        venue_id: int,
        payload: dict[str, Any],
        business_id: int | None = None,
        auth_context: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        validation_errors = validate_venue_update(payload)
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        try:
            access = resolve_business_access(auth_context, business_id)
            before = await _api_get(
                api_client,
                API_PATH_VENUE_DETAILS.format(venue_id=venue_id),
                access,
            )
            ensure_venue_owned_by_business(before, access["business_id"])

            after = await _api_patch(
                api_client,
                API_PATH_UPDATE_VENUE_DETAILS.format(venue_id=venue_id),
                payload,
                access,
            )

            fields = ("name", "avatar", "banner", "description", "latitude", "longitude", "tagIds")
            return {
                "ok": True,
                "error": None,
                "validation_errors": [],
                "changed_fields": changed_fields(before, after, fields),
                "business_id": access["business_id"],
                "venue_id": venue_id,
                "result": after,
            }
        except (AccessError, ForbiddenError, RuntimeError) as e:
            return _tool_error(str(e))

    @mcp.tool(
        name=TOOL_UPDATE_VENUE_HOURS,
        description="Update venue hours.",
        exclude_args=["auth_context"],
    )
    async def update_venue_hours(
        venue_id: int,
        payload: dict[str, Any],
        business_id: int | None = None,
        auth_context: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        validation_errors = validate_venue_hours(payload)
        if validation_errors:
            return _tool_error("validation failed", validation_errors=validation_errors)

        try:
            access = resolve_business_access(auth_context, business_id)

            before = await _api_get(
                api_client,
                API_PATH_VENUE_DETAILS.format(venue_id=venue_id),
                access,
            )
            ensure_venue_owned_by_business(before, access["business_id"])

            after = await _api_patch(
                api_client,
                API_PATH_UPDATE_VENUE_HOURS.format(venue_id=venue_id),
                payload,
                access,
            )

            return {
                "ok": True,
                "error": None,
                "validation_errors": [],
                "changed_fields": ["days"],
                "business_id": access["business_id"],
                "venue_id": venue_id,
                "result": after,
            }
        except (AccessError, ForbiddenError, RuntimeError) as e:
            return _tool_error(str(e))


async def _api_get(
    api_client: BusinessAPIClient,
    path: str,
    access: dict[str, Any],
    params: dict[str, Any] | None = None,
) -> dict[str, Any]:
    result = await api_client.get(
        path,
        auth_token=access["auth_token"],
        request_id=access.get("request_id"),
        params=params,
    )
    return _unwrap(result)


async def _api_patch(
    api_client: BusinessAPIClient,
    path: str,
    payload: dict[str, Any],
    access: dict[str, Any],
) -> dict[str, Any]:
    result = await api_client.patch(
        path,
        json_data=payload,
        auth_token=access["auth_token"],
        request_id=access.get("request_id"),
    )
    return _unwrap(result)


def _unwrap(result: dict[str, Any]) -> dict[str, Any]:
    if result.get("is_error") is True:
        message = result.get("error_message", "unknown business api error")
        raise RuntimeError(f"Business API error: {message}")
    data = result.get("data")
    if data is None:
        return {}
    if isinstance(data, dict):
        return data
    return {"value": data}


def _tool_error(
    message: str,
    validation_errors: list[dict[str, str]] | None = None,
    changed: list[str] | None = None,
) -> dict[str, Any]:
    return {
        "ok": False,
        "error": message,
        "validation_errors": validation_errors or [],
        "changed_fields": changed or [],
        "result": None,
    }