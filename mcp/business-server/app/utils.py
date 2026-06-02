from typing import Any
from datetime import datetime

ValidationErrors = list[dict[str, str]]

class ForbiddenError(RuntimeError):
    pass


def validate_profile_update(payload: dict[str, Any]) -> ValidationErrors:
    allowed = {"name", "avatar", "banner", "description"}
    return _validate_update_payload(payload, allowed)


def validate_venue_update(payload: dict[str, Any]) -> ValidationErrors:
    allowed = {"name", "avatar", "banner", "description", "latitude", "longitude", "tagIds"}
    errors = _validate_update_payload(payload, allowed)

    lat = payload.get("latitude")
    if lat is not None and (not isinstance(lat, (int, float)) or lat < -90 or lat > 90):
        errors.append({"field": "latitude", "message": "latitude must be between -90 and 90"})

    lon = payload.get("longitude")
    if lon is not None and (not isinstance(lon, (int, float)) or lon < -180 or lon > 180):
        errors.append({"field": "longitude", "message": "longitude must be between -180 and 180"})

    tag_ids = payload.get("tagIds")
    if tag_ids is not None:
        if not isinstance(tag_ids, list) or any(not isinstance(x, int) for x in tag_ids):
            errors.append({"field": "tagIds", "message": "tagIds must be list[int]"})
        elif len(tag_ids) > 5:
            errors.append({"field": "tagIds", "message": "location can have at most 5 tags"})

    return errors


def validate_venue_hours(payload: dict[str, Any]) -> ValidationErrors:
    errors: ValidationErrors = []

    days = payload.get("days")
    if not isinstance(days, list) or len(days) == 0:
        return [{"field": "days", "message": "days must be a non-empty list"}]

    seen: set[int] = set()
    for idx, day in enumerate(days):
        if not isinstance(day, dict):
            errors.append({"field": f"days[{idx}]", "message": "must be an object"})
            continue

        weekday = day.get("weekday")
        if not isinstance(weekday, int) or weekday < 1 or weekday > 7:
            errors.append({"field": f"days[{idx}].weekday", "message": "weekday must be integer 1..7"})
            continue

        if weekday in seen:
            errors.append({"field": f"days[{idx}].weekday", "message": "duplicate weekday"})
        seen.add(weekday)

        open_time = day.get("openTime")
        close_time = day.get("closeTime")

        # closed day
        if open_time is None and close_time is None:
            continue

        # partial pair
        if open_time is None or close_time is None:
            errors.append({"field": f"days[{idx}]", "message": "both openTime and closeTime must be provided together"})
            continue

        # format checks
        if not isinstance(open_time, str) or not isinstance(close_time, str):
            errors.append({"field": f"days[{idx}]", "message": "openTime and closeTime must be strings in HH:MM format"})
            continue

        try:
            open_dt = datetime.strptime(open_time, "%H:%M")
        except ValueError:
            errors.append({"field": f"days[{idx}].openTime", "message": "openTime must be HH:MM"})
            continue

        try:
            close_dt = datetime.strptime(close_time, "%H:%M")
        except ValueError:
            errors.append({"field": f"days[{idx}].closeTime", "message": "closeTime must be HH:MM"})
            continue

        if not open_dt < close_dt:
            errors.append({"field": f"days[{idx}]", "message": "openTime must be before closeTime"})

    return errors


def changed_fields(before: dict[str, Any], after: dict[str, Any], fields: tuple[str, ...]) -> list[str]:
    out: list[str] = []
    for field in fields:
        if before.get(field) != after.get(field):
            out.append(field)
    return out


def ensure_venue_owned_by_business(venue_data: dict[str, Any], business_id: int) -> None:
    brand_id = None

    brand = venue_data.get("brand")
    if isinstance(brand, dict):
        brand_id = brand.get("id")

    brand_id = brand_id or venue_data.get("brandId") or venue_data.get("brand_id") or venue_data.get("parentId") or venue_data.get("parent_id")

    if brand_id is None:
        raise ForbiddenError("cannot verify venue ownership")

    if int(brand_id) != int(business_id):
        raise ForbiddenError("unauthorized access to another business venue")


def _validate_update_payload(payload: dict[str, Any], allowed: set[str]) -> ValidationErrors:
    errors: ValidationErrors = []

    if not isinstance(payload, dict) or len(payload) == 0:
        return [{"field": "payload", "message": "payload is required"}]

    for key in payload.keys():
        if key not in allowed:
            errors.append({"field": key, "message": "unknown field"})

    has_any = any(payload.get(k) is not None for k in allowed)
    if not has_any:
        errors.append({"field": "payload", "message": "at least one updatable field is required"})

    return errors
