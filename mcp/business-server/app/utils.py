from typing import Any
from datetime import datetime

ValidationErrors = list[dict[str, str]]


class ForbiddenError(RuntimeError):
    pass


def validate_profile_update(payload: dict[str, Any]) -> ValidationErrors:
    allowed = {"name", "avatar", "banner", "description"}
    return _validate_update_payload(payload, allowed)


def validate_venue_update(payload: dict[str, Any]) -> ValidationErrors:
    allowed = {
        "name",
        "avatar",
        "banner",
        "description",
        "latitude",
        "longitude",
        "tagIds",
    }
    errors = _validate_update_payload(payload, allowed)

    lat = payload.get("latitude")
    if lat is not None and (not isinstance(lat, (int, float)) or lat < -90 or lat > 90):
        errors.append(
            {"field": "latitude", "message": "latitude must be between -90 and 90"}
        )

    lon = payload.get("longitude")
    if lon is not None and (
        not isinstance(lon, (int, float)) or lon < -180 or lon > 180
    ):
        errors.append(
            {"field": "longitude", "message": "longitude must be between -180 and 180"}
        )

    tag_ids = payload.get("tagIds")
    if tag_ids is not None:
        if not isinstance(tag_ids, list) or any(
            not isinstance(x, int) for x in tag_ids
        ):
            errors.append({"field": "tagIds", "message": "tagIds must be list[int]"})
        elif len(tag_ids) > 5:
            errors.append(
                {"field": "tagIds", "message": "location can have at most 5 tags"}
            )

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
            errors.append(
                {
                    "field": f"days[{idx}].weekday",
                    "message": "weekday must be integer 1..7",
                }
            )
            continue

        if weekday in seen:
            errors.append(
                {"field": f"days[{idx}].weekday", "message": "duplicate weekday"}
            )
        seen.add(weekday)

        open_time = day.get("openTime")
        close_time = day.get("closeTime")

        # closed day
        if open_time is None and close_time is None:
            continue

        # partial pair
        if open_time is None or close_time is None:
            errors.append(
                {
                    "field": f"days[{idx}]",
                    "message": "both openTime and closeTime must be provided together",
                }
            )
            continue

        # format checks
        if not isinstance(open_time, str) or not isinstance(close_time, str):
            errors.append(
                {
                    "field": f"days[{idx}]",
                    "message": "openTime and closeTime must be strings in HH:MM format",
                }
            )
            continue

        try:
            open_dt = datetime.strptime(open_time, "%H:%M")
        except ValueError:
            errors.append(
                {"field": f"days[{idx}].openTime", "message": "openTime must be HH:MM"}
            )
            continue

        try:
            close_dt = datetime.strptime(close_time, "%H:%M")
        except ValueError:
            errors.append(
                {
                    "field": f"days[{idx}].closeTime",
                    "message": "closeTime must be HH:MM",
                }
            )
            continue

        if not open_dt < close_dt:
            errors.append(
                {
                    "field": f"days[{idx}]",
                    "message": "openTime must be before closeTime",
                }
            )

    return errors


def changed_fields(
    before: dict[str, Any], after: dict[str, Any], fields: tuple[str, ...]
) -> list[str]:
    out: list[str] = []
    for field in fields:
        if before.get(field) != after.get(field):
            out.append(field)
    return out


def ensure_venue_owned_by_business(
    venue_data: dict[str, Any], business_id: int
) -> None:
    brand_id = None

    brand = venue_data.get("brand")
    if isinstance(brand, dict):
        brand_id = brand.get("id")

    brand_id = (
        brand_id
        or venue_data.get("brandId")
        or venue_data.get("brand_id")
        or venue_data.get("parentId")
        or venue_data.get("parent_id")
    )

    if brand_id is None:
        raise ForbiddenError("cannot verify venue ownership")

    if int(brand_id) != int(business_id):
        raise ForbiddenError("unauthorized access to another business venue")

def extract_venue_hours_days(venue_data: dict[str, Any]) -> list[dict[str, Any]] | None:
    for key in ("days", "hours", "venueHours", "venue_hours"):
        if key not in venue_data:
            continue
        value = venue_data.get(key)
        if isinstance(value, list):
            return value

        if isinstance(value, dict):
            nested_days = value.get("days")
            if isinstance(nested_days, list):
                return nested_days

    return None


def normalize_venue_hours_days(days: list[dict[str, Any]] | None) -> list[dict[str, Any]]:
    if not isinstance(days, list):
        return []

    normalized: list[dict[str, Any]] = []
    for day in days:
        if not isinstance(day, dict):
            continue

        weekday = day.get("weekday")
        if not isinstance(weekday, int):
            continue

        normalized.append(
            {
                "weekday": weekday,
                "openTime": day.get("openTime"),
                "closeTime": day.get("closeTime"),
            }
        )

    return sorted(normalized, key=lambda item: item["weekday"])


