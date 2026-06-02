import pytest

from app.utils import (
    ForbiddenError,
    ensure_venue_owned_by_business,
    validate_profile_update,
    validate_venue_hours,
    validate_venue_update,
)


def test_validate_profile_update_unknown_field():
    errs = validate_profile_update({"badField": "x"})
    assert any(e["field"] == "badField" for e in errs)


def test_validate_venue_update_tag_ids_limit():
    errs = validate_venue_update({"tagIds": [1, 2, 3, 4, 5, 6]})
    assert any(e["field"] == "tagIds" for e in errs)


def test_validate_venue_hours_partial_pair_error():
    errs = validate_venue_hours({"days": [{"weekday": 1, "openTime": "09:00", "closeTime": None}]})
    assert any("both openTime and closeTime" in e["message"] for e in errs)


def test_validate_venue_hours_invalid_order_error():
    errs = validate_venue_hours({"days": [{"weekday": 1, "openTime": "18:00", "closeTime": "09:00"}]})
    assert any("openTime must be before closeTime" in e["message"] for e in errs)


def test_validate_venue_hours_closed_day_allowed():
    errs = validate_venue_hours({"days": [{"weekday": 7, "openTime": None, "closeTime": None}]})
    assert errs == []


def test_ensure_venue_owned_by_business_forbidden():
    with pytest.raises(ForbiddenError):
        ensure_venue_owned_by_business({"brand": {"id": 999}}, business_id=10)


def test_ensure_venue_owned_by_business_accepts_parent_id():
    ensure_venue_owned_by_business({"parentId": 10}, business_id=10)