def build_venue_hours_preview(
    current_days: list[dict[str, Any]] | None,
    proposed_days: list[dict[str, Any]] | None,
) -> dict[str, Any]:
    current_norm = normalize_venue_hours_days(current_days)
    proposed_norm = normalize_venue_hours_days(proposed_days)

    current_map = {
        day["weekday"]: (day.get("openTime"), day.get("closeTime"))
        for day in current_norm
    }
    proposed_map = {
        day["weekday"]: (day.get("openTime"), day.get("closeTime"))
        for day in proposed_norm
    }

    counts = {
        "added": 0,
        "removed": 0,
        "updated": 0,
        "opened": 0,
        "closed": 0,
    }
    day_changes: list[dict[str, Any]] = []

    for weekday in sorted(set(current_map) | set(proposed_map)):
        before = current_map.get(weekday)
        after = proposed_map.get(weekday)

        if before == after:
            status = "unchanged"
        elif before is None:
            status = "added"
            counts["added"] += 1
        elif after is None:
            status = "removed"
            counts["removed"] += 1
        elif _is_closed_pair(before) and not _is_closed_pair(after):
            status = "opened"
            counts["opened"] += 1
        elif not _is_closed_pair(before) and _is_closed_pair(after):
            status = "closed"
            counts["closed"] += 1
        else:
            status = "updated"
            counts["updated"] += 1

        day_changes.append(
            {
                "weekday": weekday,
                "status": status,
                "before": _pair_to_dict(before),
                "after": _pair_to_dict(after),
            }
        )

    changed_fields = ["days"] if any(counts.values()) else []

    return {
        "current_days": current_norm,
        "preview_days": proposed_norm,
        "day_changes": day_changes,
        "summary": _build_venue_hours_summary(counts),
        "changed_fields": changed_fields,
    }


def _is_closed_pair(value: tuple[Any, Any] | None) -> bool:
    if value is None:
        return False
    return value[0] is None and value[1] is None


def _pair_to_dict(value: tuple[Any, Any] | None) -> dict[str, Any] | None:
    if value is None:
        return None

    return {
        "openTime": value[0],
        "closeTime": value[1],
    }


def _build_venue_hours_summary(counts: dict[str, int]) -> str:
    parts: list[str] = []

    if counts["added"]:
        parts.append(f'{counts["added"]} day(s) added')
    if counts["removed"]:
        parts.append(f'{counts["removed"]} day(s) removed')
    if counts["updated"]:
        parts.append(f'{counts["updated"]} day(s) updated')
    if counts["opened"]:
        parts.append(f'{counts["opened"]} day(s) opened')
    if counts["closed"]:
        parts.append(f'{counts["closed"]} day(s) closed')

    if not parts:
        return "No venue-hours changes detected."

    return ", ".join(parts)

def validate_food_box_creation(payload: dict[str, Any]) -> ValidationErrors:
    """Validate food box creation payload.
    
    Required fields: venue_id, price_full, expires_at, quantity, image_base64
    Optional fields: category_id, price_discount
    """
    required_fields = {"venue_id", "price_full", "expires_at", "quantity", "image_base64"}
    errors: ValidationErrors = []

    for field in required_fields:
        if field not in payload or payload[field] is None:
            errors.append({"field": field, "message": f"{field} is required"})

    allowed = required_fields | {"category_id", "price_discount"}
    for key in payload.keys():
        if key not in allowed:
            errors.append({"field": key, "message": "unknown field"})

    venue_id = payload.get("venue_id")
    if venue_id is not None:
        if not isinstance(venue_id, int) or venue_id <= 0:
            errors.append({"field": "venue_id", "message": "venue_id must be a positive integer"})

    category_id = payload.get("category_id")
    if category_id is not None:
        if not isinstance(category_id, int) or category_id <= 0:
            errors.append({"field": "category_id", "message": "category_id must be a positive integer"})

    full_price = payload.get("price_full")
    if full_price is not None:
        try:
            price_val = float(full_price)
            if price_val <= 0:
                errors.append({"field": "price_full", "message": "price_full must be a positive number"})
        except (TypeError, ValueError):
            errors.append({"field": "price_full", "message": "price_full must be a valid number"})

    discount_price = payload.get("price_discount")
    if discount_price is not None:
        try:
            price_val = float(discount_price)
            if price_val < 0:
                errors.append({"field": "price_discount", "message": "price_discount must be a non-negative number"})
        except (TypeError, ValueError):
            errors.append({"field": "price_discount", "message": "price_discount must be a valid number"})

    expires_at = payload.get("expires_at")
    if expires_at is not None:
        if not isinstance(expires_at, str):
            errors.append({"field": "expires_at", "message": "expires_at must be an ISO 8601 datetime string"})

    quantity = payload.get("quantity")
    if quantity is not None:
        if not isinstance(quantity, int) or quantity < 1 or quantity > 1000:
            errors.append({"field": "quantity", "message": "quantity must be an integer between 1 and 1000"})

    image_base64 = payload.get("image_base64")
    if image_base64 is not None:
        if not isinstance(image_base64, str) or len(image_base64) == 0:
            errors.append({"field": "image_base64", "message": "image_base64 must be a non-empty string"})
        elif len(image_base64) > 10_000_000:  # ~7.5MB max
            errors.append({"field": "image_base64", "message": "image_base64 is too large (max ~7.5MB)"})

    return errors


def validate_food_box_update(payload: dict[str, Any]) -> ValidationErrors:
    """Validate food box update payload.
    
    Optional fields: category_id, price_full, price_discount, expires_at
    """
    errors: ValidationErrors = []

    if not isinstance(payload, dict) or len(payload) == 0:
        return [{"field": "payload", "message": "payload is required"}]

    allowed = {"category_id", "price_full", "price_discount", "expires_at"}
    for key in payload.keys():
        if key not in allowed:
            errors.append({"field": key, "message": "unknown field"})

    category_id = payload.get("category_id")
    if category_id is not None:
        if not isinstance(category_id, int) or category_id <= 0:
            errors.append({"field": "category_id", "message": "category_id must be a positive integer"})

    full_price = payload.get("price_full")
    if full_price is not None:
        try:
            price_val = float(full_price)
            if price_val <= 0:
                errors.append({"field": "price_full", "message": "price_full must be a positive number"})
        except (TypeError, ValueError):
            errors.append({"field": "price_full", "message": "price_full must be a valid number"})

    discount_price = payload.get("price_discount")
    if discount_price is not None:
        try:
            price_val = float(discount_price)
            if price_val < 0:
                errors.append({"field": "price_discount", "message": "price_discount must be a non-negative number"})
        except (TypeError, ValueError):
            errors.append({"field": "price_discount", "message": "price_discount must be a valid number"})

    expires_at = payload.get("expires_at")
    if expires_at is not None:
        if not isinstance(expires_at, str):
            errors.append({"field": "expires_at", "message": "expires_at must be an ISO 8601 datetime string"})

    has_any = any(payload.get(k) is not None for k in allowed)
    if not has_any:
        errors.append({"field": "payload", "message": "at least one updatable field is required"})

    return errors

def _validate_update_payload(payload: dict[str, Any], allowed: set[str]) -> ValidationErrors:
    errors: ValidationErrors = []

    if not isinstance(payload, dict) or len(payload) == 0:
        return [{"field": "payload", "message": "payload is required"}]

    for key in payload.keys():
        if key not in allowed:
            errors.append({"field": key, "message": "unknown field"})

    has_any = any(payload.get(k) is not None for k in allowed)
    if not has_any:
        errors.append(
            {"field": "payload", "message": "at least one updatable field is required"}
        )

    return errors


def validate_date_range(
    start_date: str, end_date: str, max_days: int = 90
) -> str | None:
    """
    Validates that the start and end dates are in correct format,
    chronological order, and within the maximum allowed days.
    """
    try:
        start_obj = datetime.strptime(start_date, "%Y-%m-%d")
        end_obj = datetime.strptime(end_date, "%Y-%m-%d")
    except ValueError:
        return "Dates must be in YYYY-MM-DD format."

    if start_obj > end_obj:
        return "start_date cannot be after end_date."

    delta_days = (end_obj - start_obj).days
    if delta_days > max_days:
        return f"Date range exceeds maximum allowed period of {max_days} days."

    return None


def validate_discovery_coords(lat: Any, lon: Any) -> ValidationErrors:
    errors: ValidationErrors = []
    if (
        lat is None
        or isinstance(lat, bool)
        or not isinstance(lat, (int, float))
        or lat < -90
        or lat > 90
    ):
        errors.append(
            {"field": "lat", "message": "latitude must be a number between -90 and 90"}
        )
    if (
        lon is None
        or isinstance(lon, bool)
        or not isinstance(lon, (int, float))
        or lon < -180
        or lon > 180
    ):
        errors.append(
            {
                "field": "lon",
                "message": "longitude must be a number between -180 and 180",
            }
        )
    return errors


def validate_pagination(skip: Any, limit: Any) -> ValidationErrors:
    errors: ValidationErrors = []
    if isinstance(skip, bool) or not isinstance(skip, int) or skip < 0:
        errors.append(
            {"field": "skip", "message": "skip must be a non-negative integer"}
        )
    if (
        isinstance(limit, bool)
        or not isinstance(limit, int)
        or limit < 1
        or limit > 100
    ):
        errors.append(
            {"field": "limit", "message": "limit must be an integer between 1 and 100"}
        )
    return errors
